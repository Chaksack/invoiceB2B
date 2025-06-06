<template>
  <main class="font-inter">
    <Toaster richColors position="top-right" /> <section>
    <div class="flex mb-2 gap-8 pt-10 items-center py-2 px-4 mx-auto max-w-screen-xl ">
      <div>
        <h3 v-if="isLoadingUser && !user">Loading user...</h3>
        <h3 v-else-if="userError">{{ userError }}</h3>
        <div v-if="user">
          <h2 class="text-2xl tracking-tight font-bold">Welcome,
            <span class=" bg-gradient-to-r from-sky-500 via-purple-500 to-pink-500 bg-clip-text text-transparent ">
                {{ user.firstName }} {{ user.lastName }}
              </span>
          </h2>
          <p class="mt-2">Submit, track, and fund your business invoices with ease.</p>
          <p class="mt-1 text-sm">Company: {{ user.companyName }} | KYC Status: <span :class="kycStatusClass(user.kycStatus)">{{ user.kycStatus }}</span></p>
        </div>
        <div v-else-if="!isLoadingUser && !userError && !user">
          <h2 class="text-2xl tracking-tight font-bold">Welcome!</h2>
          <p class="mt-2">Could not load user details. Please ensure you are logged in.</p>
        </div>
      </div>

      <div class="ml-auto">
        <Dialog v-model:open="isUploadDialogOpen">
          <DialogTrigger as-child>
            <Button @click="isUploadDialogOpen = true" class=" bg-primary rounded-md hover:bg-gray-200 transition-colors">
              <Upload class="mr-2 h-4 w-4" />Upload Invoice
            </Button>
          </DialogTrigger>
          <DialogContent class="sm:max-w-[600px] shadow-lg rounded-lg">
            <div class="">
              <FileUpload class="">
                <FileUploadGrid />
                <p v-if="fileUploadError" class="text-red-500 text-sm mt-1">{{ fileUploadError }}</p>
              </FileUpload>
            </div>
            <DialogFooter>
              <Button @click="submitInvoice" :disabled="isUploadingInvoice || !selectedFile" class="rounded-md">
                <span v-if="isUploadingInvoice">Uploading...</span>
                <span v-else>Upload Invoice</span>
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  </section>

    <div v-if="user && user.kycStatus !== 'approved' && user.kycStatus !== 'Verified'" class="relative bg-red-500 isolate flex items-center gap-x-6 overflow-hidden bg-gray-50 px-6 py-2.5 sm:px-3.5 sm:before:flex-1 mb-6">
      <div class="flex flex-wrap items-center gap-x-4 gap-y-2">
        <p class="text-sm/6 text-white flex items-center">
          <FileWarning class="mr-2 h-5 w-5" /> <strong class="font-semibold">Compliance Form Incomplete</strong>
          <svg viewBox="0 0 2 2" class="mx-2 inline size-0.5 fill-current" aria-hidden="true"><circle cx="1" cy="1" r="1" /></svg>
          To get more out of your score kindly complete the compliance forms.
        </p>
        <Dialog v-model:open="isUploadDialogOpen"> <DialogTrigger as-child>
          <Button @click="isUploadDialogOpen = true" class="text-white bg-black rounded-md hover:bg-gray-800 transition-colors">
            Complete now
          </Button>
        </DialogTrigger>
          <DialogContent class="sm:max-w-[425px] shadow-lg rounded-lg">
            <DialogHeader>
              <DialogTitle>Complete KYC</DialogTitle>
              <DialogDescription>
                Upload your business registration files (PDF, CSV, XLSX formats supported).
              </DialogDescription>
            </DialogHeader>
            <div class="grid w-full max-w-sm items-center gap-1.5 py-4">
              <Label for="kycFile">Business Registration</Label> <Input id="kycFile" type="file" @change="handleFileSelect" accept=".pdf,.csv,.xlsx,.xls" class="rounded-md"/> <p v-if="fileUploadError" class="text-red-500 text-sm mt-1">{{ fileUploadError }}</p>
            </div>
            <DialogFooter>
              <Button @click="submitInvoice" :disabled="isUploadingInvoice || !selectedFile" class="rounded-md">
                <span v-if="isUploadingInvoice">Uploading...</span>
                <span v-else>Upload Business files</span>
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
      <div class="flex flex-1 justify-end">
        <button type="button" class="-m-3 p-3 focus-visible:outline-offset-[-4px]" @click="dismissComplianceBanner">
          <span class="sr-only">Dismiss</span>
          <CircleX class="size-5 text-white" aria-hidden="true" />
        </button>
      </div>
    </div>

    <div class=" mx-auto mt-2 grid w-full max-w-6xl flex-1 auto-rows-max gap-6 lg:gap-8 px-4">
      <section class="flex justify-center w-full">
        <div class="grid w-full max-w-4xl items-start justify-items-center gap-4 sm:grid-cols-2 md:grid-cols-3">
          <Card class="w-full rounded-lg shadow-md">
            <CardHeader class="p-4">
              <CardTitle class="flex items-center text-blue-500 text-2xl font-bold mt-1">
                <FileText class="mr-4 h-7 w-7"/>
                {{ summaryStats.totalInvoices }}
              </CardTitle>
              <CardDescription class="font-medium text-lg ">Total Invoices</CardDescription>
              <CardDescription class="font-medium text-xs ">Submitted to date</CardDescription>
            </CardHeader>
          </Card>
          <Card class="w-full rounded-lg shadow-md">
            <CardHeader class="p-4">
              <CardTitle class="flex items-center text-2xl text-green-600 font-bold mt-1">
                <CircleCheck class="mr-4 h-7 w-7 fill-current" stroke="#ffffff" stroke-width="1"/>
                {{ summaryStats.approvedInvoices }}
              </CardTitle>
              <CardDescription class="font-medium text-lg ">Approved</CardDescription>
              <CardDescription class="font-medium text-xs ">Eligible for funding</CardDescription>
            </CardHeader>
          </Card>
          <Card class=" w-full rounded-lg shadow-md">
            <CardHeader class="p-4">
              <CardTitle class="flex text-2xl font-bold text-yellow-500 mt-1 items-center">
                <HandCoins class="mr-4 h-7 w-7 items-center"/> GHS {{ summaryStats.totalFunded.toLocaleString() }}
              </CardTitle>
              <CardDescription class="font-medium text-lg ">Total Funded</CardDescription>
              <CardDescription class="font-medium text-xs">Received from bank</CardDescription>
            </CardHeader>
          </Card>
        </div>
      </section>
    </div>

    <section class="">
      <div class="py-8 px-4 mx-auto max-w-screen-xl lg:py-16 lg:px-6">
        <div class="flex items-center mb-6 gap-4">
          <h2 class="text-xl tracking-tight font-semibold ">Recent Invoices</h2>
          <div class="ml-auto flex items-center gap-2">
            <!-- Search Input -->
            <div class="relative w-full max-w-sm">
              <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input v-model="searchQuery" placeholder="Search by customer, invoice #..." class="pl-9 rounded-md" />
            </div>
            <!-- Status Filter -->
            <Select v-model="statusFilter">
              <SelectTrigger class="w-[180px] rounded-md">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="option in statusOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <div v-if="isLoadingInvoices" class="text-center py-10">
          <p>Loading invoices...</p>
        </div>
        <div v-else-if="invoicesError" class="text-center py-10 text-red-500">
          <p>{{ invoicesError }}</p>
        </div>
        <div v-else-if="filteredInvoices.length === 0" class="text-center py-10 flex flex-col items-center">
          <Ban class="w-12 h-12 mb-4 text-muted-foreground" />
          <p class=" text-2xl font-semibold">No Invoices Found</p>
          <p v-if="searchQuery || statusFilter !== 'all'" class="text-muted-foreground">Try adjusting your search or filters to find what you're looking for.</p>
          <p v-else class="text-muted-foreground">Upload your first invoice to get started!</p>
        </div>

        <div v-else class="space-y-4">
          <div v-for="invoice in filteredInvoices" :key="invoice.id" class="rounded-lg p-4 shadow-lg hover:shadow-xl transition-shadow">
            <Accordion type="single" class="w-full" collapsible>
              <AccordionItem class="border-none" :value="`invoice-${invoice.id}`">
                <AccordionTrigger class="mx-2 sm:mx-8 mb-2 items-start hover:no-underline focus:no-underline">
                  <div class="grid gap-2 md:flex w-full items-center">
                    <Avatar class="relative overflow-visible h-10 w-10">
                      <AvatarImage :src="invoice.customer?.avatarUrl || ''" class="rounded-full" alt="Customer Avatar" />
                      <AvatarFallback class="bg-gray-250 rounded-full flex items-center justify-center">
                        {{ invoice.customerName ? invoice.customerName.substring(0, 2).toUpperCase() : 'N/A' }}
                      </AvatarFallback>
                    </Avatar>
                    <div class="text-md font-semibold text-left">
                      {{ invoice.customerName || 'N/A' }}
                    </div>
                    <div class="text-md text-muted-foreground text-left">
                      #{{ invoice.invoiceNumber || invoice.id }}
                    </div>
                    <Badge :class="getInvoiceStatusBadgeClass(invoice.status)" class="py-1 px-3 rounded-full text-xs">
                      <component :is="getInvoiceStatusIcon(invoice.status)" class="w-4 h-4 mr-1"/>
                      {{ formatStatus(invoice.status) }}
                    </Badge>
                  </div>
                  <div class="gap-4 ml-auto md:flex items-center text-right">
                    <div class="text-sm">
                      Invoice Amount:
                      <div class="text-xl text-black font-semibold">
                        {{ invoice.currency || 'GHS' }} {{ (invoice.totalAmount || 0).toLocaleString() }}
                      </div>
                    </div>
                  </div>
                </AccordionTrigger>
                <AccordionContent class="border-t pt-4">
                  <div class="py-6 px-4 mx-auto max-w-screen-xl">
                    <div class="w-full mb-8">
                      <div class="flex items-start">
                        <template v-for="(step, index) in predefinedSteps" :key="step.step">
                          <div class="flex flex-col items-center w-1/5">
                            <div
                                class="w-10 h-10 rounded-full flex items-center justify-center"
                                :class="{
                                'bg-blue-600 text-white': getStepForInvoiceStatus(invoice.status) === step.step,
                                'bg-green-500 text-white': getStepForInvoiceStatus(invoice.status) > step.step,
                                'bg-gray-300 text-gray-600': getStepForInvoiceStatus(invoice.status) < step.step
                              }"
                            >
                              <component :is="step.icon" class="w-5 h-5" />
                            </div>
                            <p
                                class="mt-2 text-xs text-center font-medium"
                                :class="{
                                'text-blue-600': getStepForInvoiceStatus(invoice.status) === step.step,
                                'text-green-500': getStepForInvoiceStatus(invoice.status) > step.step,
                                'text-gray-500': getStepForInvoiceStatus(invoice.status) < step.step
                              }"
                            >
                              {{ step.title }}
                            </p>
                            <p class="mt-1 text-xs text-center px-1" style="min-height: 2.5em;"> {{ step.description }}
                            </p>
                          </div>
                          <div
                              v-if="index < predefinedSteps.length - 1"
                              class="flex-1 h-1 mt-5"
                              :class="{
                              'bg-green-500': getStepForInvoiceStatus(invoice.status) > step.step,
                              'bg-gray-300': getStepForInvoiceStatus(invoice.status) <= step.step
                            }"
                          ></div>
                        </template>
                      </div>
                    </div>
                    <h2 class="mt-4 mb-4 text-xl tracking-tight font-semibold ">Invoice Details</h2>
                    <div class="space-y-8 md:grid md:grid-cols-3 lg:grid-cols-3 md:gap-12 md:space-y-0">
                      <div class="col-span-2 max-w-screen-lg  sm:text-lg">
                        <h3 class="font-semibold text-black">{{ invoice.customerName || 'N/A' }}</h3>
                        <p class="text-sm">{{ invoice.customerAddress?.street || 'Street not available' }}</p>
                        <p class="text-sm">{{ invoice.customerAddress?.city || 'City not available' }}</p>
                        <p class="text-sm">{{ invoice.customerAddress?.zipCode }} {{ invoice.customerAddress?.country }}</p>
                      </div>
                      <div class="col-span-1 md:grid md:grid-cols-2 lg:grid-cols-2 text-sm">
                        <p class="font-semibold">Invoice #</p>
                        <p class=" font-medium">{{ invoice.invoiceNumber || invoice.id }}</p>
                        <p class="">Invoice Date </p>
                        <p class=" font-medium">{{ formatDate(invoice.invoiceDate) }}</p>
                        <p class="">Terms</p>
                        <p class=" font-medium">{{ invoice.paymentTerms || 'N/A' }}</p>
                        <p class="">Due Date</p>
                        <p class=" font-medium">{{ formatDate(invoice.dueDate) }}</p>
                      </div>
                    </div>

                    <Table v-if="invoice.items && invoice.items.length > 0" class="mt-6 w-full min-w-[768px]">
                      <TableHeader>
                        <TableRow>
                          <TableHead>ID</TableHead>
                          <TableHead>ITEM & DESCRIPTION</TableHead>
                          <TableHead class="text-right">UNIT PRICE ({{ invoice.currency || 'GHS' }})</TableHead>
                          <TableHead class="text-right">QUANTITY</TableHead>
                          <TableHead class="text-right">AMOUNT ({{ invoice.currency || 'GHS' }})</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        <TableRow v-for="item in invoice.items" :key="item.id">
                          <TableCell>{{ item.id }}</TableCell>
                          <TableCell>
                            <div>{{ item.name }}</div>
                            <div class="text-xs font-light text-muted-foreground">{{ item.description }}</div>
                          </TableCell>
                          <TableCell class="text-right">{{ item.unitPrice?.toLocaleString() || '0.00' }}</TableCell>
                          <TableCell class="text-right">{{ item.quantity }}</TableCell>
                          <TableCell class="text-right">{{ item.amount?.toLocaleString() || '0.00' }}</TableCell>
                        </TableRow>
                      </TableBody>
                    </Table>
                    <p v-else class="mt-6">No items listed for this invoice.</p>


                    <div class="mt-6 max-w-screen-lg bg-gray-100 p-4 rounded-md  sm:text-lg">
                      <div class="flex col-span-1 md:grid md:grid-cols-3 lg:grid-cols-3">
                        <div class="col-span-2">
                          <p class="text-sm font-semibold">Invoice Totals</p>
                        </div>
                        <div class="col-span-1 text-sm">
                          <div class="flex items-center">
                            <p>Sub Total</p>
                            <p class=" font-semibold ml-auto">{{ invoice.currency || 'GHS' }} {{ invoice.subTotalAmount?.toLocaleString() || '0.00' }}</p>
                          </div>
                          <div class="flex items-center">
                            <p>Tax Rate ({{ invoice.taxRatePercentage || 0 }}%)</p>
                            <p class=" font-semibold ml-auto">{{ invoice.currency || 'GHS' }} {{ invoice.taxAmount?.toLocaleString() || '0.00' }}</p>
                          </div>
                          <hr class="my-1 border-gray-300">
                          <div class="flex items-center font-bold">
                            <p>Total</p>
                            <p class=" ml-auto">{{ invoice.currency || 'GHS' }} {{ invoice.totalAmount?.toLocaleString() || '0.00' }}</p>
                          </div>
                          <div class="flex items-center text-green-600">
                            <p>Balance Due</p>
                            <p class="font-semibold ml-auto">{{ invoice.currency || 'GHS' }} {{ invoice.balanceDue?.toLocaleString() || '0.00' }}</p>
                          </div>
                        </div>
                      </div>
                    </div>

                    <div class="flex justify-end mt-6 space-x-2">
                      <Dialog v-if="invoice.status === 'DISBURSED' || invoice.status === 'PAID'"> <DialogTrigger as-child>
                        <Button @click="fetchReceiptDetails(invoice.id)" variant="outline" class=" border-gray-700 hover:bg-gray-200 px-6 text-sm rounded-md">
                          <Eye class="mr-2 h-4 w-4" /> View Receipt
                        </Button>
                      </DialogTrigger>
                        <DialogContent class="sm:max-w-[625px] shadow-lg rounded-lg">
                          <DialogHeader class=" p-6 rounded-t-lg">
                            <DialogTitle class="">InvoiceB2B Platform</DialogTitle>
                          </DialogHeader>
                          <DialogTitle class="px-6 py-4 text-xl font-semibold">Transfer Receipt</DialogTitle>
                          <div v-if="isLoadingReceipt" class="p-6 text-center">Loading receipt...</div>
                          <div v-else-if="receiptError" class="p-6 text-center text-red-500">{{ receiptError }}</div>
                          <div v-else-if="currentReceipt" class="col-span-1 p-6 md:grid md:grid-cols-2 lg:grid-cols-2 gap-x-4 gap-y-2 text-sm">
                            <p class="">Reference Number:</p><p class=" font-medium">{{ currentReceipt.referenceNumber }}</p>
                            <p class="">Transfer to:</p><p class=" font-medium">{{ currentReceipt.transferTo }}</p>
                            <p class="">Account Type:</p><p class=" font-medium">{{ currentReceipt.accountType }}</p>
                            <p class="">Account Number:</p><p class=" font-medium">{{ currentReceipt.accountNumber }}</p>
                            <p class="">Account Name:</p><p class=" font-medium">{{ currentReceipt.accountName }}</p>
                            <p class="">Amount:</p><p class=" font-medium">{{ invoice.currency || 'GHS' }} {{ currentReceipt.amount?.toLocaleString() }}</p>
                            <p class="">Transfer Date:</p><p class=" font-medium">{{ formatDate(currentReceipt.transferDate) }}</p>
                            <p class="">Purpose:</p><p class=" font-medium col-span-2">{{ currentReceipt.purpose }}</p>
                          </div>
                          <DialogDescription class="px-6 pb-4 text-xs ">
                            This is a computer-generated receipt; no signature is required.
                            Electronic receipts may not have official legal effect. You may go to a branch to get a paper receipt.
                          </DialogDescription>
                          <DialogFooter class="p-6 border-t">
                            <Button @click="triggerDownloadReceipt(invoice.id)" class=" bg-blue-700 hover:bg-blue-800 px-6 text-sm rounded-md">
                              <Download class="mr-2 h-4 w-4" /> Download Receipt
                            </Button>
                          </DialogFooter>
                        </DialogContent>
                      </Dialog>
                      <Button v-else @click="triggerDownloadReceipt(invoice.id)" class=" bg-blue-700 hover:bg-blue-800 px-6 text-sm rounded-md">
                        <Download class="mr-2 h-4 w-4" /> Download Invoice PDF
                      </Button>
                    </div>
                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import axios from 'axios'
