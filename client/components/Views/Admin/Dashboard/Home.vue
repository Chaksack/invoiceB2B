<template>
  <main class="flex-1 space-y-4 p-4 md:p-8 pt-6 bg-gray-100 font-inter">
    <Toaster richColors position="top-right" />
    <div class="flex items-center justify-between space-y-2">
      <h2 class="text-3xl font-bold tracking-tight text-gray-800">Admin Dashboard</h2>
      <div class="flex items-center space-x-2">
        <Button variant="outline" size="sm">
          <Download class="mr-2 h-4 w-4" />
          Download Report
        </Button>
      </div>
    </div>

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Total Disbursed
          </CardTitle>
          <DollarSign class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div v-if="isLoadingStats" class="h-8 w-3/4 bg-gray-200 animate-pulse rounded"></div>
          <div v-else class="text-2xl font-bold">GHS {{ stats.totalRevenue?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00' }}</div>
          <p v-if="!isLoadingStats" class="text-xs text-muted-foreground">
            {{ stats.revenueChange || '+0.0' }}% from last month
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">
            Total Users
          </CardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div v-if="isLoadingStats" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
          <div v-else class="text-2xl font-bold">{{ stats.totalUsers?.toLocaleString() || '0' }}</div>
          <p v-if="!isLoadingStats" class="text-xs text-muted-foreground">
            {{ stats.usersChange || '+0' }} new users this month
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Invoices</CardTitle>
          <FileText class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div v-if="isLoadingStats" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
          <div v-else class="text-2xl font-bold">{{ stats.totalInvoices?.toLocaleString() || '0' }}</div>
          <p v-if="!isLoadingStats" class="text-xs text-muted-foreground">
            {{ stats.invoicesChange || '+0' }} from last month
          </p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Pending KYC</CardTitle>
          <AlertCircle class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div v-if="isLoadingStats" class="h-8 w-1/4 bg-gray-200 animate-pulse rounded"></div>
          <div v-else class="text-2xl font-bold">{{ stats.pendingKyc?.toLocaleString() || '0' }}</div>
          <p v-if="!isLoadingStats" class="text-xs text-muted-foreground">
            Users awaiting KYC approval
          </p>
        </CardContent>
      </Card>
    </div>

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
      <Card class="col-span-4">
        <CardHeader>
          <CardTitle>Invoice Overview</CardTitle>
          <CardDescription>Monthly invoice volume and status distribution.</CardDescription>
        </CardHeader>
        <CardContent class="pl-2">
          <div v-if="isLoadingCharts" class="h-[350px] w-full bg-gray-200 animate-pulse rounded flex items-center justify-center text-gray-500">
            Loading chart data...
          </div>
          <div v-else class="h-[350px] w-full bg-gray-50 rounded flex items-center justify-center text-gray-500 border">
            Chart: Invoice Volume & Status (e.g., using Recharts)
          </div>
        </CardContent>
      </Card>
      <Card class="col-span-4 lg:col-span-3">
        <CardHeader>
          <CardTitle>User Registrations</CardTitle>
          <CardDescription>New user sign-ups over the past 6 months.</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="isLoadingCharts" class="h-[350px] w-full bg-gray-200 animate-pulse rounded flex items-center justify-center text-gray-500">
            Loading chart data...
          </div>
          <div v-else class="h-[350px] w-full bg-gray-50 rounded flex items-center justify-center text-gray-500 border">
            Chart: User Registrations Trend (e.g., using Recharts)
          </div>
        </CardContent>
      </Card>
    </div>

    <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card class="lg:col-span-2">
        <CardHeader>
          <CardTitle>Recent Activities</CardTitle>
          <CardDescription>Latest user and invoice activities.</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="isLoadingActivities" class="space-y-4">
            <div v-for="i in 5" :key="i" class="flex items-center space-x-3 p-2 bg-gray-200 animate-pulse rounded h-12"></div>
          </div>
          <div v-else-if="recentActivities.length === 0" class="text-sm text-gray-500">No recent activities.</div>
          <div v-else class="space-y-4 max-h-96 overflow-y-auto">
            <div v-for="(activity, index) in recentActivities" :key="index" class="flex items-center space-x-3 p-2 hover:bg-gray-50 rounded-md transition-colors">
              <Avatar class="h-9 w-9">
                <AvatarImage :src="activity.avatarUrl" :alt="activity.actorName" />
                <AvatarFallback :class="getActivityIconBg(activity.type)">
                  <component :is="getActivityIcon(activity.type)" class="h-4 w-4 text-white" />
                </AvatarFallback>
              </Avatar>
              <div class="flex-1 text-sm">
                <p class="font-medium text-gray-800">{{ activity.description }}</p>
                <p class="text-xs text-gray-500">
                  <span v-if="activity.actorName">{{ activity.actorName }} &bull; </span>
                  {{ formatRelativeTime(activity.timestamp) }}
                </p>
              </div>
              <Badge v-if="activity.tag" :variant="getActivityTagVariant(activity.tag)">{{ activity.tag }}</Badge>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Quick Links</CardTitle>
          <CardDescription>Frequently accessed admin sections.</CardDescription>
        </CardHeader>
        <CardContent class="grid gap-2">
          <NuxtLink to="/admin/users">
            <Button variant="outline" class="w-full justify-start">
              <Users class="mr-2 h-4 w-4" /> Manage Users
            </Button>
          </NuxtLink>
          <NuxtLink to="/admin/invoices">
            <Button variant="outline" class="w-full justify-start">
              <FileText class="mr-2 h-4 w-4" /> Manage Invoices
            </Button>
          </NuxtLink>
          <NuxtLink to="/admin/settings">
            <Button variant="outline" class="w-full justify-start">
              <Settings class="mr-2 h-4 w-4" /> System Settings
            </Button>
          </NuxtLink>
          <NuxtLink to="/admin/reports">
            <Button variant="outline" class="w-full justify-start">
              <BarChart class="mr-2 h-4 w-4" /> View Reports
            </Button>
          </NuxtLink>
        </CardContent>
      </Card>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { Download, DollarSign, Users, FileText, AlertCircle, Settings, BarChart, UserPlus, FilePlus, CheckCircle2, XCircle } from 'lucide-vue-next';
import { Button } from '~/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '~/components/ui/card';
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar';
import { Badge } from '~/components/ui/badge';
import { Toaster, toast } from 'vue-sonner';
import { useCookie } from '#app';
import axios from 'axios';

definePageMeta({
  layout: 'admin',
  middleware: 'auth'
});

const API_BASE_URL = 'http://localhost:3000/api/v1'; // Replace with your actual API base URL
const tokenCookie = useCookie('token');
const authToken = tokenCookie.value || null;

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    ...(authToken ? { Authorization: `Bearer ${authToken}` } : {}),
    'Content-Type': 'application/json',
  },
});

