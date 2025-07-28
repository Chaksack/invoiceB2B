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
    name: 'First National Bank', 
    rate: '2.5%', 
    term: '30 days',
    amount: 'GHC 12,250.00',
    description: 'Standard financing with 2.5% fee'
  },
  { 
    id: 2, 
    name: 'Global Finance Partners', 
    rate: '2.0%', 
    term: '45 days',
    amount: 'GHC 12,300.00',
    description: 'Extended term financing with competitive rates'
  },
  { 
    id: 3, 
    name: 'Metropolis Credit Union', 
    rate: '1.8%', 
    term: '30 days',
    amount: 'GHC 12,350.00',
    description: 'Premium client rate with fast processing'
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
const invoices = ref([]);

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
  </div>
</template>