import {
  Upload, Download, Eye, CreditCard, CircleCheckBig, FileWarning, FileText,
  CircleCheck, BookUser, HandCoins, Clock, AlertCircle, CircleX, Ban, Search
} from 'lucide-vue-next'
import { useCookie } from '#app';
import { Toaster, toast } from 'vue-sonner'

// Shadcn-vue components
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Label } from '@/components/ui/label'

// API Configuration
const API_BASE_URL = 'http://localhost:3000/api/v1'
const tokenCookie = useCookie('token');
const authToken = tokenCookie.value || null;

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    ...(authToken ? { 'Authorization': `Bearer ${authToken}` } : {}),
    'Content-Type': 'application/json'
  }
})

// Reactive State
const user = ref<any>(null)
const isLoadingUser = ref(false)
const userError = ref<string | null>(null)

const invoices = ref<any[]>([])
const isLoadingInvoices = ref(false)
const invoicesError = ref<string | null>(null)

const selectedFile = ref<File | null>(null)
const isUploadingInvoice = ref(false)
const fileUploadError = ref<string | null>(null)
const isUploadDialogOpen = ref(false);

const currentReceipt = ref<any>(null)
const isLoadingReceipt = ref(false)
const receiptError = ref<string | null>(null)

const searchQuery = ref('');
const statusFilter = ref('all');

