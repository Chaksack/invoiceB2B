<script setup>
import { ref } from 'vue';
import ForgotPassword from '~/components/Views/Auth/ForgotPassword.vue';
const loading = ref(false);

definePageMeta({
  layout: 'auth'
});
</script>

<template>
    <ForgotPassword />
</template>