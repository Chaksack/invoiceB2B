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
  Settings2,
  SquareTerminal,
} from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute } from '#app'

// Importing all the necessary child components for the sidebar
import NavMain from '@/components/NavMain.vue'
import NavSecondary from '@/components/NavSecondary.vue'
import NavUser from '@/components/NavUser.vue'
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
  user: {
    name: 'Andrew Chakdahah',
    email: 'm@example.com',
    avatar: '/avatars/shadcn.jpg',
  },
  navMain: [
    {
      title: 'Dashboard',
      url: '/admin',
      icon: SquareTerminal,
    },
    {
      title: 'Invoices',
      url: '/admin/invoices',
      icon: Bot,
    },
    {
      title: 'Users',
      url: '/admin/users',
      icon: BookOpen,
    },
    {
      title: 'Staffs',
      url: '/admin/staff', // Corrected URL from '/admin.staff'
      icon: BookOpen,
    },
  ],
  navSecondary: [
    {
      title: 'Settings',
      url: '/settings/profile',
      icon: Settings2,
    },
    {
      title: 'Help Center',
      url: '/help',
      icon: LifeBuoy,
    },
    {
      title: 'Notifications',
      url: '/notifications',
      icon: Bell,
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
  isActive: route.path === item.url,
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
              <div class="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary">
                <Command class="size-4" />
              </div>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">Profundr Inc</span>
                <span class="truncate text-xs">Demo</span>
              </div>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="navMainItems" />
      <NavSecondary :items="navSecondaryItems" class="mt-auto" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="data.user" />
    </SidebarFooter>
  </Sidebar>
</template>