const statusOptions = [
  { value: 'all', label: 'All Statuses' },
  { value: 'SUBMITTED', label: 'Submitted' },
  { value: 'PENDING REVIEW', label: 'Under Review' },
  { value: 'APPROVED', label: 'Approved' },
  { value: 'DISBURSED', label: 'Disbursed' },
  { value: 'PAID', label: 'Paid' },
  { value: 'REJECTED', label: 'Rejected' },
];


// --- User Profile ---
const fetchUser = async () => {
  if (!authToken) {
    userError.value = "Authentication token not found. Please log in.";
    isLoadingUser.value = false;
    return;
  }
  isLoadingUser.value = true
  userError.value = null
  try {
    const response = await apiClient.get('/user/profile')
    user.value = response.data
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      console.error('API Error (User Profile):', err.response.status, err.response.data)
      if (err.response.status === 401 || err.response.status === 403) {
        userError.value = `Authentication error: ${err.response.data?.message || 'Please log in again.'}`;
      } else {
        userError.value = `Error ${err.response.status}: ${err.response.data?.message || 'Failed to fetch user details'}`
      }
    } else {
      console.error('Failed to fetch user details:', err)
      userError.value = 'An unexpected error occurred while fetching user details.'
    }
  } finally {
    isLoadingUser.value = false
  }
}

