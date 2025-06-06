<script setup lang="ts">
import { h } from 'vue'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'
import { toast } from 'vue-sonner'

import { cn } from '@/lib/utils'
import { Button, buttonVariants } from '@/components/ui/button'
import { FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Separator } from '@/components/ui/separator'

// Define the form schema using Zod for validation
const appearanceFormSchema = toTypedSchema(z.object({
  theme: z.enum(['light', 'dark'], {
    required_error: 'Please select a theme.',
  }),
  font: z.enum(['inter', 'manrope', 'system'], {
    invalid_type_error: 'Select a font',
    required_error: 'Please select a font.',
  }),
}))

// Get the color mode controller from Nuxt
const color = useColorMode()

// Set up the form using VeeValidate's composable
const { handleSubmit } = useForm({
  validationSchema: appearanceFormSchema,
  // Sync initial form value with the app's current theme for better UX
  initialValues: {
    theme: color.value,
    font: 'inter',
  },
})

// Define the submission handler
const onSubmit = handleSubmit((values) => {
  // Show a toast notification with the submitted values
  toast({
    title: 'You submitted the following values:',
    description: h('pre', { class: 'mt-2 w-[340px] rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white' }, JSON.stringify(values, null, 2))),
  })

  // Update the application's color theme based on form input
  color.preference = values.theme
})
</script>

<template>
  <div>
    <h3 class="text-lg font-medium">
      Appearance
    </h3>
    <p class="text-sm text-muted-foreground">
      Customize the appearance of the app. Automatically switch between day and night themes.
    </p>
  </div>
  <Separator />
  <form class="space-y-8" @submit="onSubmit">
    <FormField v-slot="{ field }" name="font">
      <FormItem>
        <FormLabel>Font</FormLabel>
        <div class="relative w-[200px]">
          <FormControl>
            <select
                :class="cn(
                buttonVariants({ variant: 'outline' }),
                'w-[200px] appearance-none bg-transparent font-normal',
              )"
                v-bind="field"
            >
              <option value="inter">
                Inter
              </option>
              <option value="manrope">
                Manrope
              </option>
              <option value="system">
                System
              </option>
            </select>
          </FormControl>
          <!-- FIX: Replaced <Icon> with an inline SVG for portability -->
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 15 15" class="pointer-events-none absolute right-3 top-2.5 h-4 w-4 opacity-50"><path fill="currentColor" d="M3.135 6.158a.5.5 0 0 1 .707-.023L7.5 9.565l3.658-3.43a.5.5 0 0 1 .684.728l-4 3.75a.5.5 0 0 1-.684 0l-4-3.75a.5.5 0 0 1-.023-.707Z"></path></svg>
        </div>
        <FormDescription>
          Set the font you want to use in the dashboard.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ componentField }" name="theme">
      <FormItem class="space-y-1">
        <FormLabel>Theme</FormLabel>
        <FormDescription>
          Select the theme for the dashboard.
        </FormDescription>
        <FormMessage />

        <RadioGroup
            class="grid max-w-md grid-cols-2 gap-8 pt-2"
            v-bind="componentField"
        >
          <FormItem>
            <FormLabel class="[&:has([data-state=checked])>div]:border-primary">
              <FormControl>
                <RadioGroupItem value="light" class="sr-only" />
              </FormControl>
              <div class="items-center rounded-md border-2 border-muted p-1 hover:border-accent">
                <div class="space-y-2 rounded-sm bg-[#ecedef] p-2">
                  <div class="space-y-2 rounded-md bg-white p-2 shadow-sm">
                    <div class="h-2 w-20 rounded-lg bg-[#ecedef]" />
                    <div class="h-2 w-[100px] rounded-lg bg-[#ecedef]" />
                  </div>
                  <div class="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm">
                    <div class="h-4 w-4 rounded-full bg-[#ecedef]" />
                    <div class="h-2 w-[100px] rounded-lg bg-[#ecedef]" />
                  </div>
                  <div class="flex items-center space-x-2 rounded-md bg-white p-2 shadow-sm">
                    <div class="h-4 w-4 rounded-full bg-[#ecedef]" />
                    <div class="h-2 w-[100px] rounded-lg bg-[#ecedef]" />
                  </div>
                </div>
              </div>
              <span class="block w-full p-2 text-center font-normal">
                Light
              </span>
            </FormLabel>
          </FormItem>
          <FormItem>
            <FormLabel class="[&:has([data-state=checked])>div]:border-primary">
              <FormControl>
                <RadioGroupItem value="dark" class="sr-only" />
              </FormControl>
              <div class="items-center rounded-md border-2 border-muted bg-popover p-1 hover:bg-accent hover:text-accent-foreground">
                <div class="space-y-2 rounded-sm bg-slate-950 p-2">
                  <div class="space-y-2 rounded-md bg-slate-800 p-2 shadow-sm">
                    <div class="h-2 w-20 rounded-lg bg-slate-400" />
                    <div class="h-2 w-[100px] rounded-lg bg-slate-400" />
                  </div>
                  <div class="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm">
                    <div class="h-4 w-4 rounded-full bg-slate-400" />
                    <div class="h-2 w-[100px] rounded-lg bg-slate-400" />
                  </div>
                  <div class="flex items-center space-x-2 rounded-md bg-slate-800 p-2 shadow-sm">
                    <div class="h-4 w-4 rounded-full bg-slate-400" />
                    <div class="h-2 w-[100px] rounded-lg bg-slate-400" />
                  </div>
                </div>
              </div>
              <span class="block w-full p-2 text-center font-normal">
                Dark
              </span>
            </FormLabel>
          </FormItem>
        </RadioGroup>
      </FormItem>
    </FormField>

    <div class="flex justify-start">
      <Button type="submit">
        Update preferences
      </Button>
    </div>
  </form>
</template>
