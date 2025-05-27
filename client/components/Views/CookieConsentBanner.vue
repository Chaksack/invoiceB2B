<template>
  <div v-if="!consentGiven" class="relative isolate flex cookie-banner text-white items-center gap-x-6 overflow-hidden bg-gray-50 px-6 py-2.5 sm:px-3.5 sm:before:flex-1">
    <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
      <p class="text-sm/6">
        <strong class="font-semibold">Invoicefnd</strong><svg viewBox="0 0 2 2" class="mx-2 inline size-0.5 fill-current" aria-hidden="true"><circle cx="1" cy="1" r="1" /></svg>This website uses cookies. Do you accept?
      </p>
      <button @click="acceptCookies" class="flex-none rounded-full bg-white px-3.5 py-1 text-sm font-semibold text-black shadow-xs hover:bg-gray-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900">Accept <span aria-hidden="true">&rarr;</span></button>
    </div>
    <div class="flex flex-1 justify-end">
      <button type="button" class="-m-3 p-3 focus-visible:outline-offset-[-4px]">
        <span class="sr-only">Dismiss</span>
        <XMarkIcon class="size-5 text-white" aria-hidden="true" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { XMarkIcon } from '@heroicons/vue/20/solid'

const consentGiven = ref(false);

const emit = defineEmits(['consentGiven']);

const acceptCookies = () => {
  consentGiven.value = true;
  emit('consentGiven');
};
</script>

<style scoped>
.cookie-banner {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  background-color: #005234;
  padding: 10px;
  text-align: center;
}
</style>