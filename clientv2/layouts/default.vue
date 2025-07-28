<template>
  <div class="min-h-screen bg-background">
    <header class="bg-card shadow-sm">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <h1 class="text-xl font-bold">Invoice B2B</h1>
        <nav>
          <!-- Navigation links would go here -->
        </nav>
      </div>
    </header>
    
    <main class="container mx-auto px-4 py-8">
      <slot />
    </main>
    
    <footer class="bg-card shadow-sm mt-auto">
      <div class="container mx-auto px-4 py-4 text-center text-muted-foreground">
        <p>&copy; 2023 Invoice B2B. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>