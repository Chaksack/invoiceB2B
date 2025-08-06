<template>
  <div>
    <Toaster 
      richColors 
      position="top-right" 
      :expand="false"
      closeButton
      :duration="4000"
      theme="light"
    />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>
<script setup lang="ts">
import { Toaster } from 'vue-sonner';
import { useHead } from '#app';

// Add viewport meta tag for responsive design
useHead({
  meta: [
    { name: 'viewport', content: 'width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no' }
  ]
});
</script>