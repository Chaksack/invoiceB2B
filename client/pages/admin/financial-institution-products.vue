<template>
  <div>
    <ViewsAdminDashboardHeader />
    <main class="flex-1 space-y-4 p-4 md:p-8 pt-6 bg-gray-100 font-inter">
      <Toaster richColors position="top-right" />
      <div class="flex items-center justify-between space-y-2">
        <h2 class="text-3xl font-bold tracking-tight text-gray-800">Financial Institution Products</h2>
        <div class="flex items-center space-x-2">
          <Button @click="openCreateModal" variant="default" size="sm">
            <PlusCircle class="mr-2 h-4 w-4" />
            Add Product
          </Button>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">
              Total Products
            </CardTitle>
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
            <CardTitle class="text-sm font-medium">
              Active Products
            </CardTitle>
            <CheckCircle2 class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ activeProducts }}</div>
            <p class="text-xs text-muted-foreground">
              Currently active products
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
        <Card>
          <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle class="text-sm font-medium">Total Terms</CardTitle>
            <FileText class="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div v-if="isLoading" class="h-8 w-1/2 bg-gray-200 animate-pulse rounded"></div>
            <div v-else class="text-2xl font-bold">{{ totalTerms }}</div>
            <p class="text-xs text-muted-foreground">
              Financing terms configured
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
            placeholder="Search products..."
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
          <Select v-model="statusFilter" @update:modelValue="fetchProducts">
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

      <!-- Products Table -->
      <Card>
        <CardHeader>
          <CardTitle>Financial Institution Products</CardTitle>
          <CardDescription>Manage financing products offered by your financial institution partners.</CardDescription>
        </CardHeader>
        <CardContent>
          <div v-if="isLoading" class="space-y-4">
            <div v-for="i in 5" :key="i" class="h-12 bg-gray-200 animate-pulse rounded"></div>
          </div>
          <div v-else-if="products.length === 0" class="text-center py-8 text-gray-500">
            No financial institution products found.
          </div>
          <div v-else>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Institution</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Terms</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead class="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="product in products" :key="product.id">
                  <TableCell class="font-medium">{{ product.name }}</TableCell>
                  <TableCell>{{ product.institutionName }}</TableCell>
                  <TableCell>{{ product.type }}</TableCell>
                  <TableCell>
                    <Badge :variant="product.isActive ? 'default' : 'secondary'">
                      {{ product.isActive ? 'Active' : 'Inactive' }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ product.termCount || 0 }}</TableCell>
                  <TableCell>{{ formatDate(product.createdAt) }}</TableCell>
                  <TableCell class="text-right">
                    <Button variant="ghost" size="icon" @click="viewProduct(product)">
                      <Eye class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" @click="editProduct(product)">
                      <Edit class="h-4 w-4" />
                    </Button>
                    <Button variant="ghost" size="icon" @click="confirmDelete(product)">
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

      <!-- Create/Edit Product Modal -->
      <Dialog v-model:open="showModal">
        <DialogContent class="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>{{ isEditing ? 'Edit Product' : 'Create Product' }}</DialogTitle>
            <DialogDescription>
              {{ isEditing ? 'Update the details of this financial product.' : 'Add a new financial product to the platform.' }}
            </DialogDescription>
          </DialogHeader>
          <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="name" class="text-right">Name</Label>
              <Input id="name" v-model="formData.name" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="institution" class="text-right">Institution</Label>
              <Select v-model="formData.financialInstitutionId" class="col-span-3">
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
              <Label for="type" class="text-right">Type</Label>
              <Select v-model="formData.type" class="col-span-3">
                <SelectTrigger>
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="INVOICE_FINANCING">Invoice Financing</SelectItem>
                  <SelectItem value="SUPPLY_CHAIN_FINANCING">Supply Chain Financing</SelectItem>
                  <SelectItem value="WORKING_CAPITAL">Working Capital</SelectItem>
                  <SelectItem value="TRADE_FINANCING">Trade Financing</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="description" class="text-right">Description</Label>
              <Input id="description" v-model="formData.description" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="minAmount" class="text-right">Min Amount</Label>
              <Input id="minAmount" v-model="formData.minAmount" type="number" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="maxAmount" class="text-right">Max Amount</Label>
              <Input id="maxAmount" v-model="formData.maxAmount" type="number" class="col-span-3" />
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
            <Button @click="saveProduct">Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <!-- Delete Confirmation Dialog -->
      <Dialog v-model:open="showDeleteDialog">
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirm Deletion</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete {{ selectedProduct?.name }}? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" @click="showDeleteDialog = false">Cancel</Button>
            <Button variant="destructive" @click="deleteProduct">Delete</Button>
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
const products = ref<any[]>([]);
const institutions = ref<any[]>([]);
const totalProducts = ref(0);
const activeProducts = ref(0);
const totalInstitutions = ref(0);
const totalTerms = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const totalPages = ref(1);
const searchQuery = ref('');
const statusFilter = ref('');
const institutionFilter = ref('');
const showModal = ref(false);
const isEditing = ref(false);
const selectedProduct = ref<any>(null);
const showDeleteDialog = ref(false);

