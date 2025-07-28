<script setup>
import { onMounted } from 'vue';
import { navigateTo } from '#app';

// Redirect to login page
onMounted(() => {
  navigateTo('/login');
});

definePageMeta({
  layout: 'default'
});
</script>

<template>
  <div class="flex items-center justify-center h-screen">
    <p>Redirecting to login page...</p>
  </div>
</template>
