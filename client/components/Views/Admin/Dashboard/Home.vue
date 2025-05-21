<script setup lang="ts">
import { ref } from 'vue';
// Assuming AlertCircle might be used elsewhere or was part of a larger context.
// If not used in this specific template, it could be removed.
import { AlertCircle } from 'lucide-vue-next';

// Assuming these components are correctly imported and registered in your project.
// For example, globally or via a plugin.
// import { Avatar, AvatarFallback } from '@/components/ui/avatar'; // Example if they were missing
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '~/components/ui/card'
// import {
//   Table,
//   TableBody,
//   TableCaption,
//   TableCell,
//   TableHead,
//   TableHeader,
//   TableRow,
// } from '@/components/ui/table' // Table components are imported but not used in the template.

const isLoading = ref(false); // isLoading is defined but not used in the template.

const usages = [
  {
    username: 'johndoe',
    email: 'johndoe@test.com',
    avatarfallback: 'JD',
    count: '15',
  },
  {
    username: 'janedoe',
    email: 'janedoe@test.com',
    avatarfallback: 'JD',
    count: '35',
  },
  {
    username: 'johnsmith',
    email: 'johnsmith@test.com',
    avatarfallback: 'JS',
    count: '67',
  },
  {
    username: 'janesmith',
    email: 'janesmith@test.com',
    avatarfallback: 'JS',
    count: '32',
  },
];

// Placeholder for Avatar and AvatarFallback if they are not globally registered
// This is just for making the snippet runnable if these components are local.
// In a real project, these would come from your UI library (e.g., @/components/ui/avatar)
const Avatar = { template: '<div class="avatar-placeholder"><slot /></div>' };
const AvatarFallback = { template: '<div class="avatar-fallback-placeholder"><slot /></div>' };

</script>

<template>
  <div class="flex min-h-screen bg-white text-black justify-center font-sans">
    <main class="flex-1 items-start gap-4 p-4 md:p-6 lg:p-8">
      <div class="mx-auto grid w-full max-w-6xl flex-1 auto-rows-max gap-6 lg:gap-8">

        <section class="flex justify-center w-full">
          <div class="grid w-full max-w-4xl items-start justify-items-center gap-4 sm:grid-cols-2 md:grid-cols-3">
            <Card class="bg-gray-900 w-full shadow-lg rounded-lg text-white dark:bg-gray-800 dark:text-white">
              <CardHeader class="p-4">
                <CardDescription class="font-medium text-sm text-white dark:text-gray-400">INVOICES</CardDescription>
                <CardTitle class="text-2xl font-bold mt-1">
                  10
                </CardTitle>
              </CardHeader>
            </Card>
            <Card class="bg-gray-900 w-full shadow-lg rounded-lg text-white dark:bg-gray-800 dark:text-white">
              <CardHeader class="p-4">
                <CardDescription class="font-medium text-sm text-white dark:text-gray-400">DISBURSED</CardDescription>
                <CardTitle class="text-2xl font-bold mt-1">
                  GHS 15,500
                </CardTitle>
              </CardHeader>
            </Card>
            <Card class="bg-gray-900 w-full shadow-lg rounded-lg text-white dark:bg-gray-800 dark:text-white">
              <CardHeader class="p-4">
                <CardDescription class="font-medium text-sm text-white dark:text-gray-400">LOAN DUE AMOUNT</CardDescription>
                <CardTitle class="text-2xl font-bold mt-1">
                  GHS 18,000
                </CardTitle>
              </CardHeader>
            </Card>
          </div>
        </section>

        <div class="grid gap-6 lg:gap-8 md:grid-cols-2">
          <Card class="bg-gray-900 shadow-lg rounded-lg text-white">
            <CardHeader class="p-4">
              <CardTitle class="text-lg font-semibold">
                Exceptions
              </CardTitle>
              {/* Placeholder content */}
              <CardContent class="pt-2 text-gray-400">
                <p>Last 24h: 3</p>
                <p>Critical: 1</p>
              </CardContent>
            </CardHeader>
          </Card>
          <Card class="bg-gray-900 shadow-lg rounded-lg text-white">
            <CardHeader class="p-4">
              <CardTitle class="text-lg font-semibold">
                Routes
              </CardTitle>
              {/* Placeholder content */}
              <CardContent class="pt-2 text-gray-400">
                <p>Monitored: 25</p>
                <p>Avg. RPS: 300</p>
              </CardContent>
            </CardHeader>
          </Card>
        </div>

      </div>
    </main>
  </div>
</template>

<style scoped>
/* Scoped styles for the component if needed */
/* Fallback for Avatar if component is not fully styled by library */
.avatar-placeholder {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 9999px; /* full */
  background-color: #cbd5e1; /* bg-slate-300 */
  color: #1e293b; /* text-slate-700 */
}
.avatar-fallback-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

/* Ensure full height for the main layout */
.flex.min-h-screen {
  min-height: 100vh;
}

/* General font for readability */
.font-sans {
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
}
</style>