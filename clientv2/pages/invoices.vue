<script setup>
import { ref, computed } from 'vue';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogClose } from '@/components/ui/dialog';
import { toast } from 'vue-sonner';

// Dialog state
const isDialogOpen = ref(false);
const dialogStep = ref('upload'); // 'upload', 'processing', 'financing'
const uploadedFiles = ref([]);
const isProcessing = ref(false);
const financingOptions = ref([
  { 
    id: 1, 
    name: 'Standard Financing Package', 
    rate: '2.5%', 
    term: '30 days',
    amount: 'GHC 12,250.00',
    description: 'Basic financing package with standard rates and terms'
  },
  { 
    id: 2, 
    name: 'Extended Financing Package', 
    rate: '2.0%', 
    term: '45 days',
    amount: 'GHC 12,300.00',
    description: 'Extended term financing with competitive rates and longer payment period'
  },
  { 
    id: 3, 
    name: 'Premium Financing Package', 
    rate: '1.8%', 
    term: '30 days',
    amount: 'GHC 12,350.00',
    description: 'Premium package with lowest rates and expedited processing'
  }
]);
const selectedFinancingOption = ref(null);
const errors = ref({});

// Handle file upload
const handleFileUpload = (event) => {
  const files = Array.from(event.target.files);
  uploadedFiles.value = [...uploadedFiles.value, ...files];
};

// Remove file
const removeFile = (index) => {
  uploadedFiles.value.splice(index, 1);
};

// Process invoice
const processInvoice = () => {
  // Validate
  if (uploadedFiles.value.length === 0) {
    errors.value.files = 'At least one invoice file is required';
    return;
  }
  
  // Clear errors
  errors.value = {};
  
  // Show processing state
  dialogStep.value = 'processing';
  isProcessing.value = true;
  
  // Simulate processing delay
  setTimeout(() => {
    isProcessing.value = false;
    dialogStep.value = 'financing';
  }, 2000);
};

// Select financing option
const selectFinancing = (option) => {
  selectedFinancingOption.value = option;
};

// Complete financing
const completeFinancing = () => {
  if (!selectedFinancingOption.value) {
    errors.value.financing = 'Please select a financing option';
    return;
  }
  
  // Clear errors
  errors.value = {};
  
  // Show processing
  isProcessing.value = true;
  
  // Simulate API call
  setTimeout(() => {
    isProcessing.value = false;
    // Reset and close dialog
    resetDialog();
    isDialogOpen.value = false;
    
    // Show success message
    toast.success('Invoice submitted for financing successfully!');
  }, 1500);
};

// Save financing for later
const saveForLater = () => {
  if (!selectedFinancingOption.value) {
    errors.value.financing = 'Please select a financing option';
    return;
  }
  
  // Clear errors
  errors.value = {};
  
  // Show processing
  isProcessing.value = true;
  
  // Simulate API call
  setTimeout(() => {
    isProcessing.value = false;
    // Close dialog but don't reset (could be implemented to keep state)
    isDialogOpen.value = false;
    
    // Show success message
    toast.success('Financing option saved for later application!');
    
    // In a real implementation, we would save the selected financing option
    // to the user's profile or a database for later retrieval
  }, 1500);
};

// Reset dialog state
const resetDialog = () => {
  dialogStep.value = 'upload';
  uploadedFiles.value = [];
  isProcessing.value = false;
  selectedFinancingOption.value = null;
  errors.value = {};
};

// Open dialog
const openSubmitDialog = () => {
  resetDialog();
  isDialogOpen.value = true;
};

