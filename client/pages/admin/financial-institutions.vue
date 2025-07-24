<template>
  <div>
    <ViewsAdminDashboardHeader />
    <main class="flex-1 space-y-4 p-4 md:p-8 pt-6 bg-gray-100 font-inter">
      <Toaster richColors position="top-right" />
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight text-gray-800">Financial Institutions</h2>
        <div class="flex items-center space-x-2">
          <Button @click="openCreateModal" variant="default" size="sm">
            <PlusCircle class="mr-2 h-4 w-4" />
            Add Institution
          </Button>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">
              Total Institutions
            </CardTitle>
            <Building2 class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ totalInstitutions }}</div>
            <p class="text-xs text-muted-foreground">
              Financial partners on the platform
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">
              Active Institutions
            </CardTitle>
            <CheckCircle2 class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ activeInstitutions }}</div>
            <p class="text-xs text-muted-foreground">
              Currently active partners
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Total Products</CardTitle>
            <Package class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ totalProducts }}</div>
            <p class="text-xs text-muted-foreground">
              Financing products available
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Total Terms</CardTitle>
            <FileText class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ totalTerms }}</div>
            <p class="text-xs text-muted-foreground">
              Financing terms configured
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- Search and Filter -->
      <div class="flex flex-col sm:flex-row gap-4 items-center justify-between">
        <div class="relative w-full sm:w-64">
          <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            v-model="searchQuery"
            placeholder="Search institutions..."
            class="pl-8 w-full"
            @input="handleSearch"
          />
        </div>
        <div class="flex gap-2 w-full sm:w-auto">
          <Select v-model="statusFilter" @update:modelValue="fetchInstitutions">
            <SelectTrigger class="w-full sm:w-[180px]">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">All Statuses</SelectItem>
              <SelectItem value="true">Active</SelectItem>
              <SelectItem value="false">Inactive</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <!-- Institutions Table -->
      <Card>
        <CardHeader>
          <CardTitle>Financial Institutions</CardTitle>
          <CardDescription>Manage your financial institution partners.</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="isLoading" class="space-y-4">
            <div v-for="i in 5" :key="i" class="h-12 bg-gray-200 animate-pulse rounded"></div>
          </div>
          <div v-else-if="institutions.length === 0" class="text-center py-8 text-gray-500">
            No financial institutions found.
          </div>
          <div v-else>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Code</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Products</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead class="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="institution in institutions" :key="institution.id">
                  <TableCell class="font-medium">{{ institution.name }}</TableCell>
                  <TableCell>{{ institution.code }}</TableCell>
                  <TableCell>
                    <Badge :variant="institution.isActive ? 'default' : 'secondary'">
                      {{ institution.isActive ? 'Active' : 'Inactive' }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ institution.productCount || 0 }}</TableCell>
                  <TableCell>{{ formatDate(institution.createdAt) }}</TableCell>
                  <TableCell class="text-right">
                    <Button variant="ghost" size="icon" @click="viewInstitution(institution)">
                      <Eye class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" @click="editInstitution(institution)">
                      <Edit class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" @click="confirmDelete(institution)">
                      <Trash2 class="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
            
            <!-- Pagination -->
            <div class="flex items-center justify-end space-x-2 py-4">
              <Button
                variant="outline"
                size="sm"
                :disabled="currentPage <= 1"
                @click="changePage(currentPage - 1)"
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                :disabled="currentPage >= totalPages"
                @click="changePage(currentPage + 1)"
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Create/Edit Institution Modal -->
      <Dialog v-model:open="showModal">
        <DialogContent class="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{{ isEditing ? 'Edit Institution' : 'Create Institution' }}</DialogTitle>
            <DialogDescription>
              {{ isEditing ? 'Update the details of this financial institution.' : 'Add a new financial institution to the platform.' }}
            </DialogDescription>
          </DialogHeader>
          <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="name" class="text-right">Name</Label>
              <Input id="name" v-model="formData.name" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="code" class="text-right">Code</Label>
              <Input id="code" v-model="formData.code" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="description" class="text-right">Description</Label>
              <Input id="description" v-model="formData.description" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="status" class="text-right">Status</Label>
              <Select v-model="formData.isActive" class="col-span-3">
                <SelectTrigger>
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="true">Active</SelectItem>
                  <SelectItem :value="false">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" @click="showModal = false">Cancel</Button>
            <Button @click="saveInstitution">Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <!-- Delete Confirmation Dialog -->
      <Dialog v-model:open="showDeleteDialog">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm Deletion</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {{ selectedInstitution?.name }}? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" @click="showDeleteDialog = false">Cancel</Button>
            <Button variant="destructive" @click="deleteInstitution">Delete</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { 
  Building2, CheckCircle2, Package, FileText, PlusCircle, 
  Search, Eye, Edit, Trash2 
} from 'lucide-vue-next';
import { Button } from '~/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import { Badge } from '~/components/ui/badge';
import { 
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow 
} from '~/components/ui/table';
import { 
  Dialog, DialogContent, DialogDescription, DialogFooter, 
  DialogHeader, DialogTitle 
} from '~/components/ui/dialog';
import { Label } from '~/components/ui/label';
import { 
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue 
} from '~/components/ui/select';
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

