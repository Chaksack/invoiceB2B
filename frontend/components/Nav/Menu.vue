<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import {
  Dialog,
  DialogPanel,
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
  Menu,
  MenuButton,
  MenuItem,
  MenuItems,
} from '@headlessui/vue'
import {
  Bars3Icon,
  XMarkIcon,
} from '@heroicons/vue/24/outline'
import { ChevronDownIcon } from '@heroicons/vue/20/solid'
import { NavigationMenu, NavigationMenuItem, NavigationMenuList } from "~/components/ui/navigation-menu/index.js";

const mobileMenuOpen = ref(false)
const isScrolled = ref(false)

// --- Mock Authentication Logic ---
// In a real Nuxt app, you would likely use a composable from an auth module
// (e.g., nuxt-auth, supabase) instead of useState.
const user = useState('user', () => null) // Default state: user is not logged in

// Mock user data for demonstration purposes
const loggedInUser = {
  name: 'Andrew C.',
  email: 'm@example.com',
  avatar: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
}

// Mock login/logout functions
function login() {
  user.value = loggedInUser;
}

function logout() {
  user.value = null;
  mobileMenuOpen.value = false
  // Optional: Redirect to home or login page
  navigateTo('/')
}
// --- End Mock Auth Logic ---


const products = [
  { name: 'Cloud',description: 'Cloud Migrations, Management and Infrastrcture', href: '/cloud' },
  { name: 'Artificial Intelligence', description: 'Cloud Migrations, Management and Infrastrcture', href: '/ai' },
  { name: 'Cybersecurity',description: 'Cloud Migrations, Management and Infrastrcture', href: '/cybersecurity' },
  { name: 'Application Development',description: 'Cloud Migrations, Management and Infrastrcture', href: '/application' },
  { name: 'Digital Marketing',description: 'Cloud Migrations, Management and Infrastrcture', href: '/marketing' },
  { name: 'Project Management',description: 'Cloud Migrations, Management and Infrastrcture', href: '/projectmanagement' },
]