// --- Invoices ---
const fetchInvoices = async () => {
  if (!authToken) {
    invoicesError.value = "Authentication token not found. Please log in.";
    isLoadingInvoices.value = false;
    return;
  }
  isLoadingInvoices.value = true;
  invoicesError.value = null;
  try {
    const response = await apiClient.get('/invoices');
    console.log('Raw invoice API response:', response);

    let invoiceDataArray = [];
    if (Array.isArray(response.data)) {
      invoiceDataArray = response.data;
    } else if (response.data && Array.isArray(response.data.data)) {
      invoiceDataArray = response.data.data;
    } else if (response.data && Array.isArray(response.data.invoices)) {
      invoiceDataArray = response.data.invoices;
    } else if (response.status === 200 && response.data && typeof response.data === 'object' && !Array.isArray(response.data)) {
      console.warn("Invoice API returned 200 OK with a non-array payload. Assuming empty list.", response.data);
      invoiceDataArray = [];
    } else if (!Array.isArray(response.data)) {
      console.error("Invoice API response.data is not an array.", response.data);
      throw new Error("Invoice data from API is not in expected array format.");
    }

    invoices.value = invoiceDataArray.map((inv: any) => {
      let parsedJsonData: any = null;
      if (inv.jsonData && typeof inv.jsonData === 'string') {
        try {
          parsedJsonData = JSON.parse(inv.jsonData);
        } catch (e) {
          console.error(`Failed to parse jsonData for invoice ID ${inv.id}:`, e, inv.jsonData);
        }
      }

      // Initialize with base inv values or defaults
      let customerName = inv.customer?.name || inv.companyName || 'Unknown Customer';
      let invoiceNumber = inv.invoiceNumber || String(inv.id);
      let invoiceDate = inv.invoiceDate || inv.createdAt;
      let items: any[] = inv.items || [];
      let currency = inv.currency || 'GHS';
      let totalAmount = parseFloat(inv.totalAmount) || 0;
      let subTotalAmount = parseFloat(inv.subTotal || inv.subTotalAmount) || 0;
      let taxAmount = parseFloat(inv.tax || inv.taxAmount) || 0;
      let balanceDue = parseFloat(inv.balanceDue);
      let taxRatePercentage = parseFloat(inv.taxRate) || 0;
      let paymentTerms = inv.terms || inv.paymentTerms;
      let dueDate = inv.dueDate;


      if (parsedJsonData) {
        customerName = parsedJsonData.billedTo || customerName;
        invoiceNumber = parsedJsonData.extractedInvoiceNumber || invoiceNumber;
        invoiceDate = parsedJsonData.invoiceDate || invoiceDate;
        currency = parsedJsonData.extractedCurrency || currency;
        paymentTerms = parsedJsonData.paymentTerms || paymentTerms;
        dueDate = parsedJsonData.dueDate || dueDate;

        // Prioritize totals from parsedJsonData
        if (typeof parsedJsonData.total !== 'undefined') {
          totalAmount = parseFloat(parsedJsonData.total) || 0;
        } else if (typeof parsedJsonData.grandTotal !== 'undefined') {
          totalAmount = parseFloat(parsedJsonData.grandTotal) || 0;
        }

        if (typeof parsedJsonData.subtotal !== 'undefined') {
          subTotalAmount = parseFloat(parsedJsonData.subtotal) || 0;
        }
        if (typeof parsedJsonData.tax !== 'undefined') {
          taxAmount = parseFloat(parsedJsonData.tax) || 0;
        } else if (typeof parsedJsonData.taxAmount !== 'undefined') {
          taxAmount = parseFloat(parsedJsonData.taxAmount) || 0;
        }

        balanceDue = (typeof parsedJsonData.balanceDue !== 'undefined') ? (parseFloat(parsedJsonData.balanceDue) || totalAmount) : totalAmount;

        if (typeof parsedJsonData.taxRatePercentage !== 'undefined') {
          taxRatePercentage = parseFloat(parsedJsonData.taxRatePercentage) || 0;
        } else if (subTotalAmount > 0 && taxAmount > 0) { // Calculate if not present
          taxRatePercentage = parseFloat(((taxAmount / subTotalAmount) * 100).toFixed(2));
        }
        // else, it keeps the value from inv.taxRate or 0 from initialization

        if (parsedJsonData.lineItems && Array.isArray(parsedJsonData.lineItems)) {
          items = parsedJsonData.lineItems.map((item: any, index: number) => ({
            id: item.id || `jsonItem-${inv.id}-${index + 1}`,
            name: item.item || item.name || 'N/A',
            description: item.description || '',
            unitPrice: parseFloat(item.unitPrice) || 0,
            quantity: parseInt(item.quantity, 10) || 0,
            amount: parseFloat(item.total || item.amount) || 0,
          }));
        }
      }

      if (isNaN(balanceDue)) {
        balanceDue = totalAmount;
      }

      // Normalize status to uppercase and replace underscores with spaces
      const normalizedStatus = inv.status ? String(inv.status).trim().replace(/_/g, ' ').toUpperCase() : 'UNKNOWN';

      return {
        ...inv,
        id: inv.id,
        customerName: customerName,
        invoiceNumber: invoiceNumber,
        totalAmount: totalAmount,
        subTotalAmount: subTotalAmount,
        taxAmount: taxAmount,
        balanceDue: balanceDue,
        taxRatePercentage: taxRatePercentage,
        items: items,
        status: normalizedStatus, // Use the consistently formatted status
        invoiceDate: invoiceDate,
        dueDate: dueDate,
        paymentTerms: paymentTerms,
        customerAddress: inv.customer?.address || (parsedJsonData?.customerAddress) || {},
        currency: currency,
      };
    }).sort((a: any, b: any) => {
      const dateA = a.invoiceDate ? new Date(a.invoiceDate).getTime() : 0;
      const dateB = b.invoiceDate ? new Date(b.invoiceDate).getTime() : 0;
      return dateB - dateA;
    });

  } catch (err: any) {
    console.error('Error during fetchInvoices processing:', err);
    if (err.response) {
      console.error('Axios error response data:', err.response.data);
      console.error('Axios error response status:', err.response.status);
    }
    if (axios.isAxiosError(err) && err.response) {
      if (err.response.status === 401 || err.response.status === 403) {
        invoicesError.value = `Authentication error: ${err.response.data?.message || 'Please log in again.'}`;
      } else {
        invoicesError.value = `Error ${err.response.status}: ${err.response.data?.message || 'Failed to fetch invoices'}`;
      }
    } else {
      invoicesError.value = `An unexpected error occurred: ${err.message || 'Failed to process invoice data.'}`;
    }
  } finally {
    isLoadingInvoices.value = false;
  }
};

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    selectedFile.value = target.files[0];
    console.log('File selected:', selectedFile.value);
    fileUploadError.value = null;
  } else {
    selectedFile.value = null;
  }
};

