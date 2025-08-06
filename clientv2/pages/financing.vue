<script setup>
import { ref, computed } from 'vue';
import { toast } from 'vue-sonner';

// This data would typically come from an API
const invoices = ref([
  { id: 'INV-2025-001', client: 'Acme Corp', amount: 'GHC 12,500.00', dueDate: '2025-08-25', status: 'Pending', selected: false },
  { id: 'INV-2025-002', client: 'Globex Inc', amount: 'GHC 8,750.00', dueDate: '2025-08-24', status: 'Pending', selected: false },
  { id: 'INV-2025-006', client: 'Oscorp Industries', amount: 'GHC 7,200.00', dueDate: '2025-08-15', status: 'Pending', selected: false },
]);

// Financing options
const financingOptions = ref([
  { 
    id: 1, 
    name: 'Standard Financing', 
    description: 'Get up to 80% of your invoice value upfront with the remaining balance paid when your customer pays.', 
    advanceRate: 0.8, 
    fee: 0.03, 
    processingTime: '1-2 business days',
    selected: true
  },
  { 
    id: 2, 
    name: 'Express Financing', 
    description: 'Get up to 75% of your invoice value within 24 hours with the remaining balance paid when your customer pays.', 
    advanceRate: 0.75, 
    fee: 0.04, 
    processingTime: 'Same day',
    selected: false
  },
  { 
    id: 3, 
    name: 'Full Financing', 
    description: 'Get up to 90% of your invoice value upfront with the remaining balance paid when your customer pays.', 
    advanceRate: 0.9, 
    fee: 0.05, 
    processingTime: '2-3 business days',
    selected: false
  }
]);

// Form data
const financingForm = ref({
  bankName: '',
  accountNumber: '',
  routingNumber: '',
  accountType: 'checking',
  termsAccepted: false
});

// Form validation
const errors = ref({});
const isSubmitting = ref(false);

// Account type options
const accountTypes = [
  { value: 'checking', label: 'Checking Account' },
  { value: 'savings', label: 'Savings Account' }
];

// Select financing option
const selectFinancingOption = (optionId) => {
  financingOptions.value.forEach(option => {
    option.selected = option.id === optionId;
  });
};

// Toggle invoice selection
const toggleInvoiceSelection = (invoiceId) => {
  const invoice = invoices.value.find(inv => inv.id === invoiceId);
  if (invoice) {
    invoice.selected = !invoice.selected;
  }
};

// Select all invoices
const selectAllInvoices = () => {
  const allSelected = invoices.value.every(invoice => invoice.selected);
  invoices.value.forEach(invoice => {
    invoice.selected = !allSelected;
  });
};

// Computed properties
const selectedInvoices = computed(() => {
  return invoices.value.filter(invoice => invoice.selected);
});

const totalInvoiceAmount = computed(() => {
  return selectedInvoices.value.reduce((total, invoice) => {
    // Remove GHC and commas for calculation
    const amount = parseFloat(invoice.amount.replace('GHC', '').replace(',', ''));
    return total + amount;
  }, 0).toFixed(2);
});

const selectedFinancingOption = computed(() => {
  return financingOptions.value.find(option => option.selected);
});

const advanceAmount = computed(() => {
  if (!selectedFinancingOption.value) return 0;
  return (parseFloat(totalInvoiceAmount.value) * selectedFinancingOption.value.advanceRate).toFixed(2);
});

const financingFee = computed(() => {
  if (!selectedFinancingOption.value) return 0;
  return (parseFloat(totalInvoiceAmount.value) * selectedFinancingOption.value.fee).toFixed(2);
});

const netAmount = computed(() => {
  return (parseFloat(advanceAmount.value) - parseFloat(financingFee.value)).toFixed(2);
});

