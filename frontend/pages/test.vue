<template>
  <div class="min-h-screen bg-gray-100 flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8 font-inter">
    <div class="max-w-md w-full bg-white p-8 rounded-2xl shadow-xl space-y-8">
      <div class="text-center">
        <h2 class="mt-6 text-3xl font-extrabold text-gray-900 rounded-lg">
          Invoice Financing Platform
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          Connects businesses with financial institutions
        </p>
      </div>

      <div class="mt-8 space-y-6">
        <div>
          <button @click="testApi"
                  class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 shadow-md hover:shadow-lg rounded-xl">
            Test Backend API
          </button>
        </div>
        <p v-if="apiResponse" class="mt-2 text-center text-sm text-gray-700 p-3 bg-gray-50 rounded-lg shadow-sm">
          API Response: {{ apiResponse }}
        </p>
        <p v-if="apiError" class="mt-2 text-center text-sm text-red-600 p-3 bg-red-50 rounded-lg shadow-sm">
          API Error: {{ apiError }}
        </p>
      </div>

      <div class="text-center text-sm text-gray-500">
        <p>Get started by setting up your PostgreSQL database and populating the `.env` file.</p>
        <p class="mt-2">For development, run: `npm run dev`</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { NuxtAxiosInstance } from '@nuxtjs/axios'; // Import NuxtAxiosInstance type

declare module 'vue/types/vue' {
  interface Vue {
    $axios: NuxtAxiosInstance;
  }
}

interface DataType {
  apiResponse: string;
  apiError: string;
}

export default Vue.extend({
  data(): DataType {
    return {
      apiResponse: '',
      apiError: '',
    };
  },
  methods: {
    async testApi() {
      this.apiResponse = '';
      this.apiError = '';
      try {
        const response = await this.$axios.$get('/api/hello');
        this.apiResponse = response.message;
      } catch (error: any) {
        console.error('API Test Error:', error);
        this.apiError = error.response ? error.response.data.message : error.message;
      }
    },
  },
});
</script>

<style>
/* You can define global styles or import Tailwind directives here */
/* In assets/css/main.css or directly in a <style> tag */
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap');

body {
  font-family: 'Inter', sans-serif;
}
</style>