// State
const isLoading = ref(true);
const institutions = ref<any[]>([]);
const totalInstitutions = ref(0);
const activeInstitutions = ref(0);
const totalProducts = ref(0);
const totalTerms = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const searchQuery = ref('');
const statusFilter = ref('');
const showModal = ref(false);
const isEditing = ref(false);
const selectedInstitution = ref<any>(null);
const showDeleteDialog = ref(false);

const formData = ref({
  name: '',
  code: '',
  description: '',
  isActive: true
});

// Fetch financial institutions
const fetchInstitutions = async () => {
  isLoading.value = true;
  try {
    let url = `/admin/financial-institutions?page=${currentPage.value}&pageSize=${pageSize.value}`;
    
    if (searchQuery.value) {
      url += `&name=${searchQuery.value}`;
    }
    
    if (statusFilter.value) {
      url += `&is_active=${statusFilter.value}`;
    }
    
    const response = await apiClient.get(url);
    institutions.value = response.data.financialInstitutions || [];
    totalInstitutions.value = response.data.total || 0;
    totalPages.value = Math.ceil(totalInstitutions.value / pageSize.value);
    
    // Count active institutions
    activeInstitutions.value = institutions.value.filter(inst => inst.isActive).length;
    
    // Fetch additional stats
    await fetchStats();
    
  } catch (error) {
    console.error("Failed to fetch financial institutions:", error);
    toast.error("Could not load financial institutions.");
  } finally {
    isLoading.value = false;
  }
};

// Fetch additional statistics
const fetchStats = async () => {
  try {
    // In a real implementation, you would fetch these from the API
    // For now, we'll use mock data
    totalProducts.value = 24;
    totalTerms.value = 48;
  } catch (error) {
    console.error("Failed to fetch stats:", error);
  }
};

// Handle search
const handleSearch = () => {
  currentPage.value = 1; // Reset to first page when searching
  fetchInstitutions();
};

// Pagination
const changePage = (page: number) => {
  currentPage.value = page;
  fetchInstitutions();
};

// Format date
const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  });
};

// Modal functions
const openCreateModal = () => {
  isEditing.value = false;
  formData.value = {
    name: '',
    code: '',
    description: '',
    isActive: true
  };
  showModal.value = true;
};

const editInstitution = (institution: any) => {
  isEditing.value = true;
  selectedInstitution.value = institution;
  formData.value = {
    name: institution.name,
    code: institution.code,
    description: institution.description || '',
    isActive: institution.isActive
  };
  showModal.value = true;
};

const viewInstitution = (institution: any) => {
  // In a real implementation, you would navigate to a detail page
  // For now, we'll just show a toast
  toast.info(`Viewing details for ${institution.name}`);
};

const saveInstitution = async () => {
  try {
    if (isEditing.value && selectedInstitution.value) {
      // Update existing institution
      await apiClient.put(`/admin/financial-institutions/${selectedInstitution.value.id}`, formData.value);
      toast.success(`${formData.value.name} updated successfully!`);
    } else {
      // Create new institution
      await apiClient.post('/admin/financial-institutions', formData.value);
      toast.success(`${formData.value.name} created successfully!`);
    }
    showModal.value = false;
    fetchInstitutions();
  } catch (error) {
    console.error("Failed to save financial institution:", error);
    let errorMessage = "Could not save financial institution.";
    if (axios.isAxiosError(error) && error.response) {
      errorMessage += ` (Status: ${error.response.status})`;
    }
    toast.error(errorMessage);
  }
};

// Delete functions
const confirmDelete = (institution: any) => {
  selectedInstitution.value = institution;
  showDeleteDialog.value = true;
};

const deleteInstitution = async () => {
  if (!selectedInstitution.value) return;
  
  try {
    await apiClient.delete(`/admin/financial-institutions/${selectedInstitution.value.id}`);
    toast.success(`${selectedInstitution.value.name} deleted successfully!`);
    showDeleteDialog.value = false;
    fetchInstitutions();
  } catch (error) {
    console.error("Failed to delete financial institution:", error);
    let errorMessage = "Could not delete financial institution.";
    if (axios.isAxiosError(error) && error.response) {
      errorMessage += ` (Status: ${error.response.status})`;
    }
    toast.error(errorMessage);
  }
};

// Initialize
onMounted(() => {
  if (authToken) {
    fetchInstitutions();
  } else {
    toast.error("Authentication token not found. Please log in.");
    isLoading.value = false;
  }
});
</script>