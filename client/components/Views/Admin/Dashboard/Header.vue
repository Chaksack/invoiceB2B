<template>
  <header class="sticky top-0 flex h-12 items-center justify-between text-white bg-indigo-500 px-6">
    <h1 class="text-xl font-semibold font-mono">
      InvoiceB2B
    </h1>
    <NavigationMenu>
      <NavigationMenuList>
        <LazyNuxtLink to="/admin/home">
        <NavigationMenuItem class="px-2 text-xs" >Dashboard
        </NavigationMenuItem>
        </LazyNuxtLink>
        <LazyNuxtLink to="/admin/invoices">
          <NavigationMenuItem class="px-2 text-xs" >Invoices
          </NavigationMenuItem>
        </LazyNuxtLink>
        <LazyNuxtLink to="/admin/customers">
          <NavigationMenuItem class="px-2 text-xs" >Customers
          </NavigationMenuItem>
        </LazyNuxtLink>
<!--        <NavigationMenuItem>-->
<!--          <NavigationMenuTrigger class="bg-transparent text-xs">Invoices</NavigationMenuTrigger>-->
<!--          <NavigationMenuContent class="bg-gray-900 text-white md:w-[400px] lg:w-[500px] ">-->
<!--            <LazyNuxtLink to="/admin/invoices/pending">-->
<!--            <NavigationMenuLink>Pending</NavigationMenuLink>-->
<!--            </LazyNuxtLink>-->
<!--            <LazyNuxtLink to="">-->
<!--            <NavigationMenuLink>Approved</NavigationMenuLink>-->
<!--            </LazyNuxtLink>-->
<!--            <LazyNuxtLink to="">-->
<!--            <NavigationMenuLink>Disbursed</NavigationMenuLink>-->
<!--            </LazyNuxtLink>-->
<!--          </NavigationMenuContent>-->
<!--        </NavigationMenuItem>-->
<!--<NavigationMenuItem>-->
<!--        <NavigationMenuTrigger class="bg-transparent text-xs">Loans</NavigationMenuTrigger>-->
<!--          <NavigationMenuContent class="bg-gray-900 text-white md:w-[400px] lg:w-[500px] ">-->
<!--            <LazyNuxtLink to="">-->
<!--            <NavigationMenuLink>Active</NavigationMenuLink>-->
<!--            </LazyNuxtLink>-->
<!--            <LazyNuxtLink to="">-->
<!--            <NavigationMenuLink>Paid</NavigationMenuLink>-->
<!--            </LazyNuxtLink>-->
<!--          </NavigationMenuContent>-->
<!--        </NavigationMenuItem>-->

        <LazyNuxtLink to="">
        <NavigationMenuItem class="px-2 text-xs">Profile
        </NavigationMenuItem>
        </LazyNuxtLink>
      </NavigationMenuList>
    </NavigationMenu>
    <div class="flex items-center">
      <p class="pr-2">Andrew Chakdahah</p>
      <Avatar class="relative overflow-visible">
        <AvatarFallback class="text-black"> AC </AvatarFallback>
      </Avatar>
    </div>
  </header>
</template>
<script setup lang="ts">
import {Avatar, AvatarFallback} from "~/components/ui/avatar";
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuIndicator,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  NavigationMenuViewport,
} from '~/components/ui/navigation-menu'
</script>