<template>
  <div>
    <ViewsAdminDashboardHeader />
    <main class="flex-1 space-y-4 p-4 md:p-8 pt-6 bg-gray-100 font-inter">
      <Toaster richColors position="top-right" />
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight text-gray-800">Financial Institution Terms</h2>
        <div class="flex items-center space-x-2">
          <Button @click="openCreateModal" variant="default" size="sm">
            <PlusCircle class="mr-2 h-4 w-4" />
            Add Term
          </Button>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">
              Total Terms
            </CardTitle>
            <FileText class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ totalTerms }}</div>
            <p class="text-xs text-muted-foreground">
              Financing terms available
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">
              Active Terms
            </CardTitle>
            <CheckCircle2 class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ activeTerms }}</div>
            <p class="text-xs text-muted-foreground">
              Currently active terms
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Total Products</CardTitle>
            <Package class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ totalProducts }}</div>
            <p class="text-xs text-muted-foreground">
              Financing products available
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Total Institutions</CardTitle>
            <Building2 class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ totalInstitutions }}</div>
            <p class="text-xs text-muted-foreground">
              Financial partners on the platform
            </p>
          </CardContent>
        </Card>
      </div>

      <!-- Search and Filter -->
      <div class="flex flex-col sm:flex-row gap-4 items-center justify-between">
        <div class="relative w-full sm:w-64">
          <Search class="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            v-model="searchQuery"
            placeholder="Search terms..."
            class="pl-8 w-full"
            @input="handleSearch"
          />
        </div>
        <div class="flex gap-2 w-full sm:w-auto">
          <Select v-model="institutionFilter" @update:modelValue="fetchProducts">
            <SelectTrigger class="w-full sm:w-[180px]">
              <SelectValue placeholder="Institution" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">All Institutions</SelectItem>
              <SelectItem v-for="institution in institutions" :key="institution.id" :value="institution.id.toString()">
                {{ institution.name }}
              </SelectItem>
            </SelectContent>
          </Select>
          <Select v-model="productFilter" @update:modelValue="fetchTerms">
            <SelectTrigger class="w-full sm:w-[180px]">
              <SelectValue placeholder="Product" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">All Products</SelectItem>
              <SelectItem v-for="product in products" :key="product.id" :value="product.id.toString()">
                {{ product.name }}
              </SelectItem>
            </SelectContent>
          </Select>
          <Select v-model="statusFilter" @update:modelValue="fetchTerms">
            <SelectTrigger class="w-full sm:w-[180px]">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">All Statuses</SelectItem>
              <SelectItem value="true">Active</SelectItem>
              <SelectItem value="false">Inactive</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <!-- Terms Table -->
      <Card>
        <CardHeader>
          <CardTitle>Financial Institution Terms</CardTitle>
          <CardDescription>Manage financing terms offered by your financial institution partners.</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="isLoading" class="space-y-4">
            <div v-for="i in 5" :key="i" class="h-12 bg-gray-200 animate-pulse rounded"></div>
          </div>
          <div v-else-if="terms.length === 0" class="text-center py-8 text-gray-500">
            No financial institution terms found.
          </div>
          <div v-else>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Institution</TableHead>
                  <TableHead>Product</TableHead>
                  <TableHead>Duration</TableHead>
                  <TableHead>Interest Rate</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead class="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="term in terms" :key="term.id">
                  <TableCell class="font-medium">{{ term.name }}</TableCell>
                  <TableCell>{{ term.institutionName }}</TableCell>
                  <TableCell>{{ term.productName }}</TableCell>
                  <TableCell>{{ term.durationMonths }} months</TableCell>
                  <TableCell>{{ term.interestRate }}%</TableCell>
                  <TableCell>
                    <Badge :variant="term.isActive ? 'default' : 'secondary'">
                      {{ term.isActive ? 'Active' : 'Inactive' }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ formatDate(term.createdAt) }}</TableCell>
                  <TableCell class="text-right">
                    <Button variant="ghost" size="icon" @click="viewTerm(term)">
                      <Eye class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" @click="editTerm(term)">
                      <Edit class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" @click="confirmDelete(term)">
                      <Trash2 class="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
            
            <!-- Pagination -->
            <div class="flex items-center justify-end space-x-2 py-4">
              <Button
                variant="outline"
                size="sm"
                :disabled="currentPage <= 1"
                @click="changePage(currentPage - 1)"
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                :disabled="currentPage >= totalPages"
                @click="changePage(currentPage + 1)"
              >
                Next
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Create/Edit Term Modal -->
      <Dialog v-model:open="showModal">
        <DialogContent class="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{{ isEditing ? 'Edit Term' : 'Create Term' }}</DialogTitle>
            <DialogDescription>
              {{ isEditing ? 'Update the details of this financing term.' : 'Add a new financing term to the platform.' }}
            </DialogDescription>
          </DialogHeader>
          <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="name" class="text-right">Name</Label>
              <Input id="name" v-model="formData.name" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="institution" class="text-right">Institution</Label>
              <Select v-model="formData.financialInstitutionId" class="col-span-3" @update:modelValue="loadProductsForInstitution">
                <SelectTrigger>
                  <SelectValue placeholder="Select institution" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="institution in institutions" :key="institution.id" :value="institution.id">
                    {{ institution.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="product" class="text-right">Product</Label>
              <Select v-model="formData.productId" class="col-span-3">
                <SelectTrigger>
                  <SelectValue placeholder="Select product" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="product in filteredProducts" :key="product.id" :value="product.id">
                    {{ product.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="durationMonths" class="text-right">Duration (months)</Label>
              <Input id="durationMonths" v-model="formData.durationMonths" type="number" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="interestRate" class="text-right">Interest Rate (%)</Label>
              <Input id="interestRate" v-model="formData.interestRate" type="number" step="0.01" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="processingFee" class="text-right">Processing Fee (%)</Label>
              <Input id="processingFee" v-model="formData.processingFee" type="number" step="0.01" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="description" class="text-right">Description</Label>
              <Input id="description" v-model="formData.description" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="status" class="text-right">Status</Label>
              <Select v-model="formData.isActive" class="col-span-3">
                <SelectTrigger>
                  <SelectValue placeholder="Select status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="true">Active</SelectItem>
                  <SelectItem :value="false">Inactive</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" @click="showModal = false">Cancel</Button>
            <Button @click="saveTerm">Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <!-- Delete Confirmation Dialog -->
      <Dialog v-model:open="showDeleteDialog">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm Deletion</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {{ selectedTerm?.name }}? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" @click="showDeleteDialog = false">Cancel</Button>
            <Button variant="destructive" @click="deleteTerm">Delete</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { 
  Building2, CheckCircle2, Package, FileText, PlusCircle, 
  Search, Eye, Edit, Trash2 
} from 'lucide-vue-next';
import { Button } from '~/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import { Badge } from '~/components/ui/badge';
import { 
  Table, TableBody, TableCell, TableHead, TableHeader, TableRow 
} from '~/components/ui/table';
import { 
  Dialog, DialogContent, DialogDescription, DialogFooter, 
  DialogHeader, DialogTitle 
} from '~/components/ui/dialog';
import { Label } from '~/components/ui/label';
import { 
  Select, SelectContent, SelectItem, SelectTrigger, SelectValue 
} from '~/components/ui/select';
import { Toaster, toast } from 'vue-sonner';
import { useCookie } from '#app';
import axios from 'axios';

definePageMeta({
  layout: 'admin',
  middleware: 'auth'
});

const API_BASE_URL = 'http://localhost:3000/api/v1'; // Replace with your actual API base URL
const tokenCookie = useCookie('token');
const authToken = tokenCookie.value || null;

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    ...(authToken ? { Authorization: `Bearer ${authToken}` } : {}),
    'Content-Type': 'application/json',
  },
});

// State
const isLoading = ref(true);
const terms = ref<any[]>([]);
const institutions = ref<any[]>([]);
const products = ref<any[]>([]);
const filteredProducts = ref<any[]>([]);
const totalTerms = ref(0);
const activeTerms = ref(0);
const totalProducts = ref(0);
const totalInstitutions = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const searchQuery = ref('');
const statusFilter = ref('');
const institutionFilter = ref('');
const productFilter = ref('');
const showModal = ref(false);
const isEditing = ref(false);
const selectedTerm = ref<any>(null);
const showDeleteDialog = ref(false);

const formData = ref({
  name: '',
  financialInstitutionId: 0,
  productId: 0,
  durationMonths: 0,
  interestRate: 0,
  processingFee: 0,
  description: '',
  isActive: true
});

// Fetch financial institution terms
const fetchTerms = async () => {
  isLoading.value = true;
  try {
    let url = `/admin/financial-institutions/terms?page=${currentPage.value}&pageSize=${pageSize.value}`;
    
    if (searchQuery.value) {
      url += `&name=${searchQuery.value}`;
    }
    
    if (statusFilter.value) {
      url += `&is_active=${statusFilter.value}`;
    }
    
    if (institutionFilter.value) {
      url += `&financial_institution_id=${institutionFilter.value}`;
    }
    
    if (productFilter.value) {
      url += `&product_id=${productFilter.value}`;
    }
    
    const response = await apiClient.get(url);
    terms.value = response.data.terms || [];
    totalTerms.value = response.data.total || 0;
    totalPages.value = Math.ceil(totalTerms.value / pageSize.value);
    
    // Count active terms
    activeTerms.value = terms.value.filter(term => term.isActive).length;
    
    // Fetch additional stats
    await fetchStats();
    
  } catch (error) {
    console.error("Failed to fetch financial institution terms:", error);
    toast.error("Could not load financial institution terms.");
  } finally {
    isLoading.value = false;
  }
};

// Fetch financial institutions for dropdown
const fetchInstitutions = async () => {
  try {
    const response = await apiClient.get('/admin/financial-institutions?pageSize=100');
    institutions.value = response.data.financialInstitutions || [];
    totalInstitutions.value = response.data.total || 0;
  } catch (error) {
    console.error("Failed to fetch financial institutions:", error);
    toast.error("Could not load financial institutions.");
  }
};

// Fetch products for dropdown
const fetchProducts = async () => {
  try {
    const response = await apiClient.get('/admin/financial-institutions/products?pageSize=100');
    products.value = response.data.products || [];
    totalProducts.value = response.data.total || 0;
  } catch (error) {
    console.error("Failed to fetch financial institution products:", error);
    toast.error("Could not load financial institution products.");
  }
};

// Load products for selected institution
const loadProductsForInstitution = () => {
  if (formData.value.financialInstitutionId) {
    filteredProducts.value = products.value.filter(
      product => product.financialInstitutionId === formData.value.financialInstitutionId
    );
    // Reset product selection if current selection is not valid for this institution
    const productExists = filteredProducts.value.some(p => p.id === formData.value.productId);
    if (!productExists) {
      formData.value.productId = 0;
    }
  } else {
    filteredProducts.value = products.value;
  }
};

// Fetch additional statistics
const fetchStats = async () => {
  try {
    // In a real implementation, you would fetch these from the API
    // For now, we'll use mock data or data we already have
    totalProducts.value = products.value.length;
  } catch (error) {
    console.error("Failed to fetch stats:", error);
  }
};

// Handle search
const handleSearch = () => {
  currentPage.value = 1; // Reset to first page when searching
  fetchTerms();
};

// Pagination
const changePage = (page: number) => {
  currentPage.value = page;
  fetchTerms();
};

// Format date
const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  });
};

