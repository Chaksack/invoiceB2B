<template>
  <main class="font-inter">
    <Toaster richColors position="top-right" />
    <section>
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
              <DialogHeader>
                <DialogTitle>Upload New Invoice</DialogTitle>
                <DialogDescription>
                  Upload your invoice file (PDF, CSV, XLSX formats supported).
                </DialogDescription>
              </DialogHeader>
              <div class="grid w-full max-w-sm items-center gap-1.5 py-4">
                <Label for="invoiceFile">Invoice File</Label>
                <Input id="invoiceFile" type="file" @change="handleFileSelect" accept=".pdf,.csv,.xlsx,.xls" class="rounded-md"/>
                <p v-if="fileUploadError" class="text-red-500 text-sm mt-1">{{ fileUploadError }}</p>
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
        <Dialog v-model:open="isKycDialogOpen"> <DialogTrigger as-child>
          <Button @click="isKycDialogOpen = true" class="text-white bg-black rounded-md hover:bg-gray-800 transition-colors">
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
              <Label for="kycFile">Business Registration</Label> <Input id="kycFile" type="file" @change="handleKycFileSelect" accept=".pdf,.csv,.xlsx,.xls" class="rounded-md"/> <p v-if="kycFileUploadError" class="text-red-500 text-sm mt-1">{{ kycFileUploadError }}</p>
            </div>
            <DialogFooter>
              <Button @click="submitKycFile" :disabled="isUploadingKyc || !selectedKycFile" class="rounded-md">
                <span v-if="isUploadingKyc">Uploading...</span>
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


    <!-- New Post-Upload Invoice Details & Providers Dialog -->
    <Dialog v-model:open="isPostUploadDetailsDialogOpen">
      <DialogContent class="sm:max-w-[1000px] shadow-lg rounded-lg">
        <DialogHeader>
          <DialogTitle>Invoice Submitted Successfully!</DialogTitle>
          <DialogDescription>
            Review the extracted details and choose a financial provider to fund your invoice.
          </DialogDescription>
        </DialogHeader>

        <div v-if="uploadedInvoiceDetails" class="grid md:grid-cols-2 gap-8 py-4 px-6">
          <!-- Left Column: Full Invoice Details -->
          <div class="left-column space-y-4">
            <h3 class="text-lg font-semibold border-b pb-2">Extracted Invoice Details</h3>
            <!-- Customer Info -->
            <div class="bg-gray-50 p-3 rounded-md">
              <p class="font-semibold text-black">{{ uploadedInvoiceDetails.customerName || 'N/A' }}</p>
              <p class="text-sm">{{ uploadedInvoiceDetails.customerAddress?.street || 'Street not available' }}</p>
              <p class="text-sm">{{ uploadedInvoiceDetails.customerAddress?.city || 'City not available' }}</p>
              <p class="text-sm">{{ uploadedInvoiceDetails.customerAddress?.zipCode }} {{ uploadedInvoiceDetails.customerAddress?.country }}</p>
            </div>

            <!-- Invoice Header Details -->
            <div class="grid grid-cols-2 gap-x-4 gap-y-1 text-sm">
              <div><span class="font-semibold">Invoice #:</span> {{ uploadedInvoiceDetails.invoiceNumber || 'N/A' }}</div>
              <div><span class="font-semibold">Invoice Date:</span> {{ formatDate(uploadedInvoiceDetails.invoiceDate) }}</div>
              <div><span class="font-semibold">Due Date:</span> {{ formatDate(uploadedInvoiceDetails.dueDate) }}</div>
              <div><span class="font-semibold">Terms:</span> {{ uploadedInvoiceDetails.paymentTerms || 'N/A' }}</div>
            </div>

            <!-- Items Table (if any) -->
            <Table v-if="uploadedInvoiceDetails.items && uploadedInvoiceDetails.items.length > 0" class="mt-4 w-full text-sm">
              <TableHeader>
                <TableRow>
                  <TableHead>ITEM & DESCRIPTION</TableHead>
                  <TableHead class="text-right">QTY</TableHead>
                  <TableHead class="text-right">AMOUNT</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="item in uploadedInvoiceDetails.items" :key="item.id">
                  <TableCell>
                    <div>{{ item.name }}</div>
                    <div class="text-xs font-light text-muted-foreground">{{ item.description }}</div>
                  </TableCell>
                  <TableCell class="text-right">{{ item.quantity }}</TableCell>
                  <TableCell class="text-right">{{ uploadedInvoiceDetails.currency || 'GHS' }} {{ item.amount?.toLocaleString() || '0.00' }}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
            <p v-else class="mt-4 text-sm text-muted-foreground">No items found for this invoice.</p>

            <!-- Totals -->
            <div class="bg-gray-100 p-3 rounded-md text-sm mt-4">
              <div class="flex justify-between items-center">
                <p>Sub Total:</p>
                <p class="font-semibold">{{ uploadedInvoiceDetails.currency || 'GHS' }} {{ uploadedInvoiceDetails.subTotalAmount?.toLocaleString() || '0.00' }}</p>
              </div>
              <div class="flex justify-between items-center">
                <p>Tax ({{ uploadedInvoiceDetails.taxRatePercentage || 0 }}%):</p>
                <p class="font-semibold">{{ uploadedInvoiceDetails.currency || 'GHS' }} {{ uploadedInvoiceDetails.taxAmount?.toLocaleString() || '0.00' }}</p>
              </div>
              <hr class="my-2 border-gray-300">
              <div class="flex justify-between items-center font-bold">
                <p>Total Amount:</p>
                <p>{{ uploadedInvoiceDetails.currency || 'GHS' }} {{ uploadedInvoiceDetails.totalAmount?.toLocaleString() || '0.00' }}</p>
              </div>
              <div class="flex justify-between items-center text-green-600 font-bold">
                <p>Balance Due:</p>
                <p>{{ uploadedInvoiceDetails.currency || 'GHS' }} {{ uploadedInvoiceDetails.balanceDue?.toLocaleString() || '0.00' }}</p>
              </div>
            </div>
          </div>

          <!-- Right Column: Financial Providers -->
          <div class="right-column space-y-4">
            <h3 class="text-lg font-semibold border-b pb-2">Available Financial Providers</h3>
            <div class="grid grid-cols-1 gap-4">
              <Card v-for="provider in financialProviders" :key="provider.id"
                    class="flex flex-col p-4 rounded-lg shadow-sm hover:shadow-md transition-shadow cursor-pointer"
                    :class="{ 'border-2 border-blue-500': selectedProvider?.id === provider.id }"
                    @click="selectedProvider = provider">
                <div class="flex items-center mb-2">
                  <img :src="provider.logo" :alt="provider.name" class="h-10 w-10 mr-4 rounded-full object-cover">
                  <div>
                    <p class="font-semibold text-lg">{{ provider.name }}</p>
                    <p class="text-sm text-muted-foreground">{{ provider.description }}</p>
                  </div>
                </div>
                <div class="text-sm border-t pt-2 mt-2">
                  <p class="font-medium">Terms:</p>
                  <ul class="list-disc list-inside text-muted-foreground">
                    <li v-for="(term, idx) in provider.terms" :key="idx">{{ term }}</li>
                  </ul>
                </div>
              </Card>
              <p v-if="financialProviders.length === 0" class="col-span-2 text-center text-muted-foreground">No financial providers available at this moment. Please check back later.</p>
            </div>
            <p v-if="providerSelectionError" class="text-red-500 text-sm mt-2">{{ providerSelectionError }}</p>
          </div>
        </div>
        <div v-else class="text-center py-4 text-muted-foreground">
          No invoice details to display.
        </div>

        <DialogFooter class="px-6 pb-4">
          <Button variant="outline" @click="isPostUploadDetailsDialogOpen = false" class="rounded-md">Close</Button>
          <Button @click="fundInvoice" :disabled="!selectedProvider || isFundingInvoice" class="bg-green-600 hover:bg-green-700 rounded-md">
            <span v-if="isFundingInvoice">Funding...</span>
            <span v-else>Fund Invoice Now</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue'
