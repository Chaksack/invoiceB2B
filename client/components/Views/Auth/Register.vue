<template>
  <div class="flex h-screen flex-col lg:grid lg:grid-cols-[1fr_1.8fr] overflow-hidden">
    <!-- Left side - Form -->
    <div class="flex flex-col items-center justify-center p-6 sm:p-12">
      <div class="mt-10 max-w-[800px]">
        <h1 class="text-2xl font-bold mb-4">
          Sign Up Account.
        </h1>
        <p class="text-muted-foreground mb-6">
          Enter your personal data to create your account.
        </p>
        <div class="grid gap-4">
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="first_name">First Name</Label>
            </div>
            <Input
                id="first_name"
                type="name"
                placeholder="John"
                required
            />
          </div>

          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="last_name">Last Name</Label>
            </div>
            <Input
                id="last_name"
                type="name"
                placeholder="Doe"
                required
            />
          </div>

          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="company">Company Name</Label>
            </div>
            <Input
                id="company"
                type="text"
                placeholder="Company Ltd"
                required
            />
          </div>
          <!-- Email Input -->
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="email">Email Address</Label>
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
        <NuxtLink to="/home">
          <button class="py-2 mt-4 rounded-lg bg-bgC text-white bg-black w-full font-medium mb-4 hover:bg-gray-500">
            Create Account
          </button>
          <p class="text-muted-foreground mb-6">
            Already have an InvoiceFnd account? <a href="/login" class="text-black text-semibold">Login</a>
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