const submitInvoice = async () => {
  if (!authToken) {
    toast.error("Authentication token not found. Please log in.");
    isUploadingInvoice.value = false;
    return;
  }
  if (!selectedFile.value) {
    fileUploadError.value = "Please select a file to upload.";
    return;
  }

  console.log('Submitting invoice with file:', selectedFile.value);

  isUploadingInvoice.value = true;
  fileUploadError.value = null;
  const formData = new FormData();
  formData.append('invoiceFile', selectedFile.value);

  try {
    const uploadApiClient = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        ...(authToken ? { 'Authorization': `Bearer ${authToken}` } : {}),
      }
    });
    const response = await uploadApiClient.post('/invoices', formData);

    console.log('Invoice uploaded:', response.data);
    toast.success('Invoice uploaded successfully!');
    await fetchInvoices();
    selectedFile.value = null;
    isUploadDialogOpen.value = false;
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      console.error('API Error (Upload):', err.response.status, err.response.data);
      const errorMessage = err.response.data?.message || err.response.data?.error || 'Failed to upload invoice';
      if (err.response.status === 401 || err.response.status === 403) {
        toast.error(`Authentication error: ${errorMessage}`);
      } else {
        if (err.response.status === 400 && (typeof errorMessage === 'string' && errorMessage.toLowerCase().includes("invoice file"))) {
          fileUploadError.value = `Upload Error: ${errorMessage}`;
        } else {
          toast.error(`Upload Error ${err.response.status}: ${errorMessage}`);
        }
      }
    } else {
      console.error('Failed to upload invoice:', err);
      toast.error('An unexpected error occurred during upload.');
    }
  } finally {
    isUploadingInvoice.value = false;
  }
};


