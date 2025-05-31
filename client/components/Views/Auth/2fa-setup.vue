<template>
  <div class="max-w-md mx-auto p-6 mt-10 border rounded shadow">
    <h2 class="text-xl font-bold mb-4">Two-Factor Authentication Setup</h2>

    <div v-if="!is2faEnabled">
      <button
          @click="generate2fa"
          class="mb-4 bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
          :disabled="isLoading"
      >
        Enable 2FA
      </button>

      <div v-if="qrCodeUrl" class="mb-4 text-center">
        <p class="mb-2">Scan this QR code with your Authenticator app:</p>
        <img :src="qrCodeUrl" alt="2FA QR Code" class="mx-auto" />
      </div>

      <form v-if="qrCodeUrl" @submit.prevent="verifyCode">
        <label for="code" class="block mb-1">Enter the 6-digit code from your app:</label>
        <input
            v-model="code"
            id="code"
            type="text"
            maxlength="6"
            required
            class="border p-2 w-full mb-4 text-center tracking-widest"
        />
        <button
            type="submit"
            class="w-full bg-green-600 text-white py-2 rounded hover:bg-green-700"
            :disabled="isLoading"
        >
          Verify & Enable 2FA
        </button>
      </form>
    </div>

    <div v-else>
      <p class="mb-4 text-green-700 font-semibold">Two-Factor Authentication is currently <strong>enabled</strong>.</p>
      <button
          @click="disable2fa"
          class="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
          :disabled="isLoading"
      >
        Disable 2FA
      </button>
    </div>

    <p v-if="message" :class="messageClass" class="mt-4">{{ message }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import axios from 'axios';

const qrCodeUrl = ref('');
const code = ref('');
const isLoading = ref(false);
const message = ref('');
const messageClass = ref('');
const is2faEnabled = ref(false);

const fetch2faStatus = async () => {
  try {
    const response = await axios.get('http://localhost:3000/api/v1/auth/2fa/status', { withCredentials: true });
    is2faEnabled.value = response.data.enabled;
  } catch (err) {
    console.error('Failed to fetch 2FA status', err);
  }
};

const generate2fa = async () => {
  isLoading.value = true;
  message.value = '';
  messageClass.value = '';

  try {
    const response = await axios.post('http://localhost:3000/api/v1/auth/2fa/generate', null, { withCredentials: true });
    qrCodeUrl.value = response.data.qrCodeUrl;
  } catch (err: any) {
    message.value = err.response?.data?.message || 'Failed to generate 2FA secret';
    messageClass.value = 'text-red-600';
  } finally {
    isLoading.value = false;
  }
};

const verifyCode = async () => {
  if (code.value.length !== 6) {
    message.value = 'Please enter a valid 6-digit code';
    messageClass.value = 'text-red-600';
    return;
  }

  isLoading.value = true;
  message.value = '';
  messageClass.value = '';

  try {
    const response = await axios.post(
        'http://localhost:3000/api/v1/auth/2fa/enable',
        { code: code.value },
        { withCredentials: true }
    );
    is2faEnabled.value = true;
    qrCodeUrl.value = '';
    code.value = '';
    message.value = 'Two-Factor Authentication has been enabled successfully!';
    messageClass.value = 'text-green-600';
  } catch (err: any) {
    message.value = err.response?.data?.message || 'Invalid 2FA code';
    messageClass.value = 'text-red-600';
  } finally {
    isLoading.value = false;
  }
};

const disable2fa = async () => {
  isLoading.value = true;
  message.value = '';
  messageClass.value = '';

  try {
    await axios.post('http://localhost:3000/api/v1/auth/2fa/disable', null, { withCredentials: true });
    is2faEnabled.value = false;
    message.value = 'Two-Factor Authentication has been disabled.';
    messageClass.value = 'text-green-600';
  } catch (err: any) {
    message.value = err.response?.data?.message || 'Failed to disable 2FA';
    messageClass.value = 'text-red-600';
  } finally {
    isLoading.value = false;
  }
};

fetch2faStatus();
</script>
