<template>
  <main class="flex flex-col font-inter z-50 items-start gap-4 p-4 sm:px-6 sm:py-2 md:gap-8">
    <div class="container w-full mx-auto p-4 md:p-6 lg:p-8">

      <div class=" p-6 rounded-lg shadow-md">
        <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between mb-4">
          <Input
              type="text"
              v-model="searchTerm"
              placeholder="Search by name, email, company, ID..."
              class="w-full sm:max-w-md py-2 px-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
          />
        </div>

        <div v-if="isLoadingCustomers" class="flex flex-col justify-center items-center h-64">
          <div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-primary"></div>
          <p class="mt-3">Loading Customers...</p>
        </div>
        <div v-else-if="customersError" class="text-center py-10 text-red-500">
          <AlertCircle class="w-12 h-12 mx-auto mb-2 text-red-400" />
          <p class="text-xl font-semibold">Error Loading Customers</p>
          <p>{{ customersError }}</p>
        </div>
        <div v-else-if="filteredCustomers.length === 0" class="text-center py-10 ">
          <Ban class="w-12 h-12 mx-auto mb-2 " />
          <p class="text-xl font-semibold">No Customers Found</p>
          <p>{{ searchTerm ? 'Try adjusting your search.' : 'There are no customers to display for the current page.' }}</p>
        </div>
        <div v-else class="overflow-x-auto">
          <Table class="w-full min-w-[800px]">
            <TableHeader>
              <TableRow>
                <TableHead>ID</TableHead>
                <TableHead>Avatar</TableHead>
                <TableHead>Name</TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Company</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                  v-for="customer in filteredCustomers"
                  :key="customer.id"
                  @click="openCustomerDetails(customer)"
                  class="cursor-pointer hover:bg-gray-50"
              >
                <TableCell class="font-medium">{{ customer.id }}</TableCell>
                <TableCell>
                  <Avatar class="relative bg-gray-200 overflow-visible h-10 w-10">
                    <AvatarImage class="rounded-full" :src="customer.avatarUrl || ''" :alt="`${customer.firstName} ${customer.lastName}`" />
                    <AvatarFallback class="bg-blue-500 text-white">
                      {{ customer.firstName ? customer.firstName.substring(0, 1).toUpperCase() : '' }}{{ customer.lastName ? customer.lastName.substring(0, 1).toUpperCase() : 'N/A' }}
                    </AvatarFallback>
                  </Avatar>
                </TableCell>
                <TableCell>{{ customer.firstName }} {{ customer.lastName }}</TableCell>
                <TableCell>{{ customer.email }}</TableCell>
                <TableCell>{{ customer.companyName || 'N/A' }}</TableCell>
                <TableCell>
                  <Badge :variant="customer.isActive ? 'default' : 'destructive'">
                    {{ customer.isActive ? 'Active' : 'Inactive' }}
                  </Badge>
                </TableCell>
                <TableCell>
                  <Button variant="outline" size="sm" @click.stop="openCustomerDetails(customer)">
                    <Eye class="w-4 h-4 mr-1" /> View Details
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
        <div v-if="totalPages > 1" class="mt-6 flex justify-center items-center space-x-2">
          <Button @click="prevPage" :disabled="currentPage === 1" variant="outline" size="sm" class="text-xs px-2.5 py-1">Previous</Button>
          <span v-for="pageNumber in paginationNumbers" :key="pageNumber">
                <Button
                    v-if="pageNumber !== '...'"
                    @click="goToPage(pageNumber as number)"
                    :variant="currentPage === pageNumber ? 'default' : 'outline'"
                    size="sm"
                    class="text-xs px-3 py-1"
                >
                    {{ pageNumber }}
                </Button>
                <span v-else class="text-xs px-2 py-1">...</span>
            </span>
          <Button @click="nextPage" :disabled="currentPage === totalPages" variant="outline" size="sm" class="text-xs px-2.5 py-1">Next</Button>
        </div>
        <div v-if="totalPages > 0" class="mt-2 text-center text-xs text-gray-500">
          Showing {{ filteredCustomers.length }} of {{ totalCustomers }} customers
        </div>

      </div>
    </div>

    <Dialog :open="isCustomerDetailsOpen" @update:open="isCustomerDetailsOpen = $event">
      <DialogContent
          class="fixed ml-auto right-0 h-[100vh] w-1/2 flex flex-col transform transition-transform duration-300 ease-out translate-x-full md:translate-x-0 shadow-lg overflow-y-auto p-6"
          :show-close-button="true"
          @escape-key-down="closeCustomerDetails"
          @pointer-down-outside="closeCustomerDetails"
      >
        <DialogHeader class="p-6 border-b sticky top-0 bg-white z-10">
          <DialogTitle class="text-xl font-semibold">Customer KYC Details</DialogTitle>
          <DialogDescription v-if="selectedCustomer">
            For {{ selectedCustomer.firstName }} {{ selectedCustomer.lastName }} (User ID: {{ selectedCustomer.id }})
          </DialogDescription>
        </DialogHeader>

        <div class="p-6 flex-1 overflow-y-auto">
          <div v-if="isLoadingCustomerDetails" class="flex justify-center items-center h-40">
            <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-600"></div>
          </div>
          <div v-else-if="customerDetailsError" class="text-red-500 p-4 bg-red-50 rounded-md">
            <AlertCircle class="w-6 h-6 inline mr-2" /> Error: {{ customerDetailsError }}
          </div>
          <div v-else-if="selectedCustomerDetails">
            <div class="space-y-6">
              <div>
                <h3 class="text-md font-semibold text-gray-700 mb-2">KYC Overview</h3>
                <div class="bg-gray-50 p-4 rounded-md text-sm space-y-2">
                  <p><strong>KYC Record ID:</strong> {{ selectedCustomerDetails.id }}</p>
                  <p><strong>User Email:</strong> {{ selectedCustomerDetails.userEmail }}</p>
                  <p><strong>Current KYC Status:</strong>
                    <Badge :variant="getKycStatusBadgeVariant(selectedCustomerDetails.status)">
                      {{ formatKycStatus(selectedCustomerDetails.status) }}
                    </Badge>
                  </p>
                  <p><strong>Submitted At:</strong> {{ formatDate(selectedCustomerDetails.submittedAt, true) || 'N/A' }}</p>
                  <p><strong>Reviewed At:</strong> {{ formatDate(selectedCustomerDetails.reviewedAt, true) || 'N/A' }}</p>
                  <p><strong>Reviewed By:</strong> {{ selectedCustomerDetails.reviewedByEmail || 'N/A' }}</p>
                  <p v-if="selectedCustomerDetails.remarks"><strong>Previous Remarks:</strong> {{ selectedCustomerDetails.remarks }}</p>
                  <p><strong>KYC Record Created:</strong> {{ formatDate(selectedCustomerDetails.createdAt, true) || 'N/A' }}</p>
                  <p><strong>KYC Record Updated:</strong> {{ formatDate(selectedCustomerDetails.updatedAt, true) || 'N/A' }}</p>
                </div>
              </div>

              <div v-if="selectedCustomerDetails.parsedDocumentsInfo">
                <h3 class="text-md font-semibold text-gray-700 mb-2">Document Information</h3>
                <div class="bg-gray-50 p-4 rounded-md text-sm space-y-2">
                  <p><strong>Full Name (as per doc):</strong> {{ selectedCustomerDetails.parsedDocumentsInfo.full_name || 'N/A' }}</p>
                  <p><strong>Date of Birth (as per doc):</strong> {{ formatDate(selectedCustomerDetails.parsedDocumentsInfo.date_of_birth) || 'N/A' }}</p>
                  <p><strong>Address (as per doc):</strong> {{ selectedCustomerDetails.parsedDocumentsInfo.address || 'N/A' }}</p>
                  <hr class="my-2">
                  <p><strong>Document Type:</strong> {{ selectedCustomerDetails.parsedDocumentsInfo.document_type || 'N/A' }}</p>
                  <p><strong>Document Number:</strong> {{ selectedCustomerDetails.parsedDocumentsInfo.document_number || 'N/A' }}</p>
                  <p><strong>Document Issue Date:</strong> {{ formatDate(selectedCustomerDetails.parsedDocumentsInfo.document_issue_date) || 'N/A' }}</p>
                  <p><strong>Document Expiry Date:</strong> {{ formatDate(selectedCustomerDetails.parsedDocumentsInfo.document_expiry_date) || 'N/A' }}</p>
                  <div v-if="selectedCustomerDetails.parsedDocumentsInfo.document_urls && selectedCustomerDetails.parsedDocumentsInfo.document_urls.length > 0">
                    <p><strong>Document URLs:</strong></p>
                    <ul class="list-disc list-inside pl-4">
                      <li v-for="(url, index) in selectedCustomerDetails.parsedDocumentsInfo.document_urls" :key="index">
                        <a :href="url" target="_blank" class="text-blue-600 hover:underline">View Document {{ index + 1 }}</a>
                      </li>
                    </ul>
                  </div>
                  <p v-if="selectedCustomerDetails.parsedDocumentsInfo.additional_notes"><strong>Additional Notes:</strong> {{ selectedCustomerDetails.parsedDocumentsInfo.additional_notes }}</p>
                </div>
              </div>
              <div v-else>
                <p class="text-sm text-gray-500">No detailed document information available or failed to parse.</p>
              </div>


              <div class="pt-4 border-t mt-6">
                <h3 class="text-md font-semibold text-gray-700 mb-2">Update KYC Status</h3>
                <div class="space-y-3">
                  <div>
                    <Label for="kycStatus" class="text-sm">New Status</Label>
                    <select v-model="kycReviewData.status" id="kycStatus" class="mt-1 block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm">
                      <option value="" disabled>Select status</option>
                      <option value="PENDING">Pending</option>
                      <option value="VERIFIED">Verified / Approved</option>
                      <option value="REJECTED">Rejected</option>
                      <option value="ACTION_REQUIRED">Action Required</option>
                    </select>
                  </div>
                  <div>
                    <Label for="kycRemarks" class="text-sm">Remarks</Label>
                    <textarea v-model="kycReviewData.remarks" id="kycRemarks" rows="3" class="mt-1 block w-full py-2 px-3 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm" placeholder="Reason for status change... (Required for Rejected/Action Required)"></textarea>
                  </div>
                  <Button @click="submitKycReview" :disabled="isUpdatingKycStatus || !kycReviewData.status" class="w-full sm:w-auto">
                    {{ isUpdatingKycStatus ? 'Updating...' : 'Submit Review' }}
                  </Button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center text-gray-500 py-10">
            <p>No KYC details available for this customer or an error occurred.</p>
          </div>
        </div>
        <DialogFooter class="p-4 border-t sticky bottom-0 bg-white z-10">
          <Button variant="outline" @click="closeCustomerDetails">Close</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import axios from 'axios';
