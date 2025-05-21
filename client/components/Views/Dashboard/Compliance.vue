<script setup lang="ts">
import { Button } from '~/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select'
import { Stepper, StepperDescription, StepperItem, StepperSeparator, StepperTitle, StepperTrigger } from '~/components/ui/stepper'
import { toTypedSchema } from '@vee-validate/zod'
import { Check, Circle, Dot } from 'lucide-vue-next'
import { h, ref } from 'vue'
import * as z from 'zod'
import {toast } from 'vue-sonner';


const formSchema = [
  z.object({
    phoneNumber: z.string().regex(/^\d+$/, 'Invalid phone number'),
    dob: z.string(),
    occupation: z.string(),
    income: z.string(),
    address: z.string(),
    digitalAddress: z.string(),
  }),
  z.object({
    idType: z.union([z.literal('ghcard'), z.literal('passport')]),
    idNumber: z.string(),
  }),
  z.object({
    consent: z.preprocess(val => val === "true" || val === true, z.boolean().refine(val => val === true, {
      message: "You must agree to the terms.",
    }))
  })
]

const stepIndex = ref(1)
const steps = [
  {
    step: 1,
    title: 'Personal details',
    description: 'Provide your personal information',
  },
  {
    step: 2,
    title: 'Identification',
    description: 'Provide your identification',
  },
  {
    step: 3,
    title: 'Confirmation',
    description: 'Kindly give consent to your information',
  },
]

function onSubmit(values: any) {
  toast.success({
    title: 'You submitted the following values:',
    description: h('pre', { class: 'mt-2 w-[340px] rounded-md bg-slate-950 p-4' }, h('code', { class: 'text-white bg-customIndigo' }, JSON.stringify(values, null, 2))),
  })
}
</script>

<template>
  <div class="mt-10 px-4 lg:px-6 py-2.5" >
    <div class="flex flex-wrap justify-start mx-auto max-w-screen-xl">
      <Form
          v-slot="{ meta, values, validate }"
          as="" keep-values :validation-schema="toTypedSchema(formSchema[stepIndex - 1])"
      >
        <Stepper v-slot="{ isNextDisabled, isPrevDisabled, nextStep, prevStep }" v-model="stepIndex" class="block w-full">
          <form
              @submit="(e) => {
          e.preventDefault()
          validate()

          if (stepIndex === steps.length && meta.valid) {
            onSubmit(values)
          }
        }"
          >
            <div class="flex w-full flex-start gap-2">
              <StepperItem
                  v-for="step in steps"
                  :key="step.step"
                  v-slot="{ state }"
                  class="relative flex w-full flex-col items-center justify-center text-black"
                  :step="step.step"
              >
                <StepperSeparator
                    v-if="step.step !== steps[steps.length - 1].step"
                    class="absolute left-[calc(50%+20px)] right-[calc(-50%+10px)] top-5 block h-0.5 shrink-0 rounded-full bg-muted group-data-[state=completed]:bg-green-800"
                />

                <StepperTrigger as-child>
                  <Button
                      :variant="state === 'completed' || state === 'active' ? 'default' : 'outline'"
                      size="icon"
                      class="z-10 rounded-full shrink-0 bg-black text-white "
                      :class="[state === 'active' && 'ring-2 ring-ring ring-offset-2 ring-offset-background']"
                      :disabled="state !== 'completed' && !meta.valid"
                  >
                    <Check v-if="state === 'completed'" class="size-5" />
                    <Circle v-if="state === 'active'" />
                    <Dot v-if="state === 'inactive'" />
                  </Button>
                </StepperTrigger>

                <div class="mt-5 flex flex-col items-center text-center">
                  <StepperTitle
                      :class="[state === 'active' && 'text-black']"
                      class="text-sm font-semibold transition lg:text-base"
                  >
                    {{ step.title }}
                  </StepperTitle>
                  <StepperDescription
                      :class="[state === 'active' && 'text-black']"
                      class="sr-only text-xs text-muted-foreground transition md:not-sr-only lg:text-sm"
                  >
                    {{ step.description }}
                  </StepperDescription>
                </div>
              </StepperItem>
            </div>

            <div class="flex flex-col gap-4 mt-8">
              <template v-if="stepIndex === 1">
                <div class="space-y-8 md:grid md:grid-cols-2 lg:grid-cols-2 md:gap-12 md:space-y-0">
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="phoneNumber">
                      <FormItem>
                        <FormLabel>Phone Number</FormLabel>
                        <FormControl>
                          <Input type="text" v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="dob">
                      <FormItem>
                        <FormLabel>Date of Birth</FormLabel>
                        <FormControl>
                          <Input type="text" placeholder="yyyy/mm/dd" v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="occupation">
                      <FormItem>
                        <FormLabel>Occupation</FormLabel>
                        <FormControl>
                          <Input type="text " v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="income">
                      <FormItem>
                        <FormLabel>Monthly Income</FormLabel>
                        <FormControl>
                          <Input type="text" prefix="GHC" v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="address">
                      <FormItem>
                        <FormLabel>Residential Address</FormLabel>
                        <FormControl>
                          <Input type="text" v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="digitalAddress">
                      <FormItem>
                        <FormLabel>Digital Address</FormLabel>
                        <FormControl>
                          <Input type="text" v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                </div>
              </template>

              <template v-if="stepIndex === 2">
                <div class="space-y-8 md:grid md:grid-cols-2 lg:grid-cols-2 md:gap-12 md:space-y-0">
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="idType" required>
                      <FormItem>
                        <FormLabel>Identification Type</FormLabel>

                        <Select v-bind="componentField">
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Select id type" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            <SelectGroup>
                              <SelectItem value="ghcard">
                                National Id
                              </SelectItem>
                              <SelectItem value="passport">
                                Passport
                              </SelectItem>
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                  <div class="grid gap-2 ">
                    <FormField v-slot="{ componentField }" name="idNumber">
                      <FormItem>
                        <FormLabel>ID Number</FormLabel>
                        <FormControl>
                          <Input type="text" v-bind="componentField" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    </FormField>
                  </div>
                </div>

              </template>

              <template v-if="stepIndex === 3">
                <FormField v-slot="{ value, handleChange }" name="consent">
                  <FormItem class="flex items-center space-x-2">
                    <FormControl>
                      <input
                          type="checkbox"
                          :checked="value"
                          @change="handleChange($event.target.checked)"
                          class="size-5"
                      />
                    </FormControl>
                    <FormLabel> I consent to providing my information </FormLabel>
                    <FormMessage />
                  </FormItem>
                </FormField>


              </template>
            </div>

            <div class="flex items-center justify-between mt-4">
              <Button :disabled="isPrevDisabled" class="border-black" variant="outline" size="sm" @click="prevStep()">
                Back
              </Button>
              <div class="flex items-center gap-3">
                <Button class="bg-black text-white" v-if="stepIndex !== 3" :type="meta.valid ? 'button' : 'submit'" :disabled="isNextDisabled" size="sm" @click="meta.valid && nextStep()">
                  Next
                </Button>
                <NuxtLink to="/admin/home">
                  <Button
                      class="bg-green-800 text-white" v-if="stepIndex === 3" size="sm" type="submit"
                  >
                    Submit
                  </Button>
                </NuxtLink>

              </div>
            </div>
          </form>
        </Stepper>
      </Form>
    </div>
  </div>
</template>