// This data would typically come from an API
const invoices = ref([
  { 
    id: 'INV-2025-001', 
    client: 'Acme Corp', 
    clientAddress: '123 Main St, Metropolis, NY 10001',
    clientEmail: 'billing@acmecorp.com',
    clientPhone: '(555) 123-4567',
    amount: 'GHC 12,500.00', 
    issueDate: '2025-07-25', 
    dueDate: '2025-08-25', 
    status: 'Financed',
    items: [
      { description: 'Web Development Services', quantity: 1, unitPrice: 'GHC 10,000.00', total: 'GHC 10,000.00' },
      { description: 'UI/UX Design', quantity: 1, unitPrice: 'GHC 2,500.00', total: 'GHC 2,500.00' }
    ],
    financingBank: {
      name: 'First National Bank',
      address: '789 Finance Ave, Metropolis, NY 10001',
      accountNumber: 'FNBX-7890-1234',
      contactPerson: 'Jane Smith',
      contactEmail: 'j.smith@fnb.com',
      contactPhone: '(555) 987-6543',
      financingTerms: '30-day net financing with 2.5% fee'
    }
  },
  { 
    id: 'INV-2025-002', 
    client: 'Globex Inc', 
    clientAddress: '456 Tech Blvd, Silicon Valley, CA 94025',
    clientEmail: 'accounts@globex.com',
    clientPhone: '(555) 234-5678',
    amount: 'GHC 8,750.00', 
    issueDate: '2025-07-24', 
    dueDate: '2025-08-24', 
    status: 'Pending',
    items: [
      { description: 'Mobile App Development', quantity: 1, unitPrice: 'GHC 7,500.00', total: 'GHC 7,500.00' },
      { description: 'Technical Support (10 hours)', quantity: 10, unitPrice: 'GHC 125.00', total: 'GHC 1,250.00' }
    ]
  },
  { 
    id: 'INV-2025-003', 
    client: 'Stark Industries', 
    clientAddress: '200 Park Avenue, New York, NY 10017',
    clientEmail: 'finance@stark.com',
    clientPhone: '(555) 345-6789',
    amount: 'GHC 15,200.00', 
    issueDate: '2025-07-22', 
    dueDate: '2025-08-22', 
    status: 'Financed',
    items: [
      { description: 'AI Integration Services', quantity: 1, unitPrice: 'GHC 12,000.00', total: 'GHC 12,000.00' },
      { description: 'Hardware Components', quantity: 8, unitPrice: 'GHC 400.00', total: 'GHC 3,200.00' }
    ],
    financingBank: {
      name: 'Metropolis Credit Union',
      address: '100 Finance Street, Metropolis, NY 10001',
      accountNumber: 'MCU-1234-5678',
      contactPerson: 'Robert Johnson',
      contactEmail: 'r.johnson@mcu.com',
      contactPhone: '(555) 456-7890',
      financingTerms: '45-day financing with 2% fee'
    }
  },
  { 
    id: 'INV-2025-004', 
    client: 'Wayne Enterprises', 
    clientAddress: '1007 Mountain Drive, Gotham, NJ 08302',
    clientEmail: 'accounting@wayne.com',
    clientPhone: '(555) 456-7890',
    amount: 'GHC 9,800.00', 
    issueDate: '2025-07-20', 
    dueDate: '2025-08-20', 
    status: 'Paid',
    items: [
      { description: 'Security System Implementation', quantity: 1, unitPrice: 'GHC 8,500.00', total: 'GHC 8,500.00' },
      { description: 'Staff Training', quantity: 13, unitPrice: 'GHC 100.00', total: 'GHC 1,300.00' }
    ]
  },
  { 
    id: 'INV-2025-005', 
    client: 'Umbrella Corp', 
    clientAddress: '544 Raccoon Ave, Raccoon City, CA 90210',
    clientEmail: 'finance@umbrella.com',
    clientPhone: '(555) 567-8901',
    amount: 'GHC 11,350.00', 
    issueDate: '2025-07-18', 
    dueDate: '2025-08-18', 
    status: 'Financed',
    items: [
      { description: 'Biotech Consulting Services', quantity: 1, unitPrice: 'GHC 9,500.00', total: 'GHC 9,500.00' },
      { description: 'Laboratory Equipment', quantity: 1, unitPrice: 'GHC 1,850.00', total: 'GHC 1,850.00' }
    ],
    financingBank: {
      name: 'Global Finance Partners',
      address: '555 Banking Lane, New York, NY 10004',
      accountNumber: 'GFP-5678-9012',
      contactPerson: 'Michael Chen',
      contactEmail: 'm.chen@gfp.com',
      contactPhone: '(555) 678-9012',
      financingTerms: '60-day financing with 3% fee'
    }
  },
  { 
    id: 'INV-2025-006', 
    client: 'Oscorp Industries', 
    clientAddress: '120 Science Blvd, New York, NY 10016',
    clientEmail: 'billing@oscorp.com',
    clientPhone: '(555) 678-9012',
    amount: 'GHC 7,200.00', 
    issueDate: '2025-07-15', 
    dueDate: '2025-08-15', 
    status: 'Pending',
    items: [
      { description: 'Research & Development Support', quantity: 1, unitPrice: 'GHC 6,000.00', total: 'GHC 6,000.00' },
      { description: 'Technical Documentation', quantity: 12, unitPrice: 'GHC 100.00', total: 'GHC 1,200.00' }
    ]
  },
  { 
    id: 'INV-2025-007', 
    client: 'LexCorp', 
    clientAddress: '1000 Corporate Drive, Metropolis, DE 19901',
    clientEmail: 'accounts@lexcorp.com',
    clientPhone: '(555) 789-0123',
    amount: 'GHC 18,900.00', 
    issueDate: '2025-07-12', 
    dueDate: '2025-08-12', 
    status: 'Financed',
    items: [
      { description: 'Enterprise Software Implementation', quantity: 1, unitPrice: 'GHC 15,000.00', total: 'GHC 15,000.00' },
      { description: 'Custom Development', quantity: 26, unitPrice: 'GHC 150.00', total: 'GHC 3,900.00' }
    ],
    financingBank: {
      name: 'Metropolis Bank & Trust',
      address: '800 Financial District, Metropolis, DE 19901',
      accountNumber: 'MBT-9012-3456',
      contactPerson: 'Sarah Williams',
      contactEmail: 's.williams@mbt.com',
      contactPhone: '(555) 890-1234',
      financingTerms: '30-day financing with 2.75% fee'
    }
  },
  { 
    id: 'INV-2025-008', 
    client: 'Cyberdyne Systems', 
    clientAddress: '18144 El Camino Real, Sunnyvale, CA 94087',
    clientEmail: 'payments@cyberdyne.com',
    clientPhone: '(555) 890-1234',
    amount: 'GHC 14,500.00', 
    issueDate: '2025-07-10', 
    dueDate: '2025-08-10', 
    status: 'Paid',
    items: [
      { description: 'AI Algorithm Development', quantity: 1, unitPrice: 'GHC 12,000.00', total: 'GHC 12,000.00' },
      { description: 'System Integration', quantity: 25, unitPrice: 'GHC 100.00', total: 'GHC 2,500.00' }
    ]
  },
  { 
    id: 'INV-2025-009', 
    client: 'Massive Dynamic', 
    clientAddress: '1 Massive Dynamic Way, New York, NY 10018',
    clientEmail: 'accounting@massivedynamic.com',
    clientPhone: '(555) 901-2345',
    amount: 'GHC 22,750.00', 
    issueDate: '2025-07-08', 
    dueDate: '2025-08-08', 
    status: 'Financed',
    items: [
      { description: 'Quantum Computing Research', quantity: 1, unitPrice: 'GHC 20,000.00', total: 'GHC 20,000.00' },
      { description: 'Specialized Equipment', quantity: 1, unitPrice: 'GHC 2,750.00', total: 'GHC 2,750.00' }
    ],
    financingBank: {
      name: 'Atlantic Financial Services',
      address: '200 Wall Street, New York, NY 10005',
      accountNumber: 'AFS-3456-7890',
      contactPerson: 'David Thompson',
      contactEmail: 'd.thompson@afs.com',
      contactPhone: '(555) 012-3456',
      financingTerms: '45-day financing with 2.25% fee'
    }
  },
  { 
    id: 'INV-2025-010', 
    client: 'Soylent Corp', 
    clientAddress: '101 Green Street, New York, NY 10013',
    clientEmail: 'finance@soylent.com',
    clientPhone: '(555) 012-3456',
    amount: 'GHC 5,800.00', 
    issueDate: '2025-07-05', 
    dueDate: '2025-08-05', 
    status: 'Rejected',
    items: [
      { description: 'Food Processing Consultation', quantity: 1, unitPrice: 'GHC 4,500.00', total: 'GHC 4,500.00' },
      { description: 'Market Research', quantity: 13, unitPrice: 'GHC 100.00', total: 'GHC 1,300.00' }
    ]
  }
]);