// Modal functions
const openCreateModal = () => {
  isEditing.value = false;
  formData.value = {
    name: '',
    financialInstitutionId: 0,
    productId: 0,
    durationMonths: 0,
    interestRate: 0,
    processingFee: 0,
    description: '',
    isActive: true
  };
  filteredProducts.value = products.value;
  showModal.value = true;
};

const editTerm = (term: any) => {
  isEditing.value = true;
  selectedTerm.value = term;
  formData.value = {
    name: term.name,
    financialInstitutionId: term.financialInstitutionId,
    productId: term.productId,
    durationMonths: term.durationMonths || 0,
    interestRate: term.interestRate || 0,
    processingFee: term.processingFee || 0,
    description: term.description || '',
    isActive: term.isActive
  };
  
  // Filter products for this institution
  loadProductsForInstitution();
  
  showModal.value = true;
};

const viewTerm = (term: any) => {
  // In a real implementation, you would navigate to a detail page
  // For now, we'll just show a toast
  toast.info(`Viewing details for ${term.name}`);
};

const saveTerm = async () => {
  try {
    if (isEditing.value && selectedTerm.value) {
      // Update existing term
      await apiClient.put(`/admin/financial-institutions/terms/${selectedTerm.value.id}`, formData.value);
      toast.success(`${formData.value.name} updated successfully!`);
    } else {
      // Create new term
      await apiClient.post('/admin/financial-institutions/terms', formData.value);
      toast.success(`${formData.value.name} created successfully!`);
    }
    showModal.value = false;
    fetchTerms();
  } catch (error) {
    console.error("Failed to save financial institution term:", error);
    let errorMessage = "Could not save financial institution term.";
    if (axios.isAxiosError(error) && error.response) {
      errorMessage += ` (Status: ${error.response.status})`;
    }
    toast.error(errorMessage);
  }
};

// Delete functions
const confirmDelete = (term: any) => {
  selectedTerm.value = term;
  showDeleteDialog.value = true;
};

const deleteTerm = async () => {
  if (!selectedTerm.value) return;
  
  try {
    await apiClient.delete(`/admin/financial-institutions/terms/${selectedTerm.value.id}`);
    toast.success(`${selectedTerm.value.name} deleted successfully!`);
    showDeleteDialog.value = false;
    fetchTerms();
  } catch (error) {
    console.error("Failed to delete financial institution term:", error);
    let errorMessage = "Could not delete financial institution term.";
    if (axios.isAxiosError(error) && error.response) {
      errorMessage += ` (Status: ${error.response.status})`;
    }
    toast.error(errorMessage);
  }
};

// Initialize
onMounted(() => {
  if (authToken) {
    Promise.all([
      fetchInstitutions(),
      fetchProducts()
    ]).then(() => {
      fetchTerms();
    });
  } else {
    toast.error("Authentication token not found. Please log in.");
    isLoading.value = false;
  }
});
</script>