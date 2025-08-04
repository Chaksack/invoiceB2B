<script setup>
import { ref } from 'vue';

// Form data
const accountSettings = ref({
  companyName: 'Acme Business Solutions',
  taxId: '12-3456789',
  industry: 'technology',
  address: {
    street: '123 Business Ave',
    city: 'San Francisco',
    state: 'CA',
    zip: '94107',
    country: 'United States'
  }
});

const notificationSettings = ref({
  email: true,
  sms: false,
  push: true,
  invoiceReminders: true,
  paymentReceipts: true,
  financingOffers: true,
  marketingUpdates: false
});

const securitySettings = ref({
  twoFactorEnabled: true,
  lastPasswordChange: '2025-06-15',
  sessionTimeout: '30'
});

// Industry options
const industryOptions = [
  { value: 'technology', label: 'Technology' },
  { value: 'healthcare', label: 'Healthcare' },
  { value: 'retail', label: 'Retail' },
  { value: 'manufacturing', label: 'Manufacturing' },
  { value: 'finance', label: 'Finance' },
  { value: 'education', label: 'Education' },
  { value: 'other', label: 'Other' }
];

// Session timeout options
const timeoutOptions = [
  { value: '15', label: '15 minutes' },
  { value: '30', label: '30 minutes' },
  { value: '60', label: '1 hour' },
  { value: '120', label: '2 hours' },
  { value: '240', label: '4 hours' }
];

// Active tab
const activeTab = ref('account');

// Save settings
const saveSettings = (section) => {
  // Simulate API call
  setTimeout(() => {
    alert(`${section} settings saved successfully!`);
  }, 500);
};

definePageMeta({
  layout: 'dashboard'
});
</script>