// Filter and sort state
const searchQuery = ref('');
const statusFilter = ref('All');
const sortBy = ref('issueDate');
const sortDirection = ref('desc');

// Status options for filter
const statusOptions = ['All', 'Pending', 'Financed', 'Paid', 'Rejected'];

// Computed filtered and sorted invoices
const filteredInvoices = computed(() => {
  let result = [...invoices.value];
  
  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(invoice => 
      invoice.id.toLowerCase().includes(query) || 
      invoice.client.toLowerCase().includes(query)
    );
  }
  
  // Apply status filter
  if (statusFilter.value !== 'All') {
    result = result.filter(invoice => invoice.status === statusFilter.value);
  }
  
  // Apply sorting
  result.sort((a, b) => {
    let comparison = 0;
    if (sortBy.value === 'amount') {
      // Remove GHC and commas for numeric comparison
      const amountA = parseFloat(a.amount.replace('GHC', '').replace(',', ''));
      const amountB = parseFloat(b.amount.replace('GHC', '').replace(',', ''));
      comparison = amountA - amountB;
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

// Selected invoice for details view
const selectedInvoice = ref(null);

const viewInvoiceDetails = (invoice) => {
  selectedInvoice.value = invoice;
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
        <h1 class="text-2xl font-bold text-gray-900">Invoices</h1>
        <p class="text-gray-500 mt-1">Manage and track all your invoices</p>
      </div>
      <div class="mt-4 md:mt-0">
        <button 
          @click="openSubmitDialog"
          class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary-700 transition-colors flex items-center"
        >
          <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="mr-2">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          Submit New Invoice
        </button>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="bg-card rounded-lg shadow-sm border border-border p-4 mb-6">
      <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div class="relative flex-1">
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Search invoices..." 
            class="w-full py-2 pl-10 pr-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
          />
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
        </div>
        
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="w-full sm:w-auto">
            <select 
              v-model="statusFilter"
              class="w-full py-2 px-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            >
              <option v-for="status in statusOptions" :key="status" :value="status">
                {{ status }} Invoices
              </option>
            </select>
          </div>
          
          <div class="flex gap-2">
            <button class="px-3 py-2 border border-input rounded-md hover:bg-muted/20 transition-colors">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"></path>
              </svg>
            </button>
            <button class="px-3 py-2 border border-input rounded-md hover:bg-muted/20 transition-colors">
              <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Two-column layout: Invoice List and Details -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left Column: Invoice List -->
      <div class="lg:col-span-5 bg-card rounded-lg shadow-sm border border-border overflow-hidden">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="bg-muted/50">
                <th 
                  @click="toggleSort('id')" 
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70"
                >
                  <div class="flex items-center">
                    Invoice ID
                    <svg v-if="sortBy === 'id'" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="ml-1">
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
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-muted/70 hidden md:table-cell"
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
              <tr 
                v-for="invoice in filteredInvoices" 
                :key="invoice.id" 
                class="hover:bg-muted/20 cursor-pointer"
                :class="{'bg-primary-50 border-l-4 border-primary font-medium': selectedInvoice && selectedInvoice.id === invoice.id}"
                @click="viewInvoiceDetails(invoice)"
              >
                <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ invoice.id }}</td>
                <td class="px-4 py-3 text-sm text-gray-500">{{ invoice.client }}</td>
                <td class="px-4 py-3 text-sm text-gray-500 hidden md:table-cell">{{ invoice.amount }}</td>
                <td class="px-4 py-3 text-sm">
                  <span :class="[
                    'px-2 py-1 text-xs font-medium rounded-full',
                    invoice.status === 'Financed' ? 'bg-blue-100 text-blue-800' : 
                    invoice.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 
                    invoice.status === 'Paid' ? 'bg-green-100 text-green-800' :
                    'bg-red-100 text-red-800'
                  ]">
                    {{ invoice.status }}
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-right">
                  <span class="text-primary hover:text-primary-700 font-medium">
                    View
                  </span>
                </td>
              </tr>
              <tr v-if="filteredInvoices.length === 0">
                <td colspan="5" class="px-4 py-6 text-center text-gray-500">
                  No invoices found matching your filters.
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
                Showing <span class="font-medium">1</span> to <span class="font-medium">{{ filteredInvoices.length }}</span> of <span class="font-medium">{{ filteredInvoices.length }}</span> results
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

      <!-- Right Column: Invoice Details -->
      <div class="lg:col-span-7">
        <div v-if="selectedInvoice" class="bg-card rounded-lg shadow-sm border border-border h-full overflow-auto">
          <!-- Professional Invoice Layout -->
          <div class="p-6 print:p-0">
            <!-- Invoice Header -->
            <div class="flex justify-between items-start mb-8 border-b border-gray-200 pb-4">
              <div>
                <div class="flex items-center mb-2">
                  <div class="bg-primary h-10 w-10 rounded-md flex items-center justify-center text-white mr-3">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div>
                    <h2 class="text-xl font-bold text-gray-900">InvoiceB2B Finance</h2>
                    <p class="text-sm text-gray-500">123 Finance Street, Suite 100</p>
                    <p class="text-sm text-gray-500">New York, NY 10004</p>
                  </div>
                </div>
                <p class="text-sm text-gray-500">contact@invoiceb2b.com | (555) 123-4567</p>
              </div>
              <div class="text-right">
                <h1 class="text-2xl font-bold text-gray-900 mb-1">INVOICE</h1>
                <p class="text-lg font-semibold text-gray-700">{{ selectedInvoice.id }}</p>
                <div class="mt-2">
                  <span :class="[
                    'px-3 py-1 text-sm font-medium rounded-full',
                    selectedInvoice.status === 'Financed' ? 'bg-blue-100 text-blue-800' : 
                    selectedInvoice.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 
                    selectedInvoice.status === 'Paid' ? 'bg-green-100 text-green-800' :
                    'bg-red-100 text-red-800'
                  ]">
                    {{ selectedInvoice.status }}
                  </span>
                </div>
              </div>
            </div>
            
            <!-- Invoice Information -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
              <!-- Client Information -->
              <div>
                <h3 class="text-sm font-medium text-gray-500 uppercase mb-3">Bill To:</h3>
                <h4 class="text-base font-bold text-gray-900 mb-1">{{ selectedInvoice.client }}</h4>
                <p class="text-sm text-gray-700 mb-1">{{ selectedInvoice.clientAddress }}</p>
                <p class="text-sm text-gray-700 mb-1">{{ selectedInvoice.clientEmail }}</p>
                <p class="text-sm text-gray-700">{{ selectedInvoice.clientPhone }}</p>
              </div>
              
              <!-- Invoice Details -->
              <div class="md:text-right">
                <div class="mb-2">
                  <h3 class="text-sm font-medium text-gray-500 uppercase mb-1">Invoice Date:</h3>
                  <p class="text-sm text-gray-900">{{ selectedInvoice.issueDate }}</p>
                </div>
                <div class="mb-2">
                  <h3 class="text-sm font-medium text-gray-500 uppercase mb-1">Due Date:</h3>
                  <p class="text-sm text-gray-900">{{ selectedInvoice.dueDate }}</p>
                </div>
                <div>
                  <h3 class="text-sm font-medium text-gray-500 uppercase mb-1">Total Amount:</h3>
                  <p class="text-lg font-bold text-gray-900">{{ selectedInvoice.amount }}</p>
                </div>
              </div>
            </div>
            
            <!-- Invoice Items -->
            <div class="mb-8">
              <h3 class="text-base font-semibold text-gray-900 mb-3">Invoice Items</h3>
              <div class="overflow-x-auto">
                <table class="w-full">
                  <thead>
                    <tr class="bg-gray-50 text-left">
                      <th class="px-4 py-2 text-xs font-medium text-gray-500 uppercase">Description</th>
                      <th class="px-4 py-2 text-xs font-medium text-gray-500 uppercase text-center">Quantity</th>
                      <th class="px-4 py-2 text-xs font-medium text-gray-500 uppercase text-right">Unit Price</th>
                      <th class="px-4 py-2 text-xs font-medium text-gray-500 uppercase text-right">Total</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-gray-200">
                    <tr v-for="(item, index) in selectedInvoice.items" :key="index" class="text-sm">
                      <td class="px-4 py-3 text-gray-900">{{ item.description }}</td>
                      <td class="px-4 py-3 text-gray-700 text-center">{{ item.quantity }}</td>
                      <td class="px-4 py-3 text-gray-700 text-right">{{ item.unitPrice }}</td>
                      <td class="px-4 py-3 font-medium text-gray-900 text-right">{{ item.total }}</td>
                    </tr>
                  </tbody>
                  <tfoot>
                    <tr class="bg-gray-50">
                      <td colspan="3" class="px-4 py-3 text-sm font-medium text-gray-500 text-right">Total:</td>
                      <td class="px-4 py-3 text-base font-bold text-gray-900 text-right">{{ selectedInvoice.amount }}</td>
                    </tr>
                  </tfoot>
                </table>
              </div>
            </div>
            
            <!-- Financing Information (if applicable) -->
            <div v-if="selectedInvoice.status === 'Financed' && selectedInvoice.financingBank" class="mb-8 bg-blue-50 p-4 rounded-lg border border-blue-200">
              <h3 class="text-base font-semibold text-blue-800 mb-3">Financing Information</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <h4 class="text-sm font-medium text-blue-700 mb-1">Financing Bank:</h4>
                  <p class="text-sm text-gray-800 mb-3">{{ selectedInvoice.financingBank.name }}</p>
                  
                  <h4 class="text-sm font-medium text-blue-700 mb-1">Bank Address:</h4>
                  <p class="text-sm text-gray-800 mb-3">{{ selectedInvoice.financingBank.address }}</p>
                  
                  <h4 class="text-sm font-medium text-blue-700 mb-1">Account Number:</h4>
                  <p class="text-sm text-gray-800">{{ selectedInvoice.financingBank.accountNumber }}</p>
                </div>
                <div>
                  <h4 class="text-sm font-medium text-blue-700 mb-1">Contact Person:</h4>
                  <p class="text-sm text-gray-800 mb-1">{{ selectedInvoice.financingBank.contactPerson }}</p>
                  <p class="text-sm text-gray-800 mb-1">{{ selectedInvoice.financingBank.contactEmail }}</p>
                  <p class="text-sm text-gray-800 mb-3">{{ selectedInvoice.financingBank.contactPhone }}</p>
                  
                  <h4 class="text-sm font-medium text-blue-700 mb-1">Financing Terms:</h4>
                  <p class="text-sm text-gray-800">{{ selectedInvoice.financingBank.financingTerms }}</p>
                </div>
              </div>
            </div>
            
            <!-- Payment Terms & Notes -->
            <div class="mb-8">
              <h3 class="text-base font-semibold text-gray-900 mb-2">Payment Terms</h3>
              <p class="text-sm text-gray-700 mb-4">
                Payment is due by {{ selectedInvoice.dueDate }}. Please include the invoice number with your payment.
                For questions concerning this invoice, please contact our accounts department.
              </p>
              
              <h3 class="text-base font-semibold text-gray-900 mb-2">Notes</h3>
              <p class="text-sm text-gray-700">
                Thank you for your business! We appreciate the opportunity to work with {{ selectedInvoice.client }}.
              </p>
            </div>
            
            <!-- Footer -->
            <div class="border-t border-gray-200 pt-4 text-center text-xs text-gray-500">
              <p>InvoiceB2B Finance, Inc. | Tax ID: XX-XXXXXXX</p>
              <p>This is a computer-generated document. No signature is required.</p>
            </div>
            
            <!-- Action Buttons -->
            <div class="mt-6 flex flex-col sm:flex-row gap-3 justify-end">
              <button class="px-4 py-2 border border-input rounded-md text-sm font-medium text-gray-700 hover:bg-muted/20">
                Download PDF
              </button>
              <button v-if="selectedInvoice.status === 'Pending'" class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700">
                Apply for Financing
              </button>
            </div>
          </div>
        </div>
        <div v-else class="bg-card rounded-lg shadow-sm border border-border h-full flex items-center justify-center p-8">
          <div class="text-center">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No invoice selected</h3>
            <p class="mt-1 text-sm text-gray-500">Select an invoice from the list to view its details.</p>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Invoice Upload Dialog -->
  <Dialog v-model:open="isDialogOpen" @update:open="(val) => !val && resetDialog()">
    <DialogContent class="sm:max-w-md md:max-w-lg">
      <DialogHeader>
        <DialogTitle>
          <span v-if="dialogStep === 'upload'">Submit New Invoice</span>
          <span v-else-if="dialogStep === 'processing'">Processing Invoice</span>
          <span v-else-if="dialogStep === 'financing'">Financing Options</span>
        </DialogTitle>
      </DialogHeader>
      
      <!-- Upload Step -->
      <div v-if="dialogStep === 'upload'" class="py-6">
        <p class="text-sm text-gray-500 mb-6">
          Upload your invoice file to get financing options from our partner banks.
        </p>
        
        <!-- File Upload -->
        <div class="mb-6">
          <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center" :class="{ 'border-red-500': errors.files }">
            <input 
              type="file" 
              id="fileUpload" 
              @change="handleFileUpload" 
              multiple 
              class="hidden" 
              accept=".pdf,.jpg,.jpeg,.png,.doc,.docx,.xls,.xlsx"
            />
            <label for="fileUpload" class="cursor-pointer">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
              </svg>
              <p class="mt-2 text-sm text-gray-600">
                <span class="font-medium text-primary hover:text-primary-700">Click to upload</span> or drag and drop
              </p>
              <p class="mt-1 text-xs text-gray-500">
                PDF, JPG, PNG, DOC, DOCX, XLS, XLSX up to 10MB
              </p>
            </label>
            <p v-if="errors.files" class="mt-2 text-sm text-red-600">{{ errors.files }}</p>
          </div>
          
          <!-- File List -->
          <div v-if="uploadedFiles.length > 0" class="mt-4">
            <h4 class="text-sm font-medium text-gray-700 mb-2">Uploaded Files</h4>
            <ul class="space-y-2">
              <li v-for="(file, index) in uploadedFiles" :key="index" class="flex items-center justify-between p-2 bg-muted/20 rounded-md">
                <div class="flex items-center">
                  <svg class="h-5 w-5 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                  </svg>
                  <span class="text-sm truncate max-w-xs">{{ file.name }}</span>
                </div>
                <button 
                  type="button" 
                  @click="removeFile(index)" 
                  class="text-gray-500 hover:text-red-500"
                >
                  <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </li>
            </ul>
          </div>
        </div>
        
        <div class="flex justify-end gap-3">
          <DialogClose class="px-4 py-2 border border-input rounded-md text-sm font-medium text-gray-700 hover:bg-muted/20">
            Cancel
          </DialogClose>
          <button 
            @click="processInvoice"
            class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700"
          >
            Continue
          </button>
        </div>
      </div>
      
      <!-- Processing Step -->
      <div v-else-if="dialogStep === 'processing'" class="py-6">
        <div class="flex flex-col items-center justify-center py-10">
          <svg class="animate-spin h-12 w-12 text-primary mb-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Processing Your Invoice</h3>
          <p class="text-sm text-gray-500 text-center max-w-sm">
            We're analyzing your invoice and finding the best financing options for you. This will only take a moment.
          </p>
        </div>
      </div>
      
      <!-- Financing Options Step -->
      <div v-else-if="dialogStep === 'financing'" class="py-6">
        <p class="text-sm text-gray-500 mb-6">
          Select a financing option that works best for your business needs.
        </p>
        
        <div class="space-y-4 mb-6">
          <div 
            v-for="option in financingOptions" 
            :key="option.id"
            @click="selectFinancing(option)"
            class="border rounded-lg p-4 cursor-pointer transition-colors"
            :class="selectedFinancingOption && selectedFinancingOption.id === option.id ? 'border-primary bg-primary/5' : 'border-border hover:border-primary/50'"
          >
            <div class="flex justify-between items-start">
              <div>
                <h3 class="font-medium text-gray-900">{{ option.name }}</h3>
                <p class="text-sm text-gray-500 mt-1">{{ option.description }}</p>
              </div>
              <div class="text-right">
                <p class="font-bold text-gray-900">{{ option.amount }}</p>
                <p class="text-sm text-gray-500">{{ option.rate }} / {{ option.term }}</p>
              </div>
            </div>
          </div>
        </div>
        
        <p v-if="errors.financing" class="text-sm text-red-600 mb-4">{{ errors.financing }}</p>
        
        <div class="flex justify-end gap-3">
          <button 
            @click="dialogStep = 'upload'"
            class="px-4 py-2 border border-input rounded-md text-sm font-medium text-gray-700 hover:bg-muted/20"
          >
            Back
          </button>
          <button 
            @click="saveForLater"
            class="px-4 py-2 border border-primary text-primary rounded-md text-sm font-medium hover:bg-primary/10"
            :disabled="isProcessing"
          >
            <span v-if="isProcessing" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-primary" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Processing...
            </span>
            <span v-else>Save for Later</span>
          </button>
          <button 
            @click="completeFinancing"
            class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700"
            :disabled="isProcessing"
          >
            <span v-if="isProcessing" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Processing...
            </span>
            <span v-else>Complete Financing</span>
          </button>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>