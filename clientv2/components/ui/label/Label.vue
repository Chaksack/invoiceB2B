<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { Label as RekaLabel, type LabelProps } from 'reka-ui'
import { cn } from '@/lib/utils'

const props = defineProps<LabelProps & { class?: HTMLAttributes['class'] }>()
</script>

<template>
  <RekaLabel
    data-slot="label"
    :class="cn('text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70', props.class)"
  >
    <slot />
  </RekaLabel>
</template>