<template>
  <div>
    <!-- Page Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Account Settings</h1>
      <p class="text-gray-500 mt-1">Manage your account preferences and settings</p>
    </div>

    <!-- Settings Tabs -->
    <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
      <div class="border-b border-border">
        <nav class="flex -mb-px">
          <button 
            @click="activeTab = 'account'" 
            class="py-4 px-6 text-sm font-medium border-b-2 focus:outline-none"
            :class="activeTab === 'account' ? 'border-primary text-primary' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
          >
            Account
          </button>
          <button 
            @click="activeTab = 'notifications'" 
            class="py-4 px-6 text-sm font-medium border-b-2 focus:outline-none"
            :class="activeTab === 'notifications' ? 'border-primary text-primary' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
          >
            Notifications
          </button>
          <button 
            @click="activeTab = 'security'" 
            class="py-4 px-6 text-sm font-medium border-b-2 focus:outline-none"
            :class="activeTab === 'security' ? 'border-primary text-primary' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
          >
            Security
          </button>
        </nav>
      </div>
      
      <!-- Account Settings Tab -->
      <div v-if="activeTab === 'account'" class="p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Business Information</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div>
            <label for="companyName" class="block text-sm font-medium text-gray-700 mb-1">Company Name</label>
            <input 
              id="companyName" 
              v-model="accountSettings.companyName" 
              type="text" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <label for="taxId" class="block text-sm font-medium text-gray-700 mb-1">Tax ID / EIN</label>
            <input 
              id="taxId" 
              v-model="accountSettings.taxId" 
              type="text" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <label for="industry" class="block text-sm font-medium text-gray-700 mb-1">Industry</label>
            <select 
              id="industry" 
              v-model="accountSettings.industry" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            >
              <option v-for="option in industryOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
        </div>
        
        <h2 class="text-lg font-medium text-gray-900 mb-4">Business Address</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div class="md:col-span-2">
            <label for="street" class="block text-sm font-medium text-gray-700 mb-1">Street Address</label>
            <input 
              id="street" 
              v-model="accountSettings.address.street" 
              type="text" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <label for="city" class="block text-sm font-medium text-gray-700 mb-1">City</label>
            <input 
              id="city" 
              v-model="accountSettings.address.city" 
              type="text" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <label for="state" class="block text-sm font-medium text-gray-700 mb-1">State / Province</label>
            <input 
              id="state" 
              v-model="accountSettings.address.state" 
              type="text" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <label for="zip" class="block text-sm font-medium text-gray-700 mb-1">ZIP / Postal Code</label>
            <input 
              id="zip" 
              v-model="accountSettings.address.zip" 
              type="text" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <label for="country" class="block text-sm font-medium text-gray-700 mb-1">Country</label>
            <input 
              id="country" 
              v-model="accountSettings.address.country" 
              type="text" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
        </div>
        
        <div class="flex justify-end">
          <button 
            @click="saveSettings('Account')" 
            class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          >
            Save Changes
          </button>
        </div>
      </div>
      
      <!-- Notifications Tab -->
      <div v-if="activeTab === 'notifications'" class="p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Notification Preferences</h2>
        
        <div class="space-y-4 mb-6">
          <div class="flex items-center justify-between py-2 border-b border-border">
            <div>
              <h3 class="text-sm font-medium text-gray-900">Email Notifications</h3>
              <p class="text-xs text-gray-500">Receive notifications via email</p>
            </div>
            <div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.email" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
          
          <div class="flex items-center justify-between py-2 border-b border-border">
            <div>
              <h3 class="text-sm font-medium text-gray-900">SMS Notifications</h3>
              <p class="text-xs text-gray-500">Receive notifications via text message</p>
            </div>
            <div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.sms" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
          
          <div class="flex items-center justify-between py-2 border-b border-border">
            <div>
              <h3 class="text-sm font-medium text-gray-900">Push Notifications</h3>
              <p class="text-xs text-gray-500">Receive notifications on your device</p>
            </div>
            <div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.push" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
        </div>
        
        <h2 class="text-lg font-medium text-gray-900 mb-4">Notification Types</h2>
        
        <div class="space-y-4 mb-6">
          <div class="flex items-center justify-between py-2 border-b border-border">
            <div>
              <h3 class="text-sm font-medium text-gray-900">Invoice Reminders</h3>
              <p class="text-xs text-gray-500">Receive reminders about upcoming invoice due dates</p>
            </div>
            <div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.invoiceReminders" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
          
          <div class="flex items-center justify-between py-2 border-b border-border">
            <div>
              <h3 class="text-sm font-medium text-gray-900">Payment Receipts</h3>
              <p class="text-xs text-gray-500">Receive notifications when payments are processed</p>
            </div>
            <div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.paymentReceipts" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
          
          <div class="flex items-center justify-between py-2 border-b border-border">
            <div>
              <h3 class="text-sm font-medium text-gray-900">Financing Offers</h3>
              <p class="text-xs text-gray-500">Receive notifications about new financing opportunities</p>
            </div>
            <div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.financingOffers" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
          
          <div class="flex items-center justify-between py-2 border-b border-border">
            <div>
              <h3 class="text-sm font-medium text-gray-900">Marketing Updates</h3>
              <p class="text-xs text-gray-500">Receive marketing communications and updates</p>
            </div>
            <div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notificationSettings.marketingUpdates" class="sr-only peer">
                <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
              </label>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end">
          <button 
            @click="saveSettings('Notification')" 
            class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          >
            Save Changes
          </button>
        </div>
      </div>
      
      <!-- Security Tab -->
      <div v-if="activeTab === 'security'" class="p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Security Settings</h2>
        
        <div class="space-y-6 mb-6">
          <div>
            <div class="flex items-center justify-between py-2 border-b border-border">
              <div>
                <h3 class="text-sm font-medium text-gray-900">Two-Factor Authentication</h3>
                <p class="text-xs text-gray-500">Add an extra layer of security to your account</p>
              </div>
              <div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" v-model="securitySettings.twoFactorEnabled" class="sr-only peer">
                  <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary"></div>
                </label>
              </div>
            </div>
            <div v-if="securitySettings.twoFactorEnabled" class="mt-2 p-3 bg-muted/20 rounded-md">
              <p class="text-xs text-gray-700">Two-factor authentication is enabled. You'll be asked for a verification code when signing in from a new device.</p>
            </div>
          </div>
          
          <div>
            <h3 class="text-sm font-medium text-gray-900 mb-2">Password</h3>
            <div class="flex items-center justify-between">
              <div>
                <p class="text-xs text-gray-500">Last changed: {{ securitySettings.lastPasswordChange }}</p>
              </div>
              <button class="text-sm text-primary hover:text-primary-700 font-medium">
                Change Password
              </button>
            </div>
          </div>
          
          <div>
            <h3 class="text-sm font-medium text-gray-900 mb-2">Session Timeout</h3>
            <p class="text-xs text-gray-500 mb-2">Automatically log out after a period of inactivity</p>
            <select 
              v-model="securitySettings.sessionTimeout" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            >
              <option v-for="option in timeoutOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
          
          <div>
            <h3 class="text-sm font-medium text-gray-900 mb-2">Active Sessions</h3>
            <p class="text-xs text-gray-500 mb-2">Devices currently logged into your account</p>
            <div class="bg-muted/10 rounded-md p-3">
              <div class="flex items-center justify-between py-2 border-b border-border">
                <div>
                  <p class="text-sm font-medium">Current Session</p>
                  <p class="text-xs text-gray-500">Chrome on macOS - San Francisco, CA</p>
                </div>
                <div class="text-xs text-green-600 font-medium">
                  Active Now
                </div>
              </div>
              <div class="flex items-center justify-between py-2">
                <div>
                  <p class="text-sm font-medium">Mobile App</p>
                  <p class="text-xs text-gray-500">iOS - Last active 2 days ago</p>
                </div>
                <button class="text-xs text-red-600 hover:text-red-800 font-medium">
                  Revoke
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end">
          <button 
            @click="saveSettings('Security')" 
            class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
          >
            Save Changes
          </button>
        </div>
      </div>
    </div>
  </div>
</template>