const isLoadingStats = ref(true);
const isLoadingCharts = ref(true);
const isLoadingActivities = ref(true);

const stats = ref({
  totalRevenue: 0,      // Will be populated by API: totalDisbursedAmount
  revenueChange: '+0.0',// Mocked for now
  totalUsers: 0,        // Will be populated by API
  usersChange: '+0',    // Mocked for now
  totalInvoices: 0,     // Will be populated by API
  invoicesChange: '+0', // Mocked for now
  pendingKyc: 0,        // Will be populated by API
});

const recentActivities = ref<any[]>([]);

const fetchDashboardStats = async () => {
  isLoadingStats.value = true;
  try {
    const response = await apiClient.get('/admin/dashboard/analytics'); // Updated endpoint
    const data = response.data;

    stats.value = {
      totalRevenue: data.totalDisbursedAmount || 0, // Using totalDisbursedAmount for totalRevenue
      revenueChange: stats.value.revenueChange, // Keep mock or implement calculation
      totalUsers: data.totalUsers || 0,
      usersChange: stats.value.usersChange, // Keep mock or implement calculation
      totalInvoices: data.invoiceStats?.totalInvoices || 0,
      invoicesChange: stats.value.invoicesChange, // Keep mock or implement calculation
      pendingKyc: data.kycStats?.pendingReview || 0,
    };
    toast.success("Dashboard statistics loaded successfully!");

  } catch (error) {
    console.error("Failed to fetch dashboard stats:", error);
    let errorMessage = "Could not load dashboard statistics.";
    if (axios.isAxiosError(error) && error.response) {
      errorMessage += ` (Status: ${error.response.status})`;
    }
    toast.error(errorMessage);
    // Fallback to some default/zero values if API fails
    stats.value = {
      totalRevenue: 0,
      revenueChange: 'N/A',
      totalUsers: 0,
      usersChange: 'N/A',
      totalInvoices: 0,
      invoicesChange: 'N/A',
      pendingKyc: 0,
    };
  } finally {
    isLoadingStats.value = false;
  }
};

