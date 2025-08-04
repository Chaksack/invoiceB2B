<script setup>
import { ref, computed } from 'vue';

// This data would typically come from an API
const payments = ref([
  { 
    id: 'PAY-2025-001', 
    invoiceId: 'INV-2025-001', 
    client: 'Acme Corp', 
    amount: '$10,000.00', 
    date: '2025-07-25', 
    type: 'Advance', 
    status: 'Completed',
    method: 'Bank Transfer'
  },
  { 
    id: 'PAY-2025-002', 
    invoiceId: 'INV-2025-003', 
    client: 'Stark Industries', 
    amount: '$12,160.00', 
    date: '2025-07-22', 
    type: 'Advance', 
    status: 'Completed',
    method: 'Bank Transfer'
  },
  { 
    id: 'PAY-2025-003', 
    invoiceId: 'INV-2025-005', 
    client: 'Umbrella Corp', 
    amount: '$9,080.00', 
    date: '2025-07-18', 
    type: 'Advance', 
    status: 'Completed',
    method: 'Bank Transfer'
  },
  { 
    id: 'PAY-2025-004', 
    invoiceId: 'INV-2025-001', 
    client: 'Acme Corp', 
    amount: '$2,500.00', 
    date: '2025-08-25', 
    type: 'Balance', 
    status: 'Pending',
    method: 'Automatic'
  },
  { 
    id: 'PAY-2025-005', 
    invoiceId: 'INV-2025-003', 
    client: 'Stark Industries', 
    amount: '$3,040.00', 
    date: '2025-08-22', 
    type: 'Balance', 
    status: 'Pending',
    method: 'Automatic'
  },
  { 
    id: 'PAY-2025-006', 
    invoiceId: 'INV-2025-005', 
    client: 'Umbrella Corp', 
    amount: '$2,270.00', 
    date: '2025-08-18', 
    type: 'Balance', 
    status: 'Pending',
    method: 'Automatic'
  },
  { 
    id: 'PAY-2025-007', 
    invoiceId: 'INV-2025-007', 
    client: 'LexCorp', 
    amount: '$15,120.00', 
    date: '2025-07-12', 
    type: 'Advance', 
    status: 'Completed',
    method: 'Bank Transfer'
  },
  { 
    id: 'PAY-2025-008', 
    invoiceId: 'INV-2025-007', 
    client: 'LexCorp', 
    amount: '$3,780.00', 
    date: '2025-08-12', 
    type: 'Balance', 
    status: 'Pending',
    method: 'Automatic'
  },
  { 
    id: 'PAY-2025-009', 
    invoiceId: 'INV-2025-009', 
    client: 'Massive Dynamic', 
    amount: '$18,200.00', 
    date: '2025-07-08', 
    type: 'Advance', 
    status: 'Completed',
    method: 'Bank Transfer'
  },
  { 
    id: 'PAY-2025-010', 
    invoiceId: 'INV-2025-009', 
    client: 'Massive Dynamic', 
    amount: '$4,550.00', 
    date: '2025-08-08', 
    type: 'Balance', 
    status: 'Pending',
    method: 'Automatic'
  }
]);

// Filter and sort state
const searchQuery = ref('');
const statusFilter = ref('All');
const typeFilter = ref('All');
const dateRange = ref({
  start: '',
  end: ''
});
const sortBy = ref('date');
const sortDirection = ref('desc');

// Status and type options for filters
const statusOptions = ['All', 'Completed', 'Pending', 'Failed'];
const typeOptions = ['All', 'Advance', 'Balance', 'Fee'];

// Computed filtered and sorted payments
const filteredPayments = computed(() => {
  let result = [...payments.value];
  
  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(payment => 
      payment.id.toLowerCase().includes(query) || 
      payment.invoiceId.toLowerCase().includes(query) ||
      payment.client.toLowerCase().includes(query)
    );
  }
  
  // Apply status filter
  if (statusFilter.value !== 'All') {
    result = result.filter(payment => payment.status === statusFilter.value);
  }
  
  // Apply type filter
  if (typeFilter.value !== 'All') {
    result = result.filter(payment => payment.type === typeFilter.value);
  }
  
  // Apply date range filter
  if (dateRange.value.start) {
    result = result.filter(payment => new Date(payment.date) >= new Date(dateRange.value.start));
  }
  if (dateRange.value.end) {
    result = result.filter(payment => new Date(payment.date) <= new Date(dateRange.value.end));
  }
  
  // Apply sorting
  result.sort((a, b) => {
    let comparison = 0;
    if (sortBy.value === 'amount') {
      // Remove $ and commas for numeric comparison
      const amountA = parseFloat(a.amount.replace('$', '').replace(',', ''));
      const amountB = parseFloat(b.amount.replace('$', '').replace(',', ''));
      comparison = amountA - amountB;
    } else if (sortBy.value === 'date') {
      comparison = new Date(a.date) - new Date(b.date);
    } else {
      comparison = a[sortBy.value].localeCompare(b[sortBy.value]);
    }
    
    return sortDirection.value === 'asc' ? comparison : -comparison;
  });
  
  return result;
});

