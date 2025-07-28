<script setup>
import { ref } from 'vue';

// User profile data
const userProfile = ref({
  firstName: 'John',
  lastName: 'Doe',
  email: 'john.doe@example.com',
  phone: '(555) 123-4567',
  jobTitle: 'Finance Manager',
  department: 'Finance',
  bio: 'Experienced finance professional with over 10 years in the industry. Specializing in invoice management and financial operations.',
  avatar: null,
  avatarUrl: '/placeholder-avatar.jpg'
});

// Form validation
const errors = ref({});
const isSubmitting = ref(false);

// Department options
const departmentOptions = [
  { value: 'finance', label: 'Finance' },
  { value: 'accounting', label: 'Accounting' },
  { value: 'operations', label: 'Operations' },
  { value: 'sales', label: 'Sales' },
  { value: 'marketing', label: 'Marketing' },
  { value: 'executive', label: 'Executive' },
  { value: 'other', label: 'Other' }
];

// Handle avatar upload
const handleAvatarUpload = (event) => {
  const file = event.target.files[0];
  if (file) {
    userProfile.value.avatar = file;
    
    // Create a preview URL
    const reader = new FileReader();
    reader.onload = (e) => {
      userProfile.value.avatarUrl = e.target.result;
    };
    reader.readAsDataURL(file);
  }
};

// Remove avatar
const removeAvatar = () => {
  userProfile.value.avatar = null;
  userProfile.value.avatarUrl = '/placeholder-avatar.jpg';
};

// Validate form
const validateForm = () => {
  const newErrors = {};
  
  if (!userProfile.value.firstName) {
    newErrors.firstName = 'First name is required';
  }
  
  if (!userProfile.value.lastName) {
    newErrors.lastName = 'Last name is required';
  }
  
  if (!userProfile.value.email) {
    newErrors.email = 'Email is required';
  } else if (!/^\S+@\S+\.\S+$/.test(userProfile.value.email)) {
    newErrors.email = 'Please enter a valid email address';
  }
  
  if (userProfile.value.phone && !/^(\+\d{1,3})?\s?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$/.test(userProfile.value.phone)) {
    newErrors.phone = 'Please enter a valid phone number';
  }
  
  errors.value = newErrors;
  return Object.keys(newErrors).length === 0;
};

// Save profile
const saveProfile = () => {
  if (!validateForm()) return;
  
  isSubmitting.value = true;
  
  // Simulate API call
  setTimeout(() => {
    isSubmitting.value = false;
    alert('Profile updated successfully!');
  }, 1000);
};

definePageMeta({
  layout: 'dashboard'
});
</script>

