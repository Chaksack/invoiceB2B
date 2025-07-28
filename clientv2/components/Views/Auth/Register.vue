<script setup>
import { ref } from 'vue';
import { useRouter } from '#app';
import { toast } from 'vue-sonner';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from '~/components/ui/form';

const router = useRouter();
const isSubmitting = ref(false);

// Form data
const formData = ref({
  firstName: '',
  lastName: '',
  email: '',
  password: ''
});

// Handle form submission
const handleSubmit = (values) => {
  isSubmitting.value = true;
  
  // In a real application, you would submit the form data to an API
  // For now, we'll simulate a successful registration
  setTimeout(() => {
    isSubmitting.value = false;
    toast.success('Registration successful! Redirecting to onboarding...');
    
    // Redirect to onboarding page
    setTimeout(() => {
      router.push('/onboarding');
    }, 1500);
  }, 1000);
};
</script>

<template>
  <div class="min-h-screen flex flex-col lg:flex-row font-inter">
    <!-- Left Section: Primary Background with Illustrations -->
    <div class="relative hidden lg:flex bg-primary text-white p-4 sm:p-6 md:p-8 lg:w-1/2 flex-col justify-center items-center overflow-hidden">
      <!-- Background patterns -->
      <div class="absolute top-0 left-0 w-full h-full opacity-10">
        <svg class="absolute top-8 left-8 hidden sm:block" width="50" height="50" viewBox="0 0 50 50" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="25" cy="25" r="23" stroke="white" stroke-width="2" />
        </svg>
        <svg class="absolute bottom-16 right-16 hidden sm:block" width="100" height="100" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="10" y="10" width="80" height="80" rx="10" stroke="white" stroke-width="2" />
        </svg>
        <div class="absolute top-1/4 right-1/4 w-4 h-4 bg-white rounded-full hidden sm:block"></div>
        <div class="absolute bottom-1/3 left-1/3 w-6 h-6 bg-white rounded-full hidden sm:block"></div>
      </div>

      <!-- Logo -->
      <div class="w-full flex justify-center sm:justify-start sm:absolute sm:top-8 sm:left-8 items-center space-x-2 mb-6 sm:mb-0">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 2L2 12L12 22L22 12L12 2Z" fill="white"/>
        </svg>
        <span class="text-2xl font-bold">InvoiceFlow</span>
      </div>

      <!-- Illustrative Cards -->
      <div class="relative mt-8 sm:mt-12 lg:mt-16 mb-8 w-full max-w-sm sm:max-w-md lg:max-w-xl mx-auto">
        <!-- Main Card -->
        <div class="bg-white text-gray-800 p-4 sm:p-6 rounded-xl shadow-xl z-10 relative">
          <div class="flex justify-between items-center mb-4">
            <div class="text-sm sm:text-lg font-semibold">Invoices Financed</div>
            <div class="text-sm sm:text-lg font-semibold">Funds Available</div>
          </div>
          <div class="flex justify-between items-center mb-6">
            <div class="text-lg sm:text-2xl font-bold">$124,908.00</div>
            <div class="text-lg sm:text-2xl font-bold">$85,750.00</div>
          </div>
          <!-- Graph Placeholder -->
          <div class="h-16 sm:h-24 bg-gray-100 rounded-lg flex items-center justify-center text-gray-400 mb-6">
            <svg class="w-full h-full" viewBox="0 0 300 100" preserveAspectRatio="none">
              <polyline fill="none" stroke="var(--primary)" stroke-width="2" points="0,80 50,50 100,70 150,30 200,60 250,40 300,20" />
              <polyline fill="none" stroke="#EF4444" stroke-width="2" points="0,20 50,40 100,30 150,60 200,40 250,70 300,80" />
            </svg>
          </div>
          <div class="flex items-center mb-4">
            <div class="w-5 h-5 sm:w-6 sm:h-6 bg-green-100 rounded-full flex items-center justify-center mr-2 sm:mr-3">
              <svg class="w-3 h-3 sm:w-4 sm:h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
            </div>
            <div>
              <div class="text-sm sm:text-base font-medium">Invoice financing approved</div>
              <div class="text-xs sm:text-sm text-gray-500">Today at 11:15 AM</div>
            </div>
          </div>
          <div class="border-t border-gray-200 pt-4">
            <div class="flex justify-between items-center mb-2">
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 12</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 15</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 18</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 21</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 24</div>
            </div>
            <div class="flex justify-between items-center text-xs sm:text-sm">
              <div class="font-medium text-gray-800">$35,798.00</div>
              <div class="font-medium text-gray-800">$23.10</div>
              <div class="font-medium text-gray-800">$23.10</div>
              <div class="font-medium text-gray-800">$23.10</div>
              <div class="font-medium text-gray-800">$23.10</div>
            </div>
          </div>
        </div>

        <!-- Overlapping Card 1 -->
        <div class="absolute bg-white text-gray-800 p-3 sm:p-4 rounded-xl shadow-lg -top-6 sm:-top-10 left-1/3 sm:left-1/4 transform -translate-x-1/2 z-0 opacity-90">
          <div class="font-semibold text-xs sm:text-sm mb-1 sm:mb-2">Invoice Financed</div>
          <div class="text-sm sm:text-lg font-bold text-green-500">+$42,500.00</div>
        </div>

        <!-- Overlapping Card 2 -->
        <div class="absolute bg-white text-gray-800 p-3 sm:p-4 rounded-xl shadow-lg bottom-6 sm:bottom-10 right-0 transform translate-x-1/6 sm:translate-x-1/4 z-0 opacity-90">
          <div class="flex items-center justify-between text-xs sm:text-sm mb-1">
            <span class="font-medium">Financing Fee</span>
            <span class="text-red-500 font-bold">-$1,275</span>
          </div>
          <div class="flex items-center justify-between text-xs sm:text-sm mb-1">
            <span class="font-medium">Service Charge</span>
            <span class="text-red-500 font-bold">-$425</span>
          </div>
          <div class="flex items-center justify-between text-xs sm:text-sm">
            <span class="font-medium">Processing Fee</span>
            <span class="text-red-500 font-bold">-$212</span>
          </div>
          <div class="text-xs text-gray-500 mt-1 sm:mt-2">Today at 9:24 AM</div>
        </div>
      </div>

      <!-- Main Marketing Text -->
      <div class="text-center px-4 mt-8 sm:mt-10">
        <h1 class="text-3xl sm:text-4xl lg:text-5xl font-extrabold mb-3 sm:mb-4">Unlock Your Business Cash Flow</h1>
        <p class="text-base sm:text-lg leading-relaxed max-w-md mx-auto">
          InvoiceFlow helps your business convert unpaid invoices into immediate working capital. Get paid within 24 hours, bridge cash flow gaps, and focus on growing your business while we handle your invoice financing needs.
        </p>
      </div>
    </div>

    <!-- Right Section: Sign Up Form -->
    <div class="bg-white p-4 sm:p-6 md:p-8 lg:p-12 w-full lg:w-1/2 flex flex-col justify-between">
      <div class="flex flex-col items-center w-full max-w-md lg:max-w-lg mx-auto p-4 sm:p-6 md:p-8 my-auto">
        <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2 text-center">Business Financing Made Simple</h2>
        <p class="text-gray-600 mb-6 sm:mb-8 text-center text-sm sm:text-base">Turn your invoices into immediate capital</p>

        <!-- Social Sign-up Buttons -->
        <div class="flex flex-col space-y-3 w-full mb-6 sm:mb-8">
          <Button variant="outline" class="w-full py-2 sm:py-3 px-4 sm:px-6 border border-gray-300 text-gray-700 hover:bg-gray-50 text-sm sm:text-base flex items-center justify-center">
            <svg class="w-4 h-4 sm:w-5 sm:h-5 mr-2" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
              <path d="M12.24 10.29C11.66 10.11 10.99 10 10.24 10C8.04 10 6.27 11.17 5.24 13.06L2.5 10.18C3.89 7.7 6.81 6 10.24 6C12.55 6 14.5 6.8 15.9 8.1L18.78 5.22C16.89 3.32 14.07 2 10.24 2C6.18 2 2.65 4.38 0.76 7.74L3.5 10.62C4.47 8.84 6.18 7.6 8.24 7.6C9.18 7.6 10.04 7.78 10.81 8.16L12.24 10.29Z" fill="#EA4335"/>
              <path d="M23.5 12.24C23.5 11.44 23.42 10.66 23.27 9.9L12 9.9V14.1H18.47C18.18 15.63 17.21 16.92 16.03 17.77L18.9 20.64C20.67 19.06 21.84 16.73 22.5 14.1L23.5 12.24Z" fill="#4285F4"/>
              <path d="M12 22C15.24 22 17.96 20.9 19.96 19.04L17.09 16.17C16.27 16.74 15.29 17.2 14.19 17.44C13.09 17.68 11.96 17.8 10.84 17.8C7.45 17.8 4.54 15.82 3.03 12.96L0.29 15.84C2.18 19.1 6.06 21.5 10.84 21.5C10.96 21.5 11.08 21.5 11.2 21.5C11.32 21.5 11.44 21.5 11.56 21.5L12 22Z" fill="#34A853"/>
              <path d="M22.28 12.24C22.28 11.44 22.2 10.66 22.05 9.9L12 9.9V14.1H18.47C18.18 15.63 17.21 16.92 16.03 17.77L18.9 20.64C20.67 19.06 21.84 16.73 22.5 14.1L23.5 12.24Z" fill="#FBBC05"/>
            </svg>
            Sign Up with Google
          </Button>
          <Button variant="outline" class="w-full py-2 sm:py-3 px-4 sm:px-6 border border-gray-300 text-gray-700 hover:bg-gray-50 text-sm sm:text-base flex items-center justify-center">
            <svg class="w-4 h-4 sm:w-5 sm:h-5 mr-2" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" fill="#000000"/>
              <path d="M15.5 10.5C15.5 9.67157 14.8284 9 14 9C13.1716 9 12.5 9.67157 12.5 10.5C12.5 11.3284 13.1716 12 14 12C14.8284 12 15.5 11.3284 15.5 10.5Z" fill="#FFFFFF"/>
              <path d="M10.5 10.5C10.5 9.67157 9.82843 9 9 9C8.17157 9 7.5 9.67157 7.5 10.5C7.5 11.3284 8.17157 12 9 12C9.82843 12 10.5 11.3284 10.5 10.5Z" fill="#FFFFFF"/>
              <path d="M12 16C10.6716 16 9.5 16.6716 9.5 17.5C9.5 18.3284 10.6716 19 12 19C13.3284 19 14.5 18.3284 14.5 17.5C14.5 16.6716 13.3284 16 12 16Z" fill="#FFFFFF"/>
            </svg>
            Sign Up with Apple
          </Button>
        </div>

        <!-- Separator -->
        <div class="relative w-full text-center text-gray-500 mb-6 sm:mb-8">
          <div class="absolute inset-0 flex items-center">
            <span class="w-full border-t border-gray-300"></span>
          </div>
          <div class="relative bg-white px-4 inline-block text-sm">Or with email</div>
        </div>

        <!-- Form Fields -->
        <Form :initial-values="formData" @submit="handleSubmit" class="w-full">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4 w-full mb-3 sm:mb-4">
            <FormField v-slot="{ componentField }" name="firstName">
              <FormItem>
                <FormLabel for="firstName" class="sr-only">First Name</FormLabel>
                <FormControl>
                  <Input 
                    id="firstName" 
                    v-bind="componentField"
                    placeholder="Alesia" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base border border-gray-300 focus:ring-primary focus:border-primary" 
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField v-slot="{ componentField }" name="lastName">
              <FormItem>
                <FormLabel for="lastName" class="sr-only">Last Name</FormLabel>
                <FormControl>
                  <Input 
                    id="lastName" 
                    v-bind="componentField"
                    placeholder="Karapova" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base border border-gray-300 focus:ring-primary focus:border-primary" 
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>
          
          <FormField v-slot="{ componentField }" name="email">
            <FormItem class="w-full mb-3 sm:mb-4">
              <FormLabel for="email" class="sr-only">Email</FormLabel>
              <FormControl>
                <Input 
                  id="email" 
                  v-bind="componentField"
                  type="email" 
                  placeholder="alesiakara@po" 
                  class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base border border-gray-300 focus:ring-primary focus:border-primary" 
                  required
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
          
          <FormField v-slot="{ componentField }" name="password">
            <FormItem class="w-full mb-4 sm:mb-6">
              <FormLabel for="password" class="sr-only">Password</FormLabel>
              <FormControl>
                <Input 
                  id="password" 
                  v-bind="componentField"
                  type="password" 
                  placeholder="Password" 
                  class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base border border-gray-300 focus:ring-primary focus:border-primary" 
                  required
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Agreement Text -->
          <p class="text-xs sm:text-sm text-gray-600 mb-4 sm:mb-6 text-center">
            By creating an account, you agree to our
            <NuxtLink to="/settings" class="text-primary hover:underline">Privacy Policy</NuxtLink> and
            <NuxtLink to="/settings" class="text-primary hover:underline">Electronics Communication Policy</NuxtLink>.
          </p>

          <!-- Sign Up Button -->
          <Button 
            type="submit" 
            class="w-full bg-primary text-white py-2 sm:py-3 px-4 sm:px-6 text-sm sm:text-base hover:bg-primary-700 focus:ring-4 focus:ring-primary"
            :disabled="isSubmitting"
          >
            <span v-if="isSubmitting">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white inline-block" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Signing Up...
            </span>
            <span v-else>Sign Up</span>
          </Button>
        </Form>

        <!-- Sign In Link -->
        <p class="mt-4 sm:mt-6 text-center text-xs sm:text-sm text-gray-600">
          Already have an account? <NuxtLink to="/login" class="text-primary hover:underline font-medium">Sign In</NuxtLink>
        </p>
      </div>

      <!-- Footer -->
      <div class="mt-6 flex justify-center sm:justify-between text-xs sm:text-sm text-gray-500 w-full max-w-xs sm:max-w-sm md:max-w-md mx-auto px-4">
        <NuxtLink to="/settings" class="hover:underline mr-4 sm:mr-0">Privacy Policy</NuxtLink>
        <span>Copyright 2025</span>
      </div>
    </div>
  </div>
</template>

<style>
/* Ensure the Inter font is applied globally if not already configured in Tailwind */
body {
  font-family: 'Inter', sans-serif;
}
</style>
