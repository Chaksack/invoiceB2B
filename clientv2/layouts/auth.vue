<template>
  <div class="min-h-screen bg-background flex items-center justify-center">
    <div class="w-full max-w-md">
      <div class="text-center mb-6">
        <h1 class="text-2xl font-bold">Invoice B2B</h1>
        <p class="text-muted-foreground">Your Business Invoice Solution</p>
      </div>
      
      <slot />
      
      <div class="text-center mt-6 text-sm text-muted-foreground">
        <p>&copy; 2023 Invoice B2B. All rights reserved.</p>
      </div>
    </div>
  </div>
</template>