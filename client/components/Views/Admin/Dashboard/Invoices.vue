<template>
  <main class="bg-gray-100 font-inter min-h-screen">
    <Toaster richColors position="top-right" />

    <div class="container mx-auto p-4 md:p-6 lg:p-8">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-semibold text-gray-800">Invoice</h1>
      </div>

      <div class="md:flex md:space-x-6">
        <div class="md:w-1/3 bg-white p-6 rounded-lg shadow-md mb-6 md:mb-0 flex flex-col">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-700">All Invoices</h2>
          </div>

          <div class="mb-4 space-y-3">
            <div>
              <Input
                  type="text"
                  v-model="searchTerm"
                  placeholder="Search invoices..."
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 text-sm"
              />
            </div>
            <div>
              <Label class="text-xs text-gray-600">Filter by Status:</Label>
              <div class="flex flex-wrap gap-2 mt-1">
                <Button
                    v-for="status in availableStatuses"
                    :key="status.value"
                    @click="toggleStatusFilter(status.value)"
                    :variant="selectedStatusFilters.includes(status.value) ? 'default' : 'outline'"
                    class="text-xs px-2.5 py-1 rounded-full"
                    :class="{
                    'bg-blue-600 text-white hover:bg-blue-700': selectedStatusFilters.includes(status.value),
                    'border-gray-300 text-gray-700 hover:bg-gray-100': !selectedStatusFilters.includes(status.value)
                  }"
                >
                  {{ status.label }}
                </Button>
                <Button
                    v-if="selectedStatusFilters.length > 0"
                    @click="clearStatusFilters"
                    variant="ghost"
                    class="text-xs px-2.5 py-1 rounded-full text-red-500 hover:bg-red-50"
                >
                  Clear
                </Button>
              </div>
            </div>
          </div>

          <div class="flex-grow overflow-hidden">
            <div v-if="isLoadingInvoices" class="text-center py-10 text-gray-500">
              <p>Loading invoices...</p>
            </div>
            <div v-else-if="invoicesError" class="text-center py-10 text-red-500">
              <p>{{ invoicesError }}</p>
            </div>
            <div v-else-if="paginatedInvoices.length === 0 && invoices.length > 0" class="text-center py-6 text-gray-500">
              <p class="text-md font-semibold">No invoices match your criteria.</p>
              <p class="text-sm">Try adjusting your search or filters.</p>
            </div>
            <div v-else-if="paginatedInvoices.length === 0" class="text-center py-10 text-gray-500 flex flex-col items-center">
              <Ban class="w-12 h-12 mb-4 text-gray-400" />
              <p class="text-xl font-semibold">No invoices found.</p>
              <p class="text-sm">Upload your first invoice to get started!</p>
            </div>
            <div v-else class="space-y-1 max-h-[calc(60vh-100px)] overflow-y-auto pr-1">
              <div
                  v-for="invoice in paginatedInvoices"
                  :key="invoice.id"
                  @click="selectInvoice(invoice)"
                  :class="[
                  'p-3 border-l-4 rounded-r-md cursor-pointer hover:bg-gray-50 transition-all',
                  selectedInvoiceId === invoice.id ? 'border-blue-500 bg-blue-50' : 'border-transparent'
                ]"
              >
                <div class="flex items-center">
                  <input type="checkbox" class="form-checkbox h-4 w-4 text-blue-600 mr-3 rounded border-gray-300 focus:ring-blue-500" @click.stop />
                  <div class="flex-grow">
                    <div class="flex justify-between items-start">
                      <span class="font-semibold text-gray-800 text-sm">{{ invoice.customerName || 'N/A' }}</span>
                      <span class="font-semibold text-gray-800 text-sm">{{ invoice.currencySymbol || '$' }}{{ (invoice.totalAmount || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</span>
                    </div>
                    <div class="flex justify-between items-center mt-1">
                      <span class="text-xs text-gray-500">
                        {{ invoice.invoiceNumber || invoice.id }} &bull; {{ formatDate(invoice.invoiceDate, 'dd/MM/yyyy') }}
                      </span>
                      <Badge :class="getListInvoiceStatusBadgeClass(invoice.statusForBadge)" class="py-0.5 px-2 rounded-full text-xs font-medium">
                        {{ formatListStatus(invoice.statusForBadge) }}
                      </Badge>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div v-if="totalPages > 1" class="mt-auto pt-4 flex justify-center items-center space-x-2">
            <Button @click="prevPage" :disabled="currentPage === 1" variant="outline" size="sm" class="text-xs px-2.5 py-1">Previous</Button>
            <span class="text-xs text-gray-600">Page {{ currentPage }} of {{ totalPages }}</span>
            <Button @click="nextPage" :disabled="currentPage === totalPages" variant="outline" size="sm" class="text-xs px-2.5 py-1">Next</Button>
          </div>
        </div>

        <div class="md:w-2/3 bg-white p-6 rounded-lg shadow-md">
          <div v-if="isLoadingInvoices && invoices.length === 0" class="text-center py-10 text-gray-500">
            <p>Loading details...</p>
          </div>
          <div v-else-if="!computedSelectedInvoice && paginatedInvoices.length > 0" class="text-center py-10 text-gray-500">
            <p>Select an invoice to view details.</p>
          </div>
          <div v-else-if="!computedSelectedInvoice && invoices.length > 0 && paginatedInvoices.length === 0" class="text-center py-10 text-gray-500">
            <p>No invoices match your current filters. Clear filters or select from available invoices if any.</p>
          </div>
          <div v-else-if="!computedSelectedInvoice && invoices.length === 0" class="text-center py-10 text-gray-500">
            <p>No invoices available to display.</p>
          </div>
          <div v-else-if="computedSelectedInvoice">
            <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
              <h2 class="text-2xl font-bold text-gray-800">{{ computedSelectedInvoice.invoiceNumber || computedSelectedInvoice.id }}</h2>
              <div class="flex space-x-2">

                <DropdownMenu>
                  <DropdownMenuTrigger as-child>
                    <Button variant="outline" class="text-sm px-3 py-1.5 rounded-md bg-blue-500 text-white border border-gray-300 hover:bg-ble-300">
                      Actions <ChevronDown class="ml-1 h-4 w-4 inline-block" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent class="w-40 bg-white shadow-lg rounded-md border border-gray-200">
                    <DropdownMenuItem @click="updateInvoiceStatus(computedSelectedInvoice.id, 'APPROVED')" class="text-sm px-3 py-2 hover:bg-gray-100 cursor-pointer">Approve</DropdownMenuItem>
                    <DropdownMenuItem @click="updateInvoiceStatus(computedSelectedInvoice.id, 'DISBURSED')" class="text-sm px-3 py-2 hover:bg-gray-100 cursor-pointer">Mark Disbursed</DropdownMenuItem>
                    <DropdownMenuItem @click="updateInvoiceStatus(computedSelectedInvoice.id, 'PAID')" class="text-sm px-3 py-2 hover:bg-gray-100 cursor-pointer">Mark Paid</DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem @click="updateInvoiceStatus(computedSelectedInvoice.id, 'REJECTED')" class="text-sm px-3 py-2 hover:bg-red-50 text-red-600 cursor-pointer">Reject</DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>

                <Button @click="isUploadReceiptDialogOpen = true" variant="outline" class="text-sm px-3 py-1.5 rounded-md border-gray-300 hover:bg-gray-100">
                  <UploadCloud class="mr-1 h-4 w-4 inline" /> Upload Receipt
                </Button>
                <Button @click="triggerDownloadPDF(computedSelectedInvoice.id, 'invoice')" class="text-sm px-3 py-1.5 rounded-md border border-gray-300 bg-white hover:bg-gray-100 text-blue-600 border-blue-500">
                  <Download class="mr-1 h-4 w-4 inline" /> Download PDF
                </Button>
              </div>
            </div>

            <div class="w-full mb-8">
              <div class="flex items-start">
                <template v-for="(step, index) in predefinedSteps" :key="step.step">
                  <div class="flex flex-col items-center w-1/5">
                    <div
                        class="w-10 h-10 rounded-full flex items-center justify-center"
                        :class="{
                          'bg-blue-600 text-white': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) === step.step,
                          'bg-green-500 text-white': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) > step.step,
                          'bg-gray-300 text-gray-600': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) < step.step
                        }"
                    >
                      <component :is="step.icon" class="w-5 h-5" />
                    </div>
                    <p
                        class="mt-2 text-xs text-center font-medium"
                        :class="{
                          'text-blue-600': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) === step.step,
                          'text-green-500': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) > step.step,
                          'text-gray-500': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) < step.step
                        }"
                    >
                      {{ step.title }}
                    </p>
                    <p class="mt-1 text-xs text-center text-gray-400 px-1" style="min-height: 2.5em;"> {{ step.description }}
                    </p>
                  </div>
                  <div
                      v-if="index < predefinedSteps.length - 1"
                      class="flex-1 h-1 mt-5"
                      :class="{
                        'bg-green-500': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) > step.step,
                        'bg-gray-300': getStepForInvoiceStatus(computedSelectedInvoice.originalStatus) <= step.step
                      }"
                  ></div>
                </template>
              </div>
            </div>

            <div class="grid md:grid-cols-2 gap-6 mb-6">
              <div>
                <div :class="getDetailStatusClass(computedSelectedInvoice.statusForBadge)" class="text-xs font-semibold px-2 py-1 inline-block rounded mb-4">
                  STATUS : {{ computedSelectedInvoice.statusForBadge ? computedSelectedInvoice.statusForBadge.toUpperCase() : 'N/A' }}
                </div>
                <h3 class="text-xs text-gray-500 mb-0.5">Invoice: {{ computedSelectedInvoice.invoiceNumber || computedSelectedInvoice.id }}</h3>
                <p class="text-xs text-gray-500 mb-1">Issue Date: {{ formatDate(computedSelectedInvoice.invoiceDate) }}</p>
                <p class="text-xs text-gray-500">Due Date: {{ formatDate(computedSelectedInvoice.dueDate) }}</p>

                <div class="mt-4 pt-4 border-t border-gray-100">
                  <p class="text-xs text-gray-500 mb-1">CLIENT</p>
                  <h4 class="font-semibold text-gray-700">{{ computedSelectedInvoice.customerName || 'N/A' }}</h4>
                  <p class="text-xs text-gray-600">{{ computedSelectedInvoice.customerAddress?.street || 'Street not available' }}</p>
                  <p class="text-xs text-gray-600">{{ computedSelectedInvoice.customerAddress?.city }}{{ computedSelectedInvoice.customerAddress?.zipCode ? ', ' + computedSelectedInvoice.customerAddress.zipCode : '' }}</p>
                  <p class="text-xs text-gray-600">{{ computedSelectedInvoice.customerAddress?.country }}</p>
                </div>
              </div>
              <div class="text-right">
                <div v-if="isLoadingUserProfile" class="text-xs text-gray-500">Loading sender info...</div>
                <div v-else-if="userProfileError" class="text-xs text-red-500">{{ userProfileError }}</div>
                <div v-else-if="userProfile.value">
                  <img v-if="userProfile.value.companyLogoUrl" :src="userProfile.value.companyLogoUrl" alt="Sender Company Logo" class="h-10 mb-2 inline-block" onerror="this.style.display='none'">
                  <h3 class="font-semibold text-gray-700 text-md">{{ userProfile.value.companyName || userProfile.value.name || 'Your Company' }}</h3>
                  <p class="text-xs text-gray-600">{{ userProfile.value.companyAddress?.street || 'Street not set' }}</p>
                  <p class="text-xs text-gray-600">
                    {{ userProfile.value.companyAddress?.city || 'City not set' }}{{ userProfile.value.companyAddress?.state ? ', ' + userProfile.value.companyAddress.state : '' }} {{ userProfile.value.companyAddress?.zipCode || '' }}
                  </p>
                  <p class="text-xs text-gray-600">{{ userProfile.value.companyAddress?.country || 'Country not set' }}</p>
                </div>
                <div v-else class="text-xs text-gray-500">Sender information not available.</div>


                <div class="mt-4 pt-4 border-t border-gray-100 text-right">
                  <p class="text-xs text-gray-500">TOTAL AMOUNT</p>
                  <p class="text-2xl font-bold text-gray-800 mb-1">{{ computedSelectedInvoice.currencySymbol || '$' }}{{ (computedSelectedInvoice.totalAmount || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</p>
                  <p class="text-xs text-gray-500">BALANCE DUE</p>
                  <p class="text-lg font-semibold text-blue-600 mb-1">{{ computedSelectedInvoice.currencySymbol || '$' }}{{ (computedSelectedInvoice.balanceDue || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</p>
                  <p class="text-xs text-gray-500">AMOUNT PAID</p>
                  <p class="text-sm text-gray-700">{{ computedSelectedInvoice.currencySymbol || '$' }}{{ (computedSelectedInvoice.amountPaid || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</p>
                </div>
              </div>
            </div>

            <h3 class="text-md font-semibold text-gray-700 mb-3">Items</h3>
            <Table v-if="computedSelectedInvoice.items && computedSelectedInvoice.items.length > 0" class="w-full text-sm">
              <TableHeader class="bg-gray-50">
                <TableRow>
                  <TableHead class="text-left text-gray-600 font-medium py-2 px-3">#</TableHead>
                  <TableHead class="text-left text-gray-600 font-medium py-2 px-3">ITEM & DESCRIPTION</TableHead>
                  <TableHead class="text-right text-gray-600 font-medium py-2 px-3">QTY.</TableHead>
                  <TableHead class="text-right text-gray-600 font-medium py-2 px-3">RATE ({{ computedSelectedInvoice.currencySymbol || '$' }})</TableHead>
                  <TableHead class="text-right text-gray-600 font-medium py-2 px-3">AMOUNT ({{ computedSelectedInvoice.currencySymbol || '$' }})</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="(item, index) in computedSelectedInvoice.items" :key="item.id || index" class="border-b border-gray-100">
                  <TableCell class="py-2 px-3 text-gray-500">{{ index + 1 }}</TableCell>
                  <TableCell class="py-2 px-3">
                    <div class="text-gray-800 font-medium">{{ item.name }}</div>
                    <div class="text-xs text-gray-500">{{ item.description }}</div>
                  </TableCell>
                  <TableCell class="text-right py-2 px-3 text-gray-700">{{ item.quantity }}</TableCell>
                  <TableCell class="text-right py-2 px-3 text-gray-700">{{ item.unitPrice?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00' }}</TableCell>
                  <TableCell class="text-right py-2 px-3 text-gray-700 font-medium">{{ item.amount?.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00' }}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
            <p v-else class="mt-4 text-gray-500 text-sm">No items listed for this invoice.</p>

            <div class="mt-8 flex justify-end">
              <div class="w-full max-w-xs space-y-1 text-sm">
                <div class="flex justify-between">
                  <span class="text-gray-600">Sub Total:</span>
                  <span class="text-gray-800 font-medium">{{ computedSelectedInvoice.currencySymbol || '$' }}{{ (computedSelectedInvoice.subTotalAmount || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</span>
                </div>
                <div v-if="computedSelectedInvoice.taxAmount > 0" class="flex justify-between">
                  <span class="text-gray-600">Tax ({{ computedSelectedInvoice.taxRatePercentage || 0 }}%):</span>
                  <span class="text-gray-800 font-medium">{{ computedSelectedInvoice.currencySymbol || '$' }}{{ (computedSelectedInvoice.taxAmount || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</span>
                </div>
                <div class="flex justify-between border-t border-gray-200 pt-1 mt-1">
                  <span class="text-gray-800 font-semibold">Total:</span>
                  <span class="text-gray-800 font-semibold">{{ computedSelectedInvoice.currencySymbol || '$' }}{{ (computedSelectedInvoice.totalAmount || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</span>
                </div>
                <div class="flex justify-between text-blue-600">
                  <span class="font-semibold">Balance Due:</span>
                  <span class="font-semibold">{{ computedSelectedInvoice.currencySymbol || '$' }}{{ (computedSelectedInvoice.balanceDue || 0).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</span>
                </div>
              </div>
            </div>

          </div>
          <div v-else-if="!isLoadingInvoices && invoices.length > 0 && !computedSelectedInvoice" class="text-center py-10 text-gray-500">
            <p>Select an invoice from the list to see the details.</p>
          </div>
        </div>
      </div>
    </div>

    <Dialog :open="isUploadDialogOpen" @update:open="isUploadDialogOpen = $event">
      <DialogContent class="sm:max-w-[425px] bg-white rounded-lg shadow-xl">
        <DialogHeader class="p-6 border-b border-gray-200">
          <DialogTitle class="text-lg font-semibold text-gray-800">Upload New Invoice</DialogTitle>
          <DialogDescription class="text-sm text-gray-500 mt-1">
            Select a PDF or image file for your invoice.
          </DialogDescription>
        </DialogHeader>
        <div class="p-6 space-y-4">
          <div>
            <Label for="invoiceFile" class="text-sm font-medium text-gray-700">Invoice File</Label>
            <Input
                id="invoiceFile"
                type="file"
                @change="handleFileSelect"
                class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                accept=".pdf,.jpg,.jpeg,.png"
            />
            <p v-if="selectedFile" class="mt-1 text-xs text-gray-500">Selected: {{ selectedFile.name }}</p>
            <p v-if="fileUploadError" class="mt-1 text-xs text-red-500">{{ fileUploadError }}</p>
          </div>
        </div>
        <DialogFooter class="p-6 border-t border-gray-200 flex justify-end space-x-2">
          <Button variant="outline" @click="isUploadDialogOpen = false" class="px-4 py-2 text-sm rounded-md border-gray-300 hover:bg-gray-50">Cancel</Button>
          <Button @click="submitInvoice" :disabled="isUploadingInvoice || !selectedFile" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm rounded-md disabled:opacity-50">
            {{ isUploadingInvoice ? 'Uploading...' : 'Upload Invoice' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <Dialog :open="isUploadReceiptDialogOpen" @update:open="isUploadReceiptDialogOpen = $event">
      <DialogContent class="sm:max-w-[425px] bg-white rounded-lg shadow-xl">
        <DialogHeader class="p-6 border-b border-gray-200">
          <DialogTitle class="text-lg font-semibold text-gray-800">Upload Receipt</DialogTitle>
          <DialogDescription class="text-sm text-gray-500 mt-1">
            Select a PDF or image file for the receipt. (Invoice: {{ computedSelectedInvoice?.invoiceNumber || computedSelectedInvoice?.id }})
          </DialogDescription>
        </DialogHeader>
        <div class="p-6 space-y-4">
          <div>
            <Label for="receiptFile" class="text-sm font-medium text-gray-700">Receipt File</Label>
            <Input
                id="receiptFile"
                type="file"
                @change="handleReceiptFileSelect"
                class="mt-1 block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                accept=".pdf,.jpg,.jpeg,.png"
            />
            <p v-if="selectedReceiptFile" class="mt-1 text-xs text-gray-500">Selected: {{ selectedReceiptFile.name }}</p>
            <p v-if="receiptUploadError" class="mt-1 text-xs text-red-500">{{ receiptUploadError }}</p>
          </div>
        </div>
        <DialogFooter class="p-6 border-t border-gray-200 flex justify-end space-x-2">
          <Button variant="outline" @click="isUploadReceiptDialogOpen = false" class="px-4 py-2 text-sm rounded-md border-gray-300 hover:bg-gray-50">Cancel</Button>
          <Button @click="submitReceipt" :disabled="isUploadingReceipt || !selectedReceiptFile || !computedSelectedInvoice" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm rounded-md disabled:opacity-50">
            {{ isUploadingReceipt ? 'Uploading...' : 'Upload Receipt' }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h, watch } from 'vue' // Added watch
import axios from 'axios'
import {
  Upload, Download, Ban,
  CreditCard, CircleCheckBig, BookUser, HandCoins, Clock, ChevronDown, UploadCloud
} from 'lucide-vue-next'
import { useCookie } from '#app';
import { Toaster, toast } from 'vue-sonner'

// Shadcn-vue components
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
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
const invoices = ref<any[]>([])
const isLoadingInvoices = ref(false)
const invoicesError = ref<string | null>(null)

const selectedInvoiceId = ref<string | number | null>(null);

// User Profile State
const userProfile = ref<any>(null);
const isLoadingUserProfile = ref(false);
const userProfileError = ref<string | null>(null);


// Invoice Upload State
const selectedFile = ref<File | null>(null)
const isUploadingInvoice = ref(false)
const fileUploadError = ref<string | null>(null)
const isUploadDialogOpen = ref(false);

// Receipt Upload State
const selectedReceiptFile = ref<File | null>(null);
const isUploadingReceipt = ref(false);
const receiptUploadError = ref<string | null>(null);
const isUploadReceiptDialogOpen = ref(false);


const searchTerm = ref('');
const selectedStatusFilters = ref<string[]>([]);

const availableStatuses = ref([
  { label: 'Pending', value: 'PENDING' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Disbursed', value: 'DISBURSED' },
  { label: 'Paid', value: 'PAID' },
  { label: 'Rejected', value: 'REJECTED' },
]);

// Pagination State
const currentPage = ref(1);
const itemsPerPage = ref(10);

const predefinedSteps = [
  { step: 1, title: 'Submitted', description: 'Invoice submitted for processing.', icon: BookUser },
  { step: 2, title: 'Under Review', description: 'Invoice is being reviewed.', icon: Clock },
  { step: 3, title: 'Approved', description: 'Invoice has been approved.', icon: CircleCheckBig },
  { step: 4, title: 'Disbursed', description: 'Funds have been disbursed.', icon: CreditCard },
  { step: 5, title: 'Paid', description: 'Invoice has been paid.', icon: HandCoins },
];

const invoiceStatusToStepMap: Record<string, number> = {
  'SUBMITTED': 1, 'PENDING ADMIN REVIEW': 2, 'PENDING REVIEW': 2, 'PENDING APPROVAL': 2, 'UNDER REVIEW': 2,
  'APPROVED': 3, 'DISBURSED': 4, 'PAID': 5,
  'REJECTED': 0, 'CANCELLED': 0, 'PROCESSING FAILED': 0, 'UNKNOWN': 0, 'SENT': 1, 'DRAFT': 0,
};

const getStepForInvoiceStatus = (status: string | undefined): number => {
  if (!status) return 0;
  const upperStatus = status.toUpperCase();
  const step = invoiceStatusToStepMap[upperStatus];
  if (typeof step === 'undefined') {
    const partialMatchKey = Object.keys(invoiceStatusToStepMap).find(key => upperStatus.includes(key));
    if (partialMatchKey) return invoiceStatusToStepMap[partialMatchKey];
    console.warn(`Unknown status encountered for stepper: "${status}". Defaulting to step 0.`);
    return 0;
  }
  return step;
};

// --- User Profile ---
const fetchUserProfile = async () => {
  if (!authToken) {
    userProfileError.value = "Authentication token not found for user profile.";
    isLoadingUserProfile.value = false;
    return;
  }
  isLoadingUserProfile.value = true;
  userProfileError.value = null;
  try {
    const response = await apiClient.get('/admin/user/{{ $userId }}');
    userProfile.value = response.data;
    console.log('User profile fetched:', userProfile.value);
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      console.error('API Error (User Profile):', err.response.status, err.response.data);
      if (err.response.status === 401 || err.response.status === 403) {
        userProfileError.value = `Authentication error: ${err.response.data?.message || 'Please log in again.'}`;
      } else {
        userProfileError.value = `Error ${err.response.status}: ${err.response.data?.message || 'Failed to fetch user profile'}`;
      }
    } else {
      console.error('Failed to fetch user profile:', err);
      userProfileError.value = 'An unexpected error occurred while fetching user profile.';
    }
  } finally {
    isLoadingUserProfile.value = false;
  }
};


const fetchInvoices = async () => {
  if (!authToken) {
    invoicesError.value = "Authentication token not found. Please log in.";
    isLoadingInvoices.value = false; return;
  }
  isLoadingInvoices.value = true; invoicesError.value = null;
  try {
    // If your API supports pagination: `/admin/invoices?page=${currentPage.value}&limit=${itemsPerPage.value}`
    // And the response should include totalItems/totalPages.
    // For now, fetching all and paginating client-side.
    const response = await apiClient.get('/admin/invoices');
    let invoiceDataArray = [];
    if (Array.isArray(response.data)) invoiceDataArray = response.data;
    else if (response.data && Array.isArray(response.data.data)) invoiceDataArray = response.data.data;
    else if (response.data && Array.isArray(response.data.invoices)) invoiceDataArray = response.data.invoices;
    else if (response.status === 200 && response.data && typeof response.data === 'object' && !Array.isArray(response.data)) invoiceDataArray = [];
    else throw new Error("Invoice data from API is not in expected array format.");

    invoices.value = invoiceDataArray.map((inv: any) => {
      let parsedJsonData: any = null;
      if (inv.jsonData && typeof inv.jsonData === 'string') {
        try { parsedJsonData = JSON.parse(inv.jsonData); }
        catch (e) { console.error(`Failed to parse jsonData for invoice ID ${inv.id}:`, e, inv.jsonData); }
      }

      let customerName = inv.customer?.name || inv.companyName || 'Unknown Customer';
      let invoiceNumber = inv.invoiceNumber || String(inv.id);
      let invoiceDate = inv.invoiceDate || inv.createdAt;
      let items: any[] = inv.items || [];
      let currency = inv.currency || 'USD';
      let currencySymbol = inv.currencySymbol || '$';
      let totalAmount = parseFloat(inv.totalAmount) || 0;
      let subTotalAmount = parseFloat(inv.subTotal || inv.subTotalAmount) || 0;
      let taxAmount = parseFloat(inv.tax || inv.taxAmount) || 0;
      let balanceDue = parseFloat(inv.balanceDue);
      let taxRatePercentage = parseFloat(inv.taxRate) || 0;
      let paymentTerms = inv.terms || inv.paymentTerms;
      let dueDate = inv.dueDate;
      let amountPaid = parseFloat(inv.amountPaid) || 0;

      if (parsedJsonData) {
        customerName = parsedJsonData.billedTo || parsedJsonData.customerName || customerName;
        invoiceNumber = parsedJsonData.extractedInvoiceNumber || parsedJsonData.invoiceNumber || invoiceNumber;
        invoiceDate = parsedJsonData.invoiceDate || invoiceDate;
        currency = parsedJsonData.extractedCurrency || parsedJsonData.currency || currency;
        currencySymbol = parsedJsonData.currencySymbol || currencySymbol;
        paymentTerms = parsedJsonData.paymentTerms || paymentTerms;
        dueDate = parsedJsonData.dueDate || dueDate;
        if (typeof parsedJsonData.total !== 'undefined') totalAmount = parseFloat(parsedJsonData.total) || 0;
        else if (typeof parsedJsonData.grandTotal !== 'undefined') totalAmount = parseFloat(parsedJsonData.grandTotal) || 0;
        if (typeof parsedJsonData.subtotal !== 'undefined') subTotalAmount = parseFloat(parsedJsonData.subtotal) || 0;
        if (typeof parsedJsonData.tax !== 'undefined') taxAmount = parseFloat(parsedJsonData.tax) || 0;
        else if (typeof parsedJsonData.taxAmount !== 'undefined') taxAmount = parseFloat(parsedJsonData.taxAmount) || 0;
        balanceDue = (typeof parsedJsonData.balanceDue !== 'undefined') ? (parseFloat(parsedJsonData.balanceDue) || totalAmount) : totalAmount;
        amountPaid = (typeof parsedJsonData.amountPaid !== 'undefined') ? (parseFloat(parsedJsonData.amountPaid) || 0) : amountPaid;
        if (typeof parsedJsonData.taxRatePercentage !== 'undefined') taxRatePercentage = parseFloat(parsedJsonData.taxRatePercentage) || 0;
        else if (subTotalAmount > 0 && taxAmount > 0 && subTotalAmount !== 0) taxRatePercentage = parseFloat(((taxAmount / subTotalAmount) * 100).toFixed(2));
        if (parsedJsonData.lineItems && Array.isArray(parsedJsonData.lineItems)) {
          items = parsedJsonData.lineItems.map((item: any, index: number) => ({
            id: item.id || `jsonItem-${inv.id}-${index + 1}`, name: item.item || item.name || 'N/A', description: item.description || '',
            unitPrice: parseFloat(item.unitPrice || item.rate) || 0, quantity: parseInt(item.quantity || item.qty, 10) || 0,
            amount: parseFloat(item.total || item.amount) || 0,
          }));
        }
      }
      if (currency === 'USD' && currencySymbol !== '$') currencySymbol = '$';
      if (currency === 'GHS' && currencySymbol !== 'GH₵') currencySymbol = 'GH₵';
      if (isNaN(balanceDue)) balanceDue = totalAmount - amountPaid;
      if (isNaN(subTotalAmount) && items.length > 0 && isNaN(totalAmount)) subTotalAmount = items.reduce((sum, item) => sum + (item.amount || 0), 0);
      if (isNaN(totalAmount) && !isNaN(subTotalAmount) && !isNaN(taxAmount)) {
        totalAmount = subTotalAmount + taxAmount;
        if(isNaN(balanceDue)) balanceDue = totalAmount - amountPaid;
      }
      const originalFullStatus = inv.status ? String(inv.status).trim().replace(/_/g, ' ').toUpperCase() : 'UNKNOWN';
      let filterableStatus = 'UNKNOWN';
      if (originalFullStatus.includes('PENDING') || originalFullStatus.includes('REVIEW') || originalFullStatus === 'SUBMITTED' || originalFullStatus === 'SENT') filterableStatus = 'PENDING';
      else if (originalFullStatus === 'APPROVED') filterableStatus = 'APPROVED';
      else if (originalFullStatus === 'DISBURSED') filterableStatus = 'DISBURSED';
      else if (originalFullStatus === 'PAID') filterableStatus = 'PAID';
      else if (originalFullStatus === 'REJECTED' || originalFullStatus === 'CANCELLED' || originalFullStatus === 'PROCESSING FAILED') filterableStatus = 'REJECTED';
      let statusForBadgeDisplay = filterableStatus;
      if (originalFullStatus === 'SENT' && filterableStatus === 'PENDING') statusForBadgeDisplay = 'SENT';
      if (originalFullStatus === 'DRAFT') statusForBadgeDisplay = 'DRAFT';

      return {
        ...inv, id: inv.id, customerName, invoiceNumber, totalAmount, subTotalAmount, taxAmount, balanceDue, taxRatePercentage, items,
        status: filterableStatus, statusForBadge: statusForBadgeDisplay, originalStatus: originalFullStatus,
        invoiceDate, dueDate, paymentTerms, customerAddress: inv.customer?.address || (parsedJsonData?.customerAddress) || {},
        currency, currencySymbol, amountPaid,
      };
    }).sort((a: any, b: any) => (a.invoiceDate ? new Date(a.invoiceDate).getTime() : 0) < (b.invoiceDate ? new Date(b.invoiceDate).getTime() : 0) ? 1 : -1);

    if (invoices.value.length > 0 && !selectedInvoiceId.value) selectInvoice(invoices.value[0]);
  } catch (err: any) {
    console.error('Error during fetchInvoices processing:', err);
    if (err.response) {
      if (err.response.status === 401 || err.response.status === 403) invoicesError.value = `Authentication error: ${err.response.data?.message || 'Please log in again.'}`;
      else invoicesError.value = `Error ${err.response.status}: ${err.response.data?.message || 'Failed to fetch invoices'}`;
    } else invoicesError.value = `An unexpected error occurred: ${err.message || 'Failed to process invoice data.'}`;
  } finally { isLoadingInvoices.value = false; }
};

const selectInvoice = (invoice: any) => { selectedInvoiceId.value = invoice.id; };

const computedSelectedInvoice = computed(() => {
  if (!selectedInvoiceId.value) return null;
  return invoices.value.find(inv => inv.id === selectedInvoiceId.value) || null;
});

const filteredInvoices = computed(() => {
  let items = invoices.value;
  if (searchTerm.value.trim() !== '') {
    const lowerSearchTerm = searchTerm.value.toLowerCase();
    items = items.filter(invoice =>
        (invoice.customerName && invoice.customerName.toLowerCase().includes(lowerSearchTerm)) ||
        (invoice.invoiceNumber && invoice.invoiceNumber.toLowerCase().includes(lowerSearchTerm)) ||
        (invoice.id && String(invoice.id).toLowerCase().includes(lowerSearchTerm)) ||
        (invoice.totalAmount && String(invoice.totalAmount).includes(lowerSearchTerm))
    );
  }
  if (selectedStatusFilters.value.length > 0) {
    items = items.filter(invoice => selectedStatusFilters.value.includes(invoice.status));
  }
  return items;
});

// Pagination Computed Properties
const totalPages = computed(() => {
  return Math.ceil(filteredInvoices.value.length / itemsPerPage.value);
});

const paginatedInvoices = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value;
  const end = start + itemsPerPage.value;
  return filteredInvoices.value.slice(start, end);
});

// Pagination Methods
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
  }
};
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};

// Watch for changes in filters that might affect pagination
watch([searchTerm, selectedStatusFilters], () => {
  currentPage.value = 1;
});


const toggleStatusFilter = (statusValue: string) => {
  const index = selectedStatusFilters.value.indexOf(statusValue);
  if (index > -1) selectedStatusFilters.value.splice(index, 1);
  else selectedStatusFilters.value.push(statusValue);
};
const clearStatusFilters = () => { selectedStatusFilters.value = []; };

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) { selectedFile.value = target.files[0]; fileUploadError.value = null; }
  else { selectedFile.value = null; }
};


const handleReceiptFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) { selectedReceiptFile.value = target.files[0]; receiptUploadError.value = null; }
  else { selectedReceiptFile.value = null; }
};

const submitReceipt = async () => {
  if (!computedSelectedInvoice.value) { toast.error("No invoice selected to upload receipt for."); return; }
  if (!authToken) { toast.error("Authentication token not found. Please log in."); isUploadingReceipt.value = false; return; }
  if (!selectedReceiptFile.value) { receiptUploadError.value = "Please select a receipt file to upload."; return; }

  isUploadingReceipt.value = true; receiptUploadError.value = null;
  const formData = new FormData();
  formData.append('receiptFile', selectedReceiptFile.value);

  try {
    const uploadApiClient = axios.create({ baseURL: API_BASE_URL, headers: { ...(authToken ? { 'Authorization': `Bearer ${authToken}` } : {}), }});
    await uploadApiClient.post(`/admin/invoices/${computedSelectedInvoice.value.id}/receipt`, formData);
    toast.success(`Receipt uploaded successfully for invoice ${computedSelectedInvoice.value.invoiceNumber || computedSelectedInvoice.value.id}!`);
    selectedReceiptFile.value = null;
    isUploadReceiptDialogOpen.value = false;
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      const errorMessage = err.response.data?.message || err.response.data?.error || 'Failed to upload receipt';
      if (err.response.status === 401 || err.response.status === 403) toast.error(`Authentication error: ${errorMessage}`);
      else if (err.response.status === 400 && (typeof errorMessage === 'string' && errorMessage.toLowerCase().includes("receipt file"))) receiptUploadError.value = `Upload Error: ${errorMessage}`;
      else toast.error(`Upload Error ${err.response.status}: ${errorMessage}`);
    } else toast.error('An unexpected error occurred during receipt upload.');
    console.error('Failed to upload receipt:', err);
  } finally { isUploadingReceipt.value = false; }
};

const triggerDownloadPDF = async (invoiceId: string | number, documentType: 'invoice' | 'receipt' = 'invoice') => {
  if (!authToken) { toast.error("Authentication token not found. Please log in."); return; }
  toast.info("Preparing download...");

  const endpoint = `/admin/invoices/${invoiceId}/pdf`;

  try {
    const response = await apiClient.get(endpoint, { responseType: 'blob' });
    const blob = new Blob([response.data], { type: response.headers['content-type'] || 'application/pdf' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    const contentDisposition = response.headers['content-disposition'];
    let filename = `${documentType}-${invoiceId}.pdf`;
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename="?(.+)"?/i);
      if (filenameMatch && filenameMatch.length === 2) filename = filenameMatch[1];
    }
    link.download = filename;
    document.body.appendChild(link); link.click(); document.body.removeChild(link);
    URL.revokeObjectURL(link.href);
    toast.success("Download started: " + filename);
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      if (err.response.data instanceof Blob && err.response.data.type === "application/json") {
        const errorText = await err.response.data.text();
        try { const errorJson = JSON.parse(errorText); toast.error(`Download Error ${err.response.status}: ${errorJson.message || 'Failed to download.'}`); }
        catch (parseError) { toast.error(`Download Error ${err.response.status}: Failed to download. Invalid error format.`); }
      } else { const errorMessage = err.response.data?.message || 'Failed to download document.'; toast.error(`Download Error ${err.response.status}: ${errorMessage}`);}
    } else { toast.error('An unexpected error occurred while downloading the document.'); }
    console.error('Failed to download document:', err);
  }
};

const updateInvoiceStatus = async (invoiceId: string | number, newStatus: string) => {
  if (!authToken) { toast.error("Authentication token not found. Please log in."); return; }
  toast.info(`Updating status to ${newStatus}...`);
  try {
    // Using PUT as per the previous context for admin status update
    await apiClient.put(`/admin/invoices/${invoiceId}/status`, { status: newStatus });
    toast.success(`Invoice status updated to ${newStatus}!`);
    await fetchInvoices();
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      const errorMessage = err.response.data?.message || `Failed to update status to ${newStatus}`;
      toast.error(`Error ${err.response.status}: ${errorMessage}`);
    } else {
      toast.error('An unexpected error occurred while updating status.');
    }
    console.error(`Failed to update status for invoice ${invoiceId}:`, err);
  }
};

const formatListStatus = (status: string): string => {
  if (!status || status === 'UNKNOWN') return 'Unknown';
  return status.toUpperCase();
};

const getListInvoiceStatusBadgeClass = (status: string): string => {
  const upperStatus = status ? status.toUpperCase() : 'UNKNOWN';
  if (['PAID', 'APPROVED', 'DISBURSED'].includes(upperStatus)) return 'bg-green-100 text-green-700';
  if (upperStatus === 'SENT') return 'bg-blue-100 text-blue-700';
  if (upperStatus === 'PENDING') return 'bg-yellow-100 text-yellow-700';
  if (upperStatus === 'DRAFT') return 'bg-gray-200 text-gray-700';
  if (['REJECTED', 'CANCELLED'].includes(upperStatus)) return 'bg-red-100 text-red-700';
  return 'bg-gray-100 text-gray-500';
};

const getDetailStatusClass = (status: string): string => {
  const upperStatus = status ? status.toUpperCase() : 'UNKNOWN';
  if (['PAID', 'APPROVED', 'DISBURSED'].includes(upperStatus)) return 'bg-green-100 text-green-700';
  if (upperStatus === 'SENT') return 'bg-blue-100 text-blue-700';
  if (upperStatus === 'PENDING') return 'bg-yellow-100 text-yellow-700';
  if (upperStatus === 'DRAFT') return 'bg-gray-200 text-gray-700';
  if (['REJECTED', 'CANCELLED'].includes(upperStatus)) return 'bg-red-100 text-red-700';
  return 'bg-gray-100 text-gray-500';
};

const formatDate = (dateString?: string | Date, format: 'default' | 'dd/MM/yyyy' = 'default'): string => {
  if (!dateString) return 'N/A';
  try {
    const date = new Date(dateString);
    if (format === 'dd/MM/yyyy') {
      const day = String(date.getDate()).padStart(2, '0');
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const year = date.getFullYear();
      return `${day}/${month}/${year}`;
    }
    return date.toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' });
  } catch (e) { console.warn("Invalid date string for formatDate:", dateString); return 'Invalid Date'; }
}

onMounted(() => {
  if (authToken) {
    fetchUserProfile();
    fetchInvoices();
  }
  else {
    invoicesError.value = "You are not logged in. Please log in to view invoices.";
    userProfileError.value = "Cannot fetch profile. Please log in.";
    toast.info("Please log in to access all features.", { description: "Authentication token not found." });
  }
})

</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');
</style>
