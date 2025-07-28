<script setup>
import { ref } from 'vue';

const name = ref('');
const email = ref('');
const password = ref('');
const confirmPassword = ref('');
const loading = ref(false);

const handleRegister = async () => {
  if (password.value !== confirmPassword.value) {
    // Handle password mismatch
    return;
  }

  loading.value = true;
  // Add registration logic here
  loading.value = false;
};
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-background">
    <div class="w-full max-w-md p-8 space-y-8 bg-card rounded-lg shadow-md">
      <div class="text-center">
        <h1 class="text-2xl font-bold">Create an Account</h1>
        <p class="text-muted-foreground">Register to access the platform</p>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-6">
        <div class="space-y-2">
          <label for="name" class="text-sm font-medium">Full Name</label>
          <input 
            id="name" 
            v-model="name" 
            type="text" 
            required 
            class="w-full px-3 py-2 border rounded-md"
            placeholder="Enter your full name"
          />
        </div>

        <div class="space-y-2">
          <label for="email" class="text-sm font-medium">Email</label>
          <input 
            id="email" 
            v-model="email" 
            type="email" 
            required 
            class="w-full px-3 py-2 border rounded-md"
            placeholder="Enter your email"
          />
        </div>

        <div class="space-y-2">
          <label for="password" class="text-sm font-medium">Password</label>
          <input 
            id="password" 
            v-model="password" 
            type="password" 
            required 
            class="w-full px-3 py-2 border rounded-md"
            placeholder="Create a password"
          />
        </div>

        <div class="space-y-2">
          <label for="confirmPassword" class="text-sm font-medium">Confirm Password</label>
          <input 
            id="confirmPassword" 
            v-model="confirmPassword" 
            type="password" 
            required 
            class="w-full px-3 py-2 border rounded-md"
            placeholder="Confirm your password"
          />
        </div>

        <button 
          type="submit" 
          class="w-full py-2 text-white bg-primary rounded-md hover:bg-primary-700"
          :disabled="loading"
        >
          {{ loading ? 'Registering...' : 'Register' }}
        </button>
      </form>

      <div class="text-center">
        <p class="text-sm text-muted-foreground">
          Already have an account? 
          <NuxtLink to="/login" class="text-primary hover:underline">Login</NuxtLink>
        </p>
      </div>
    </div>
  </div>
</template>
