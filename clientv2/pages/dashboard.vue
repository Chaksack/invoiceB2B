<script setup>
import { ref } from 'vue';

// This data would typically come from an API
const stats = ref([
  { title: 'Total Invoices', value: 'GHC 245,678.00', change: '+12.5%', changeType: 'positive' },
  { title: 'Financed Amount', value: 'GHC 187,432.00', change: '+8.2%', changeType: 'positive' },
  { title: 'Available Funds', value: 'GHC 58,246.00', change: '-3.1%', changeType: 'negative' },
  { title: 'Average Financing Rate', value: '4.2%', change: '-0.5%', changeType: 'positive' }
]);

const recentInvoices = ref([
  { id: 'INV-2025-001', client: 'Acme Corp', amount: 'GHC 12,500.00', date: '2025-07-25', status: 'Financed' },
  { id: 'INV-2025-002', client: 'Globex Inc', amount: 'GHC 8,750.00', date: '2025-07-24', status: 'Pending' },
  { id: 'INV-2025-003', client: 'Stark Industries', amount: 'GHC 15,200.00', date: '2025-07-22', status: 'Financed' },
  { id: 'INV-2025-004', client: 'Wayne Enterprises', amount: 'GHC 9,800.00', date: '2025-07-20', status: 'Paid' },
  { id: 'INV-2025-005', client: 'Umbrella Corp', amount: 'GHC 11,350.00', date: '2025-07-18', status: 'Financed' }
]);

const upcomingPayments = ref([
  { id: 'PAY-2025-001', invoice: 'INV-2025-001', amount: 'GHC 12,500.00', dueDate: '2025-08-25' },
  { id: 'PAY-2025-002', invoice: 'INV-2025-003', amount: 'GHC 15,200.00', dueDate: '2025-08-22' },
  { id: 'PAY-2025-003', invoice: 'INV-2025-005', amount: 'GHC 11,350.00', dueDate: '2025-08-18' }
]);

definePageMeta({
  layout: 'dashboard'
});
</script>

<template>
  <div>
    <!-- Page Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p class="text-gray-500 mt-1">Welcome back! Here's an overview of your invoice financing.</p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div v-for="stat in stats" :key="stat.title" class="bg-card rounded-lg shadow-sm p-4 border border-border">
        <div class="flex justify-between items-start">
          <div>
            <p class="text-sm font-medium text-gray-500">{{ stat.title }}</p>
            <p class="text-2xl font-bold mt-1">{{ stat.value }}</p>
          </div>
          <div :class="[
            'text-sm font-medium px-2 py-1 rounded-full',
            stat.changeType === 'positive' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
          ]">
            {{ stat.change }}
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Recent Invoices -->
      <div class="lg:col-span-2 bg-card rounded-lg shadow-sm border border-border overflow-hidden">
        <div class="p-4 border-b border-border flex justify-between items-center">
          <h2 class="font-semibold text-lg">Recent Invoices</h2>
          <NuxtLink to="/invoices" class="text-sm text-primary hover:underline">View All</NuxtLink>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="bg-muted/50">
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Invoice ID</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Client</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border">
              <tr v-for="invoice in recentInvoices" :key="invoice.id" class="hover:bg-muted/20">
                <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ invoice.id }}</td>
                <td class="px-4 py-3 text-sm text-gray-500">{{ invoice.client }}</td>
                <td class="px-4 py-3 text-sm text-gray-500">{{ invoice.amount }}</td>
                <td class="px-4 py-3 text-sm text-gray-500">{{ invoice.date }}</td>
                <td class="px-4 py-3 text-sm">
                  <span :class="[
                    'px-2 py-1 text-xs font-medium rounded-full',
                    invoice.status === 'Financed' ? 'bg-blue-100 text-blue-800' : 
                    invoice.status === 'Pending' ? 'bg-yellow-100 text-yellow-800' : 
                    'bg-green-100 text-green-800'
                  ]">
                    {{ invoice.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Upcoming Payments -->
      <div class="bg-card rounded-lg shadow-sm border border-border overflow-hidden">
        <div class="p-4 border-b border-border">
          <h2 class="font-semibold text-lg">Upcoming Payments</h2>
        </div>
        <div class="p-4 space-y-4">
          <div v-for="payment in upcomingPayments" :key="payment.id" class="border-b border-border pb-4 last:border-0 last:pb-0">
            <div class="flex justify-between items-start">
              <div>
                <p class="text-sm font-medium text-gray-900">{{ payment.invoice }}</p>
                <p class="text-xs text-gray-500 mt-1">Due on {{ payment.dueDate }}</p>
              </div>
              <p class="text-sm font-bold">{{ payment.amount }}</p>
            </div>
          </div>
        </div>
        <div class="p-4 border-t border-border">
          <NuxtLink to="/payments" class="block w-full py-2 bg-primary text-white rounded-md hover:bg-primary-700 transition-colors text-center">
            View All Payments
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="mt-6 bg-card rounded-lg shadow-sm border border-border p-4">
      <h2 class="font-semibold text-lg mb-4">Quick Actions</h2>
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4">
        <NuxtLink to="/submit-invoice" class="p-4 border border-border rounded-lg hover:bg-muted/20 transition-colors flex flex-col items-center justify-center">
          <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="text-primary mb-2">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m5 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
          </svg>
          <span class="text-sm font-medium">Submit Invoice</span>
        </NuxtLink>
        <NuxtLink to="/financing" class="p-4 border border-border rounded-lg hover:bg-muted/20 transition-colors flex flex-col items-center justify-center">
          <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="text-primary mb-2">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          <span class="text-sm font-medium">Apply for Financing</span>
        </NuxtLink>
        <NuxtLink to="/invoices" class="p-4 border border-border rounded-lg hover:bg-muted/20 transition-colors flex flex-col items-center justify-center">
          <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="text-primary mb-2">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"></path>
          </svg>
          <span class="text-sm font-medium">View Reports</span>
        </NuxtLink>
        <NuxtLink to="/payments" class="p-4 border border-border rounded-lg hover:bg-muted/20 transition-colors flex flex-col items-center justify-center">
          <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" class="text-primary mb-2">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
          </svg>
          <span class="text-sm font-medium">Payment Schedule</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>