<script setup lang="ts">
import { useRoute } from '#app'
import { NuxtLink } from '#components'
import { ChevronRight, type LucideIcon } from 'lucide-vue-next'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  SidebarGroup,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuAction,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from '@/components/ui/sidebar'
import { cn } from '@/lib/utils'

// Component props now accept nested items for sub-menus
defineProps<{
  items: {
    title: string
    url: string
    icon: LucideIcon
    isActive?: boolean // This should be calculated in the parent component
    items?: {
      title: string
      url: string
    }[]
  }[]
}>()

// Use Nuxt's composable to get the current route
const route = useRoute()
</script>

<template>
  <SidebarGroup>
    <SidebarGroupLabel>Platform</SidebarGroupLabel>
    <SidebarMenu>
      <!-- The collapsible opens if the parent or any child route is active -->
      <Collapsible v-for="item in items" :key="item.title" as-child :default-open="item.isActive">
        <SidebarMenuItem>
          <!-- The button variant changes based on the active state -->
          <SidebarMenuButton
              :variant="item.isActive ? 'primary' : 'ghost'"
              class="w-full justify-start"
              as-child
              :tooltip="item.title"
              :class="cn(
        'w-full text-left justify-start items-start',
        route.path === item.href && 'bg-primary hover:bg-primary',
      )"
          >
            <!-- Use NuxtLink for proper client-side routing -->
            <NuxtLink :to="item.url">
              <component :is="item.icon" class="mr-2 size-4" />
              <span>{{ item.title }}</span>
            </NuxtLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </Collapsible>
    </SidebarMenu>
  </SidebarGroup>
</template>
