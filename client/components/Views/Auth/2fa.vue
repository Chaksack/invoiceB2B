<template>
  <div class="flex flex-col items-center justify-center h-screen">
    <h1 class="text-2xl font-bold mb-4">Two-Factor Authentication</h1>
    <p class="mb-6">Please enter the 6-digit code sent to your email or authenticator app.</p>

    <form @submit.prevent="verify2fa" class="max-w-sm w-full">
      <input
          v-model="code"
          type="text"
          maxlength="6"
          placeholder="Enter 2FA code"
          class="border p-2 w-full text-center text-xl tracking-widest"
          required
      />
      <button
          type="submit"
          class="mt-4 w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
          :disabled="isLoading"
      >
        Verify
      </button>
    </form>

    <Alert v-if="errorMessage" variant="destructive" class="mt-4">
      <AlertTitle>Error</AlertTitle>
      <AlertDescription>{{ errorMessage }}</AlertDescription>
    </Alert>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import axios from 'axios';
import { useCookie, navigateTo } from '#app';
import { toast } from 'vue-sonner';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

const code = ref('');
const isLoading = ref(false);
const errorMessage = ref('');

const emailCookie = useCookie('2fa_email');
const email = emailCookie.value;

if (!email) {
  // If no email in cookie, redirect back to login
  navigateTo('/login');
}

const verify2fa = async () => {
  isLoading.value = true;
  errorMessage.value = '';

  try {
    const response = await axios.post('http://localhost:3000/api/v1/auth/2fa/verify', {
      email,
      code: code.value,
    });

    if (response.status === 200) {
      const data = response.data;
      const tokenCookie = useCookie('token');
      const roleCookie = useCookie('role');

      tokenCookie.value = data.accessToken || data.token;
      roleCookie.value = data.role;

      // Clean up 2fa email cookie
      emailCookie.value = null;

      toast.success('Two-Factor Authentication successful');

      if (data.redirectPath) {
        await navigateTo(data.redirectPath);
      } else if (data.role === 'admin') {
        await navigateTo('/admin');
      } else {
        await navigateTo('/home');
      }
    }
  } catch (error: any) {
    errorMessage.value = error.response?.data?.message || 'Invalid 2FA code';
    toast.error('2FA Error', { description: errorMessage.value });
  } finally {
    isLoading.value = false;
  }
};
</script>
