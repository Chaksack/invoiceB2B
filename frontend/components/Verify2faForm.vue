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
          Enter 2fa Authentication
        </CardTitle>
        <CardDescription>
          Enter 2fa code
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form>
          <div class="grid gap-6">
            <div class="grid gap-2">
              <PinInput
                  id="pin-input"
                  v-model="value"
                  placeholder="○"
                  @complete="handleComplete"
              >
                <PinInputGroup class="gap-1">
                  <template v-for="(id, index) in 5" :key="id">
                    <PinInputSlot
                        class="rounded-md border"
                        :index="index"
                    />
                    <template v-if="index !== 4">
                      <PinInputSeparator />
                    </template>
                  </template>
                </PinInputGroup>
              </PinInput>
            </div>

            <Button type="submit" class="w-full">
              Submit
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
