<script setup lang="ts">
import {
  BookOpen,
  Bot,
  Bell,
  Command,
  LogOut,
  Frame,
  LifeBuoy,
  Map,
  PieChart,
  Send,
  History,
  Settings2,
  SquareTerminal, MapPlus, MessageCircle,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute } from '#app'

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  type SidebarProps,
} from '@/components/ui/sidebar'

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
})

// Get the current route using Nuxt's composable to determine the active page
const route = useRoute()

// Static data for the sidebar navigation links and user info
const data = {
  navMain: [
    {
      title: 'Dashboard',
      url: '/home',
      icon: MapPlus,
    },
    {
      title: 'Invoices',
      url: '/invoices',
      icon: MessageCircle,
    },
    {
      title: 'Settings',
      url: '/settings/profile',
      icon: Settings2,
    },
  ],
  navSecondary: [
    {
      title: 'Help Center',
      url: '/help',
      icon: LifeBuoy,
    },
    {
      title: 'Log Out',
      url: '/logout',
      icon: LogOut,
    },
  ],
}

// Create computed properties to dynamically add `isActive: true` to the
// nav item that matches the current route. This is the cleanest way to handle active state.
const navMainItems = computed(() => data.navMain.map(item => ({
  ...item,
  isActive: route.path === item.url && 'bg-stone-300 hover:bg-stone-400',

})))

const navSecondaryItems = computed(() => data.navSecondary.map(item => ({
  ...item,
  isActive: route.path === item.url,
})))
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <a href="#">
              <div class="flex aspect-square size-8 items-center justify-center rounded-lg">
                <Command class="size-4" />
              </div>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">Profundr</span>
              </div>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <UserNavMain :items="navMainItems" />
      <UserNavSecondary :items="navSecondaryItems" class="mt-auto" />
    </SidebarContent>
  </Sidebar>
</template>