const formData = ref({
  name: '',
  financialInstitutionId: 0,
  type: '',
  description: '',
  minAmount: 0,
  maxAmount: 0,
  isActive: true
});

// Fetch financial institution products
const fetchProducts = async () => {
  isLoading.value = true;
  try {
    let url = `/admin/financial-institutions/products?page=${currentPage.value}&pageSize=${pageSize.value}`;
    
    if (searchQuery.value) {
      url += `&name=${searchQuery.value}`;
    }
    
    if (statusFilter.value) {
      url += `&is_active=${statusFilter.value}`;
    }
    
    if (institutionFilter.value) {
      url += `&financial_institution_id=${institutionFilter.value}`;
    }
    
    const response = await apiClient.get(url);
    products.value = response.data.products || [];
    totalProducts.value = response.data.total || 0;
    totalPages.value = Math.ceil(totalProducts.value / pageSize.value);
    
    // Count active products
    activeProducts.value = products.value.filter(prod => prod.isActive).length;
    
    // Fetch additional stats
    await fetchStats();
    
  } catch (error) {
    console.error("Failed to fetch financial institution products:", error);
    toast.error("Could not load financial institution products.");
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

// Fetch additional statistics
const fetchStats = async () => {
  try {
    // In a real implementation, you would fetch these from the API
    // For now, we'll use mock data
    totalTerms.value = 48;
  } catch (error) {
    console.error("Failed to fetch stats:", error);
  }
};

// Handle search
const handleSearch = () => {
  currentPage.value = 1; // Reset to first page when searching
  fetchProducts();
};

// Pagination
const changePage = (page: number) => {
  currentPage.value = page;
  fetchProducts();
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
    type: '',
    description: '',
    minAmount: 0,
    maxAmount: 0,
    isActive: true
  };
  showModal.value = true;
};

const editProduct = (product: any) => {
  isEditing.value = true;
  selectedProduct.value = product;
  formData.value = {
    name: product.name,
    financialInstitutionId: product.financialInstitutionId,
    type: product.type,
    description: product.description || '',
    minAmount: product.minAmount || 0,
    maxAmount: product.maxAmount || 0,
    isActive: product.isActive
  };
  showModal.value = true;
};

const viewProduct = (product: any) => {
  // In a real implementation, you would navigate to a detail page
  // For now, we'll just show a toast
  toast.info(`Viewing details for ${product.name}`);
};

const saveProduct = async () => {
  try {
    if (isEditing.value && selectedProduct.value) {
      // Update existing product
      await apiClient.put(`/admin/financial-institutions/products/${selectedProduct.value.id}`, formData.value);
      toast.success(`${formData.value.name} updated successfully!`);
    } else {
      // Create new product
      await apiClient.post('/admin/financial-institutions/products', formData.value);
      toast.success(`${formData.value.name} created successfully!`);
    }
    showModal.value = false;
    fetchProducts();
  } catch (error) {
    console.error("Failed to save financial institution product:", error);
    let errorMessage = "Could not save financial institution product.";
    if (axios.isAxiosError(error) && error.response) {
      errorMessage += ` (Status: ${error.response.status})`;
    }
    toast.error(errorMessage);
  }
};

// Delete functions
const confirmDelete = (product: any) => {
  selectedProduct.value = product;
  showDeleteDialog.value = true;
};

const deleteProduct = async () => {
  if (!selectedProduct.value) return;
  
  try {
    await apiClient.delete(`/admin/financial-institutions/products/${selectedProduct.value.id}`);
    toast.success(`${selectedProduct.value.name} deleted successfully!`);
    showDeleteDialog.value = false;
    fetchProducts();
  } catch (error) {
    console.error("Failed to delete financial institution product:", error);
    let errorMessage = "Could not delete financial institution product.";
    if (axios.isAxiosError(error) && error.response) {
      errorMessage += ` (Status: ${error.response.status})`;
    }
    toast.error(errorMessage);
  }
};

// Initialize
onMounted(() => {
  if (authToken) {
    fetchInstitutions().then(() => {
      fetchProducts();
    });
  } else {
    toast.error("Authentication token not found. Please log in.");
    isLoading.value = false;
  }
});
</script>