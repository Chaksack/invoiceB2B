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
          Create your new password
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form>
          <div class="grid gap-6">
            <div class="grid gap-6">
              <div class="grid gap-2">
                <div class="flex items-center">
                  <Label html-for="password">New Password</Label>
                </div>
                <Input id="password" type="password" required />
              </div>
              <div class="grid gap-2">
                <div class="flex items-center">
                  <Label html-for="password">Confirm Password</Label>
                </div>
                <Input id="password" type="password" required />
              </div>
              <Button type="submit" class="w-full">
                Submit
              </Button>
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
