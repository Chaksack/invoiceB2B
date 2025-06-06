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
          Reset Password
        </CardTitle>
        <CardDescription>
          Enter email address to reset password
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form>
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

              <Button type="submit" class="w-full">
                Send password reset link
              </Button>
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
