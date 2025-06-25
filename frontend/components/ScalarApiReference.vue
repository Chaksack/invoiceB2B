<template>
  <ApiReference
      :configuration="{
      spec: {
        url: '/openapi.json',
      },
    }"
  />
</template>

<script setup>
import { ApiReference } from '@scalar/api-reference'
import '@scalar/api-reference/style.css'
</script>

<style>
html, body, #__nuxt, #scalar-api-reference {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
  overflow: hidden; /* Prevent body scroll if Scalar handles its own scroll */
}
</style>