import axios from 'axios' // Keep axios import for now, but remove its use
import {
  Upload, Download, Eye, CreditCard, CircleCheckBig, FileWarning, FileText,
  CircleCheck, BookUser, HandCoins, Clock, AlertCircle, CircleX, Ban, Search, RefreshCw, Send
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

// API Configuration - Mocked or removed
// const API_BASE_URL = 'http://localhost:3000/api/v1'
// const tokenCookie = useCookie('token');
// const authToken = tokenCookie.value || null;

// const apiClient = axios.create({
//   baseURL: API_BASE_URL,
//   headers: {
//     ...(authToken ? { 'Authorization': `Bearer ${authToken}` } : {}),
//     'Content-Type': 'application/json'
//   }
// })

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

const selectedKycFile = ref<File | null>(null)
const isUploadingKyc = ref(false)
const kycFileUploadError = ref<string | null>(null)
const isKycDialogOpen = ref(false);

const currentReceipt = ref<any>(null)
const isLoadingReceipt = ref(false)
const receiptError = ref<string | null>(null)

const searchQuery = ref('');
const statusFilter = ref('all');

const isPostUploadDetailsDialogOpen = ref(false);
const uploadedInvoiceDetails = ref<any>(null);
const financialProviders = ref([
  { id: 'bank-a', name: 'ABC Bank', description: 'Fast approvals, low rates.', logo: 'https://placehold.co/40x40/FF5733/FFFFFF?text=AB',
    terms: ['Interest Rate: 8.5% p.a.', 'Loan Term: Up to 12 months', 'Fees: 1% origination fee', 'Eligibility: Min. 2 years in business'] },
  { id: 'fintech-x', name: 'SwiftCash Finance', description: 'Digital-first, instant funding.', logo: 'https://placehold.co/40x40/33FF57/FFFFFF?text=SC',
    terms: ['Interest Rate: 1.5% per month (flat)', 'Loan Term: 3-6 months', 'Fees: No hidden fees', 'Eligibility: Min. 6 months in business, online application'] },
  { id: 'credit-union-y', name: 'Community Credit Union', description: 'Personalized service, fair terms.', logo: 'https://placehold.co/40x40/3357FF/FFFFFF?text=CCU',
    terms: ['Interest Rate: 7.9% p.a.', 'Loan Term: 6-18 months', 'Fees: 0.5% processing fee', 'Eligibility: Member of credit union for 1+ year'] },
]);
const selectedProvider = ref<any>(null);
const isFundingInvoice = ref(false);
const providerSelectionError = ref<string | null>(null);

// New state for selected invoice in the main list
const selectedInvoiceForDetail = ref<any>(null);


const statusOptions = [
  { value: 'all', label: 'All Statuses' },
  { value: 'SUBMITTED', label: 'Submitted' },
  { value: 'PENDING REVIEW', label: 'Under Review' },
  { value: 'APPROVED', label: 'Approved' },
  { value: 'DISBURSED', label: 'Disbursed' },
  { value: 'PAID', label: 'Paid' },
  { value: 'REJECTED', label: 'Rejected' },
];


// --- Mock Data ---
const mockUser = {
  firstName: "John",
  lastName: "Doe",
  companyName: "Innovate Solutions Inc.",
  kycStatus: "approved" // or "pending", "incomplete" to test banner
};

const mockInvoices = [
  {
    id: "INV001",
    customerName: "Acme Corp",
    invoiceNumber: "ACME-2023-001",
    status: "APPROVED",
    totalAmount: 1500.75,
    subTotalAmount: 1400.00,
    taxAmount: 100.75,
    balanceDue: 1500.75,
    taxRatePercentage: 7.19,
    currency: "GHS",
    invoiceDate: "2023-01-15T10:00:00Z",
    dueDate: "2023-02-15T10:00:00Z",
    paymentTerms: "Net 30",
    customerAddress: { street: "123 Main St", city: "Accra", zipCode: "00233", country: "Ghana" },
    items: [
      { id: 1, name: "Product A", description: "High-quality widget", unitPrice: 500, quantity: 2, amount: 1000 },
      { id: 2, name: "Service B", description: "Consulting hours", unitPrice: 200, quantity: 2, amount: 400 },
    ],
    // Added mock bank financing details for testing
    bankFinancingDetails: { bankName: "ABC Bank", accountNumber: "9876543210", fundedAmount: 1500.75, purpose: "Invoice Funding for ACME-2023-001" }
  },
  {
    id: "INV002",
    customerName: "Globex Corporation",
    invoiceNumber: "GLBX-2023-002",
    status: "PENDING REVIEW",
    totalAmount: 2500.00,
    subTotalAmount: 2300.00,
    taxAmount: 200.00,
    balanceDue: 2500.00,
    taxRatePercentage: 8.7,
    currency: "GHS",
    invoiceDate: "2023-01-20T11:30:00Z",
    dueDate: "2023-02-20T11:30:00Z",
    paymentTerms: "Net 30",
    customerAddress: { street: "456 Market St", city: "Kumasi", zipCode: "00233", country: "Ghana" },
    items: [
      { id: 3, name: "Product C", description: "Advanced gadget", unitPrice: 1150, quantity: 2, amount: 2300 },
    ],
  },
  {
    id: "INV003",
    customerName: "Stark Industries",
    invoiceNumber: "STARK-2023-003",
    status: "DISBURSED",
    totalAmount: 5000.00,
    subTotalAmount: 4800.00,
    taxAmount: 200.00,
    balanceDue: 0.00, // Funded
    taxRatePercentage: 4.16,
    currency: "GHS",
    invoiceDate: "2023-01-05T09:00:00Z",
    dueDate: "2023-02-05T09:00:00Z",
    paymentTerms: "Net 30",
    customerAddress: { street: "789 Tech Rd", city: "Tema", zipCode: "00233", country: "Ghana" },
    items: [
      { id: 4, name: "Tech Solution", description: "Enterprise software license", unitPrice: 4800, quantity: 1, amount: 4800 },
    ],
    fundedAmount: 5000.00, // Example of funded amount for disbursed
    bankFinancingDetails: { bankName: "SwiftCash Finance", accountNumber: "5555444433", fundedAmount: 5000.00, purpose: "Invoice Funding for STARK-2023-003" }
  },
  {
    id: "INV004",
    customerName: "Wayne Enterprises",
    invoiceNumber: "WAYNE-2023-004",
    status: "REJECTED",
    totalAmount: 750.00,
    subTotalAmount: 700.00,
    taxAmount: 50.00,
    balanceDue: 750.00,
    taxRatePercentage: 7.14,
    currency: "GHS",
    invoiceDate: "2023-02-01T14:00:00Z",
    dueDate: "2023-03-01T14:00:00Z",
    paymentTerms: "Net 30",
    customerAddress: { street: "101 Gotham Blvd", city: "Cape Coast", zipCode: "00233", country: "Ghana" },
    items: [
      { id: 5, name: "Security Audit", description: "Annual security assessment", unitPrice: 700, quantity: 1, amount: 700 },
    ],
  },
  {
    id: "INV005",
    customerName: "Cyberdyne Systems",
    invoiceNumber: "CYB-2023-005",
    status: "PAID",
    totalAmount: 12000.00,
    subTotalAmount: 10000.00,
    taxAmount: 2000.00,
    balanceDue: 0.00,
    taxRatePercentage: 20.00,
    currency: "GHS",
    invoiceDate: "2023-01-25T16:00:00Z",
    dueDate: "2023-02-25T16:00:00Z",
    paymentTerms: "Net 30",
    customerAddress: { street: "12 Terminator Ave", city: "Takoradi", zipCode: "00233", country: "Ghana" },
    items: [
      { id: 6, name: "AI Integration", description: "Advanced AI system integration", unitPrice: 10000, quantity: 1, amount: 10000 },
    ],
    fundedAmount: 12000.00 // Example of funded amount for paid
  }
];

const mockReceipt = {
  referenceNumber: "TXN123456789",
  transferTo: "Innovate Solutions Inc.",
  accountType: "Bank Account",
  accountNumber: "1234567890",
  accountName: "John Doe",
  amount: 5000.00,
  transferDate: "2023-02-06T10:30:00Z",
  purpose: "Invoice Funding for GLBX-2023-002"
};


// --- User Profile ---
const fetchUser = async () => {
  isLoadingUser.value = true
  userError.value = null
  try {
    await new Promise(resolve => setTimeout(resolve, 500)); // Simulate API delay
    user.value = mockUser;
  } catch (err) {
    userError.value = 'An error occurred while fetching mock user details.'
  } finally {
    isLoadingUser.value = false
  }
}

// --- Invoices ---
const fetchInvoices = async () => {
  isLoadingInvoices.value = true;
  invoicesError.value = null;
  try {
    await new Promise(resolve => setTimeout(resolve, 800)); // Simulate API delay
    invoices.value = mockInvoices.map(inv => ({
      ...inv,
      status: String(inv.status).trim().replace(/_/g, ' ').toUpperCase() // Ensure consistent status
    })).sort((a: any, b: any) => {
      const dateA = a.invoiceDate ? new Date(a.invoiceDate).getTime() : 0;
      const dateB = b.invoiceDate ? new Date(b.invoiceDate).getTime() : 0;
      return dateB - dateA;
    });
    if (invoices.value.length > 0) {
      selectedInvoiceForDetail.value = invoices.value[0]; // Select the first invoice by default
    }
  } catch (err: any) {
    invoicesError.value = `An unexpected error occurred while fetching mock invoices: ${err.message || ''}`;
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
  if (!selectedFile.value) {
    fileUploadError.value = "Please select a file to upload.";
    return;
  }

  console.log('Simulating invoice submission with file:', selectedFile.value.name);

  isUploadingInvoice.value = true;
  fileUploadError.value = null;

  try {
    await new Promise(resolve => setTimeout(resolve, 1500)); // Simulate upload delay

    // Generate a mock invoice response based on the uploaded file
    const newInvoice = {
      id: `NEWINV-${Date.now()}`,
      customerName: `Customer ${Date.now().toString().slice(-4)}`,
      invoiceNumber: `INV-${Date.now().toString().slice(-6)}`,
      status: "SUBMITTED",
      totalAmount: Math.floor(Math.random() * 10000) + 500,
      subTotalAmount: 0, // Will be calculated below
      taxAmount: 0,     // Will be calculated below
      balanceDue: 0,    // Will be calculated below
      taxRatePercentage: 10 + Math.random() * 5, // Random tax rate
      currency: "GHS",
      invoiceDate: new Date().toISOString(),
      dueDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(), // 30 days from now
      paymentTerms: "Net 30",
      customerAddress: { street: "Mock Street", city: "Mock City", zipCode: "00000", country: "Ghana" },
      items: [
        { id: 1, name: "Service Package", description: "Standard service", unitPrice: Math.floor(Math.random() * 500) + 50, quantity: Math.floor(Math.random() * 5) + 1, amount: 0 },
        { id: 2, name: "Consulting Fee", description: "Expert consultation", unitPrice: Math.floor(Math.random() * 300) + 30, quantity: Math.floor(Math.random() * 10) + 1, amount: 0 },
      ],
    };

    // Calculate item amounts and totals
    newInvoice.items.forEach(item => {
      item.amount = item.unitPrice * item.quantity;
      newInvoice.subTotalAmount += item.amount;
    });

    newInvoice.taxAmount = newInvoice.subTotalAmount * (newInvoice.taxRatePercentage / 100);
    newInvoice.totalAmount = newInvoice.subTotalAmount + newInvoice.taxAmount;
    newInvoice.balanceDue = newInvoice.totalAmount;

    uploadedInvoiceDetails.value = newInvoice;
    toast.success('Invoice uploaded successfully!');

    // Add the new invoice to the mock data and refresh the list
    mockInvoices.unshift(uploadedInvoiceDetails.value);
    await fetchInvoices();

    isUploadDialogOpen.value = false;
    isPostUploadDetailsDialogOpen.value = true;
    selectedFile.value = null;
  } catch (err: any) {
    fileUploadError.value = `Simulated upload error: ${err.message || 'Failed to upload invoice'}`;
    toast.error('Simulated upload failed.');
  } finally {
    isUploadingInvoice.value = false;
  }
};

// --- KYC Upload Logic ---
const handleKycFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    selectedKycFile.value = target.files[0];
    console.log('KYC file selected:', selectedKycFile.value);
    kycFileUploadError.value = null;
  } else {
    selectedKycFile.value = null;
  }
};

