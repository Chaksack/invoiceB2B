<template>
  <div class="min-h-screen bg-background">
    <header class="bg-card shadow-sm fixed top-0 left-0 right-0 z-10">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <h1 class="text-xl font-bold">Invoice B2B</h1>
        <nav class="flex space-x-4">
          <NuxtLink to="/" class="text-sm font-medium hover:text-primary">Home</NuxtLink>
          <NuxtLink to="/login" class="text-sm font-medium hover:text-primary">Login</NuxtLink>
          <NuxtLink to="/register" class="text-sm font-medium hover:text-primary">Register</NuxtLink>
        </nav>
      </div>
    </header>
    
    <main class="container mx-auto px-4 py-8 mt-16">
      <slot />
    </main>
    
    <footer class="bg-card shadow-sm mt-auto">
      <div class="container mx-auto px-4 py-4 text-center text-muted-foreground">
        <p>&copy; 2023 Invoice B2B. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>