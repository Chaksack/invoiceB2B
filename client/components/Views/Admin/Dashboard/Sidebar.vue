<script setup lang="ts">
import { Bell, CircleUser, Home, LineChart, Menu, Package, Package2, Building2, Search, ShoppingCart, Users, Landmark } from 'lucide-vue-next'

</script>

<template>
  <div class="grid h-screen w-full md:grid-cols-[220px_1fr] lg:grid-cols-[300px_1fr]">
    <div class="hidden border-r bg-muted/40 md:block">
      <div class="flex bg-blue-500 h-full max-h-screen flex-col gap-2 ">
        <div class="flex h-14 items-center px-4 lg:h-[60px] lg:px-6">
          <div class="flex text-sm text-white italic font-semibold tracking-wide">
            <img src="" class="h-12"><span class="text-white ml-2 text-2xl">InvoiceFnd</span>
          </div>
        </div>
        <div class="flex-1 text-white">
          <nav class="grid gap-4 items-start text-sm">
            <NuxtLink
                to="/admin/home"
                class="flex items-center text-xs gap-4 px-4 py-2 transition-all hover:text-white"
                exact-active-class="bg-white text-blue-500"
            >
              <Home class="h-4 w-4" />
              Dashboard
            </NuxtLink>
            <NuxtLink
                to="/admin/Invoices"
                class="flex items-center text-xs gap-4 px-4 py-2 transition-all hover:text-white"
                exact-active-class="bg-white text-blue-500"
            >
              <Home class="h-4 w-4" />
              Invoices
            </NuxtLink>
            <NuxtLink
                to="/admin/customers"
                class="flex items-center text-xs gap-4 px-4 py-2 transition-all hover:text-white"
                exact-active-class="bg-white text-blue-500"
            >
              <Home class="h-4 w-4" />
              Users
            </NuxtLink>




          </nav>
        </div>

      </div>
    </div>



  </div>
</template>
