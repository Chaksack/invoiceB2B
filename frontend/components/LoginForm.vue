<script setup lang="ts">
import {Button} from "~/components/ui/button";
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'vue-sonner'
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "~/components/ui/card";
</script>

<template>
  <div class="flex flex-col gap-6">
    <Card>
      <CardHeader class="text-center">
        <CardTitle class="text-xl">
          Welcome back
        </CardTitle>
        <CardDescription>
          Log in to your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form>
          <div class="grid gap-6">
            <div class="grid gap-6">
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
                  <a
                      href="/forgot"
                      class="ml-auto text-sm underline-offset-4 hover:underline"
                  >
                    Forgot your password?
                  </a>
                </div>
                <Input id="password" type="password" required />
              </div>

              <Button type="submit" class="w-full">
                Login
              </Button>
            </div>
            <div class="text-center text-sm">
              Don't have an account?
              <a href="/register" class="underline underline-offset-4">
                Sign up
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
