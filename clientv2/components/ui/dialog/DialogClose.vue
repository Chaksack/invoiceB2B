<script setup>
import { inject } from 'vue';

// Get parent Dialog component's open state
const dialogOpen = inject('dialogOpen', null);

const closeDialog = () => {
  if (dialogOpen) {
    dialogOpen.value = false;
  }
};
</script>

<template>
  <button 
    @click="closeDialog"
    class="inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none ring-offset-background"
  >
    <slot></slot>
  </button>
</template>