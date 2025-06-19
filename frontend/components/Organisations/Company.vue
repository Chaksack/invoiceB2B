<script setup lang="ts">
import {
  AlertCircle,
  CircleUser,
  Copy,
  CreditCard,
  File,
  EllipsisVertical,
  Plus,
  Slash,
} from 'lucide-vue-next';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarImage, AvatarFallback } from '~/components/ui/avatar';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

import {Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle, DialogTrigger} from '~/components/ui/dialog';
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '~/components/ui/breadcrumb';
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover';
import {ref, onMounted} from "vue";
const isLoading = ref(true);
const users = ref([
  {
    id: '1',
    name: 'Tech Innovators Ltd.',
    email: 'info@techinnovators.com',
    phone: '+1 555-123-4567',
    status: 'active',
    create: '2024-10-26',
    industry: 'Software Development',
    location: 'Silicon Valley, CA'
  },
  {
    id: '2',
    name: 'Global Marketing Solutions',
    email: 'contact@globalmarketing.net',
    phone: '+44 20 7987 6543',
    status: 'pending',
    create: '2024-10-20',
    industry: 'Marketing and Advertising',
    location: 'London, UK'
  },
  {
    id: '3',
    name: 'Green Energy Corp.',
    email: 'inquiries@greenenergycorp.org',
    phone: '+61 2 9876 5432',
    status: 'active',
    create: '2024-10-15',
    industry: 'Renewable Energy',
    location: 'Sydney, Australia'
  },
  {
    id: '4',
    name: 'Food Distribution Network',
    email: 'sales@fooddistrib.com',
    phone: '+233 50 123 4567',
    status: 'inactive',
    create: '2024-10-10',
    industry: 'Food and Beverage',
    location: 'Accra, Ghana'
  },
]);

onMounted(() => {
  setTimeout(() => {
    isLoading.value = false;
  }, 2000);
});
</script>

<template>
  <div class="z-50">

    <main class="flex flex-col items-start gap-4 p-4 sm:px-6 sm:py-2 md:gap-8">
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink href="/dashboard">Home</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator><Slash /></BreadcrumbSeparator>
          <BreadcrumbItem>
            <BreadcrumbLink href="/organisations">Organisations</BreadcrumbLink>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      <div class="w-full overflow-x-auto">
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-4">
          <input
              type="text"
              placeholder="Search for company........"
              class="w-full sm:w-96 p-2 bg-transparent border rounded-full mb-2 sm:mb-0 sm:mr-2"
          />
          <Dialog>
            <DialogTrigger as-child>
              <Button
                  variant="outline"
                  size="sm"
                  class="h-7 gap-2 bg-orgC font-semibold text-md text-white hover:bg-blC hover:text-white rounded-lg px-3"
              >
                <Plus class="h-3.5 w-3.5" />
                New Organisation
              </Button>            </DialogTrigger>
            <DialogContent
                class="fixed ml-auto right-0 h-[100vh] w-1/2 flex flex-col transform transition-transform duration-300 ease-out translate-x-full md:translate-x-0 shadow-lg overflow-y-auto p-6"
            >
              <DialogHeader class="pb-4">
                <DialogTitle class="">New Loan Application</DialogTitle>
                <DialogDescription>
                  Add your credit or loan details here.
                </DialogDescription>
              </DialogHeader>
              <div class="py-4 flex-1">
                <form class="" @submit.prevent="createLoan">
                  <div class="grid gap-2">
                    <Label for="text" class="mb-2 mt-2">Loan Type</Label>
                    <Input id="text" type="text" v-model="loantype" placeholder="Personal loan" required />
                  </div>
                  <div class="grid gap-2 mt-2">
                    <Label for="text" class="mb-2 mt-2">Loan Purpose</Label>
                    <Input id="text" type="text" v-model="loanpurpose" placeholder="Home renovation" required />
                  </div>
                  <div class="grid grid-cols-2 gap-4">
                    <div class="grid gap-2 mt-2">
                      <Label for="number" class="mb-2 mt-2">Loan Amount</Label>
                      <Input id="number" type="number" v-model="loanamount" placeholder="1000" required />
                    </div>
                    <div class="grid gap-2 mt-2">
                      <Label for="number" class="mb-2 mt-2">Loan Term</Label>
                      <Input id="number" type="number" v-model="loanterm" placeholder="12" required />
                    </div>
                  </div>
                  <div class="grid grid-cols-2 gap-4">
                    <div class="grid gap-2 mt-2">
                      <Label for="date" class="mb-2 mt-2">Loan Start Date</Label>
                      <Input id="date" type="date" v-model="startdate" placeholder="1000" required />
                    </div>
                    <div class="grid gap-2 mt-2">
                      <Label for="date" class="mb-2 mt-2">Loan End Date</Label>
                      <Input id="date" type="date" v-model="enddate" placeholder="12" required />
                    </div>
                  </div>
                  <Alert v-if="errorMessage" variant="default" class="bg-red-500 text-white">
                    <AlertCircle class="w-4 h-4" />
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription>
                      {{ errorMessage }}
                    </AlertDescription>
                  </Alert>

                  <div class="pt-8 flex-1">
                    <Button type="submit" :disabled="isLoading">Save changes</Button>
                  </div>
                </form>

              </div>
            </DialogContent>
          </Dialog>

        </div>
        <div>
          <div v-if="isLoading" class="flex flex-col justify-center items-center h-64">
            <div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-orgC"></div>
            <p class="mt-2">Loading Data...</p>
          </div>
          <div v-else>
        <Table class="w-full min-w-[768px]">
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Avatar</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Phone</TableHead>
              <TableHead class="hidden md:table-cell">Industry</TableHead>
              <TableHead class="hidden md:table-cell">Location</TableHead>
              <TableHead class="hidden md:table-cell">Created</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>

          <TableBody>
            <TableRow class="" v-for="user in users" :key="user.id">
              <TableCell>{{ user.id }}</TableCell>
              <TableCell>
                <Avatar class="relative bg-black overflow-visible">
                  <AvatarImage class="rounded-full" src="" alt="User Avatar" />
                  <AvatarFallback class="text-white">
                    {{ user.name.substring(0, 2).toUpperCase() }}
                  </AvatarFallback>
                </Avatar>
              </TableCell>
              <TableCell>{{ user.name }}</TableCell>
              <TableCell>{{ user.email }}</TableCell>
              <TableCell>{{ user.phone }}</TableCell>
              <TableCell class="hidden md:table-cell">{{ user.industry }}</TableCell>
              <TableCell class="hidden md:table-cell">{{ user.location }}</TableCell>
              <TableCell class="hidden md:table-cell">{{ user.create }}</TableCell>
              <TableCell>
                <span
                    :class="
                    user.status === 'active' ? 'bg-green-500' : 'bg-red-500'
                  "
                    class="inline-block ml-2 px-2 py-1 rounded-full font-semibold text-white text-xs"
                >
                  {{ user.status === 'active' ? 'Active' : 'Inactive' }}
                </span>
              </TableCell>
              <TableCell>
                    <NuxtLink to="/organisation/view">
                      <button class="bg-bgC  px-2 py-1 text-white rounded-lg">View</button>
                    </NuxtLink>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
