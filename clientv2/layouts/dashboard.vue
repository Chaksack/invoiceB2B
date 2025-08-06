<script setup>
import { ref, computed } from 'vue';
import { useRouter } from '#app';
import { toast } from 'vue-sonner';
import { onMounted, onBeforeUnmount } from 'vue';

const router = useRouter();

// User data would typically come from a store or API
const user = ref({
  name: 'John Doe',
  email: 'john@example.com',
  avatar: '/placeholder-avatar.jpg'
});

// Notification state
const showNotifications = ref(false);
const notifications = ref([
  {
    id: 1,
    title: 'Invoice Approved',
    message: 'Your invoice #INV-2023-001 has been approved for financing.',
    timestamp: new Date(Date.now() - 1000 * 60 * 30), // 30 minutes ago
    read: false,
    type: 'success'
  },
  {
    id: 2,
    title: 'Payment Received',
    message: 'You have received a payment of $12,500.00 for invoice #INV-2023-002.',
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2), // 2 hours ago
    read: true,
    type: 'info'
  },
  {
    id: 3,
    title: 'Document Required',
    message: 'Please upload the missing documentation for invoice #INV-2023-003.',
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24), // 1 day ago
    read: false,
    type: 'warning'
  }
]);

// User dropdown state
const showUserMenu = ref(false);
// Mobile sidebar state
const showMobileSidebar = ref(false);

// Toggle notification panel
const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value;
  if (showUserMenu.value) showUserMenu.value = false;
  if (showMobileSidebar.value) showMobileSidebar.value = false;
};

// Toggle user menu
const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value;
  if (showNotifications.value) showNotifications.value = false;
  if (showMobileSidebar.value) showMobileSidebar.value = false;
};

// Toggle mobile sidebar
const toggleMobileSidebar = () => {
  showMobileSidebar.value = !showMobileSidebar.value;
  if (showNotifications.value) showNotifications.value = false;
  if (showUserMenu.value) showUserMenu.value = false;
};

// Close mobile sidebar when screen size changes to desktop
const handleResize = () => {
  if (window.innerWidth >= 768 && showMobileSidebar.value) {
    showMobileSidebar.value = false;
  }
};

// Close dropdowns when clicking outside
const handleClickOutside = (event) => {
  // Close notifications dropdown if clicking outside
  if (showNotifications.value && !event.target.closest('.notifications-container')) {
    showNotifications.value = false;
  }
  
  // Close user menu dropdown if clicking outside
  if (showUserMenu.value && !event.target.closest('.user-menu-container')) {
    showUserMenu.value = false;
  }
  
  // Don't close mobile sidebar when clicking inside it
  if (showMobileSidebar.value && !event.target.closest('.mobile-sidebar') && 
      !event.target.closest('.mobile-menu-button')) {
    showMobileSidebar.value = false;
  }
};

// Setup event listeners
onMounted(() => {
  window.addEventListener('resize', handleResize);
  document.addEventListener('click', handleClickOutside);
});

// Clean up event listeners
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize);
  document.removeEventListener('click', handleClickOutside);
});

// Mark notification as read
const markAsRead = (notificationId) => {
  const notification = notifications.value.find(n => n.id === notificationId);
  if (notification) {
    notification.read = true;
    toast.success('Notification marked as read');
  }
};

// Mark all notifications as read
const markAllAsRead = () => {
  notifications.value.forEach(notification => {
    notification.read = true;
  });
  toast.success('All notifications marked as read');
};

// Format relative time for notifications
const formatRelativeTime = (timestamp) => {
  const now = new Date();
  const diff = now - timestamp;
  
  // Convert to minutes
  const minutes = Math.floor(diff / (1000 * 60));
  
  if (minutes < 1) return 'Just now';
  if (minutes < 60) return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
  
  // Convert to hours
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hour${hours > 1 ? 's' : ''} ago`;
  
  // Convert to days
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days} day${days > 1 ? 's' : ''} ago`;
  
  // Format as date for older notifications
  return timestamp.toLocaleDateString();
};

// Get unread notifications
const unreadNotifications = computed(() => {
  return notifications.value.filter(n => !n.read);
});

// Logout function
const logout = () => {
  toast.success('Logged out successfully');
  setTimeout(() => {
    router.push('/login');
  }, 1000);
};

// Navigation items for the sidebar
const navigationItems = [
  { 
    title: 'Dashboard', 
    icon: 'HomeIcon', 
    to: '/dashboard' 
  },
  { 
    title: 'Invoices', 
    icon: 'FileTextIcon', 
    to: '/invoices' 
  },
  { 
    title: 'Submit Invoice', 
    icon: 'UploadIcon', 
    to: '/submit-invoice' 
  },
  { 
    title: 'Financing', 
    icon: 'DollarSignIcon', 
    to: '/financing' 
  },
  { 
    title: 'Payments', 
    icon: 'CreditCardIcon', 
    to: '/payments' 
  },
  { 
    title: 'Settings', 
    icon: 'SettingsIcon', 
    to: '/settings' 
  },
  { 
    title: 'Profile', 
    icon: 'UserIcon', 
    to: '/profile' 
  }
];
</script>