// --- Receipt Logic ---
const fetchReceiptDetails = async (invoiceId: string | number) => {
  if (!authToken) {
    receiptError.value = "Authentication token not found. Please log in.";
    isLoadingReceipt.value = false;
    return;
  }
  isLoadingReceipt.value = true
  receiptError.value = null
  currentReceipt.value = null
  try {
    const response = await apiClient.get(`/invoices/${invoiceId}/viewreceipt`)
    currentReceipt.value = response.data
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      console.error('API Error (View Receipt):', err.response.status, err.response.data)
      const errorMessage = err.response.data?.message || 'Failed to fetch receipt';
      if (err.response.status === 401 || err.response.status === 403) {
        receiptError.value = `Authentication error: ${errorMessage}`;
      } else {
        receiptError.value = `Error ${err.response.status}: ${errorMessage}`
      }
    } else {
      console.error('Failed to fetch receipt:', err)
      receiptError.value = 'An unexpected error occurred while fetching receipt details.'
    }
  } finally {
    isLoadingReceipt.value = false
  }
}

const triggerDownloadReceipt = async (invoiceId: string | number) => {
  if (!authToken) {
    toast.error("Authentication token not found. Please log in.");
    return;
  }
  toast.info("Preparing download...");
  try {
    const response = await apiClient.get(`/invoices/${invoiceId}/receipt`, {
      responseType: 'blob',
    })
    const blob = new Blob([response.data], { type: response.headers['content-type'] || 'application/pdf' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)

    const contentDisposition = response.headers['content-disposition'];
    let filename = `invoice-${invoiceId}-document.pdf`;
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename="?(.+)"?/i);
      if (filenameMatch && filenameMatch.length === 2)
        filename = filenameMatch[1];
    }
    link.download = filename;

    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(link.href)
    toast.success("Download started: " + filename);
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      console.error('API Error (Download Document):', err.response.status, err.response.data)
      if (err.response.data instanceof Blob && err.response.data.type === "application/json") {
        const errorText = await err.response.data.text();
        try {
          const errorJson = JSON.parse(errorText);
          toast.error(`Error ${err.response.status}: ${errorJson.message || 'Failed to download. Please try again.'}`);
        } catch (parseError) {
          toast.error(`Error ${err.response.status}: Failed to download. Please try again.`);
        }
      } else {
        toast.error(`Error ${err.response.status}: ${err.response.data?.message || 'Failed to download. Please try again.'}`);
      }
    } else {
      console.error('Failed to download document:', err)
      toast.error('An unexpected error occurred while downloading the document.')
    }
  }
}


