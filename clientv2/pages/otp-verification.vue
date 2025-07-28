<script setup>
import { ref } from 'vue';
import OtpVerification from '~/components/Views/Auth/OtpVerification.vue';
const loading = ref(false);

definePageMeta({
  layout: 'auth'
});
</script>

<template>
    <OtpVerification />
</template>