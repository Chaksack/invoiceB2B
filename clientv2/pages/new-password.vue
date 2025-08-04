<script setup>
import { ref } from 'vue';
import NewPassword from '~/components/Views/Auth/NewPassword.vue';
const loading = ref(false);

definePageMeta({
  layout: 'auth'
});
</script>

<template>
    <NewPassword />
</template>