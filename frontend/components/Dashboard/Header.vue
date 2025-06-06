<template>
  <header
      class=" transition-all duration-300 ease-in-out"
      :class="[
      isScrolled ? ' backdrop-blur-lg' : '',
      'sticky top-0 z-50',
    ]"
  >    <nav class="mx-auto flex max-w-7xl items-center justify-between p-6 lg:px-8" aria-label="Global">
    <div class="flex lg:flex-1">
      <NuxtLink to="/" class="-m-1.5 p-1.5">
        <h2 class="text-2xl font-bold text-primary">ProFundr</h2>
      </NuxtLink>
    </div>
    <div class="flex lg:hidden">
      <button
          type="button"
          class="-m-2.5 inline-flex items-center justify-center rounded-md p-2.5 text-gray-700"
          @click="mobileMenuOpen = true"
      >
        <span class="sr-only">Open main menu</span>
        <Bars3Icon class="size-6" aria-hidden="true" />
      </button>
    </div>

    <NavigationMenu>
      <NavigationMenuList>

        <NavigationMenuItem class="px-2 text-md font-semibold" >Home
        </NavigationMenuItem>

        <NavigationMenuItem class="px-2 text-md font-semibold" >Features
        </NavigationMenuItem>
        <NavigationMenuItem class="px-2 text-md font-semibold" >FAQs
        </NavigationMenuItem>
        <NavigationMenuItem class="px-2 text-md font-semibold" >Support
        </NavigationMenuItem>



      </NavigationMenuList>
    </NavigationMenu>

    <div class="hidden lg:flex lg:flex-1 lg:justify-end">
      <div class="flex items-center">
        <Avatar class="relative overflow-visible">
          <AvatarFallback class="text-black"> AC </AvatarFallback>
        </Avatar>
      </div>
    </div>
  </nav>

    <Dialog class="lg:hidden" @close="mobileMenuOpen = false" :open="mobileMenuOpen">
      <div class="fixed inset-0 z-10" />
      <DialogPanel
          class="fixed inset-y-0 right-0 z-10 w-full overflow-y-auto bg-black px-6 py-6 sm:max-w-sm sm:ring-1 sm:ring-gray-900/10"
      >
        <div class="flex items-center justify-between">
          <NuxtLink to="/" class="-m-1.5 p-1.5">
            <span class="sr-only">Your Company</span>
            <img
                class="h-8 w-auto"
                src="https://tailwindui.com/plus/img/logos/mark.svg?color=indigo&shade=600"
                alt="Your Company Logo"
            />
          </NuxtLink>
          <button
              type="button"
              class="-m-2.5 rounded-md p-2.5 text-gray-200"
              @click="mobileMenuOpen = false"
          >
            <span class="sr-only">Close menu</span>
            <XMarkIcon class="size-6" aria-hidden="true" />
          </button>
        </div>
        <div class="mt-6 flow-root">
          <div class="-my-6 divide-y divide-gray-500/10">
            <div class="space-y-2 py-6">
              <Disclosure as="div" class="-mx-3" v-slot="{ open }">
                <DisclosureButton
                    class="flex w-full items-center justify-between rounded-lg py-2 pl-3 pr-3.5 text-base/7 font-semibold hover:bg-gray-50"
                >
                  What we do
                  <ChevronDownIcon :class="[open ? 'rotate-180' : '', 'size-5 flex-none']" aria-hidden="true" />
                </DisclosureButton>
                <DisclosurePanel class="mt-2 space-y-2">
                  <DisclosureButton
                      v-for="item in [...products, ...callsToAction]"
                      :key="item.name"
                      as="NuxtLink"
                      :to="item.href"
                      class="block rounded-lg py-2 pl-6 pr-3 text-sm/7 font-semibold  hover:bg-gray-50"
                  >
                    {{ item.name }}
                  </DisclosureButton>
                </DisclosurePanel>
              </Disclosure>
              <NuxtLink
                  to="/about"
                  class="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold  hover:bg-gray-50"
              >
                Who we are
              </NuxtLink>
              <NuxtLink
                  to="/careers"
                  class="-mx-3 block rounded-lg px-3 py-2 text-base/7 font-semibold  hover:bg-gray-50"
              >
                Careers
              </NuxtLink>
            </div>
            <div class="py-6">
              <NuxtLink
                  to="/contact" class="text-sm/6 font-semibold bg-gradient-to-r from-green-400 via-green-600 to-syentiaGreen px-3 py-1 rounded-lg text-white">
                Contact Us <span aria-hidden="true">&rarr;</span>
              </NuxtLink>
            </div>
          </div>
        </div>
      </DialogPanel>
    </Dialog>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import {
  Dialog,
  DialogPanel,
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
  Popover,
  PopoverButton,
  PopoverGroup,
  PopoverPanel,
} from '@headlessui/vue'
import {
  ArrowPathIcon,
  Bars3Icon,
  ChartPieIcon,
  CursorArrowRaysIcon,
  FingerPrintIcon,
  SquaresPlusIcon,
  XMarkIcon,
} from '@heroicons/vue/24/outline'
import { ChevronDownIcon, PhoneIcon, PlayCircleIcon } from '@heroicons/vue/20/solid'
import {NavigationMenu, NavigationMenuItem, NavigationMenuList} from "~/components/ui/navigation-menu/index.js";

const mobileMenuOpen = ref(false)

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




