<template>
  <div class="flex h-screen flex-col lg:grid lg:grid-cols-[1fr_1.8fr] overflow-hidden">
    <!-- Left side - Form -->
    <div class="flex flex-col items-center justify-center p-6 sm:p-12">
      <div class="mt-10 max-w-[800px]">
        <h1 class="text-2xl font-bold mb-4">
          Sign in to your InvoiceB2B account.
        </h1>
        <p class="text-muted-foreground mb-6">
          Welcome back login with your credentials.
        </p>
        <div class="grid gap-4">
          <!-- Email Input -->
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Email Address</Label>
            </div>
            <Input
                id="email"
                type="email"
                placeholder="m@example.com"
                required
            />
          </div>
          <!-- Password Input -->
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Password</Label>
            </div>
            <Input
                id="password"
                type="password"
                required
            />
          </div>
        </div>
        <NuxtLink to="/">
          <p class="font-light text-xs mt-2" >
            Forgotten Password?
          </p>
        </NuxtLink>
        <NuxtLink to="/home">
          <button class="py-2 mt-4 rounded-lg bg-bgC text-white bg-black w-full font-medium mb-4 hover:bg-gray-500">
            Sign In
          </button>
          <p class="text-muted-foreground mb-6">
            New to InvoiceB2B? <a href="/" class="text-black text-semibold">Register</a>
          </p>
        </NuxtLink>
      </div>

    </div>

    <!-- Right side - Image -->
    <div class="hidden lg:flex bg-black items-center justify-start h-screen">
      <div class="flex flex-col ml-10">
<!--        <div class="flex text-sm text-white italic font-semibold tracking-wide">-->
<!--          <img src="../../../assets/logo.png" class="h-20"><span class="text-white ml-2">OBLOS</span>-->
<!--        </div>-->
        <div class="flex flex-col mt-6">
          <h2 class="text-white font-semibold text-6xl">A growth platform for small business</h2>
        </div>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
</script>