const handleScroll = () => {
  isScrolled.value = window.scrollY > 0;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<template>
  <header
      class=" transition-all duration-300 ease-in-out"
      :class="[
      isScrolled ? 'bg-white dark:bg-black backdrop-blur-lg' : 'bg-white dark:bg-black',
      'sticky top-0 z-50',
    ]"
  >
    <nav class="mx-auto flex max-w-7xl items-center justify-between p-6 lg:px-8" aria-label="Global">
      <div class="flex lg:flex-1">
        <NuxtLink to="/" class="-m-1.5 p-1.5">
          <h2 class="text-2xl font-bold text-primary">ProFundr</h2>
        </NuxtLink>
      </div>
      <div class="flex lg:hidden">
        <button
            type="button"
            class="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700 dark:text-gray-200"
            @click="mobileMenuOpen = true"
        >
          <span class="sr-only">Open main menu</span>
          <Bars3Icon class="size-6" aria-hidden="true" />
        </button>
      </div>

      <NavigationMenu class="hidden lg:block">
        <NavigationMenuList>
          <NavigationMenuItem class="px-2 text-md font-semibold"  >Home</NavigationMenuItem>
          <NavigationMenuItem class="px-2 text-md font-semibold" >Features</NavigationMenuItem>
          <NavigationMenuItem class="px-2 text-md font-semibold" >FAQs</NavigationMenuItem>
          <NavigationMenuItem class="px-2 text-md font-semibold" >Support</NavigationMenuItem>
        </NavigationMenuList>
      </NavigationMenu>

      <!-- Conditional Desktop Menu -->
      <div class="hidden lg:flex lg:flex-1 lg:justify-end items-center">
        <!-- Logged In View -->
        <Menu v-if="user" as="div" class="relative ml-3">
          <div>
            <MenuButton class="relative flex rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
              <span class="absolute -inset-1.5" />
              <span class="sr-only">Open user menu</span>
              <img class="h-8 w-8 rounded-full" :src="user.avatar" alt="User avatar" />
            </MenuButton>
          </div>
          <transition enter-active-class="transition ease-out duration-100" enter-from-class="transform opacity-0 scale-95" enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75" leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
            <MenuItems class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none dark:bg-gray-800 ">
              <div class="px-4 py-3">
                <p class="text-sm text-gray-900 dark:text-white">{{ user.name }}</p>
                <p class="truncate text-sm text-gray-500 dark:text-gray-400">{{ user.email }}</p>
              </div>
              <MenuItem v-slot="{ active }">
                <a href="/settings/profile" :class="[active ? 'bg-primary' : '', 'block px-4 py-2 text-sm']">Profile</a>
              </MenuItem>
              <MenuItem v-slot="{ active }">
                <a href="/settings/account" :class="[active ? 'bg-primary' : '', 'block px-4 py-2 text-sm']">Settings</a>
              </MenuItem>
              <MenuItem v-slot="{ active }">
                <a href="/" @click="logout" :class="[active ? 'bg-primary' : '', 'block px-4 py-2 text-sm']">Sign out</a>
              </MenuItem>
            </MenuItems>
          </transition>
        </Menu>
        <!-- Logged Out View -->
        <div v-else>
          <button @click="login" class="text-sm/6 mr-2 font-semibold border border-primary text-primary dark:text-white px-4 py-1 rounded-full ">
            Log in
          </button>
          <NuxtLink to="/register" class="text-sm/6 font-semibold shadow-lg bg-primary px-4 py-1 rounded-full text-white">
            Create an account
          </NuxtLink>
        </div>
      </div>
    </nav>

    <!-- Conditional Mobile Menu -->
    <Dialog class="lg:hidden" @close="mobileMenuOpen = false" :open="mobileMenuOpen">
      <div class="fixed inset-0 z-50" />
      <DialogPanel class="fixed inset-y-0 right-0 z-50 w-full overflow-y-auto bg-black px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-white/10">
        <div class="flex items-center justify-between">
          <NuxtLink to="/" class="-m-1.5 p-1.5">
            <h2 class="text-2xl font-bold text-primary">ProFundr</h2>
          </NuxtLink>
          <button type="button" class="-m-2.5 rounded-md p-2.5 text-gray-200" @click="mobileMenuOpen = false">
            <span class="sr-only">Close menu</span>
            <XMarkIcon class="size-6" aria-hidden="true" />
          </button>
        </div>
        <div class="mt-6 flow-root">
          <div class="-my-6 divide-y divide-gray-500/20">
            <div class="space-y-2 py-6">
              <NuxtLink v-for="item in products" :key="item.name" :to="item.href" class="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-white hover:bg-gray-800">{{ item.name }}</NuxtLink>
            </div>
            <!-- Mobile Logged In/Out View -->
            <div class="py-6">
              <div v-if="user" class="space-y-2">
                <div class="flex items-center px-3">
                  <div class="flex-shrink-0">
                    <img class="h-10 w-10 rounded-full" :src="user.avatar" alt="">
                  </div>
                  <div class="ml-3">
                    <div class="text-base font-medium leading-none text-white">{{ user.name }}</div>
                    <div class="text-sm font-medium leading-none text-gray-400">{{ user.email }}</div>
                  </div>
                </div>
                <div class="mt-3 space-y-1">
                  <a href="#" class="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-white hover:bg-gray-800">Profile</a>
                  <a href="#" class="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-white hover:bg-gray-800">Settings</a>
                  <a href="#" @click="logout" class="-mx-3 block rounded-lg px-3 py-2 text-base font-semibold leading-7 text-white hover:bg-gray-800">Sign out</a>
                </div>
              </div>
              <div v-else>
                <button @click="login(); mobileMenuOpen = false" class="-mx-3 block rounded-lg px-3 py-2.5 text-base font-semibold leading-7 text-white hover:bg-gray-800">
                  Log in
                </button>
              </div>
            </div>
          </div>
        </div>
      </DialogPanel>
    </Dialog>
  </header>
</template>