// Validate form
const validateForm = () => {
  const newErrors = {};
  
  if (selectedInvoices.value.length === 0) {
    newErrors.invoices = 'Please select at least one invoice to finance';
  }
  
  if (!financingForm.value.bankName) {
    newErrors.bankName = 'Bank name is required';
  }
  
  if (!financingForm.value.accountNumber) {
    newErrors.accountNumber = 'Account number is required';
  } else if (!/^\d{8,17}$/.test(financingForm.value.accountNumber)) {
    newErrors.accountNumber = 'Please enter a valid account number (8-17 digits)';
  }
  
  if (!financingForm.value.routingNumber) {
    newErrors.routingNumber = 'Routing number is required';
  } else if (!/^\d{9}$/.test(financingForm.value.routingNumber)) {
    newErrors.routingNumber = 'Please enter a valid 9-digit routing number';
  }
  
  if (!financingForm.value.termsAccepted) {
    newErrors.termsAccepted = 'You must accept the terms and conditions';
  }
  
  errors.value = newErrors;
  return Object.keys(newErrors).length === 0;
};

// Submit form
const submitFinancingRequest = () => {
  if (!validateForm()) return;
  
  isSubmitting.value = true;
  
  // Simulate API call
  setTimeout(() => {
    isSubmitting.value = false;
    // Show success message
    toast.success('Financing request submitted successfully!');
  }, 1500);
};

definePageMeta({
  layout: 'dashboard'
});
</script>

