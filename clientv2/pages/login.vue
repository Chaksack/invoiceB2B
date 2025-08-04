<script setup>
import { ref } from 'vue';
import Login from '~/components/Views/Auth/Login.vue';
const loading = ref(false);

definePageMeta({
  layout: 'auth'
});
</script>

<template>
    <Login />
</template>
