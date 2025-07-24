<template>
  <header class="sticky top-0 flex h-12 items-center justify-between text-white bg-indigo-500 px-4 md:px-6">
    <div class="flex items-center">
      <Button variant="ghost" size="icon" class="mr-2 md:hidden text-white" @click="toggleMobileMenu">
        <Menu class="h-5 w-5" />
      </Button>
      <h1 class="text-xl font-semibold font-mono">
        InvoiceB2B
      </h1>
    </div>

    <!-- Desktop Navigation -->
    <NavigationMenu class="hidden md:block">
      <NavigationMenuList>
        <LazyNuxtLink to="/admin/home">
        <NavigationMenuItem class="px-2 text-sm" >Dashboard
        </NavigationMenuItem>
        </LazyNuxtLink>
        <LazyNuxtLink to="/admin/invoices">
          <NavigationMenuItem class="px-2 text-sm" >Invoices
          </NavigationMenuItem>
        </LazyNuxtLink>
        <LazyNuxtLink to="/admin/customers">
          <NavigationMenuItem class="px-2 text-sm" >Customers
          </NavigationMenuItem>
        </LazyNuxtLink>
        <LazyNuxtLink to="/admin/financial-institutions">
          <NavigationMenuItem class="px-2 text-sm" >Financial Institutions
          </NavigationMenuItem>
        </LazyNuxtLink>
        <LazyNuxtLink to="">
        <NavigationMenuItem class="px-2 text-sm">Profile
        </NavigationMenuItem>
        </LazyNuxtLink>
      </NavigationMenuList>
    </NavigationMenu>

    <!-- Mobile Menu -->
    <Sheet v-model:open="mobileMenuOpen">
      <SheetContent side="left" class="w-[250px] sm:w-[300px] bg-indigo-500 text-white">
        <SheetHeader>
          <SheetTitle class="text-white">Menu</SheetTitle>
        </SheetHeader>
        <div class="flex flex-col space-y-4 mt-6">
          <NuxtLink to="/admin/home" class="px-2 py-2 hover:bg-indigo-600 rounded-md" @click="mobileMenuOpen = false">
            Dashboard
          </NuxtLink>
          <NuxtLink to="/admin/invoices" class="px-2 py-2 hover:bg-indigo-600 rounded-md" @click="mobileMenuOpen = false">
            Invoices
          </NuxtLink>
          <NuxtLink to="/admin/customers" class="px-2 py-2 hover:bg-indigo-600 rounded-md" @click="mobileMenuOpen = false">
            Customers
          </NuxtLink>
          <NuxtLink to="/admin/financial-institutions" class="px-2 py-2 hover:bg-indigo-600 rounded-md" @click="mobileMenuOpen = false">
            Financial Institutions
          </NuxtLink>
          <NuxtLink to="" class="px-2 py-2 hover:bg-indigo-600 rounded-md" @click="mobileMenuOpen = false">
            Profile
          </NuxtLink>
        </div>
      </SheetContent>
    </Sheet>
    <div class="flex items-center">
      <div class="relative mr-4">
        <Button variant="ghost" size="icon" class="relative text-white" @click="toggleNotifications">
          <Bell class="h-5 w-5" />
          <span v-if="unreadNotifications.length > 0" class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
            {{ unreadNotifications.length > 9 ? '9+' : unreadNotifications.length }}
          </span>
        </Button>

        <div v-if="showNotifications" class="absolute right-0 mt-2 w-80 bg-white rounded-md shadow-lg overflow-hidden z-50">
          <div class="p-2 bg-indigo-500 text-white flex justify-between items-center">
            <h3 class="text-sm font-medium">Notifications</h3>
            <Button variant="ghost" size="sm" class="text-white h-6 px-2" @click="markAllAsRead">
              Mark all as read
            </Button>
          </div>
          <div class="max-h-96 overflow-y-auto">
            <div v-if="notifications.length === 0" class="p-4 text-center text-gray-500">
              No notifications
            </div>
            <div v-else>
              <div v-for="(notification, index) in notifications" :key="index" 
                   class="p-3 border-b last:border-b-0 hover:bg-gray-50 transition-colors"
                   :class="{ 'bg-blue-50': !notification.read }">
                <div class="flex items-start">
                  <div class="flex-shrink-0 mr-3">
                    <div :class="getNotificationIconBg(notification.type)" class="h-8 w-8 rounded-full flex items-center justify-center">
                      <component :is="getNotificationIcon(notification.type)" class="h-4 w-4 text-white" />
                    </div>
                  </div>
                  <div class="flex-1">
                    <p class="text-sm font-medium text-gray-900">{{ notification.title }}</p>
                    <p class="text-xs text-gray-500">{{ notification.message }}</p>
                    <p class="text-xs text-gray-400 mt-1">{{ formatRelativeTime(notification.timestamp) }}</p>
                  </div>
                  <Button v-if="!notification.read" variant="ghost" size="sm" class="ml-2 h-6 w-6 p-0" 
                          @click.stop="markAsRead(index)">
                    <Check class="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
          <div class="p-2 bg-gray-50 text-center">
            <NuxtLink to="/admin/notifications" class="text-xs text-indigo-600 hover:text-indigo-800">
              View all notifications
            </NuxtLink>
          </div>
        </div>
      </div>

      <p class="pr-2">{{ user.fullName || 'User' }}</p>
      <Avatar class="relative overflow-visible">
        <AvatarFallback class="text-black"> {{ user.initials || 'U' }} </AvatarFallback>
      </Avatar>
    </div>
    <Toaster richColors position="top-right" />
  </header>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { Avatar, AvatarFallback } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "~/components/ui/sheet";