const fetchChartData = async () => {
  isLoadingCharts.value = true;
  try {
    await new Promise(resolve => setTimeout(resolve, 1500));
    toast.info("Chart data would be loaded here for Recharts or similar library.");
  } catch (error) {
    console.error("Failed to fetch chart data:", error);
    toast.error("Could not load chart data.");
  } finally {
    isLoadingCharts.value = false;
  }
}

const fetchRecentActivities = async () => {
  isLoadingActivities.value = true;
  try {
    await new Promise(resolve => setTimeout(resolve, 1200));
    recentActivities.value = [
      { type: 'new_user', description: 'New user registered: Alice Wonderland', actorName: 'System', timestamp: new Date(Date.now() - 3600000 * 2), avatarUrl: 'https://placehold.co/40x40/7c3aed/ffffff?text=AW', tag: 'User' },
      { type: 'invoice_submitted', description: 'Invoice INV-2025-001 submitted', actorName: 'Bob The Builder', timestamp: new Date(Date.now() - 3600000 * 5), avatarUrl: 'https://placehold.co/40x40/f59e0b/ffffff?text=BB', tag: 'Invoice' },
      { type: 'kyc_approved', description: 'KYC approved for Charlie Brown', actorName: 'Admin User', timestamp: new Date(Date.now() - 3600000 * 24), avatarUrl: 'https://placehold.co/40x40/10b981/ffffff?text=CB', tag: 'KYC' },
      { type: 'invoice_paid', description: 'Invoice INV-2024-980 marked as paid', actorName: 'Finance Bot', timestamp: new Date(Date.now() - 3600000 * 48), avatarUrl: 'https://placehold.co/40x40/3b82f6/ffffff?text=FB', tag: 'Payment' },
      { type: 'kyc_rejected', description: 'KYC rejected for Diana Prince', actorName: 'Admin User', timestamp: new Date(Date.now() - 3600000 * 72), avatarUrl: 'https://placehold.co/40x40/ef4444/ffffff?text=DP', tag: 'KYC' },
    ];
  } catch (error) {
    console.error("Failed to fetch recent activities:", error);
    toast.error("Could not load recent activities.");
  } finally {
    isLoadingActivities.value = false;
  }
};

const formatRelativeTime = (timestamp: string | Date): string => {
  const date = new Date(timestamp);
  const now = new Date();
  const seconds = Math.round((now.getTime() - date.getTime()) / 1000);
  const minutes = Math.round(seconds / 60);
  const hours = Math.round(minutes / 60);
  const days = Math.round(hours / 24);

  if (seconds < 60) return `${seconds} sec ago`;
  if (minutes < 60) return `${minutes} min ago`;
  if (hours < 24) return `${hours} hr ago`;
  return `${days} day(s) ago`;
};

const getActivityIcon = (type: string) => {
  switch (type) {
    case 'new_user': return UserPlus;
    case 'invoice_submitted': return FilePlus;
    case 'kyc_approved': return CheckCircle2;
    case 'invoice_paid': return CheckCircle2;
    case 'kyc_rejected': return XCircle;
    default: return AlertCircle;
  }
};
const getActivityIconBg = (type: string) => {
  switch (type) {
    case 'new_user': return 'bg-purple-500';
    case 'invoice_submitted': return 'bg-amber-500';
    case 'kyc_approved': return 'bg-green-500';
    case 'invoice_paid': return 'bg-blue-500';
    case 'kyc_rejected': return 'bg-red-500';
    default: return 'bg-gray-400';
  }
};

const getActivityTagVariant = (tag: string): 'default' | 'secondary' | 'destructive' | 'outline' => {
  if (tag === 'User') return 'secondary';
  if (tag === 'Invoice') return 'outline';
  if (tag === 'KYC') return 'default';
  if (tag === 'Payment') return 'default';
  return 'outline';
}

onMounted(() => {
  if (authToken) {
    fetchDashboardStats();
    fetchChartData();
    fetchRecentActivities();
  } else {
    toast.error("Authentication token not found. Please log in.");
    isLoadingStats.value = false;
    isLoadingCharts.value = false;
    isLoadingActivities.value = false;
  }
});
</script>

<style scoped>
/* Add any page-specific styles here */
</style>
