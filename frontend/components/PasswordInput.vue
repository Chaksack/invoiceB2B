<script setup lang="ts">
import { ref, computed } from 'vue'
import type { HTMLAttributes } from 'vue'
import { Eye, EyeOff, Check, Circle } from 'lucide-vue-next'
import { cn } from '~/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Popover, PopoverContent, PopoverAnchor } from '@/components/ui/popover'

// Use defineModel for the simplest, most modern v-model implementation.
const modelValue = defineModel<string>()

const props = defineProps<{
  class?: HTMLAttributes['class']
  disabled?: boolean
  autocomplete?: string
  placeholder?: string
}>()

const showPassword = ref(false)
const inputFocused = ref(false)

// --- Password Checker Logic ---
const passwordRequirements = [
  { id: 'length', text: 'At least 8 characters', regex: /.{8,}/ },
  { id: 'uppercase', text: 'At least one uppercase letter', regex: /[A-Z]/ },
  { id: 'lowercase', text: 'At least one lowercase letter', regex: /[a-z]/ },
  { id: 'number', text: 'At least one number', regex: /[0-9]/ },
  { id: 'special', text: 'At least one special character', regex: /[^A-Za-z0-9]/ },
]

// This computed property checks the password against the requirements in real-time.
const validationState = computed(() => {
  const password = modelValue.value || ''
  return passwordRequirements.map(req => ({
    ...req,
    met: req.regex.test(password),
  }))
})
</script>

<template>
  <Popover :open="inputFocused && modelValue && modelValue.length > 0">
    <PopoverAnchor as-child>
      <div class="relative">
        <!-- The v-model correctly binds to the modelValue defined above -->
        <Input
            v-model="modelValue"
            :type="showPassword ? 'text' : 'password'"
            :class="cn('pr-10', props.class)"
            :placeholder="props.placeholder || 'Enter your password'"
            :disabled="props.disabled"
            :autocomplete="props.autocomplete"
            @focus="inputFocused = true"
            @blur="inputFocused = false"
        />
        <Button
            type="button"
            variant="ghost"
            size="icon"
            class="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
            :disabled="props.disabled"
            @click="showPassword = !showPassword"
            @mousedown.prevent
        >
          <!-- Use the imported Eye and EyeOff icons for clarity -->
          <EyeOff v-if="showPassword" class="size-4" aria-hidden="true" />
          <Eye v-else class="size-4" aria-hidden="true" />
          <span class="sr-only">
            {{ showPassword ? "Show password" : "Hide password" }}
          </span>
        </Button>
      </div>
    </PopoverAnchor>
    <PopoverContent side="right" :side-offset="8" class="w-72">
      <div class="space-y-2">
        <p class="text-sm font-medium text-foreground">Password must contain:</p>
        <ul class="space-y-2">
          <!-- The list items dynamically update based on the validationState -->
          <li v-for="req in validationState" :key="req.id" :class="['flex items-center text-sm transition-colors', req.met ? 'text-emerald-500' : 'text-muted-foreground']">
            <component :is="req.met ? Check : Circle" class="mr-2 h-3.5 w-3.5" />
            <span>{{ req.text }}</span>
          </li>
        </ul>
      </div>
    </PopoverContent>
  </Popover>
</template>
