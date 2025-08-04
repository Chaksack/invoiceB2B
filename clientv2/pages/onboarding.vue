<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from '#app';
import { Stepper, StepperDescription, StepperItem, StepperSeparator, StepperTitle, StepperTrigger } from '~/components/ui/stepper';
import { Button } from '~/components/ui/button';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '~/components/ui/form';
import { Input } from '~/components/ui/input';
import { Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from '~/components/ui/select';
import { Check, Circle, Dot } from 'lucide-vue-next';
import { toast } from 'vue-sonner';

definePageMeta({
  layout: 'auth'
});

const router = useRouter();
const stepIndex = ref(1);

// Define the steps for the onboarding process
const steps = [
  {
    step: 1,
    title: 'Personal Information',
    description: 'Tell us about yourself',
  },
  {
    step: 2,
    title: 'Business Details',
    description: 'Information about your business',
  },
  {
    step: 3,
    title: 'KYC Verification',
    description: 'Upload your identification documents',
  },
  {
    step: 4,
    title: 'Banking Information',
    description: 'Add your banking details',
  },
  {
    step: 5,
    title: 'Review & Submit',
    description: 'Review your information',
  },
];

// Form data for all steps
const formData = reactive({
  // Personal Information
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  dateOfBirth: '',
  
  // Business Details
  companyName: '',
  businessType: '',
  industry: '',
  registrationNumber: '',
  taxId: '',
  businessAddress: '',
  
  // KYC Verification
  idType: '',
  idNumber: '',
  idDocument: null,
  proofOfAddress: null,
  
  // Banking Information
  bankName: '',
  accountNumber: '',
  routingNumber: '',
  accountType: '',
  
  // Terms and Consent
  termsAccepted: false
});

// Form validation state
const formValid = reactive({
  step1: false,
  step2: false,
  step3: false,
  step4: false,
  step5: false
});

// Validate current step
const validateCurrentStep = () => {
  switch(stepIndex.value) {
    case 1:
      formValid.step1 = !!formData.firstName && !!formData.lastName && !!formData.email && !!formData.phone;
      return formValid.step1;
    case 2:
      formValid.step2 = !!formData.companyName && !!formData.businessType && !!formData.industry && !!formData.businessAddress;
      return formValid.step2;
    case 3:
      formValid.step3 = !!formData.idType && !!formData.idNumber;
      return formValid.step3;
    case 4:
      formValid.step4 = !!formData.bankName && !!formData.accountNumber && !!formData.routingNumber && !!formData.accountType;
      return formValid.step4;
    case 5:
      formValid.step5 = formData.termsAccepted;
      return formValid.step5;
    default:
      return false;
  }
};

// Handle next step
const nextStep = () => {
  if (validateCurrentStep()) {
    if (stepIndex.value < steps.length) {
      stepIndex.value++;
    } else {
      submitForm();
    }
  } else {
    toast.error('Please fill in all required fields');
  }
};

// Handle previous step
const prevStep = () => {
  if (stepIndex.value > 1) {
    stepIndex.value--;
  }
};

// Submit the form
const submitForm = () => {
  // In a real application, you would submit the form data to an API
  console.log('Form submitted:', formData);
  toast.success('Onboarding completed successfully!');
  
  // Redirect to dashboard
  setTimeout(() => {
    router.push('/dashboard');
  }, 1500);
};
</script>

<template>
  <div class="min-h-screen flex flex-col lg:flex-row font-inter">
    <!-- Left Section: Primary Background with Illustrations -->
    <div class="relative hidden lg:flex bg-primary text-white p-4 sm:p-6 md:p-8 flex-col justify-center items-center overflow-hidden rounded-xl lg:rounded-r-xl lg:rounded-l-none">
      <!-- Background patterns -->
      <div class="absolute top-0 left-0 w-full h-full opacity-10">
        <svg class="absolute top-8 left-8 hidden sm:block" width="50" height="50" viewBox="0 0 50 50" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="25" cy="25" r="23" stroke="white" stroke-width="2" />
        </svg>
        <svg class="absolute bottom-16 right-16 hidden sm:block" width="100" height="100" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect x="10" y="10" width="80" height="80" rx="10" stroke="white" stroke-width="2" />
        </svg>
        <div class="absolute top-1/4 right-1/4 w-4 h-4 bg-white rounded-full hidden sm:block"></div>
        <div class="absolute bottom-1/3 left-1/3 w-6 h-6 bg-white rounded-full hidden sm:block"></div>
      </div>

      <!-- Logo -->
      <div class="w-full flex justify-center sm:justify-start sm:absolute sm:top-8 sm:left-8 items-center space-x-2 mb-6 sm:mb-0">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 2L2 12L12 22L22 12L12 2Z" fill="white"/>
        </svg>
        <span class="text-2xl font-bold">InvoiceFlow</span>
      </div>

      <!-- Illustrative Cards -->
      <div class="relative mt-8 sm:mt-12 lg:mt-16 mb-8 w-full max-w-xs sm:max-w-sm mx-auto">
        <!-- Main Card -->
        <div class="bg-white text-gray-800 p-4 sm:p-6 rounded-xl shadow-xl z-10 relative">
          <div class="flex justify-between items-center mb-4">
            <div class="text-sm sm:text-lg font-semibold">Invoices Financed</div>
            <div class="text-sm sm:text-lg font-semibold">Funds Available</div>
          </div>
          <div class="flex justify-between items-center mb-6">
            <div class="text-lg sm:text-2xl font-bold">$124,908.00</div>
            <div class="text-lg sm:text-2xl font-bold">$85,750.00</div>
          </div>
          <!-- Graph Placeholder -->
          <div class="h-16 sm:h-24 bg-gray-100 rounded-lg flex items-center justify-center text-gray-400 mb-6">
            <svg class="w-full h-full" viewBox="0 0 300 100" preserveAspectRatio="none">
              <polyline fill="none" stroke="var(--primary)" stroke-width="2" points="0,80 50,50 100,70 150,30 200,60 250,40 300,20" />
              <polyline fill="none" stroke="#EF4444" stroke-width="2" points="0,20 50,40 100,30 150,60 200,40 250,70 300,80" />
            </svg>
          </div>
          <div class="flex items-center mb-4">
            <div class="w-5 h-5 sm:w-6 sm:h-6 bg-green-100 rounded-full flex items-center justify-center mr-2 sm:mr-3">
              <svg class="w-3 h-3 sm:w-4 sm:h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
            </div>
            <div>
              <div class="text-sm sm:text-base font-medium">Invoice financing approved</div>
              <div class="text-xs sm:text-sm text-gray-500">Today at 11:15 AM</div>
            </div>
          </div>
          <div class="border-t border-gray-200 pt-4">
            <div class="flex justify-between items-center mb-2">
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 12</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 15</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 18</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 21</div>
              <div class="text-xs sm:text-sm font-medium text-gray-600">Jan 24</div>
            </div>
            <div class="flex justify-between items-center text-xs sm:text-sm">
              <div class="font-medium text-gray-800">$35,798.00</div>
              <div class="font-medium text-gray-800">$23.10</div>
              <div class="font-medium text-gray-800">$23.10</div>
              <div class="font-medium text-gray-800">$23.10</div>
              <div class="font-medium text-gray-800">$23.10</div>
            </div>
          </div>
        </div>

        <!-- Overlapping Card 1 -->
        <div class="absolute bg-white text-gray-800 p-3 sm:p-4 rounded-xl shadow-lg -top-6 sm:-top-10 left-1/3 sm:left-1/4 transform -translate-x-1/2 z-0 opacity-90">
          <div class="font-semibold text-xs sm:text-sm mb-1 sm:mb-2">Invoice Financed</div>
          <div class="text-sm sm:text-lg font-bold text-green-500">+$42,500.00</div>
        </div>

        <!-- Overlapping Card 2 -->
        <div class="absolute bg-white text-gray-800 p-3 sm:p-4 rounded-xl shadow-lg bottom-6 sm:bottom-10 right-0 transform translate-x-1/6 sm:translate-x-1/4 z-0 opacity-90">
          <div class="flex items-center justify-between text-xs sm:text-sm mb-1">
            <span class="font-medium">Financing Fee</span>
            <span class="text-red-500 font-bold">-$1,275</span>
          </div>
          <div class="flex items-center justify-between text-xs sm:text-sm mb-1">
            <span class="font-medium">Service Charge</span>
            <span class="text-red-500 font-bold">-$425</span>
          </div>
          <div class="flex items-center justify-between text-xs sm:text-sm">
            <span class="font-medium">Processing Fee</span>
            <span class="text-red-500 font-bold">-$212</span>
          </div>
          <div class="text-xs text-gray-500 mt-1 sm:mt-2">Today at 9:24 AM</div>
        </div>
      </div>

      <!-- Main Marketing Text -->
      <div class="text-center px-4 mt-8 sm:mt-10">
        <h1 class="text-2xl sm:text-3xl lg:text-4xl font-extrabold mb-3 sm:mb-4">Unlock Your Business Cash Flow</h1>
        <p class="text-sm sm:text-base leading-relaxed max-w-xs mx-auto">
          InvoiceFlow helps your business convert unpaid invoices into immediate working capital. Get paid within 24 hours, bridge cash flow gaps, and focus on growing your business while we handle your invoice financing needs.
        </p>
      </div>
    </div>

    <!-- Right Section: Onboarding Form -->
    <div class="bg-white p-4 sm:p-6 md:p-8 lg:p-12 flex flex-col justify-between">
      <div class="flex flex-col max-w-7xl p-4 sm:p-6 md:p-8 my-auto mx-auto ">
        <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2 text-center">Complete Your Profile</h2>
        <p class="text-gray-600 mb-6 sm:mb-8 text-center text-sm sm:text-base">Let's get your account set up for invoice financing</p>
        
        <!-- Stepper Component -->
        <div class="mb-10 w-full">
          <Stepper v-model="stepIndex" class="w-full">
            <div class="flex w-full justify-between">
              <StepperItem
                v-for="step in steps"
                :key="step.step"
                v-slot="{ state }"
                class="relative flex w-full flex-col items-center justify-center"
                :step="step.step"
              >
                <StepperSeparator
                  v-if="step.step !== steps[steps.length - 1].step"
                  class="absolute left-[calc(50%+24px)] right-[calc(-50%+24px)] top-6 h-1 bg-gray-200 group-data-[state=completed]:bg-primary"
                />
                
                <StepperTrigger as-child>
                  <Button
                    :variant="state === 'completed' || state === 'active' ? 'default' : 'outline'"
                    size="icon"
                    class="z-10 rounded-full shrink-0 h-12 w-12 sm:h-14 sm:w-14 shadow-md"
                    :class="[
                      state === 'completed' ? 'bg-primary text-white' : '',
                      state === 'active' ? 'bg-primary text-white ring-4 ring-primary/30 ring-offset-2' : '',
                      state === 'inactive' ? 'bg-white text-gray-400 border-2 border-gray-200' : ''
                    ]"
                  >
                    <Check v-if="state === 'completed'" class="h-6 w-6 sm:h-7 sm:w-7" />
                    <Circle v-else-if="state === 'active'" class="h-6 w-6 sm:h-7 sm:w-7" />
                    <span v-else class="text-lg font-bold">{{ step.step }}</span>
                  </Button>
                </StepperTrigger>
                
                <div class="mt-4 text-center">
                  <StepperTitle class="text-sm sm:text-base font-medium">{{ step.title }}</StepperTitle>
                  <StepperDescription class="text-xs sm:text-sm text-gray-500 hidden sm:block">{{ step.description }}</StepperDescription>
                </div>
              </StepperItem>
            </div>
          </Stepper>
        </div>
        
        <!-- Form Steps -->
        <!-- Step 1: Personal Information -->
        <div v-if="stepIndex === 1" class="space-y-8">
          <h2 class="text-xl font-semibold text-gray-900">Personal Information</h2>
          <p class="text-gray-600">Please provide your personal details</p>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <FormField name="firstName">
              <FormItem>
                <FormLabel for="firstName" class="sr-only">First Name</FormLabel>
                <FormControl>
                  <Input 
                    id="firstName" 
                    v-model="formData.firstName" 
                    placeholder="Enter your first name" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="lastName">
              <FormItem>
                <FormLabel for="lastName" class="sr-only">Last Name</FormLabel>
                <FormControl>
                  <Input 
                    id="lastName" 
                    v-model="formData.lastName" 
                    placeholder="Enter your last name" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="email">
              <FormItem>
                <FormLabel for="email" class="sr-only">Email Address</FormLabel>
                <FormControl>
                  <Input 
                    id="email" 
                    v-model="formData.email" 
                    type="email" 
                    placeholder="Enter your email address" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="phone">
              <FormItem>
                <FormLabel for="phone" class="sr-only">Phone Number</FormLabel>
                <FormControl>
                  <Input 
                    id="phone" 
                    v-model="formData.phone" 
                    placeholder="Enter your phone number" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="dateOfBirth">
              <FormItem>
                <FormLabel for="dateOfBirth" class="sr-only">Date of Birth</FormLabel>
                <FormControl>
                  <Input 
                    id="dateOfBirth" 
                    v-model="formData.dateOfBirth" 
                    type="date" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>
        </div>
        
        <!-- Step 2: Business Details -->
        <div v-if="stepIndex === 2" class="space-y-8">
          <h2 class="text-xl font-semibold text-gray-900">Business Details</h2>
          <p class="text-gray-600">Tell us about your business</p>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <FormField name="companyName">
              <FormItem>
                <FormLabel for="companyName" class="sr-only">Company Name</FormLabel>
                <FormControl>
                  <Input 
                    id="companyName" 
                    v-model="formData.companyName" 
                    placeholder="Enter your company name" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="businessType">
              <FormItem>
                <FormLabel for="businessType" class="sr-only">Business Type</FormLabel>
                <FormControl>
                  <Select v-model="formData.businessType">
                    <SelectTrigger class="w-full">
                      <SelectValue placeholder="Select business type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="sole_proprietorship">Sole Proprietorship</SelectItem>
                        <SelectItem value="partnership">Partnership</SelectItem>
                        <SelectItem value="llc">Limited Liability Company (LLC)</SelectItem>
                        <SelectItem value="corporation">Corporation</SelectItem>
                        <SelectItem value="nonprofit">Non-profit Organization</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="industry">
              <FormItem>
                <FormLabel for="industry" class="sr-only">Industry</FormLabel>
                <FormControl>
                  <Select v-model="formData.industry">
                    <SelectTrigger class="w-full">
                      <SelectValue placeholder="Select industry" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="technology">Technology</SelectItem>
                        <SelectItem value="healthcare">Healthcare</SelectItem>
                        <SelectItem value="retail">Retail</SelectItem>
                        <SelectItem value="manufacturing">Manufacturing</SelectItem>
                        <SelectItem value="finance">Finance</SelectItem>
                        <SelectItem value="education">Education</SelectItem>
                        <SelectItem value="other">Other</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="registrationNumber">
              <FormItem>
                <FormLabel for="registrationNumber" class="sr-only">Registration Number</FormLabel>
                <FormControl>
                  <Input 
                    id="registrationNumber" 
                    v-model="formData.registrationNumber" 
                    placeholder="Enter business registration number" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="taxId">
              <FormItem>
                <FormLabel for="taxId" class="sr-only">Tax ID / EIN</FormLabel>
                <FormControl>
                  <Input 
                    id="taxId" 
                    v-model="formData.taxId" 
                    placeholder="Enter tax ID or EIN" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="businessAddress" class="md:col-span-2">
              <FormItem>
                <FormLabel for="businessAddress" class="sr-only">Business Address</FormLabel>
                <FormControl>
                  <Input 
                    id="businessAddress" 
                    v-model="formData.businessAddress" 
                    placeholder="Enter business address" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>
        </div>
        
        <!-- Step 3: KYC Verification -->
        <div v-if="stepIndex === 3" class="space-y-8">
          <h2 class="text-xl font-semibold text-gray-900">KYC Verification</h2>
          <p class="text-gray-600">Upload identification documents for verification</p>
          
          <div class="space-y-6">
            <FormField name="idType">
              <FormItem>
                <FormLabel for="idType" class="sr-only">ID Type</FormLabel>
                <FormControl>
                  <Select v-model="formData.idType">
                    <SelectTrigger class="w-full">
                      <SelectValue placeholder="Select ID type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="passport">Passport</SelectItem>
                        <SelectItem value="drivers_license">Driver's License</SelectItem>
                        <SelectItem value="national_id">National ID</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="idNumber">
              <FormItem>
                <FormLabel for="idNumber" class="sr-only">ID Number</FormLabel>
                <FormControl>
                  <Input 
                    id="idNumber" 
                    v-model="formData.idNumber" 
                    placeholder="Enter ID number" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="idDocument">
              <FormItem>
                <FormLabel class="block text-sm font-medium text-gray-700 mb-1">Upload ID Document</FormLabel>
                <FormControl>
                  <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
                    <input type="file" id="idDocument" class="hidden" />
                    <label for="idDocument" class="cursor-pointer">
                      <div class="flex flex-col items-center">
                        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
                        </svg>
                        <p class="mt-2 text-sm text-gray-600">
                          <span class="font-medium text-primary hover:text-primary-700">Click to upload</span> or drag and drop
                        </p>
                        <p class="mt-1 text-xs text-gray-500">
                          PDF, JPG, PNG up to 10MB
                        </p>
                      </div>
                    </label>
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="proofOfAddress">
              <FormItem>
                <FormLabel class="block text-sm font-medium text-gray-700 mb-1">Upload Proof of Address</FormLabel>
                <FormControl>
                  <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
                    <input type="file" id="proofOfAddress" class="hidden" />
                    <label for="proofOfAddress" class="cursor-pointer">
                      <div class="flex flex-col items-center">
                        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
                        </svg>
                        <p class="mt-2 text-sm text-gray-600">
                          <span class="font-medium text-primary hover:text-primary-700">Click to upload</span> or drag and drop
                        </p>
                        <p class="mt-1 text-xs text-gray-500">
                          PDF, JPG, PNG up to 10MB
                        </p>
                      </div>
                    </label>
                  </div>
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>
        </div>
        
        <!-- Step 4: Banking Information -->
        <div v-if="stepIndex === 4" class="space-y-8">
          <h2 class="text-xl font-semibold text-gray-900">Banking Information</h2>
          <p class="text-gray-600">Provide your banking details for payments</p>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <FormField name="bankName">
              <FormItem>
                <FormLabel for="bankName" class="sr-only">Bank Name</FormLabel>
                <FormControl>
                  <Input 
                    id="bankName" 
                    v-model="formData.bankName" 
                    placeholder="Enter bank name" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="accountType">
              <FormItem>
                <FormLabel for="accountType" class="sr-only">Account Type</FormLabel>
                <FormControl>
                  <Select v-model="formData.accountType">
                    <SelectTrigger class="w-full">
                      <SelectValue placeholder="Select account type" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="checking">Checking Account</SelectItem>
                        <SelectItem value="savings">Savings Account</SelectItem>
                        <SelectItem value="business">Business Account</SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="accountNumber">
              <FormItem>
                <FormLabel for="accountNumber" class="sr-only">Account Number</FormLabel>
                <FormControl>
                  <Input 
                    id="accountNumber" 
                    v-model="formData.accountNumber" 
                    placeholder="Enter account number" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
            
            <FormField name="routingNumber">
              <FormItem>
                <FormLabel for="routingNumber" class="sr-only">Routing Number</FormLabel>
                <FormControl>
                  <Input 
                    id="routingNumber" 
                    v-model="formData.routingNumber" 
                    placeholder="Enter routing number" 
                    class="py-2 sm:py-3 px-3 sm:px-4 text-sm sm:text-base rounded-lg border border-gray-300 focus:ring-primary focus:border-primary shadow-sm w-full"
                    required
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            </FormField>
          </div>
        </div>
        
        <!-- Step 5: Review & Submit -->
        <div v-if="stepIndex === 5" class="space-y-8">
          <h2 class="text-xl font-semibold text-gray-900">Review & Submit</h2>
          <p class="text-gray-600">Please review your information before submitting</p>
          
          <div class="space-y-6">
            <!-- Personal Information Summary -->
            <div class="border rounded-lg p-4">
              <h3 class="text-lg font-medium text-gray-900 mb-3">Personal Information</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p class="text-sm font-medium text-gray-500">Full Name</p>
                  <p class="text-sm">{{ formData.firstName }} {{ formData.lastName }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Email</p>
                  <p class="text-sm">{{ formData.email }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Phone</p>
                  <p class="text-sm">{{ formData.phone }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Date of Birth</p>
                  <p class="text-sm">{{ formData.dateOfBirth }}</p>
                </div>
              </div>
            </div>
            
            <!-- Business Details Summary -->
            <div class="border rounded-lg p-4">
              <h3 class="text-lg font-medium text-gray-900 mb-3">Business Details</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p class="text-sm font-medium text-gray-500">Company Name</p>
                  <p class="text-sm">{{ formData.companyName }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Business Type</p>
                  <p class="text-sm">{{ formData.businessType }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Industry</p>
                  <p class="text-sm">{{ formData.industry }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Registration Number</p>
                  <p class="text-sm">{{ formData.registrationNumber || 'Not provided' }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Tax ID / EIN</p>
                  <p class="text-sm">{{ formData.taxId || 'Not provided' }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Business Address</p>
                  <p class="text-sm">{{ formData.businessAddress }}</p>
                </div>
              </div>
            </div>
            
            <!-- KYC Verification Summary -->
            <div class="border rounded-lg p-4">
              <h3 class="text-lg font-medium text-gray-900 mb-3">KYC Verification</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p class="text-sm font-medium text-gray-500">ID Type</p>
                  <p class="text-sm">{{ formData.idType }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">ID Number</p>
                  <p class="text-sm">{{ formData.idNumber }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">ID Document</p>
                  <p class="text-sm">{{ formData.idDocument ? 'Uploaded' : 'Not uploaded' }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Proof of Address</p>
                  <p class="text-sm">{{ formData.proofOfAddress ? 'Uploaded' : 'Not uploaded' }}</p>
                </div>
              </div>
            </div>
            
            <!-- Banking Information Summary -->
            <div class="border rounded-lg p-4">
              <h3 class="text-lg font-medium text-gray-900 mb-3">Banking Information</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <p class="text-sm font-medium text-gray-500">Bank Name</p>
                  <p class="text-sm">{{ formData.bankName }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Account Type</p>
                  <p class="text-sm">{{ formData.accountType }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Account Number</p>
                  <p class="text-sm">{{ formData.accountNumber }}</p>
                </div>
                <div>
                  <p class="text-sm font-medium text-gray-500">Routing Number</p>
                  <p class="text-sm">{{ formData.routingNumber }}</p>
                </div>
              </div>
            </div>
            
            <!-- Terms and Conditions -->
            <FormField name="termsAccepted" v-slot="{ value, handleChange }">
              <FormItem class="flex items-start mt-6">
                <div class="flex items-center h-5">
                  <FormControl>
                    <input
                      id="termsAccepted"
                      :checked="value"
                      @change="handleChange($event.target.checked)"
                      type="checkbox"
                      class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
                      required
                    />
                  </FormControl>
                </div>
                <div class="ml-3 text-sm">
                  <FormLabel for="termsAccepted" class="font-medium text-gray-700">I agree to the Terms and Conditions</FormLabel>
                  <p class="text-gray-500">By submitting this form, I confirm that all information provided is accurate and complete.</p>
                  <FormMessage />
                </div>
              </FormItem>
            </FormField>
          </div>
        </div>
      </div>
      
      <!-- Navigation Buttons -->
      <div class="flex space-x-6 justify-between mt-12 mb-6">
        <Button
          v-if="stepIndex > 1"
          @click="prevStep"
          variant="outline"
          class="w-96 py-3 sm:py-4 px-6 sm:px-8 rounded-lg border-2 border-gray-300 text-gray-700 hover:bg-gray-50 hover:border-gray-400 focus:ring-4 focus:ring-gray-200 shadow-md text-sm sm:text-base font-medium transition-all duration-200"
        >
          <span class="flex items-center justify-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M9.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L7.414 9H15a1 1 0 110 2H7.414l2.293 2.293a1 1 0 010 1.414z" clip-rule="evenodd" />
            </svg>
            Previous
          </span>
        </Button>
        <div v-else class="w-1/2"></div>
        
        <Button
          @click="nextStep"
          class="w-96 bg-primary text-white py-3 sm:py-4 px-6 sm:px-8 text-sm sm:text-base rounded-lg hover:bg-primary-700 focus:ring-4 focus:ring-primary/30 shadow-lg font-medium transition-all duration-200"
        >
          <span class="flex items-center justify-center">
            {{ stepIndex < steps.length ? 'Next' : 'Submit' }}
            <svg v-if="stepIndex < steps.length" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </span>
        </Button>
      </div>
    </div>
  </div>
</template>

<style>
/* Ensure the Inter font is applied globally if not already configured in Tailwind */
body {
  font-family: 'Inter', sans-serif;
}

/* Form styling enhancements */
.form-item {
  margin-bottom: 1.25rem;
}

/* Form validation styling */
:deep(.form-message) {
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.375rem;
  font-weight: 500;
}

:deep(.form-control) {
  position: relative;
}

:deep(.form-control input:focus),
:deep(.form-control select:focus) {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  outline: none;
}

:deep(.form-control input.error),
:deep(.form-control select.error) {
  border-color: #ef4444;
}

:deep(.form-control input.success),
:deep(.form-control select.success) {
  border-color: #10b981;
}

/* Animations and transitions */
[v-if] {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.stepper-item {
  transition: all 0.3s ease-in-out;
}

.form-field {
  transition: all 0.2s ease-in-out;
}

input, select, button {
  transition: all 0.2s ease-in-out;
}
</style>