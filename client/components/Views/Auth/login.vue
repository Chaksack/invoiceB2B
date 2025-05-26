<template>
  <!-- Loader Popup -->
  <div v-if="isLoading" class="fixed inset-0 bg-gray-800 bg-opacity-75 flex items-center justify-center">
    <div class="relative inline-block w-40 h-40">
      <!-- Spinner -->
      <div class="absolute inset-0 w-full h-full border-4 border-lime-500 border-t-transparent rounded-full animate-spin"></div>
    </div>
  </div>

  <div class="flex h-screen flex-col lg:grid lg:grid-cols-[1fr_1.8fr] overflow-hidden">
    <!-- Left side - Form -->
    <div class="flex flex-col items-center justify-center p-6 sm:p-12">
      <form class="mt-10 max-w-[800px]"  @submit.prevent="login">
        <h1 class="text-2xl font-bold mb-4">Sign in to your InvoiceFnd account.</h1>
        <p class="text-muted-foreground mb-6">Welcome back login with your credentials.</p>
        <Alert v-if="errorMessage" variant="default" class="bg-red-500 mb-4 text-white">
          <AlertCircle class="w-4 h-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription class="text-white">{{ errorMessage }}</AlertDescription>
        </Alert>
        <div class="grid gap-4">
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Email Address</Label>
            </div>
            <Input v-model="email" id="email" type="email" placeholder="m@example.com" required />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Password</Label>
            </div>
            <Input v-model="password" id="password" type="password" required />
          </div>
        </div>
        <NuxtLink to="/">
          <p class="font-light text-xs mt-2">Forgotten Password?</p>
        </NuxtLink>
        <button class="py-2 mt-4 rounded-lg bg-bgC text-white bg-black w-full font-medium mb-4 hover:bg-gray-500" :disabled="isLoading">
          Sign In
        </button>
        <p class="text-muted-foreground mb-6">
          New to InvoiceFnd? <a href="/" class="text-black text-semibold">Register</a>
        </p>
      </form>
    </div>

    <!-- Right side - Image -->
    <div class="hidden lg:flex bg-black items-center justify-start h-screen">
      <div class="flex flex-col ml-10">
        <div class="flex flex-col mt-6">
          <h2 class="text-white font-semibold text-6xl">A growth platform for small business</h2>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import axios, { AxiosError } from 'axios';
import { useRouter } from 'vue-router';
import { AlertCircle } from 'lucide-vue-next';
import { toast } from 'vue-sonner';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { useCookie } from '#app';
import { navigateTo } from '#app';

const email = ref('');
const password = ref('');
const isLoading = ref(false);
const errorMessage = ref('');
const router = useRouter();

const login = async () => {
  isLoading.value = true;
  errorMessage.value = '';

  try {
    const response = await axios.post('http://localhost:3000/api/v1/auth/login', {
      email: email.value,
      password: password.value,
    });

    const data = response.data;

    // If server says 2FA required, redirect to 2fa page with email saved in cookie/session
    if (data.twoFactorRequired) {
      const emailCookie = useCookie('2fa_email');
      emailCookie.value = email.value; // store email for 2FA step
      await navigateTo('/2fa');
      return;
    }

    // If login success without 2FA
    if (response.status === 200) {
      const tokenCookie = useCookie('token');
      const roleCookie = useCookie('role');
      tokenCookie.value = data.accessToken || data.token;
      roleCookie.value = data.role;

      toast.success('Login successful');

      if (data.redirectPath) {
        await navigateTo(data.redirectPath);
      } else if (data.role === 'admin') {
        await navigateTo('/admin');
      } else {
        await navigateTo('/home');
      }
    }
  } catch (error) {
    errorMessage.value = 'Invalid email or password';
    toast.error('Login error', { description: error.message });
  } finally {
    isLoading.value = false;
  }
};


</script>