import {
  AlertCircle, Eye, Plus, Slash, Ban,
} from 'lucide-vue-next';
import { Badge } from '~/components/ui/badge';
import { Button } from '~/components/ui/button';
import { Input } from '~/components/ui/input';
import { Avatar, AvatarImage, AvatarFallback } from '~/components/ui/avatar';
import {
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow,
} from '~/components/ui/table';
import {
  Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger,
} from '~/components/ui/dialog';
import {
  Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbSeparator,
} from '~/components/ui/breadcrumb';
import { Label } from '~/components/ui/label';
import { Toaster, toast } from 'vue-sonner';
import { useCookie } from '#app';

// API Configuration
const API_BASE_URL = 'http://localhost:3000/api/v1';
const tokenCookie = useCookie('token');
const authToken = tokenCookie.value || null;

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    ...(authToken ? { Authorization: `Bearer ${authToken}` } : {}),
    'Content-Type': 'application/json',
  },
});

// Reactive State
const customers = ref<any[]>([]);
const isLoadingCustomers = ref(true);
const customersError = ref<string | null>(null);
const searchTerm = ref('');

const selectedCustomer = ref<any | null>(null);
const selectedCustomerDetails = ref<any | null>(null);
const isLoadingCustomerDetails = ref(false);
const customerDetailsError = ref<string | null>(null);
const isCustomerDetailsOpen = ref(false);

