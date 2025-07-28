<template>
  <div>
    <Toaster richColors position="top-right" />
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>
<script setup lang="ts">
import { Toaster } from 'vue-sonner';
</script>