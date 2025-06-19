<template>
  <main class="flex flex-col z-50 items-start gap-4 p-4 sm:px-6 sm:py-2 md:gap-8">
  <div class="p-6  w-full">
    <Breadcrumb>
      <BreadcrumbList>
        <BreadcrumbItem>
          <BreadcrumbLink href="/dashboard">
            Home
          </BreadcrumbLink>
        </BreadcrumbItem>
        <BreadcrumbSeparator>
          <Slash />
        </BreadcrumbSeparator>
        <BreadcrumbItem>
          <BreadcrumbLink href="/organisations">
            Organisation
          </BreadcrumbLink>
        </BreadcrumbItem>
        <BreadcrumbSeparator>
          <Slash />
        </BreadcrumbSeparator>
        <BreadcrumbItem>
          <BreadcrumbLink href="/users">
            Company Details
          </BreadcrumbLink>
        </BreadcrumbItem>
      </BreadcrumbList>
    </Breadcrumb>
    <!-- Header -->
    <div class="flex justify-between mt-4 items-center border-b pb-4">
        <div class="flex ">
          <Avatar class="h-20 w-20 bg-grC text-white">
            <AvatarFallback>CW</AvatarFallback>
          </Avatar>
          <div class="flex flex-col ml-2">
          <h1 class="text-4xl font-semibold  text-bgC">Cameron Williamson <span class="text-orgC font-bold">GHC17,000</span></h1>
        <div class="flex mt-2">
          <p class="text-gray-600 ">Created March 15, 2025, 2:31 PM</p>
          <Badge class="bg-green-200 uppercase text-green-700 ml-2 ">
            KYC Compliance: Completed
          </Badge>
          <Badge class="bg-yellow-200 uppercase text-yellow-700 mx-2">
            Low Risk
          </Badge>
        </div>
        </div>


      </div>
      <div>
        <Button class="bg-blC mx-2">Update Limit <ChevronDown /></Button>
        <Popover>
          <PopoverTrigger>
            <Button class="bg-orgC mx-2">Change Status <ChevronDown /></Button>
          </PopoverTrigger>
          <PopoverContent class="flex flex-col gap-2 w-36 ">
              <button class="text-green-700 bg-green-200 rounded-lg ">
                Active
              </button>
            <button class="text-red-700  bg-red-200 rounded-lg ">
            Suspend
            </button>
          </PopoverContent>
        </Popover>
        <Dialog>
          <DialogTrigger as-child>
            <Button>
              <SquarePen  />
            </Button>
          </DialogTrigger>
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
    </div>

    <!-- Application Info -->
    <div class="mt-4 bg-gray-100 p-4 rounded-lg">
      <Accordion type="single" collapsible>
        <AccordionItem value="item-1">
          <AccordionTrigger>
            <div class="grid grid-cols-3 gap-6 mt-4">
              <div class="flex justify-start">
                <h2 class="text-md font-medium">Customer ID: <span class="text-orgC">#123456789</span></h2>
              </div>
              <div>
                <h2 class="text-md font-medium">Account Name: <span class="text-orgC">Cameron Williamson</span></h2>
              </div>
              <div class="ml-auto">
                <h2 class="text-md font-medium">Phone Number: <span class="text-orgC">#123456789</span></h2>
              </div>
              <div>
                <p class="text-gray-400 font-medium">Email <span class="text-bgC">cameron.w@example.com</span></p>
              </div>
              <div class="">
                <p class="text-gray-400 font-medium">Employment Status: <span class="text-bgC">Employed</span></p>
              </div>
            </div>
          </AccordionTrigger>
          <AccordionContent>
            <div class="flex justify-between items-center">
              <h2 class="text-lg font-medium">Final Review <span class="text-orgC">4/5</span></h2>
            </div>
            <div class="grid grid-cols-2 gap-4 mt-4">
              <div>
                <p class="text-gray-600">Application ID: <span class="font-medium">#6838960</span></p>
                <p class="text-gray-600">Application Date: <span class="font-medium">Oct 17, 2021, 9:48 AM</span></p>
              </div>
              <div>
                <p class="text-gray-600">Intermediary: <span class="font-medium">ABC Broker (1.5%)</span></p>
              </div>
            </div>          </AccordionContent>
        </AccordionItem>
      </Accordion>

    </div>

    <div class="mt-6">
      <div>
        <Tabs default-value="accounts">
          <TabsList class= "bg-transparent border-b-2 border-gray-200 flex justify-start space-x-4">
            <TabsTrigger value="accounts" class="">
              Accounts
            </TabsTrigger>
            <TabsTrigger value="credit-details" class="">
              Credit Details
            </TabsTrigger>
            <TabsTrigger value="documents" class="">
              Documents
            </TabsTrigger>

            <TabsTrigger value="history" class="">
              History
            </TabsTrigger>
          </TabsList>
          <TabsContent value="accounts">
          </TabsContent>
          <TabsContent value="credit-details">
            test
          </TabsContent>
          <TabsContent value="documents" class="mt-8">
            <h1 class="text-xl mb-4 font-semibold text-bgC">Documents</h1>
            <div>
              <div v-if="isLoading" class="flex flex-col justify-center items-center h-64">
                <div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-orgC"></div>
                <p class="mt-2">Loading Data...</p>
              </div>
              <div v-else>
            <Table class="w-full min-w-[768px]">
              <TableHeader>
                <TableRow class="uppercase text-xs text-gray-400">
                  <TableHead>ID</TableHead>
                  <TableHead>Name</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Size</TableHead>
                  <TableHead class="hidden md:table-cell">Created</TableHead>
                  <TableHead class="hidden md:table-cell">Last Login</TableHead>
                </TableRow>
              </TableHeader>

              <TableBody>
                <TableRow class="" v-for="user in users" :key="user.id">
                  <TableCell>{{ user.id }}</TableCell>
                  <TableCell>{{ user.name }}</TableCell>
                  <TableCell>{{ user.type }}</TableCell>
                  <TableCell>{{ user.size }}</TableCell>
                  <TableCell class="hidden md:table-cell">{{ user.create }}</TableCell>
                  <TableCell class="hidden md:table-cell">{{ user.lastLogin }}</TableCell>

                  <TableCell>
                    <NuxtLink to="/pages/User/view">
                      <button><EllipsisVertical/></button>
                    </NuxtLink>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
              </div>
            </div>
          </TabsContent>
          <TabsContent value="history">
          </TabsContent>
        </Tabs>
      </div>

    </div>
  </div>
  </main>