const kycReviewData = ref({
  status: '',
  remarks: '',
});
const isUpdatingKycStatus = ref(false);

// Pagination State
const currentPage = ref(1);
const itemsPerPage = ref(10);
const totalCustomers = ref(0);


// Fetch Customers
const fetchCustomers = async (page = 1, pageSize = itemsPerPage.value) => {
  isLoadingCustomers.value = true;
  customersError.value = null;
  try {
    const response = await apiClient.get('/admin/users', {
      params: { page, pageSize },
    });
    const fetchedUsers = response.data.users || [];
    // Sort by ID ascending (starting from 1)
    customers.value = fetchedUsers.sort((a: any, b: any) => a.id - b.id);
    totalCustomers.value = response.data.total || 0;
    currentPage.value = response.data.page || 1;
  } catch (err) {
    console.error('Failed to fetch customers:', err);
    if (axios.isAxiosError(err) && err.response) {
      customersError.value = `Error ${err.response.status}: ${err.response.data?.message || 'Could not load customers.'}`;
    } else {
      customersError.value = 'An unexpected error occurred while fetching customers.';
    }
    toast.error(customersError.value || 'Failed to load customers.');
  } finally {
    isLoadingCustomers.value = false;
  }
};

// Fetch Customer KYC Details
const fetchCustomerDetails = async (customerId: string) => {
  if (!customerId) return;
  isLoadingCustomerDetails.value = true;
  customerDetailsError.value = null;
  selectedCustomerDetails.value = null;
  try {
    const response = await apiClient.get(`/admin/users/${customerId}/kyc`);
    const kycData = response.data.kycDetails || response.data.data || response.data;

    if (kycData && typeof kycData.documentsInfo === 'string') {
      try {
        kycData.parsedDocumentsInfo = JSON.parse(kycData.documentsInfo);
      } catch (e) {
        console.error("Failed to parse documentsInfo JSON:", e);
        kycData.parsedDocumentsInfo = null;
        toast.error("Could not parse additional document information.");
      }
    } else if (kycData && typeof kycData.documentsInfo === 'object') {
      kycData.parsedDocumentsInfo = kycData.documentsInfo;
    } else if (kycData) {
      kycData.parsedDocumentsInfo = null;
    }

    selectedCustomerDetails.value = kycData;

    kycReviewData.value.status = selectedCustomerDetails.value?.status || '';
    kycReviewData.value.remarks = selectedCustomerDetails.value?.remarks || '';
  } catch (err) {
    console.error(`Failed to fetch KYC details for customer ${customerId}:`, err);
    if (axios.isAxiosError(err) && err.response) {
      customerDetailsError.value = `Error ${err.response.status}: ${err.response.data?.message || 'Could not load KYC details.'}`;
    } else {
      customerDetailsError.value = 'An unexpected error occurred.';
    }
    toast.error(customerDetailsError.value || 'Failed to load KYC details.');
  } finally {
    isLoadingCustomerDetails.value = false;
  }
};

