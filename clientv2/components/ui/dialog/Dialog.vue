<script setup>
import { ref, watch, onMounted, onBeforeUnmount, provide } from 'vue';

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:open']);

const dialogOpen = ref(props.open);

// Provide dialogOpen to child components
provide('dialogOpen', dialogOpen);

watch(() => props.open, (newValue) => {
  dialogOpen.value = newValue;
});

watch(dialogOpen, (newValue) => {
  emit('update:open', newValue);
});

// Close dialog when Escape key is pressed
const handleKeyDown = (event) => {
  if (event.key === 'Escape' && dialogOpen.value) {
    dialogOpen.value = false;
  }
};

// Add event listener when component is mounted
onMounted(() => {
  document.addEventListener('keydown', handleKeyDown);
});

// Remove event listener when component is unmounted
onBeforeUnmount(() => {
  document.removeEventListener('keydown', handleKeyDown);
});
</script>

<template>
  <div>
    <slot />
  </div>
</template>