// Toggle sort direction
const toggleSort = (column) => {
  if (sortBy.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc';
  } else {
    sortBy.value = column;
    sortDirection.value = 'desc';
  }
};

// Reset filters
const resetFilters = () => {
  searchQuery.value = '';
  statusFilter.value = 'All';
  typeFilter.value = 'All';
  dateRange.value = {
    start: '',
    end: ''
  };
};

// Selected payment for details view
const selectedPayment = ref(null);
const showDetails = ref(false);

const viewPaymentDetails = (payment) => {
  selectedPayment.value = payment;
  showDetails.value = true;
};

const closeDetails = () => {
  showDetails.value = false;
};

// Export payments
const exportPayments = () => {
  alert('Payments would be exported to CSV/Excel in a real application');
};

definePageMeta({
  layout: 'dashboard'
});
</script>

<template>
  <div>
    <!-- Page Header -->
    <div class="mb-6 flex flex-col md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Payment History</h1>
        <p class="text-gray-500 mt-1">View and track all your financing payments</p>
      </div>
      <div class="mt-4 md:mt-0">
        <button 
          @click="exportPayments"
          class="bg-white text-gray-700 border border-gray-300 px-4 py-2 rounded-md hover:bg-gray-50 transition-colors flex items-center"
        >
          <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="mr-2">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
          </svg>
          Export
        </button>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="bg-card rounded-lg shadow-sm border border-border p-4 mb-6">
      <div class="flex flex-col space-y-4">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div class="relative flex-1">
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Search payments..." 
              class="w-full py-2 pl-10 pr-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
              </svg>
            </div>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-4">
            <button 
              @click="resetFilters"
              class="px-4 py-2 border border-input rounded-md hover:bg-muted/20 transition-colors flex items-center justify-center"
            >
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
              </svg>
              Reset Filters
            </button>
            <button 
              @click="exportPayments"
              class="px-4 py-2 border border-input rounded-md hover:bg-muted/20 transition-colors flex items-center justify-center"
            >
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="mr-2">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
              </svg>
              Export
            </button>
          </div>
        </div>
        
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label for="statusFilter" class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select 
              id="statusFilter"
              v-model="statusFilter"
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            >
              <option v-for="status in statusOptions" :key="status" :value="status">
                {{ status }}
              </option>
            </select>
          </div>
          
          <div>
            <label for="typeFilter" class="block text-sm font-medium text-gray-700 mb-1">Payment Type</label>
            <select 
              id="typeFilter"
              v-model="typeFilter"
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            >
              <option v-for="type in typeOptions" :key="type" :value="type">
                {{ type }}
              </option>
            </select>
          </div>
          
          <div>
            <label for="startDate" class="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
            <input 
              id="startDate"
              v-model="dateRange.start"
              type="date" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <label for="endDate" class="block text-sm font-medium text-gray-700 mb-1">End Date</label>
            <input 
              id="endDate"
              v-model="dateRange.end"
              type="date" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Payments Table -->
    <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="bg-muted/50">
              <th 
                @click="toggleSort('id')" 
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
              >
                <div class="flex items-center">
                  Payment ID
                  <svg v-if="sortBy === 'id'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
                    <path v-if="sortDirection === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </th>
              <th 
                @click="toggleSort('invoiceId')" 
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
              >
                <div class="flex items-center">
                  Invoice ID
                  <svg v-if="sortBy === 'invoiceId'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
                    <path v-if="sortDirection === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </th>
              <th 
                @click="toggleSort('client')" 
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
              >
                <div class="flex items-center">
                  Client
                  <svg v-if="sortBy === 'client'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
                    <path v-if="sortDirection === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </th>
              <th 
                @click="toggleSort('amount')" 
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
              >
                <div class="flex items-center">
                  Amount
                  <svg v-if="sortBy === 'amount'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
                    <path v-if="sortDirection === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </th>
              <th 
                @click="toggleSort('date')" 
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
              >
                <div class="flex items-center">
                  Date
                  <svg v-if="sortBy === 'date'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
                    <path v-if="sortDirection === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </th>
              <th 
                @click="toggleSort('type')" 
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
              >
                <div class="flex items-center">
                  Type
                  <svg v-if="sortBy === 'type'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
                    <path v-if="sortDirection === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </th>
              <th 
                @click="toggleSort('status')" 
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
              >
                <div class="flex items-center">
                  Status
                  <svg v-if="sortBy === 'status'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
                    <path v-if="sortDirection === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7"></path>
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                  </svg>
                </div>
              </th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border">
            <tr v-for="payment in filteredPayments" :key="payment.id" class="hover:bg-muted/20">
              <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ payment.id }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ payment.invoiceId }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ payment.client }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ payment.amount }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ payment.date }}</td>
              <td class="px-4 py-3 text-sm">
                <span :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  payment.type === 'Advance' ? 'bg-blue-100 text-blue-800' : 
                  payment.type === 'Balance' ? 'bg-purple-100 text-purple-800' : 
                  'bg-gray-100 text-gray-800'
                ]">
                  {{ payment.type }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm">
                <span :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  payment.status === 'Completed' ? 'bg-green-100 text-green-800' : 
                  payment.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 
                  'bg-red-100 text-red-800'
                ]">
                  {{ payment.status }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-right">
                <button 
                  @click="viewPaymentDetails(payment)"
                  class="text-primary hover:text-primary-700 font-medium"
                >
                  View
                </button>
              </td>
            </tr>
            <tr v-if="filteredPayments.length === 0">
              <td colspan="8" class="px-4 py-6 text-center text-gray-500">
                No payments found matching your filters.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <!-- Pagination -->
      <div class="px-4 py-3 bg-card border-t border-border flex items-center justify-between">
        <div class="flex-1 flex justify-between sm:hidden">
          <button class="px-4 py-2 border border-input rounded-md text-sm font-medium text-gray-700 hover:bg-muted/20">
            Previous
          </button>
          <button class="ml-3 px-4 py-2 border border-input rounded-md text-sm font-medium text-gray-700 hover:bg-muted/20">
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing <span class="font-medium">1</span> to <span class="font-medium">{{ filteredPayments.length }}</span> of <span class="font-medium">{{ filteredPayments.length }}</span> results
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <button class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-input bg-background text-sm font-medium text-gray-500 hover:bg-muted/20">
                <span class="sr-only">Previous</span>
                <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                </svg>
              </button>
              <button class="relative inline-flex items-center px-4 py-2 border border-input bg-primary text-sm font-medium text-white">
                1
              </button>
              <button class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-input bg-background text-sm font-medium text-gray-500 hover:bg-muted/20">
                <span class="sr-only">Next</span>
                <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                </svg>
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Payment Details Modal -->
    <div v-if="showDetails" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-card rounded-lg shadow-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-4 border-b border-border flex justify-between items-center">
          <h3 class="text-lg font-semibold">Payment Details</h3>
          <button @click="closeDetails" class="text-gray-500 hover:text-gray-700">
            <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        <div v-if="selectedPayment" class="p-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div>
              <h4 class="text-sm font-medium text-gray-500">Payment ID</h4>
              <p class="text-lg font-semibold">{{ selectedPayment.id }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500">Status</h4>
              <p>
                <span :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  selectedPayment.status === 'Completed' ? 'bg-green-100 text-green-800' : 
                  selectedPayment.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 
                  'bg-red-100 text-red-800'
                ]">
                  {{ selectedPayment.status }}
                </span>
              </p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500">Invoice ID</h4>
              <p class="font-semibold">{{ selectedPayment.invoiceId }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500">Client</h4>
              <p>{{ selectedPayment.client }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500">Amount</h4>
              <p class="font-semibold">{{ selectedPayment.amount }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500">Date</h4>
              <p>{{ selectedPayment.date }}</p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500">Type</h4>
              <p>
                <span :class="[
                  'px-2 py-1 text-xs font-medium rounded-full',
                  selectedPayment.type === 'Advance' ? 'bg-blue-100 text-blue-800' : 
                  selectedPayment.type === 'Balance' ? 'bg-purple-100 text-purple-800' : 
                  'bg-gray-100 text-gray-800'
                ]">
                  {{ selectedPayment.type }}
                </span>
              </p>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500">Payment Method</h4>
              <p>{{ selectedPayment.method }}</p>
            </div>
          </div>
          
          <div class="border-t border-border pt-4 mb-6">
            <h4 class="text-sm font-medium text-gray-500 mb-2">Transaction Details</h4>
            <p class="text-gray-500 text-sm italic">Transaction details would be displayed here in a real application.</p>
          </div>
          
          <div class="flex flex-col sm:flex-row gap-3 justify-end">
            <button class="px-4 py-2 border border-input rounded-md text-sm font-medium text-gray-700 hover:bg-muted/20">
              Download Receipt
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>