// Open Customer Details Panel
const openCustomerDetails = (customer: any) => {
  selectedCustomer.value = customer;
  isCustomerDetailsOpen.value = true;
  fetchCustomerDetails(customer.id);
};

const closeCustomerDetails = () => {
  isCustomerDetailsOpen.value = false;
  selectedCustomer.value = null;
  selectedCustomerDetails.value = null;
  customerDetailsError.value = null;
  kycReviewData.value = { status: '', remarks: '' };
};

// Submit KYC Review
const submitKycReview = async () => {
  if (!selectedCustomer.value || !kycReviewData.value.status) {
    toast.error('Please select a new status for KYC review.');
    return;
  }
  if ((kycReviewData.value.status === 'REJECTED' || kycReviewData.value.status === 'ACTION_REQUIRED') && !kycReviewData.value.remarks.trim()) {
    toast.error('Remarks are required when rejecting or requesting action.');
    return;
  }
  isUpdatingKycStatus.value = true;
  try {
    const payload = {
      status: kycReviewData.value.status,
      remarks: kycReviewData.value.remarks || undefined,
    };
    await apiClient.put(`/admin/users/${selectedCustomer.value.id}/kyc/review`, payload);
    toast.success('KYC status updated successfully.');
    await fetchCustomers(currentPage.value, itemsPerPage.value);
    if (isCustomerDetailsOpen.value) {
      await fetchCustomerDetails(selectedCustomer.value.id);
    }
  } catch (err) {
    console.error('Failed to update KYC status:', err);
    let errorMsg = 'Failed to update KYC status.';
    if (axios.isAxiosError(err) && err.response) {
      errorMsg = `Error ${err.response.status}: ${err.response.data?.message || 'Update failed.'}`;
    }
    toast.error(errorMsg);
  } finally {
    isUpdatingKycStatus.value = false;
  }
};