<template>
  <div>
    <!-- Page Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Your Profile</h1>
      <p class="text-gray-500 mt-1">Manage your personal information and preferences</p>
    </div>

    <!-- Profile Content -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Avatar and Bio -->
      <div>
        <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
          <div class="p-4 border-b border-border">
            <h2 class="font-semibold text-lg">Profile Picture</h2>
          </div>
          
          <div class="p-6 flex flex-col items-center">
            <div class="relative mb-4">
              <div class="w-32 h-32 rounded-full overflow-hidden bg-gray-200">
                <img 
                  v-if="userProfile.avatarUrl" 
                  :src="userProfile.avatarUrl" 
                  alt="Profile picture" 
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
                  <svg class="w-16 h-16" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                    <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"></path>
                  </svg>
                </div>
              </div>
              <button 
                v-if="userProfile.avatarUrl !== '/placeholder-avatar.jpg'"
                @click="removeAvatar"
                class="absolute -top-1 -right-1 bg-red-500 text-white rounded-full p-1 shadow-md hover:bg-red-600"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </button>
            </div>
            
            <div class="w-full">
              <label 
                for="avatar-upload" 
                class="block w-full text-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 cursor-pointer"
              >
                Change Picture
              </label>
              <input 
                id="avatar-upload" 
                type="file" 
                accept="image/*" 
                @change="handleAvatarUpload" 
                class="hidden"
              />
              <p class="mt-2 text-xs text-gray-500 text-center">JPG, PNG or GIF. Max size 2MB.</p>
            </div>
          </div>
          
          <div class="p-4 border-t border-border">
            <h3 class="text-sm font-medium text-gray-700 mb-2">Bio</h3>
            <textarea 
              v-model="userProfile.bio" 
              rows="4" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
              placeholder="Tell us about yourself"
            ></textarea>
          </div>
        </div>
      </div>
      
      <!-- Right Column: Profile Information -->
      <div class="lg:col-span-2">
        <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
          <div class="p-4 border-b border-border">
            <h2 class="font-semibold text-lg">Personal Information</h2>
          </div>
          
          <form @submit.prevent="saveProfile" class="p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <div>
                <label for="firstName" class="block text-sm font-medium text-gray-700 mb-1">First Name</label>
                <input 
                  id="firstName" 
                  v-model="userProfile.firstName" 
                  type="text" 
                  class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                  :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.firstName }"
                />
                <p v-if="errors.firstName" class="mt-1 text-sm text-red-600">{{ errors.firstName }}</p>
              </div>
              
              <div>
                <label for="lastName" class="block text-sm font-medium text-gray-700 mb-1">Last Name</label>
                <input 
                  id="lastName" 
                  v-model="userProfile.lastName" 
                  type="text" 
                  class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                  :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.lastName }"
                />
                <p v-if="errors.lastName" class="mt-1 text-sm text-red-600">{{ errors.lastName }}</p>
              </div>
              
              <div>
                <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email Address</label>
                <input 
                  id="email" 
                  v-model="userProfile.email" 
                  type="email" 
                  class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                  :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.email }"
                />
                <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
              </div>
              
              <div>
                <label for="phone" class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
                <input 
                  id="phone" 
                  v-model="userProfile.phone" 
                  type="tel" 
                  class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                  :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.phone }"
                  placeholder="(555) 123-4567"
                />
                <p v-if="errors.phone" class="mt-1 text-sm text-red-600">{{ errors.phone }}</p>
              </div>
            </div>
            
            <h3 class="text-lg font-medium text-gray-900 mb-4">Work Information</h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <div>
                <label for="jobTitle" class="block text-sm font-medium text-gray-700 mb-1">Job Title</label>
                <input 
                  id="jobTitle" 
                  v-model="userProfile.jobTitle" 
                  type="text" 
                  class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                />
              </div>
              
              <div>
                <label for="department" class="block text-sm font-medium text-gray-700 mb-1">Department</label>
                <select 
                  id="department" 
                  v-model="userProfile.department" 
                  class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                >
                  <option v-for="option in departmentOptions" :key="option.value" :value="option.label">
                    {{ option.label }}
                  </option>
                </select>
              </div>
            </div>
            
            <div class="flex justify-end">
              <button 
                type="submit" 
                class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
                :disabled="isSubmitting"
              >
                <span v-if="isSubmitting" class="flex items-center justify-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Saving...
                </span>
                <span v-else>Save Profile</span>
              </button>
            </div>
          </form>
        </div>
        
        <!-- Connected Accounts Section -->
        <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden mt-6">
          <div class="p-4 border-b border-border">
            <h2 class="font-semibold text-lg">Connected Accounts</h2>
          </div>
          
          <div class="p-6">
            <div class="space-y-4">
              <div class="flex items-center justify-between py-2 border-b border-border">
                <div class="flex items-center">
                  <div class="w-10 h-10 bg-[#4285F4] rounded-full flex items-center justify-center text-white mr-3">
                    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                      <path d="M12.24 10.285C11.66 10.095 10.99 10 10.24 10C8.04 10 6.27 11.17 5.24 13.06L2.5 10.18C3.89 7.7 6.81 6 10.24 6C12.55 6 14.5 6.8 15.9 8.1L12.24 10.285Z" fill="#EA4335"/>
                      <path d="M23.5 12.24C23.5 11.44 23.42 10.66 23.27 9.9L12 9.9V14.1H18.47C18.18 15.63 17.21 16.92 16.03 17.77L18.9 20.64C20.67 19.06 21.84 16.73 22.5 14.1L23.5 12.24Z" fill="#4285F4"/>
                      <path d="M12 22C15.24 22 17.96 20.9 19.96 19.04L17.09 16.17C16.27 16.74 15.29 17.2 14.19 17.44C13.09 17.68 11.96 17.8 10.84 17.8C7.45 17.8 4.54 15.82 3.03 12.96L0.29 15.84C2.18 19.1 6.06 21.5 10.84 21.5C10.96 21.5 11.08 21.5 11.2 21.5C11.32 21.5 11.44 21.5 11.56 21.5L12 22Z" fill="#34A853"/>
                      <path d="M22.28 12.24C22.28 11.44 22.2 10.66 22.05 9.9L12 9.9V14.1H18.47C18.18 15.63 17.21 16.92 16.03 17.77L18.9 20.64C20.67 19.06 21.84 16.73 22.5 14.1L23.5 12.24Z" fill="#FBBC05"/>
                    </svg>
                  </div>
                  <div>
                    <h3 class="text-sm font-medium text-gray-900">Google</h3>
                    <p class="text-xs text-gray-500">Connected</p>
                  </div>
                </div>
                <button class="text-xs text-red-600 hover:text-red-800 font-medium">
                  Disconnect
                </button>
              </div>
              
              <div class="flex items-center justify-between py-2 border-b border-border">
                <div class="flex items-center">
                  <div class="w-10 h-10 bg-[#1DA1F2] rounded-full flex items-center justify-center text-white mr-3">
                    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                      <path d="M22.46 6C21.69 6.35 20.86 6.58 20 6.69C20.88 6.16 21.56 5.32 21.88 4.31C21.05 4.81 20.13 5.16 19.16 5.36C18.37 4.5 17.26 4 16 4C13.65 4 11.73 5.92 11.73 8.29C11.73 8.63 11.77 8.96 11.84 9.27C8.28 9.09 5.11 7.38 3 4.79C2.63 5.42 2.42 6.16 2.42 6.94C2.42 8.43 3.17 9.75 4.33 10.5C3.62 10.5 2.96 10.3 2.38 10V10.03C2.38 12.11 3.86 13.85 5.82 14.24C5.46 14.34 5.08 14.39 4.69 14.39C4.42 14.39 4.15 14.36 3.89 14.31C4.43 16 6 17.26 7.89 17.29C6.43 18.45 4.58 19.13 2.56 19.13C2.22 19.13 1.88 19.11 1.54 19.07C3.44 20.29 5.7 21 8.12 21C16 21 20.33 14.46 20.33 8.79C20.33 8.6 20.33 8.42 20.32 8.23C21.16 7.63 21.88 6.87 22.46 6Z" fill="currentColor"/>
                    </svg>
                  </div>
                  <div>
                    <h3 class="text-sm font-medium text-gray-900">Twitter</h3>
                    <p class="text-xs text-gray-500">Not connected</p>
                  </div>
                </div>
                <button class="text-xs text-primary hover:text-primary-700 font-medium">
                  Connect
                </button>
              </div>
              
              <div class="flex items-center justify-between py-2">
                <div class="flex items-center">
                  <div class="w-10 h-10 bg-[#0A66C2] rounded-full flex items-center justify-center text-white mr-3">
                    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                      <path d="M19 3H5C3.895 3 3 3.895 3 5V19C3 20.105 3.895 21 5 21H19C20.105 21 21 20.105 21 19V5C21 3.895 20.105 3 19 3ZM9 17H6.477V10H9V17ZM7.694 8.717C6.923 8.717 6.408 8.203 6.408 7.517C6.408 6.831 6.922 6.317 7.779 6.317C8.55 6.317 9.065 6.831 9.065 7.517C9.065 8.203 8.551 8.717 7.694 8.717ZM18 17H15.558V13.174C15.558 12.116 14.907 11.872 14.663 11.872C14.419 11.872 13.605 12.035 13.605 13.174C13.605 13.337 13.605 17 13.605 17H11.082V10H13.605V10.977C13.93 10.407 14.581 10 15.802 10C17.023 10 18 10.977 18 13.174V17Z" fill="currentColor"/>
                    </svg>
                  </div>
                  <div>
                    <h3 class="text-sm font-medium text-gray-900">LinkedIn</h3>
                    <p class="text-xs text-gray-500">Not connected</p>
                  </div>
                </div>
                <button class="text-xs text-primary hover:text-primary-700 font-medium">
                  Connect
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>