import { Bell, Check, UserPlus, FilePlus, CheckCircle2, XCircle, AlertCircle, Menu } from 'lucide-vue-next';
import { Toaster, toast } from 'vue-sonner';
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuIndicator,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  NavigationMenuViewport,
} from '~/components/ui/navigation-menu';
import axios from 'axios';
import { useCookie } from '#app';
import { createApiClient, userApi } from '~/lib/api';

const tokenCookie = useCookie('token');
const authToken = tokenCookie.value || null;
const apiClient = createApiClient(authToken);

// User state
const user = ref({
  fullName: '',
  email: '',
  initials: ''
});

// Fetch user profile
const fetchUserProfile = async () => {
  if (!authToken) return;

  try {
    const response = await userApi.getProfile(apiClient);
    const profile = response.data;

    if (profile) {
      // Set user data
      user.value.fullName = profile.fullName || profile.name || '';
      user.value.email = profile.email || '';

      // Generate initials from full name
      if (user.value.fullName) {
        const nameParts = user.value.fullName.split(' ');
        user.value.initials = nameParts.length > 1 
          ? (nameParts[0][0] + nameParts[nameParts.length - 1][0]).toUpperCase()
          : nameParts[0][0].toUpperCase();
      }
    }
  } catch (error) {
    console.error('Failed to fetch user profile:', error);
    // Set fallback data
    user.value = {
      fullName: 'Andrew Chakdahah',
      email: 'andrew@example.com',
      initials: 'AC'
    };
  }
};

const mobileMenuOpen = ref(false);
const showNotifications = ref(false);

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value;
};
const notifications = ref([
  {
    id: 1,
    type: 'new_user',
    title: 'New User Registration',
    message: 'Alice Wonderland has registered and submitted KYC documents.',
    timestamp: new Date(Date.now() - 3600000 * 2),
    read: false
  },
  {
    id: 2,
    type: 'invoice_submitted',
    title: 'New Invoice Submitted',
    message: 'Bob Builder has submitted invoice INV-2025-001 for review.',
    timestamp: new Date(Date.now() - 3600000 * 5),
    read: false
  },
  {
    id: 3,
    type: 'kyc_approved',
    title: 'KYC Approved',
    message: 'Charlie Brown\'s KYC documents have been approved.',
    timestamp: new Date(Date.now() - 3600000 * 24),
    read: true
  }
]);