<template>
  <div>
    <!-- Page Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Apply for Financing</h1>
      <p class="text-gray-500 mt-1">Select invoices and financing options to get immediate working capital</p>
    </div>

    <!-- Main Content -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Left Column: Invoices Selection -->
      <div class="lg:col-span-2">
        <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden mb-6">
          <div class="p-4 border-b border-border flex flex-wrap justify-between items-center gap-2">
            <h2 class="font-semibold text-lg">Select Invoices to Finance</h2>
            <button 
              @click="selectAllInvoices"
              class="text-sm text-primary hover:underline"
            >
              {{ invoices.every(invoice => invoice.selected) ? 'Deselect All' : 'Select All' }}
            </button>
          </div>
          
          <div v-if="errors.invoices" class="bg-red-50 p-3 border-b border-red-200">
            <p class="text-sm text-red-600">{{ errors.invoices }}</p>
          </div>
          
          <!-- Desktop Table View -->
          <div class="hidden md:block overflow-x-auto">
            <table class="w-full">
              <thead>
                <tr class="bg-muted/50">
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-10">
                    <input 
                      type="checkbox" 
                      :checked="invoices.every(invoice => invoice.selected)"
                      @change="selectAllInvoices"
                      class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
                    />
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Invoice ID</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Client</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Due Date</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border">
                <tr v-for="invoice in invoices" :key="invoice.id" class="hover:bg-muted/20">
                  <td class="px-4 py-3 text-sm">
                    <input 
                      type="checkbox" 
                      :checked="invoice.selected"
                      @change="toggleInvoiceSelection(invoice.id)"
                      class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
                    />
                  </td>
                  <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ invoice.id }}</td>
                  <td class="px-4 py-3 text-sm text-gray-500">{{ invoice.client }}</td>
                  <td class="px-4 py-3 text-sm text-gray-500">{{ invoice.amount }}</td>
                  <td class="px-4 py-3 text-sm text-gray-500">{{ invoice.dueDate }}</td>
                </tr>
                <tr v-if="invoices.length === 0">
                  <td colspan="5" class="px-4 py-6 text-center text-gray-500">
                    No invoices available for financing.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          
          <!-- Mobile Card View -->
          <div class="md:hidden divide-y divide-border">
            <div 
              v-for="invoice in invoices" 
              :key="invoice.id" 
              class="p-4 hover:bg-muted/20"
            >
              <div class="flex items-start mb-2">
                <div class="mr-3">
                  <input 
                    type="checkbox" 
                    :checked="invoice.selected"
                    @change="toggleInvoiceSelection(invoice.id)"
                    class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded mt-1"
                  />
                </div>
                <div class="flex-1">
                  <div class="font-medium text-gray-900">{{ invoice.id }}</div>
                  <div class="grid grid-cols-2 gap-2 mt-2 text-sm">
                    <div>
                      <div class="text-gray-500 font-medium">Client</div>
                      <div>{{ invoice.client }}</div>
                    </div>
                    <div>
                      <div class="text-gray-500 font-medium">Amount</div>
                      <div>{{ invoice.amount }}</div>
                    </div>
                    <div>
                      <div class="text-gray-500 font-medium">Due Date</div>
                      <div>{{ invoice.dueDate }}</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-if="invoices.length === 0" class="p-6 text-center text-gray-500">
              No invoices available for financing.
            </div>
          </div>
          
          <div v-if="selectedInvoices.length > 0" class="p-4 border-t border-border bg-muted/10">
            <div class="flex justify-between items-center">
              <div>
                <span class="text-sm text-gray-500">Selected:</span>
                <span class="ml-1 text-sm font-medium">{{ selectedInvoices.length }} invoice(s)</span>
              </div>
              <div>
                <span class="text-sm text-gray-500">Total Amount:</span>
                <span class="ml-1 text-sm font-bold">GHC {{ totalInvoiceAmount }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Financing Options -->
        <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
          <div class="p-4 border-b border-border">
            <h2 class="font-semibold text-lg">Select Financing Option</h2>
          </div>
          
          <div class="p-4 space-y-4">
            <div 
              v-for="option in financingOptions" 
              :key="option.id"
              class="border rounded-lg p-4 cursor-pointer"
              :class="{ 'border-primary bg-primary/5': option.selected, 'border-border': !option.selected }"
              @click="selectFinancingOption(option.id)"
            >
              <div class="flex items-start">
                <div class="flex-shrink-0 mt-0.5">
                  <div class="w-6 h-6 rounded-full border-2 flex items-center justify-center" :class="{ 'border-primary': option.selected, 'border-gray-300': !option.selected }">
                    <div v-if="option.selected" class="w-3 h-3 rounded-full bg-primary"></div>
                  </div>
                </div>
                <div class="ml-3 flex-1">
                  <div class="flex flex-wrap justify-between gap-2">
                    <h3 class="text-base font-medium text-gray-900">{{ option.name }}</h3>
                    <div class="text-sm font-medium text-primary">{{ (option.advanceRate * 100) }}% advance</div>
                  </div>
                  <p class="mt-1 text-sm text-gray-500">{{ option.description }}</p>
                  <div class="mt-2 flex flex-wrap gap-x-4 gap-y-2 text-xs text-gray-500">
                    <div>
                      <span class="font-medium">Fee:</span> {{ option.fee * 100 }}%
                    </div>
                    <div>
                      <span class="font-medium">Processing Time:</span> {{ option.processingTime }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Right Column: Summary and Payment Details -->
      <div>
        <!-- Financing Summary -->
        <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden mb-6">
          <div class="p-4 border-b border-border">
            <h2 class="font-semibold text-lg">Financing Summary</h2>
          </div>
          
          <div class="p-4 space-y-4">
            <div class="flex flex-wrap justify-between items-center gap-2">
              <span class="text-sm text-gray-500">Total Invoice Amount:</span>
              <span class="text-sm font-bold">GHC {{ totalInvoiceAmount }}</span>
            </div>
            <div class="flex flex-wrap justify-between items-center gap-2">
              <span class="text-sm text-gray-500">Advance Rate:</span>
              <span class="text-sm font-medium">{{ selectedFinancingOption ? (selectedFinancingOption.advanceRate * 100) + '%' : '0%' }}</span>
            </div>
            <div class="flex flex-wrap justify-between items-center gap-2">
              <span class="text-sm text-gray-500">Advance Amount:</span>
              <span class="text-sm font-medium">GHC {{ advanceAmount }}</span>
            </div>
            <div class="flex flex-wrap justify-between items-center gap-2">
              <span class="text-sm text-gray-500">Financing Fee:</span>
              <span class="text-sm font-medium text-red-500">-GHC {{ financingFee }}</span>
            </div>
            <div class="pt-3 border-t border-border flex flex-wrap justify-between items-center gap-2">
              <span class="text-sm font-medium text-gray-700">Net Amount:</span>
              <span class="text-base font-bold text-primary">GHC {{ netAmount }}</span>
            </div>
          </div>
        </div>
        
        <!-- Payment Details -->
        <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
          <div class="p-4 border-b border-border">
            <h2 class="font-semibold text-lg">Payment Details</h2>
          </div>
          
          <form @submit.prevent="submitFinancingRequest" class="p-4 space-y-5">
            <div>
              <label for="bankName" class="block text-sm font-medium text-gray-700 mb-1">Bank Name</label>
              <input 
                id="bankName" 
                v-model="financingForm.bankName" 
                type="text" 
                class="w-full py-3 px-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.bankName }"
                placeholder="Enter your bank name"
              />
              <p v-if="errors.bankName" class="mt-1 text-sm text-red-600">{{ errors.bankName }}</p>
            </div>
            
            <div>
              <label for="accountNumber" class="block text-sm font-medium text-gray-700 mb-1">Account Number</label>
              <input 
                id="accountNumber" 
                v-model="financingForm.accountNumber" 
                type="text" 
                inputmode="numeric"
                class="w-full py-3 px-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.accountNumber }"
                placeholder="Enter your account number"
              />
              <p v-if="errors.accountNumber" class="mt-1 text-sm text-red-600">{{ errors.accountNumber }}</p>
            </div>
            
            <div>
              <label for="routingNumber" class="block text-sm font-medium text-gray-700 mb-1">Routing Number</label>
              <input 
                id="routingNumber" 
                v-model="financingForm.routingNumber" 
                type="text" 
                inputmode="numeric"
                class="w-full py-3 px-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.routingNumber }"
                placeholder="Enter your routing number"
              />
              <p v-if="errors.routingNumber" class="mt-1 text-sm text-red-600">{{ errors.routingNumber }}</p>
            </div>
            
            <div>
              <label for="accountType" class="block text-sm font-medium text-gray-700 mb-1">Account Type</label>
              <select 
                id="accountType" 
                v-model="financingForm.accountType" 
                class="w-full py-3 px-4 rounded-md border border-input bg-background text-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
              >
                <option v-for="type in accountTypes" :key="type.value" :value="type.value">
                  {{ type.label }}
                </option>
              </select>
            </div>
            
            <div class="pt-4 border-t border-border">
              <div class="flex items-start">
                <div class="flex-shrink-0">
                  <input 
                    id="termsAccepted" 
                    v-model="financingForm.termsAccepted" 
                    type="checkbox" 
                    class="h-5 w-5 text-primary focus:ring-primary border-gray-300 rounded"
                    :class="{ 'border-red-500': errors.termsAccepted }"
                  />
                </div>
                <div class="ml-3">
                  <label for="termsAccepted" class="text-sm text-gray-600">
                    I agree to the <a href="#" class="text-primary hover:underline">Terms and Conditions</a> and <a href="#" class="text-primary hover:underline">Privacy Policy</a>
                  </label>
                  <p v-if="errors.termsAccepted" class="mt-1 text-sm text-red-600">{{ errors.termsAccepted }}</p>
                </div>
              </div>
            </div>
            
            <div class="pt-4">
              <button 
                type="submit" 
                class="w-full py-3 px-4 bg-primary text-white rounded-md text-sm font-medium hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
                :disabled="isSubmitting || selectedInvoices.length === 0"
              >
                <span v-if="isSubmitting" class="flex items-center justify-center">
                  <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Processing...
                </span>
                <span v-else>Submit Financing Request</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>