const submitKycFile = async () => {
  if (!selectedKycFile.value) {
    kycFileUploadError.value = "Please select a KYC file to upload.";
    return;
  }

  isUploadingKyc.value = true;
  kycFileUploadError.value = null;

  try {
    await new Promise(resolve => setTimeout(resolve, 1500)); // Simulate upload delay
    mockUser.kycStatus = 'pending'; // Simulate status change
    toast.success('KYC file uploaded successfully! Your status will be updated after review.');
    await fetchUser(); // Refresh user status
    selectedKycFile.value = null;
    isKycDialogOpen.value = false;
  } catch (err: any) {
    kycFileUploadError.value = `Simulated upload error: ${err.message || 'Failed to upload KYC file'}`;
    toast.error('Simulated KYC upload failed.');
  } finally {
    isUploadingKyc.value = false;
  }
};


// --- Receipt Logic ---
const fetchReceiptDetails = async (invoiceId: string | number) => {
  isLoadingReceipt.value = true
  receiptError.value = null
  currentReceipt.value = null
  try {
    await new Promise(resolve => setTimeout(resolve, 700)); // Simulate API delay
    currentReceipt.value = mockReceipt; // Use mock receipt
  } catch (err: any) {
    receiptError.value = 'An unexpected error occurred while fetching mock receipt details.'
  } finally {
    isLoadingReceipt.value = false
  }
}