// Computed property for filtered customers (client-side on current page data)
const filteredCustomers = computed(() => {
  if (!searchTerm.value) {
    return customers.value;
  }
  const lowerSearchTerm = searchTerm.value.toLowerCase();
  return customers.value.filter(customer =>
      ((customer.firstName + ' ' + customer.lastName).toLowerCase().includes(lowerSearchTerm)) ||
      (customer.email && customer.email.toLowerCase().includes(lowerSearchTerm)) ||
      (customer.companyName && customer.companyName.toLowerCase().includes(lowerSearchTerm)) ||
      (customer.id && String(customer.id).toLowerCase().includes(lowerSearchTerm))
  );
});

// Pagination Computed Properties
const totalPages = computed(() => {
  return Math.ceil(totalCustomers.value / itemsPerPage.value);
});

const paginationNumbers = computed(() => {
  const delta = 1;
  const range = [];
  const rangeWithDots: (number | string)[] = [];
  let l: number | undefined;

  if (totalPages.value <= 1) return [1];

  range.push(1);

  let left = Math.max(2, currentPage.value - delta);
  let right = Math.min(totalPages.value - 1, currentPage.value + delta);

  if (currentPage.value - delta > 2) {
    rangeWithDots.push(1);
    rangeWithDots.push('...');
  } else {
    rangeWithDots.push(1);
  }


  for (let i = left; i <= right; i++) {
    if (!rangeWithDots.includes(i)) {
      rangeWithDots.push(i);
    }
  }

  if (currentPage.value + delta < totalPages.value -1 ) {
    if (rangeWithDots[rangeWithDots.length-1] !== '...') rangeWithDots.push('...');
  }
  if (!rangeWithDots.includes(totalPages.value)) {
    rangeWithDots.push(totalPages.value);
  }

  let finalRange: (number | string)[] = [];
  let lastPushed: (number | string) | null = null;
  for(const p of rangeWithDots){
    if(p === '...' && lastPushed === '...'){
      continue;
    }
    if(typeof p === 'number' && typeof lastPushed === 'number' && p <= lastPushed){
      continue;
    }
    finalRange.push(p);
    lastPushed = p;
  }
  if (finalRange.length === 2 && finalRange[0] === 1 && finalRange[1] === '...') {
    return [1, 2];
  }


  return finalRange;
});


// Pagination Methods
const goToPage = (page: number) => {
  if (page >= 1 && page <= totalPages.value && page !== currentPage.value) {
    fetchCustomers(page, itemsPerPage.value);
  }
};
const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    fetchCustomers(currentPage.value + 1, itemsPerPage.value);
  }
};
const prevPage = () => {
  if (currentPage.value > 1) {
    fetchCustomers(currentPage.value - 1, itemsPerPage.value);
  }
};

// UI Helper Functions
const formatDate = (dateString?: string | Date, includeTime: boolean = false): string => {
  if (!dateString) return 'N/A';
  try {
    const options: Intl.DateTimeFormatOptions = {
      day: '2-digit', month: 'short', year: 'numeric',
    };
    if (includeTime) {
      options.hour = '2-digit';
      options.minute = '2-digit';
      // options.second = '2-digit'; // Optional: include seconds
    }
    return new Date(dateString).toLocaleString('en-GB', options); // Using toLocaleString for better time formatting
  } catch (e) {
    return 'Invalid Date';
  }
};

const formatKycStatus = (status?: string): string => {
  if (!status) return 'Unknown';
  return status.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
};

const getKycStatusBadgeVariant = (status?: string): 'default' | 'secondary' | 'destructive' | 'outline' => {
  const s = status?.toUpperCase();
  if (s === 'VERIFIED' || s === 'APPROVED') return 'default';
  if (s === 'PENDING' || s === 'SUBMITTED' || s === 'UNDER_REVIEW') return 'secondary';
  if (s === 'REJECTED' || s === 'ACTION_REQUIRED') return 'destructive';
  return 'outline';
};

const formatAddress = (address: any): string => {
  if (!address) return 'N/A';
  if (typeof address === 'string') return address;
  const parts = [address.street, address.city, address.state, address.zipCode, address.country].filter(Boolean);
  return parts.join(', ') || 'N/A';
};

// Lifecycle Hook
onMounted(() => {
  if (authToken) {
    fetchCustomers(currentPage.value, itemsPerPage.value);
  } else {
    customersError.value = "Authentication required. Please log in.";
    toast.error(customersError.value);
  }
});
</script>

<style>
/* Additional custom styles if needed */
.dialog-content-full-height {
  height: 100vh;
  max-height: 100vh;
  margin: 0;
  border-radius: 0;
}
</style>