<template>
  <div class="min-h-screen bg-background flex">
    <!-- Desktop Sidebar -->
    <aside class="w-64 bg-sidebar border-r border-sidebar-border hidden md:block fixed top-0 left-0 h-full z-20">
      <div class="p-4 border-b border-sidebar-border">
        <div class="flex items-center space-x-2">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 12L12 22L22 12L12 2Z" fill="var(--sidebar-primary)" />
          </svg>
          <span class="text-xl font-bold text-sidebar-foreground">InvoiceFlow</span>
        </div>
      </div>
      
      <nav class="p-2 overflow-y-auto h-[calc(100%-4rem)]">
        <ul class="space-y-1">
          <li v-for="item in navigationItems" :key="item.title">
            <NuxtLink 
              :to="item.to" 
              class="flex items-center p-2 rounded-md text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-colors"
              active-class="bg-sidebar-primary text-sidebar-primary-foreground font-medium"
              exact-active-class="bg-sidebar-primary text-sidebar-primary-foreground font-medium"
            >
              <span class="mr-2">
                <!-- Placeholder for icon -->
                <div class="w-5 h-5 flex items-center justify-center">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                  </svg>
                </div>
              </span>
              <span>{{ item.title }}</span>
            </NuxtLink>
          </li>
        </ul>
      </nav>
    </aside>

    <!-- Mobile Sidebar Overlay -->
    <div 
      v-if="showMobileSidebar" 
      class="fixed inset-0 bg-black bg-opacity-50 z-30 md:hidden"
      @click="showMobileSidebar = false"
    ></div>

    <!-- Mobile Sidebar -->
    <aside 
      v-if="showMobileSidebar"
      class="mobile-sidebar w-64 bg-sidebar border-r border-sidebar-border fixed top-0 left-0 h-full z-40 md:hidden transform transition-transform duration-300"
      :class="{ 'translate-x-0': showMobileSidebar, '-translate-x-full': !showMobileSidebar }"
    >
      <div class="p-4 border-b border-sidebar-border flex justify-between items-center">
        <div class="flex items-center space-x-2">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M12 2L2 12L12 22L22 12L12 2Z" fill="var(--sidebar-primary)" />
          </svg>
          <span class="text-xl font-bold text-sidebar-foreground">InvoiceFlow</span>
        </div>
        <button 
          @click="showMobileSidebar = false"
          class="p-1 rounded-md text-sidebar-foreground hover:bg-sidebar-accent"
        >
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      
      <nav class="p-2 overflow-y-auto h-[calc(100%-4rem)]">
        <ul class="space-y-1">
          <li v-for="item in navigationItems" :key="item.title">
            <NuxtLink 
              :to="item.to" 
              class="flex items-center p-2 rounded-md text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground transition-colors"
              active-class="bg-sidebar-primary text-sidebar-primary-foreground font-medium"
              exact-active-class="bg-sidebar-primary text-sidebar-primary-foreground font-medium"
              @click="showMobileSidebar = false"
            >
              <span class="mr-2">
                <!-- Placeholder for icon -->
                <div class="w-5 h-5 flex items-center justify-center">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                  </svg>
                </div>
              </span>
              <span>{{ item.title }}</span>
            </NuxtLink>
          </li>
        </ul>
      </nav>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col md:ml-64">
      <!-- Header -->
      <header class="bg-card shadow-sm border-b fixed top-0 right-0 left-0 md:left-64 z-10">
        <div class="container mx-auto px-4 py-3 flex justify-between items-center">
          <!-- Mobile menu button -->
          <button 
            @click="toggleMobileSidebar"
            class="mobile-menu-button md:hidden p-2 rounded-md text-gray-500 hover:text-gray-700 hover:bg-gray-100"
          >
            <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </button>
          
          <!-- Search -->
          <div class="hidden md:flex items-center flex-1 max-w-md mx-4">
            <div class="relative w-full">
              <input type="text" placeholder="Search..." class="w-full py-1.5 pl-10 pr-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary" />
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
              </div>
            </div>
          </div>
          
          <!-- User menu -->
          <div class="flex items-center space-x-4">
            <!-- Notifications -->
            <div class="notifications-container relative">
              <button 
                @click="toggleNotifications"
                class="p-1.5 rounded-md text-gray-500 hover:text-gray-700 hover:bg-gray-100 relative"
              >
                <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path>
                </svg>
                <span 
                  v-if="unreadNotifications.length > 0" 
                  class="absolute top-0 right-0 block h-4 w-4 rounded-full bg-red-500 text-white text-xs flex items-center justify-center"
                >
                  {{ unreadNotifications.length > 9 ? '9+' : unreadNotifications.length }}
                </span>
              </button>
              
              <!-- Notification Dropdown -->
              <div 
                v-if="showNotifications" 
                class="absolute right-0 mt-2 w-80 max-w-[90vw] bg-white rounded-lg shadow-lg border border-gray-200 z-50"
                style="top: 100%;"
              >
                <div class="p-3 border-b border-gray-200 flex justify-between items-center">
                  <h3 class="text-sm font-medium">Notifications</h3>
                  <button 
                    v-if="unreadNotifications.length > 0"
                    @click="markAllAsRead"
                    class="text-xs text-primary hover:text-primary-700"
                  >
                    Mark all as read
                  </button>
                </div>
                
                <div class="max-h-80 overflow-y-auto">
                  <div v-if="notifications.length === 0" class="p-4 text-center text-gray-500">
                    No notifications
                  </div>
                  
                  <div 
                    v-for="notification in notifications" 
                    :key="notification.id"
                    class="p-3 border-b border-gray-100 hover:bg-gray-50"
                    :class="{ 'bg-blue-50': !notification.read }"
                  >
                    <div class="flex">
                      <div class="flex-shrink-0 mr-3">
                        <div 
                          class="h-8 w-8 rounded-full flex items-center justify-center"
                          :class="{
                            'bg-green-500': notification.type === 'success',
                            'bg-blue-500': notification.type === 'info',
                            'bg-yellow-500': notification.type === 'warning',
                            'bg-red-500': notification.type === 'error'
                          }"
                        >
                          <svg v-if="notification.type === 'success'" class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                          </svg>
                          <svg v-else-if="notification.type === 'info'" class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                          </svg>
                          <svg v-else-if="notification.type === 'warning'" class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                          </svg>
                          <svg v-else-if="notification.type === 'error'" class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                          </svg>
                        </div>
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm font-medium text-gray-900 truncate">{{ notification.title }}</p>
                        <p class="text-xs text-gray-500 mt-1 break-words">{{ notification.message }}</p>
                        <p class="text-xs text-gray-400 mt-1">{{ formatRelativeTime(notification.timestamp) }}</p>
                      </div>
                      <button 
                        v-if="!notification.read" 
                        @click="markAsRead(notification.id)"
                        class="ml-2 text-xs text-primary hover:text-primary-700 flex-shrink-0"
                      >
                        Mark as read
                      </button>
                    </div>
                  </div>
                </div>
                
                <div class="p-2 border-t border-gray-200 text-center">
                  <NuxtLink to="/settings" class="text-xs text-primary hover:text-primary-700">
                    View all notifications
                  </NuxtLink>
                </div>
              </div>
            </div>
            
            <!-- User dropdown -->
            <div class="user-menu-container relative">
              <button 
                @click="toggleUserMenu"
                class="flex items-center space-x-2 p-1 rounded-md hover:bg-gray-100"
              >
                <div class="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                  <img v-if="user.avatar" :src="user.avatar" alt="User avatar" class="w-full h-full object-cover" />
                  <span v-else class="text-gray-700 font-medium">{{ user.name.charAt(0) }}</span>
                </div>
                <span class="hidden md:inline text-sm font-medium">{{ user.name }}</span>
                <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="hidden md:inline">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </button>
              
              <!-- User Dropdown Menu -->
              <div 
                v-if="showUserMenu" 
                class="absolute right-0 mt-2 w-48 max-w-[90vw] bg-white rounded-lg shadow-lg border border-gray-200 z-50"
                style="top: 100%;"
              >
                <div class="p-3 border-b border-gray-200">
                  <p class="text-sm font-medium text-gray-900 truncate">{{ user.name }}</p>
                  <p class="text-xs text-gray-500 truncate">{{ user.email }}</p>
                </div>
                
                <div class="py-1">
                  <NuxtLink to="/profile" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                    Your Profile
                  </NuxtLink>
                  <NuxtLink to="/settings" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
                    Settings
                  </NuxtLink>
                </div>
                
                <div class="py-1 border-t border-gray-100">
                  <button 
                    @click="logout"
                    class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100"
                  >
                    Logout
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>
      
      <!-- Main content -->
      <main class="flex-1 overflow-auto p-4 md:p-6 mt-14">
        <slot />
      </main>
      
      <!-- Footer -->
      <footer class="bg-card border-t py-4 px-6 text-center text-sm text-muted-foreground mt-auto">
        <p>&copy; 2025 InvoiceFlow. All rights reserved.</p>
      </footer>
    </div>
  </div>
</template>