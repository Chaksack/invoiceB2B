<template>
  <header class="sticky top-0 flex h-12 items-center justify-between text-white bg-black px-6">
    <h1 class="text-xl font-semibold font-mono">
      InvoiceB2B
    </h1>
    <NavigationMenu>
      <NavigationMenuList>

        <LazyNuxtLink to="/">
          <NavigationMenuItem class="px-2 text-xs" >Contact us
          </NavigationMenuItem>
        </LazyNuxtLink>

        <LazyNuxtLink to="/">
          <NavigationMenuItem class="px-2 text-xs" >About us
          </NavigationMenuItem>
        </LazyNuxtLink>
      <LazyNuxtLink to="/home">
        <NavigationMenuItem class="px-2 text-xs" >Home
        </NavigationMenuItem>
      </LazyNuxtLink>



      </NavigationMenuList>
    </NavigationMenu>



    <div class="flex items-center">
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
import {useCookie} from "#app";
import axios from "axios";
import {onMounted, ref} from "vue";


</script>