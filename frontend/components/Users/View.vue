<script setup>
import { ref, onMounted } from 'vue'
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { EllipsisVertical, ChevronDown, AlertCircle, SquarePen, Slash } from "lucide-vue-next";
import { Popover, PopoverContent, PopoverTrigger } from "~/components/ui/popover/index.js";
import { Input } from "~/components/ui/input/index.js";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs/index.js";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "~/components/ui/table/index.js";
import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar/index.js";
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "~/components/ui/dialog/index.js";
import { Badge } from "~/components/ui/badge/index.js";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'

// --- Reactive State ---

// Loan form fields
const loantype = ref('');
const loanpurpose = ref('');
const loanamount = ref('');
const loanterm = ref('');
const startdate = ref('');
const enddate = ref('');
const errorMessage = ref('');
const isFormLoading = ref(false);

// Dummy data for the documents table
const users = ref([
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

// Loading state for the main component data
const isLoading = ref(true);

// --- Lifecycle Hooks ---

// Simulate data fetching on component mount
onMounted(() => {
  setTimeout(() => {
    isLoading.value = false;
  }, 2000);
});

// --- Methods ---

/**
 * Handles the creation of a new loan application.
 * In a real app, this would make an API call.
 */
function createLoan() {
  isFormLoading.value = true;
  errorMessage.value = '';
  // Basic form validation
  if (!loantype.value || !loanamount.value) {
    errorMessage.value = "Please fill out all required fields.";
    isFormLoading.value = false;
    return;
  }
  // Simulate API call
  setTimeout(() => {
    console.log("Submitting Loan:", {
      loantype: loantype.value,
      loanpurpose: loanpurpose.value,
      loanamount: loanamount.value,
    });
    isFormLoading.value = false;
    // Here you would typically close the dialog and refresh data
  }, 1500);
}
</script>

<template>
  <main class="flex flex-col z-50 items-start gap-4 p-4 sm:px-6 sm:py-2 md:gap-8">
    <div class="p-2 md:p-6 w-full">

      <!-- Header -->
      <div class="flex flex-col md:flex-row justify-between mt-4 items-start md:items-center border-b pb-4 gap-4">
        <!-- Customer Info -->
        <div class="flex flex-col sm:flex-row items-start sm:items-center gap-4">
          <Avatar class="h-16 w-16 sm:h-20 sm:w-20 bg-gray-600 text-white">
            <AvatarFallback>CW</AvatarFallback>
          </Avatar>
          <div class="flex flex-col ml-0 sm:ml-2">
            <h1 class="text-2xl sm:text-3xl md:text-4xl font-semibold ">Cameron Williamson</h1>
            <div class="flex flex-wrap gap-2 mt-2">
              <p class="text-sm w-full sm:w-auto">Created March 15, 2025, 2:31 PM</p>
              <Badge class="bg-green-200 uppercase text-green-700">
                KYC Compliance: Completed
              </Badge>
              <Badge class="bg-yellow-200 uppercase text-yellow-700">
                Low Risk
              </Badge>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
          <Button class="bg-yellow-400 hover:bg-yellow-500 text-black w-full sm:w-auto">Update Limit <ChevronDown class="ml-2 h-4 w-4" /></Button>
          <Popover>
            <PopoverTrigger as-child>
              <Button class="bg-blue-500 text-white hover:bg-blue-600 w-full sm:w-auto">Change Status <ChevronDown class="ml-2 h-4 w-4" /></Button>
            </PopoverTrigger>
            <PopoverContent class="flex flex-col gap-2 w-36 p-2">
              <button class="text-green-700 bg-green-100 hover:bg-green-200 rounded-md p-2 text-sm">Active</button>
              <button class="text-red-700 bg-red-100 hover:bg-red-200 rounded-md p-2 text-sm">Suspend</button>
            </PopoverContent>
          </Popover>
          <Dialog>
            <DialogTrigger as-child>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Button variant="outline" class="w-full sm:w-auto"><SquarePen class="h-4 w-4" /></Button>
                  </TooltipTrigger>
                  <TooltipContent><p>Edit Customer Details</p></TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </DialogTrigger>
            <DialogContent class="fixed ml-auto right-0 top-0 h-screen w-full sm:w-3/4 md:w-1/2 flex flex-col shadow-lg overflow-y-auto p-6">
              <DialogHeader class="pb-4 border-b">
                <DialogTitle>New Loan Application</DialogTitle>
                <DialogDescription>Add your credit or loan details here.</DialogDescription>
              </DialogHeader>
              <div class="py-4 flex-1">
                <form @submit.prevent="createLoan">
                  <div class="grid gap-2"><Label for="loantype" class="mt-2">Loan Type</Label><Input id="loantype" v-model="loantype" placeholder="Personal loan" required /></div>
                  <div class="grid gap-2 mt-2"><Label for="loanpurpose" class="mt-2">Loan Purpose</Label><Input id="loanpurpose" v-model="loanpurpose" placeholder="Home renovation" required /></div>
                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div class="grid gap-2 mt-2"><Label for="loanamount" class="mt-2">Loan Amount</Label><Input id="loanamount" type="number" v-model="loanamount" placeholder="1000" required /></div>
                    <div class="grid gap-2 mt-2"><Label for="loanterm" class="mt-2">Loan Term (Months)</Label><Input id="loanterm" type="number" v-model="loanterm" placeholder="12" required /></div>
                  </div>
                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div class="grid gap-2 mt-2"><Label for="startdate" class="mt-2">Loan Start Date</Label><Input id="startdate" type="date" v-model="startdate" required /></div>
                    <div class="grid gap-2 mt-2"><Label for="enddate" class="mt-2">Loan End Date</Label><Input id="enddate" type="date" v-model="enddate" required /></div>
                  </div>
                  <Alert v-if="errorMessage" variant="destructive" class="mt-4"><AlertCircle class="w-4 h-4" /><AlertTitle>Error</AlertTitle><AlertDescription>{{ errorMessage }}</AlertDescription></Alert>
                  <div class="pt-8 flex-1"><Button type="submit" :disabled="isFormLoading">Save changes</Button></div>
                </form>
              </div>
            </DialogContent>
          </Dialog>
        </div>
      </div>

      <!-- Application Info Accordion -->
      <div class="mt-4 shadow-md p-4 rounded-lg border">
        <Accordion type="single" collapsible>
          <AccordionItem value="item-1">
            <AccordionTrigger>
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 w-full text-left">
                <h2 class="text-sm font-medium">Customer ID: <span class="font-normal ">#123456789</span></h2>
                <h2 class="text-sm font-medium">Account Name: <span class="font-normal ">Cameron Williamson</span></h2>
                <h2 class="text-sm font-medium">Phone Number: <span class="font-normal ">#123456789</span></h2>
                <p class="text-sm font-medium">Email: <span class="font-normal ">cameron.w@example.com</span></p>
                <p class="text-sm font-medium">Employment: <span class="font-normal ">Employed</span></p>
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
                <div>
                  <p class="text-sm">Application ID: <span class="font-medium">#6838960</span></p>
                  <p class="text-sm">Application Date: <span class="font-medium">Oct 17, 2021, 9:48 AM</span></p>
                </div>
                <div><p class="text-sm">Intermediary: <span class="font-medium">ABC Broker (1.5%)</span></p></div>
              </div>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>

      <!-- Tabs Section -->
      <div class="mt-6">
        <Tabs default-value="documents">
          <TabsList class="bg-transparent border-b flex-wrap h-auto justify-start">
            <TabsTrigger value="accounts">Accounts</TabsTrigger>
            <TabsTrigger value="credit-details">Credit Details</TabsTrigger>
            <TabsTrigger value="documents">Documents</TabsTrigger>
            <TabsTrigger value="history">History</TabsTrigger>
          </TabsList>

          <TabsContent value="accounts"><p class="p-4">Account details would be displayed here.</p></TabsContent>
          <TabsContent value="credit-details"><p class="p-4">Credit details would be displayed here.</p></TabsContent>
          <TabsContent value="documents" class="mt-4">
            <h1 class="text-xl mb-4 font-semibold">Documents</h1>
            <div v-if="isLoading" class="flex flex-col justify-center items-center h-64"><div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-primary"></div><p class="mt-2">Loading Data...</p></div>
            <div v-else class="overflow-x-auto">
              <Table class="w-full responsive-table">
                <TableHeader>
                  <TableRow class="uppercase text-xs text-gray-400">
                    <TableHead>ID</TableHead><TableHead>Name</TableHead><TableHead>Type</TableHead><TableHead>Size</TableHead><TableHead>Created</TableHead><TableHead>Last Login</TableHead><TableHead>Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="user in users" :key="user.id">
                    <TableCell data-label="ID:">{{ user.id }}</TableCell>
                    <TableCell data-label="Name:">{{ user.name }}</TableCell>
                    <TableCell data-label="Type:">{{ user.type }}</TableCell>
                    <TableCell data-label="Size:">{{ user.size }}</TableCell>
                    <TableCell data-label="Created:">{{ user.create }}</TableCell>
                    <TableCell data-label="Last Login:">{{ user.lastLogin }}</TableCell>
                    <TableCell>
                      <NuxtLink to="/User/view"><button class="p-2 rounded-full"><EllipsisVertical class="h-5 w-5"/></button></NuxtLink>
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </TabsContent>
          <TabsContent value="history"><p class="p-4">User history would be displayed here.</p></TabsContent>
        </Tabs>
      </div>
    </div>
  </main>
</template>

<style scoped>
/* Scoped styles for the responsive table */
@media (max-width: 768px) {
  .responsive-table thead {
    /* Hide table headers on mobile, we will use data-labels instead */
    display: none;
  }
  .responsive-table tbody,
  .responsive-table tr,
  .responsive-table td {
    /* Make table elements display as blocks for a card-like layout */
    display: block;
    width: 100%;
    text-align: right; /* Align cell content to the right */
  }
  .responsive-table tr {
    margin-bottom: 1rem;
    border: 1px solid #e2e8f0; /* Add a border to each "card" */
    border-radius: 0.5rem;
    padding: 1rem;
  }
  .responsive-table td::before {
    /* Use data-label for pseudo-headers */
    content: attr(data-label);
    font-weight: 600;
    float: left; /* Align the label to the left */
    margin-right: 0.5rem;
  }
  .responsive-table td:last-child {
    /* Adjust the actions cell for better alignment */
    text-align: right;
    padding-top: 1rem;
    border-top: 1px solid #e2e8f0;
    margin-top: 1rem;
  }
}
</style>
