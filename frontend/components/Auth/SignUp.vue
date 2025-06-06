<script setup lang="ts">
import {Button} from "~/components/ui/button";
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "~/components/ui/card";
</script>

<template>
  <div class="flex flex-col gap-6">
    <Card>
      <CardHeader class="text-center">
        <CardTitle class="text-xl">
          Welcome to Invoicing Financing
        </CardTitle>
        <CardDescription>
          Enter your details to create your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form>
          <div class="grid gap-6">
            <div class="grid gap-6">
              <div class="grid gap-2">
                <Label html-for="firstName">First Name</Label>
                <Input
                    id="firstName"
                    type="name"
                    placeholder="john"
                    required
                />
              </div>
              <div class="grid gap-2">
                <Label html-for="lastName">Last Name</Label>
                <Input
                    id="lasstName"
                    type="name"
                    placeholder="doe"
                    required
                />
              </div>
              <div class="grid gap-2">
                <Label html-for="companyName">Company Name</Label>
                <Input
                    id="companyName"
                    type="name"
                    placeholder="Acme Inc"
                    required
                />
              </div>
              <div class="grid gap-2">
                <Label html-for="email">Email</Label>
                <Input
                    id="email"
                    type="email"
                    placeholder="m@example.com"
                    required
                />
              </div>
              <div class="grid gap-2">
                <div class="flex items-center">
                  <Label html-for="password">Password</Label>
                </div>
                <PasswordInput id="password" />
              </div>
              <div class="grid gap-2">
                <div class="flex items-center">
                  <Label html-for="confirm-password">Confirm Password</Label>
                </div>
                <PasswordInput id="confirm-password" />
              </div>
              <Button type="submit" class="w-full">
                Create account
              </Button>
            </div>
            <div class="text-center text-sm">
              Already have an account?
              <a href="/" class="underline underline-offset-4">
                Sign in
              </a>
            </div>
          </div>
        </form>
      </CardContent>
    </Card>
    <div class="text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 [&_a]:hover:text-primary">
      By clicking continue, you agree to our <a href="#">Terms of Service</a>
      and <a href="#">Privacy Policy</a>.
    </div>
  </div>
</template>
