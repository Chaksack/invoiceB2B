<script setup>
// No specific script logic needed for this static UI page.
// Shadcn-vue components are imported automatically by Nuxt if configured correctly.
// For a real application, you would add form submission logic here.
</script>

<template>
  <div class="min-h-screen flex flex-col lg:flex-row font-inter">
    <!-- Left Section: Primary Background with Illustrations -->
    <div class="relative hidden lg:flex bg-primary text-white p-4 sm:p-6 md:p-8 lg:w-1/2 flex-col justify-center items-center overflow-hidden rounded-xl lg:rounded-r-xl lg:rounded-l-none">
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

    <!-- Right Section: Forgot Password Form -->
    <div class="bg-white p-4 sm:p-6 md:p-8 lg:p-12 w-full lg:w-1/2 flex flex-col justify-between">
      <div class="flex flex-col items-center w-full max-w-xs sm:max-w-sm md:max-w-md mx-auto p-4 sm:p-6 md:p-8 my-auto">
        <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2 text-center">Forgot Password</h2>
        <p class="text-gray-600 mb-6 sm:mb-8 text-center text-sm sm:text-base">Enter your email address and we'll send you a link to reset your password</p>

        <!-- Email Icon -->
        <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mb-6 sm:mb-8">
          <svg class="w-8 h-8 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
          </svg>
        </div>

        <!-- Form Fields -->
        <div class="w-full mb-6 sm:mb-8">
          <Label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email Address</Label>
          <Input id="email" type="email" placeholder="your.email@example.com" class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full" />
        </div>

        <!-- Reset Password Button -->
        <Button class="w-full bg-primary text-white py-2 sm:py-3 px-4 sm:px-6 text-sm sm:text-base rounded-lg hover:bg-primary-700 focus:ring-4 focus:ring-primary shadow-lg mb-4 sm:mb-6">
          Send Reset Link
        </Button>

        <!-- Back to Login Link -->
        <p class="text-center text-xs sm:text-sm text-gray-600">
          Remember your password? <a href="#" class="text-primary hover:underline font-medium">Back to Login</a>
        </p>
      </div>

      <!-- Footer -->
      <div class="mt-6 flex justify-center sm:justify-between text-xs sm:text-sm text-gray-500 w-full max-w-xs sm:max-w-sm md:max-w-md mx-auto px-4">
        <a href="#" class="hover:underline mr-4 sm:mr-0">Privacy Policy</a>
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