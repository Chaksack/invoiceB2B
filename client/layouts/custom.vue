<template>
    <div class="flex w-full z-50">
      <ViewsAdminDashboardSidebar class="fixed top-0 z-0 w-[calc(100%-16rem)] h-16" />
      <div class="items-center justify-center z-20">
        <slot />
      </div>
    </div>
</template>

<script setup>
</script>