// --- UI Helpers & Computed Properties ---
const filteredInvoices = computed(() => {
  let invoicesToShow = invoices.value;

  // Filter by search query
  if (searchQuery.value) {
    const lowerCaseQuery = searchQuery.value.toLowerCase();
    invoicesToShow = invoicesToShow.filter(invoice =>
        (invoice.customerName?.toLowerCase().includes(lowerCaseQuery)) ||
        (invoice.invoiceNumber?.toLowerCase().includes(lowerCaseQuery)) ||
        (String(invoice.totalAmount).includes(lowerCaseQuery))
    );
  }

  // Filter by status
  if (statusFilter.value && statusFilter.value !== 'all') {
    invoicesToShow = invoicesToShow.filter(invoice => {
      // Handle cases where a group of statuses are selected, e.g., 'Under Review'
      if (statusFilter.value === 'PENDING REVIEW') {
        return ['PENDING APPROVAL', 'UNDER REVIEW', 'PENDING ADMIN REVIEW', 'PENDING REVIEW'].includes(invoice.status);
      }
      return invoice.status === statusFilter.value
    });
  }

  return invoicesToShow;
});

const summaryStats = computed(() => {
  const approved = invoices.value.filter(inv => inv.status === 'APPROVED' || inv.status === 'DISBURSED' || inv.status === 'PAID').length
  const funded = invoices.value
      .filter(inv => inv.status === 'DISBURSED' || inv.status === 'PAID')
      .reduce((sum, inv) => sum + (parseFloat(inv.fundedAmount || inv.totalAmount) || 0), 0)

  return {
    totalInvoices: invoices.value.length,
    approvedInvoices: approved,
    totalFunded: funded,
  }
})

