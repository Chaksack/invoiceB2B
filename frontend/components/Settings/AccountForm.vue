<script setup lang="ts">
import { h, ref } from 'vue'
import { toTypedSchema } from '@vee-validate/zod'
import * as z from 'zod'
import { CalendarDate, DateFormatter, getLocalTimeZone, today } from '@internationalized/date'
import { toDate } from 'radix-vue/date'
import { Check, ChevronsUpDown } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command'
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Separator } from '@/components/ui/separator'

// State and constants
const open = ref(false)
const dateValue = ref()
const placeholder = ref()

const languages = [
  { label: 'English', value: 'en' },
  { label: 'French', value: 'fr' },
  { label: 'German', value: 'de' },
  { label: 'Indonesia', value: 'id' },
  { label: 'Spanish', value: 'es' },
  { label: 'Portuguese', value: 'pt' },
  { label: 'Russian', value: 'ru' },
  { label: 'Japanese', value: 'ja' },
  { label: 'Korean', value: 'ko' },
  { label: 'Chinese', value: 'zh' },
] as const

const df = new DateFormatter('en-US', {
  dateStyle: 'long',
})

// Form schema definition
const accountFormSchema = toTypedSchema(z.object({
  name: z
      .string({
        required_error: 'Required.',
      })
      .min(2, {
        message: 'Name must be at least 2 characters.',
      })
      .max(30, {
        message: 'Name must not be longer than 30 characters.',
      }),
  dob: z.string({ required_error: 'Please select a valid date.' }).datetime(),
  language: z.string({ required_error: 'Please select a language.' }).min(1, 'Please select a language.'),
}))

// Submission handler
async function onSubmit(values: any) {
  toast({
    title: 'You submitted the following values:',
    description: h('pre', { class: 'mt-2 w-[340px] rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white' }, JSON.stringify(values, null, 2))),
  })
}
</script>

<template>
  <div>
    <h3 class="text-lg font-medium">
      Account
    </h3>
    <p class="text-sm text-muted-foreground">
      Update your account settings. Set your preferred language and timezone.
    </p>
  </div>
  <Separator />
  <Form v-slot="{ setFieldValue }" :validation-schema="accountFormSchema" class="space-y-8" @submit="onSubmit">
    <FormField v-slot="{ componentField }" name="name">
      <FormItem>
        <FormLabel>Name</FormLabel>
        <FormControl>
          <Input type="text" placeholder="Your name" v-bind="componentField" />
        </FormControl>
        <FormDescription>
          This is the name that will be displayed on your profile and in emails.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ value }" name="dob">
      <FormItem class="flex flex-col">
        <FormLabel>Date of birth</FormLabel>
        <Popover>
          <PopoverTrigger as-child>
            <FormControl>
              <Button
                  variant="outline" :class="cn(
                  'w-[240px] justify-start text-left font-normal',
                  !value && 'text-muted-foreground',
                )"
              >
                <!-- Assuming you have an Icon component or are using an icon library -->
                <!-- <Icon name="radix-icons:calendar" class="mr-2 h-4 w-4" /> -->
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 15 15" class="mr-2 h-4 w-4 opacity-50"><path fill="currentColor" fill-rule="evenodd" d="M1 2.75A.75.75 0 0 1 1.75 2h11.5a.75.75 0 0 1 .75.75v9.5a.75.75 0 0 1-.75.75H1.75a.75.75 0 0 1-.75-.75V2.75Zm.75 1.5v1.25h11.5V4.25H1.75Zm0 2.5v5.75h11.5V6.75H1.75ZM4 8.5a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-1a.5.5 0 0 1-.5-.5v-1Zm3 0a.5.5 0 0 1 .5-.5h1a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5h-1a.5.5 0 0 1-.5-.5v-1Zm3.5-.5a.5.5 0 0 0-.5.5v1a.5.5 0 0 0 .5.5h1a.5.5 0 0 0 .5-.5v-1a.5.5 0 0 0-.5-.5h-1Z" clip-rule="evenodd"></path></svg>
                <span>{{ value ? df.format(toDate(dateValue)) : "Pick a date" }}</span>
              </Button>
            </FormControl>
          </PopoverTrigger>
          <PopoverContent class="w-auto p-0">
            <Calendar
                v-model:placeholder="placeholder"
                v-model="dateValue"
                initial-focus
                :min-value="new CalendarDate(1950, 1, 1)"
                :max-value="today(getLocalTimeZone())"
                @update:model-value="(v) => {
                if (v) {
                  dateValue = v
                  setFieldValue('dob', v.toDate(getLocalTimeZone()).toISOString())
                }
                else {
                  dateValue = undefined
                  setFieldValue('dob', undefined)
                }
              }"
            />
          </PopoverContent>
        </Popover>
        <FormDescription>
          Your date of birth is used to calculate your age.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <FormField v-slot="{ field }" name="language">
      <FormItem class="flex flex-col">
        <FormLabel>Language</FormLabel>
        <Popover v-model:open="open">
          <PopoverTrigger as-child>
            <FormControl>
              <Button
                  variant="outline"
                  role="combobox"
                  :aria-expanded="open"
                  :class="cn(
                  'w-[200px] justify-between',
                  !field.value && 'text-muted-foreground',
                )"
              >
                {{ field.value ? languages.find(
                  (language) => language.value === field.value,
              )?.label : 'Select language...' }}

                <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
              </Button>
            </FormControl>
          </PopoverTrigger>
          <PopoverContent class="w-[200px] p-0">
            <Command :model-value="field.value">
              <CommandInput placeholder="Search language..." />
              <CommandEmpty>No language found.</CommandEmpty>
              <CommandList>
                <CommandGroup>
                  <CommandItem
                      v-for="language in languages"
                      :key="language.value"
                      :value="language.value"
                      @select="(ev) => {
                      if (typeof ev.detail.value === 'string') {
                        setFieldValue('language', ev.detail.value)
                      }
                      open = false
                    }"
                  >
                    <Check
                        :class="cn(
                        'mr-2 h-4 w-4',
                        field.value === language.value ? 'opacity-100' : 'opacity-0',
                      )"
                    />
                    {{ language.label }}
                  </CommandItem>
                </CommandGroup>
              </CommandList>
            </Command>
          </PopoverContent>
        </Popover>
        <FormDescription>
          This is the language that will be used in the dashboard.
        </FormDescription>
        <FormMessage />
      </FormItem>
    </FormField>

    <div class="flex justify-start">
      <Button type="submit">
        Update account
      </Button>
    </div>
  </Form>
</template>