</template>

<script setup>
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import {EllipsisVertical, ChevronDown, AlertCircle, SquarePen , Slash} from "lucide-vue-next";
import {Popover, PopoverContent, PopoverTrigger} from "~/components/ui/popover/index.js";
import {Input} from "~/components/ui/input/index.js";
import {Tabs, TabsContent, TabsList, TabsTrigger} from "~/components/ui/tabs/index.js";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "~/components/ui/table/index.js";
import {Avatar, AvatarFallback, AvatarImage} from "~/components/ui/avatar/index.js";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "~/components/ui/dialog/index.js";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator
} from "~/components/ui/breadcrumb/index.js";
import {Badge} from "~/components/ui/badge/index.js";

const users =ref( [
  {
    id: '1',
    name: 'UtilityBill.pdf',
    type: 'pdf',
    size: '5.5mb',
    create: 'January 1, 2023',
    lastLogin: 'June 15, 2023',
  },
  {
    id: '2',
    name: 'Driverslicense.png',
    type: 'png',
    size: '5.5mb',
    create: 'January 1, 2023',
    lastLogin: 'June 15, 2023',
  },
 ]);

import {ref, onMounted} from "vue";
const isLoading = ref(true);
onMounted(() => {
  setTimeout(() => {
    isLoading.value = false;
  }, 2000);
});
</script>
