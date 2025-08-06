<script setup>
import { ref } from 'vue';
import { toast } from 'vue-sonner';

// Form data
const invoiceForm = ref({
  clientName: '',
  clientEmail: '',
  invoiceNumber: '',
  amount: '',
  issueDate: '',
  dueDate: '',
  description: '',
  terms: 'net30',
  files: []
});

// Form validation
const errors = ref({});
const isSubmitting = ref(false);

// Client suggestions (would come from API)
const clientSuggestions = ref([
  { name: 'Acme Corp', email: 'billing@acmecorp.com' },
  { name: 'Globex Inc', email: 'accounts@globexinc.com' },
  { name: 'Stark Industries', email: 'finance@stark.com' },
  { name: 'Wayne Enterprises', email: 'payments@wayne.com' },
  { name: 'Umbrella Corp', email: 'invoices@umbrella.com' }
]);

// Payment terms options
const termsOptions = [
  { value: 'net15', label: 'Net 15 - Payment due within 15 days' },
  { value: 'net30', label: 'Net 30 - Payment due within 30 days' },
  { value: 'net45', label: 'Net 45 - Payment due within 45 days' },
  { value: 'net60', label: 'Net 60 - Payment due within 60 days' },
  { value: 'immediate', label: 'Due Immediately' }
];

// Handle file upload
const handleFileUpload = (event) => {
  const files = Array.from(event.target.files);
  invoiceForm.value.files = [...invoiceForm.value.files, ...files];
};

// Remove file
const removeFile = (index) => {
  invoiceForm.value.files.splice(index, 1);
};

// Select client from suggestions
const selectClient = (client) => {
  invoiceForm.value.clientName = client.name;
  invoiceForm.value.clientEmail = client.email;
};

// Validate form
const validateForm = () => {
  const newErrors = {};
  
  if (!invoiceForm.value.clientName) newErrors.clientName = 'Client name is required';
  if (!invoiceForm.value.clientEmail) newErrors.clientEmail = 'Client email is required';
  if (!invoiceForm.value.invoiceNumber) newErrors.invoiceNumber = 'Invoice number is required';
  if (!invoiceForm.value.amount) newErrors.amount = 'Amount is required';
  else if (isNaN(parseFloat(invoiceForm.value.amount))) newErrors.amount = 'Amount must be a valid number';
  if (!invoiceForm.value.issueDate) newErrors.issueDate = 'Issue date is required';
  if (!invoiceForm.value.dueDate) newErrors.dueDate = 'Due date is required';
  if (invoiceForm.value.files.length === 0) newErrors.files = 'At least one invoice file is required';
  
  errors.value = newErrors;
  return Object.keys(newErrors).length === 0;
};

// Submit form
const submitInvoice = () => {
  if (!validateForm()) return;
  
  isSubmitting.value = true;
  
  // Simulate API call
  setTimeout(() => {
    isSubmitting.value = false;
    // Show success message and reset form
    toast.success('Invoice submitted successfully!');
    resetForm();
  }, 1500);
};

// Reset form
const resetForm = () => {
  invoiceForm.value = {
    clientName: '',
    clientEmail: '',
    invoiceNumber: '',
    amount: '',
    issueDate: '',
    dueDate: '',
    description: '',
    terms: 'net30',
    files: []
  };
  errors.value = {};
};

definePageMeta({
  layout: 'dashboard'
});
</script>