const triggerDownloadReceipt = async (invoiceId: string | number) => {
  toast.info("Preparing mock download...");
  try {
    await new Promise(resolve => setTimeout(resolve, 1000)); // Simulate download delay
    // In a real app, you'd trigger a file download from a Blob or URL
    console.log(`Simulating download for invoice ID: ${invoiceId}`);
    toast.success("Simulated download started for a mock PDF.");
  } catch (err: any) {
    toast.error('An unexpected error occurred during mock download.')
  }
}

// --- Fund Invoice (Mock) ---
const fundInvoice = async () => {
  if (!selectedProvider.value) {
    providerSelectionError.value = "Please select a financial provider.";
    return;
  }
  if (!uploadedInvoiceDetails.value || !uploadedInvoiceDetails.value.id) {
    toast.error("No invoice selected for funding. Please re-upload your invoice.");
    isPostUploadDetailsDialogOpen.value = false;
    return;
  }

  isFundingInvoice.value = true;
  providerSelectionError.value = null;
  toast.info(`Requesting funding from ${selectedProvider.value.name}...`);

  try {
    await new Promise(resolve => setTimeout(resolve, 2000)); // Simulate 2-second funding process

    // Update the status of the "uploadedInvoiceDetails" to "DISBURSED" in mockInvoices
    const invoiceIndex = mockInvoices.findIndex(inv => inv.id === uploadedInvoiceDetails.value.id);
    if (invoiceIndex !== -1) {
      mockInvoices[invoiceIndex].status = 'DISBURSED';
      mockInvoices[invoiceIndex].balanceDue = 0; // Assuming it's fully funded
      mockInvoices[invoiceIndex].fundedAmount = uploadedInvoiceDetails.value.totalAmount;
      // Add mock bank financing details to the invoice that was just funded
      mockInvoices[invoiceIndex].bankFinancingDetails = {
        bankName: selectedProvider.value.name,
        accountNumber: "XXXXXXXXXX", // Mock account number
        fundedAmount: uploadedInvoiceDetails.value.totalAmount,
        purpose: `Funding for Invoice ${uploadedInvoiceDetails.value.invoiceNumber}`
      };
    }

    console.log('Simulated Funding response: success');
    toast.success(`Funding request successful! Check your invoice status later.`);
    isPostUploadDetailsDialogOpen.value = false; // Close the dialog
    selectedProvider.value = null; // Clear selected provider
    await fetchInvoices(); // Refresh invoices to show potential status change
  } catch (error: any) {
    console.error('Error funding invoice:', error);
    toast.error('Failed to submit funding request. Please try again.');
    providerSelectionError.value = 'Failed to submit funding request.';
  } finally {
    isFundingInvoice.value = false;
  }
};


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

const selectInvoiceForDetail = (invoice: any) => {
  selectedInvoiceForDetail.value = invoice;
}

// Mock function for "Send Invoice"
const mockSendInvoice = () => {
  if (selectedInvoiceForDetail.value) {
    toast.info(`Simulating sending invoice #${selectedInvoiceForDetail.value.invoiceNumber}.`);
    // Here you would integrate actual API call to send invoice
  } else {
    toast.error("No invoice selected to send.");
  }
}

// Lifecycle Hooks
onMounted(() => {
  fetchUser();
  fetchInvoices();
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
  color: #a0e0a0;
}
</style>
