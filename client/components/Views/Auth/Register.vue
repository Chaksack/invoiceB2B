<template>
  <div class="flex h-screen flex-col lg:grid lg:grid-cols-[1fr_1.8fr] overflow-hidden">
    <!-- Left side - Form -->
    <div class="flex flex-col items-center justify-center p-6 sm:p-12">
      <form class="mt-10 max-w-[800px]" @submit.prevent="register">
        <h1 class="text-2xl font-bold mb-4">
          Sign Up Account.
        </h1>
        <p class="text-muted-foreground mb-6">
          Enter your personal data to create your account.
        </p>
        <Alert v-if="errorMessage" variant="default" class="bg-red-500 mb-4 text-white">
          <AlertCircle class="w-4 h-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription class="text-white">
            {{ errorMessage }}
          </AlertDescription>
        </Alert>
        <div class="grid gap-4">
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="firstName">First Name</Label>
            </div>
            <Input
                v-model="firstName"
                id="firstName"
                type="name"
                placeholder="John"
                required
            />
          </div>

          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="lastName">Last Name</Label>
            </div>
            <Input
                v-model="lastName"
                id="lastName"
                type="name"
                placeholder="Doe"
                required
            />
          </div>

          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="companyName">Company Name</Label>
            </div>
            <Input
                v-model="companyName"
                id="companyName"
                type="text"
                placeholder="Company Ltd"
                required
            />
          </div>
          <!-- Email Input -->
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="email">Email Address</Label>
            </div>
            <Input
                v-model="email"
                id="email"
                type="email"
                placeholder="m@example.com"
                required
            />
          </div>
          <!-- Password Input -->
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Password</Label>
            </div>
            <Input
                v-model="password"
                id="password"
                type="password"
                required
            />
          </div>
        </div>
          <button class="py-2 mt-4 rounded-lg bg-bgC text-white bg-black w-full font-medium mb-4 hover:bg-gray-500">
            Create Account
          </button>
          <p class="text-muted-foreground mb-6">
            Already have an InvoiceFnd account? <a href="/login" class="text-black text-semibold">Login</a>
          </p>
      </form>

    </div>

    <!-- Right side - Image -->
    <div class="hidden lg:flex bg-black items-center justify-start h-screen">
      <div class="flex flex-col ml-10">
        <!--        <div class="flex text-sm text-white italic font-semibold tracking-wide">-->
        <!--          <img src="../../../assets/logo.png" class="h-20"><span class="text-white ml-2">OBLOS</span>-->
        <!--        </div>-->
        <div class="flex flex-col mt-6">
          <h2 class="text-white font-semibold text-6xl">A growth platform for small business</h2>
        </div>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import axios, { AxiosError } from 'axios';
import { useRouter } from 'vue-router';
import { AlertCircle } from 'lucide-vue-next';
import { toast } from 'vue-sonner';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { navigateTo } from '#app';

const email = ref('');
const password = ref('');
const firstName = ref('');
const lastName = ref('');
const companyName = ref('');

const isLoading = ref(false);
const router = useRouter();
const errorMessage = ref('');
const successMessage = ref('');

const register = async () => {
  isLoading.value = true;

  try {
    const response = await axios.post('http://localhost:3000/api/v1/auth/register', {
      email: email.value,
      firstName: firstName.value,
      lastName: lastName.value,
      companyName: companyName.value,
      password: password.value,
    },);

    if (response.status === 200) {
      const token = response.data.token;
      localStorage.setItem('token', token);
      localStorage.setItem('email', email.value);
      toast.success('Success:', {description: 'successMessage'});
      await navigateTo('/login');

    }
  } catch (error) {
    if (error instanceof AxiosError && error.response?.status === 409) {
      toast.error('Error: ', {description: error.message});
    } else {
      errorMessage.value = 'Invalid email or password';
      toast.error('Error: ', {description: error.message});
    }
  } finally {
    isLoading.value = false;
  }
};
</script>