<template>
  <div>
    <!-- Page Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Submit Invoice</h1>
      <p class="text-gray-500 mt-1">Upload a new invoice for financing</p>
    </div>

    <!-- Form Container -->
    <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
      <div class="p-4 border-b border-border">
        <h2 class="font-semibold text-lg">Invoice Details</h2>
      </div>
      
      <form @submit.prevent="submitInvoice" class="p-4 md:p-6">
        <!-- Client Information -->
        <div class="mb-6">
          <h3 class="text-sm font-medium text-gray-700 mb-4">Client Information</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label for="clientName" class="block text-sm font-medium text-gray-700 mb-1">Client Name</label>
              <div class="relative">
                <input 
                  id="clientName" 
                  v-model="invoiceForm.clientName" 
                  type="text" 
                  class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                  :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.clientName }"
                  placeholder="Enter client name"
                />
                <div v-if="invoiceForm.clientName && clientSuggestions.length > 0" class="absolute z-10 mt-1 w-full bg-white shadow-lg rounded-md border border-gray-200 max-h-60 overflow-auto">
                  <div 
                    v-for="client in clientSuggestions.filter(c => c.name.toLowerCase().includes(invoiceForm.clientName.toLowerCase()))" 
                    :key="client.name"
                    @click="selectClient(client)"
                    class="px-4 py-2 hover:bg-muted/20 cursor-pointer"
                  >
                    <div class="font-medium">{{ client.name }}</div>
                    <div class="text-xs text-gray-500">{{ client.email }}</div>
                  </div>
                </div>
                <p v-if="errors.clientName" class="mt-1 text-sm text-red-600">{{ errors.clientName }}</p>
              </div>
            </div>
            <div>
              <label for="clientEmail" class="block text-sm font-medium text-gray-700 mb-1">Client Email</label>
              <input 
                id="clientEmail" 
                v-model="invoiceForm.clientEmail" 
                type="email" 
                class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.clientEmail }"
                placeholder="Enter client email"
              />
              <p v-if="errors.clientEmail" class="mt-1 text-sm text-red-600">{{ errors.clientEmail }}</p>
            </div>
          </div>
        </div>
        
        <!-- Invoice Details -->
        <div class="mb-6">
          <h3 class="text-sm font-medium text-gray-700 mb-4">Invoice Details</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label for="invoiceNumber" class="block text-sm font-medium text-gray-700 mb-1">Invoice Number</label>
              <input 
                id="invoiceNumber" 
                v-model="invoiceForm.invoiceNumber" 
                type="text" 
                class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.invoiceNumber }"
                placeholder="e.g. INV-2025-001"
              />
              <p v-if="errors.invoiceNumber" class="mt-1 text-sm text-red-600">{{ errors.invoiceNumber }}</p>
            </div>
            <div>
              <label for="amount" class="block text-sm font-medium text-gray-700 mb-1">Amount (GHC)</label>
              <input 
                id="amount" 
                v-model="invoiceForm.amount" 
                type="text" 
                class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.amount }"
                placeholder="e.g. 1250.00"
              />
              <p v-if="errors.amount" class="mt-1 text-sm text-red-600">{{ errors.amount }}</p>
            </div>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label for="issueDate" class="block text-sm font-medium text-gray-700 mb-1">Issue Date</label>
              <input 
                id="issueDate" 
                v-model="invoiceForm.issueDate" 
                type="date" 
                class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.issueDate }"
              />
              <p v-if="errors.issueDate" class="mt-1 text-sm text-red-600">{{ errors.issueDate }}</p>
            </div>
            <div>
              <label for="dueDate" class="block text-sm font-medium text-gray-700 mb-1">Due Date</label>
              <input 
                id="dueDate" 
                v-model="invoiceForm.dueDate" 
                type="date" 
                class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.dueDate }"
              />
              <p v-if="errors.dueDate" class="mt-1 text-sm text-red-600">{{ errors.dueDate }}</p>
            </div>
          </div>
          
          <div class="mb-4">
            <label for="terms" class="block text-sm font-medium text-gray-700 mb-1">Payment Terms</label>
            <select 
              id="terms" 
              v-model="invoiceForm.terms" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
            >
              <option v-for="option in termsOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
          
          <div>
            <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description (Optional)</label>
            <textarea 
              id="description" 
              v-model="invoiceForm.description" 
              rows="3" 
              class="w-full py-2 px-3 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
              placeholder="Enter invoice description or additional details"
            ></textarea>
          </div>
        </div>
        
        <!-- File Upload -->
        <div class="mb-6">
          <h3 class="text-sm font-medium text-gray-700 mb-4">Invoice Files</h3>
          
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
          <div v-if="invoiceForm.files.length > 0" class="mt-4">
            <h4 class="text-sm font-medium text-gray-700 mb-2">Uploaded Files</h4>
            <ul class="space-y-2">
              <li v-for="(file, index) in invoiceForm.files" :key="index" class="flex items-center justify-between p-2 bg-muted/20 rounded-md">
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
        
        <!-- Form Actions -->
        <div class="flex flex-col sm:flex-row sm:justify-end gap-3 border-t border-border pt-6">
          <button 
            type="button" 
            @click="resetForm" 
            class="px-4 py-2 border border-input rounded-md text-sm font-medium text-gray-700 hover:bg-muted/20"
          >
            Cancel
          </button>
          <button 
            type="submit" 
            class="px-4 py-2 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
            :disabled="isSubmitting"
          >
            <span v-if="isSubmitting" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Submitting...
            </span>
            <span v-else>Submit Invoice</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>