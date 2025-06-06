<script setup lang="ts">
import type { NuxtError } from '#app'

defineProps({
  error: Object as () => NuxtError
})

const handleError = () => clearError({ redirect: '/' })
</script>

<template>
  <div class="h-svh">
    <div class="m-auto h-full w-full flex flex-col items-center justify-center gap-2">
      <h1 class="text-[7rem] font-bold leading-tight">
        {{ error.statusCode }}
      </h1>
      <span class="font-medium">
        Ouch! {{ error.statusMessage || 'An error occurred.' }}
</span>
      <p class="text-center text-muted-foreground">
        {{ error.message }}
      </p>
      <div class="mt-6 flex gap-4">
        <Button variant="outline" @click="router.back()">
          Go Back
        </Button>
        <Button
            @click="handleError">
                        Back to Home
        </Button>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>