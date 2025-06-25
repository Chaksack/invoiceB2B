<template>
  <div style="width: 100%; height: 100vh;">
    <!-- Wrap ScalarApiReference in <ClientOnly> to ensure it only runs on the client side -->
    <ClientOnly>
      <ScalarApiReference />
    </ClientOnly>
  </div>
</template>

<script setup>
import ScalarApiReference from '~/components/ScalarApiReference.vue'

// You might want to set a layout for this page if you have one
definePageMeta({
  layout: 'default', // Or whatever layout you use
});
</script>

<style scoped>
/* Ensure the parent div of ScalarApiReference takes full height */
div {
  padding: 0; /* Remove padding from the main div to let Scalar take full space */
  margin: 0;
}
</style>