const predefinedSteps = [
  { step: 1, title: 'Submitted', description: 'Invoice submitted successfully.', icon: BookUser },
  { step: 2, title: 'Under Review', description: 'Invoice is being reviewed.', icon: Clock },
  { step: 3, title: 'Approved', description: 'Your invoice has been approved.', icon: CircleCheckBig },
  { step: 4, title: 'Disbursed', description: 'Amount has been disbursed.', icon: CreditCard },
  { step: 5, title: 'Paid', description: 'Invoice has been paid by customer.', icon: HandCoins },
]

const invoiceStatusToStepMap: Record<string, number> = {
  'SUBMITTED': 1,
  'PENDING ADMIN REVIEW': 2,
  'PENDING REVIEW': 2,
  'PENDING APPROVAL': 2,
  'UNDER REVIEW': 2,
  'APPROVED': 3,
  'DISBURSED': 4,
  'PAID': 5,
  'REJECTED': 0,
  'CANCELLED': 0,
  'PROCESSING FAILED': 0,
  'UNKNOWN': 0,
}

const getStepForInvoiceStatus = (status: string): number => {
  const step = invoiceStatusToStepMap[status]; // status is already normalized by fetchInvoices
  if (typeof step === 'undefined') {
    console.warn(`Unknown status encountered for stepper: "${status}" for invoice. Defaulting to step 0.`);
    return 0;
  }
  return step;
}

const formatStatus = (status: string): string => {
  if (!status || status === 'UNKNOWN') return 'Unknown';
  // Status is already uppercase with spaces, so just convert to title case for display
  return status
      .toLowerCase()
      .replace(/\b\w/g, char => char.toUpperCase());
}

const getInvoiceStatusBadgeClass = (status: string): string => { // status is already normalized
  if (status === 'APPROVED' || status === 'DISBURSED' || status === 'PAID') return 'bg-green-100 text-green-800';
  if (['PENDING APPROVAL', 'UNDER REVIEW', 'SUBMITTED', 'PENDING ADMIN REVIEW', 'PENDING REVIEW'].includes(status)) return 'bg-blue-100 text-blue-800';
  if (['REJECTED', 'CANCELLED', 'PROCESSING FAILED'].includes(status)) return 'bg-red-100 text-red-800';
  return 'bg-gray-100 text-gray-800';
}

const getInvoiceStatusIcon = (status: string) => { // status is already normalized
  if (['APPROVED', 'DISBURSED', 'PAID'].includes(status)) return CircleCheck;
  if (['PENDING APPROVAL', 'UNDER REVIEW', 'SUBMITTED', 'PENDING ADMIN REVIEW', 'PENDING REVIEW'].includes(status)) return Clock;
  if (['REJECTED', 'CANCELLED', 'PROCESSING FAILED'].includes(status)) return AlertCircle;
  return FileText;
}

const formatDate = (dateString?: string | Date): string => {
  if (!dateString) return 'N/A';
  try {
    return new Date(dateString).toLocaleDateString('en-GB', {
      day: '2-digit', month: 'short', year: 'numeric'
    });
  } catch (e) {
    console.warn("Invalid date string for formatDate:", dateString);
    return 'Invalid Date';
  }
}

const kycStatusClass = (status: string | undefined) => {
  if (!status) return 'text-gray-500';
  const s = status.toLowerCase();
  if (s === 'approved' || s === 'verified') return 'text-green-500 font-semibold';
  if (s === 'pending' || s === 'submitted') return 'text-yellow-500 font-semibold';
  if (s === 'rejected' || s === 'not submitted' || s === 'incomplete') return 'text-red-500 font-semibold';
  return 'text-gray-500';
}

const dismissComplianceBanner = (event: MouseEvent) => {
  const banner = (event.currentTarget as HTMLElement)?.closest('.relative.bg-red-500');
  if (banner) {
    banner.remove();
  }
}

// Lifecycle Hooks
onMounted(() => {
  if (authToken) {
    fetchUser();
    fetchInvoices();
  } else {
    userError.value = "You are not logged in. Please log in to view your dashboard.";
    invoicesError.value = "Please log in to view invoices.";
    toast.info("Please log in to access all features.", {
      description: "Authentication token not found."
    })
  }
})

</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');

body {
  font-family: 'Inter', sans-serif;
}

.items-placeholder {
  min-height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px dashed #e0e0e0;
  border-radius: 8px;
  color: #a0a0a0;
}

</style>
