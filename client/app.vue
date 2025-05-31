<template>
  <div>
    <NuxtLayout>
      <Toaster richColors position="top-right" />
      <NuxtLoadingIndicator />
    <NuxtPage />
    </NuxtLayout >
  </div>
</template>
<script setup lang="ts">
import { Toaster } from 'vue-sonner';
</script>