const unreadNotifications = computed(() => {
  return notifications.value.filter(notification => !notification.read);
});

const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value;
};

const markAsRead = (index) => {
  notifications.value[index].read = true;
};

const markAllAsRead = () => {
  notifications.value.forEach(notification => {
    notification.read = true;
  });
};

const formatRelativeTime = (timestamp) => {
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

const getNotificationIcon = (type) => {
  switch (type) {
    case 'new_user': return UserPlus;
    case 'invoice_submitted': return FilePlus;
    case 'kyc_approved': return CheckCircle2;
    case 'invoice_paid': return CheckCircle2;
    case 'kyc_rejected': return XCircle;
    default: return AlertCircle;
  }
};

const getNotificationIconBg = (type) => {
  switch (type) {
    case 'new_user': return 'bg-purple-500';
    case 'invoice_submitted': return 'bg-amber-500';
    case 'kyc_approved': return 'bg-green-500';
    case 'invoice_paid': return 'bg-blue-500';
    case 'kyc_rejected': return 'bg-red-500';
    default: return 'bg-gray-400';
  }
};

// Function to fetch notifications from the backend
const fetchNotifications = async () => {
  try {
    // This would be replaced with an actual API call in production
    // const response = await apiClient.get('/admin/notifications');
    // notifications.value = response.data;

    // For now, we'll use the mock data
    console.log('Fetched notifications');
  } catch (error) {
    console.error('Failed to fetch notifications:', error);
    toast.error('Could not load notifications');
  }
};

// Function to set up WebSocket connection for real-time notifications
const setupWebSocket = () => {
  // This would be implemented with a real WebSocket connection in production
  // const ws = new WebSocket('ws://localhost:3000/ws');

  // ws.onmessage = (event) => {
  //   const data = JSON.parse(event.data);
  //   if (data.type === 'notification') {
  //     notifications.value.unshift({
  //       id: data.id,
  //       type: data.notificationType,
  //       title: data.title,
  //       message: data.message,
  //       timestamp: new Date(),
  //       read: false
  //     });
  //     toast.info(data.title);
  //   }
  // };

  // return ws;

  // For now, we'll simulate a new notification every 30 seconds
  const interval = setInterval(() => {
    const types = ['new_user', 'invoice_submitted', 'kyc_approved', 'invoice_paid', 'kyc_rejected'];
    const type = types[Math.floor(Math.random() * types.length)];
    const titles = {
      'new_user': 'New User Registration',
      'invoice_submitted': 'New Invoice Submitted',
      'kyc_approved': 'KYC Approved',
      'invoice_paid': 'Invoice Paid',
      'kyc_rejected': 'KYC Rejected'
    };

    const newNotification = {
      id: Date.now(),
      type: type,
      title: titles[type],
      message: `Simulated ${type.replace('_', ' ')} notification.`,
      timestamp: new Date(),
      read: false
    };

    notifications.value.unshift(newNotification);
    toast.info(newNotification.title);
  }, 30000);

  return interval;
};

let wsConnection;

onMounted(() => {
  fetchUserProfile();
  fetchNotifications();
  wsConnection = setupWebSocket();

  // Close notifications when clicking outside
  const handleClickOutside = (event) => {
    const notificationContainer = document.querySelector('.notifications-container');
    if (showNotifications.value && notificationContainer && !notificationContainer.contains(event.target)) {
      showNotifications.value = false;
    }
  };

  document.addEventListener('click', handleClickOutside);

  onUnmounted(() => {
    document.removeEventListener('click', handleClickOutside);

    // Clean up WebSocket connection
    if (wsConnection) {
      // If using a real WebSocket
      // wsConnection.close();

      // If using the interval for simulation
      clearInterval(wsConnection);
    }
  });
});
</script>
