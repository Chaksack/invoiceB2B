<template>
  <div>
    <NuxtLayout>
      <CookieConsentBanner v-if="!cookieConsentGiven" @consentGiven="handleConsentGiven" />
      <Toaster richColors position="top-right" />
      <NuxtLoadingIndicator />
    <NuxtPage />
    </NuxtLayout >
  </div>
</template>
<script setup lang="ts">
import { Toaster } from 'vue-sonner';
import CookieConsentBanner from '~/components/Views/CookieConsentBanner.vue';
</script>