<script setup>
import { ref } from 'vue';
import Register from '~/components/Views/Auth/Register.vue';
const loading = ref(false);

definePageMeta({
  layout: 'auth'
});
</script>

<template>
    <Register />
</template>
