<template>
  <main>

    <section class="bg-customGreen border-b rounded-b-2xl ">
      <div class="flex mb-2 gap-8 pt-10 items-center py-2 px-4 mx-auto max-w-screen-xl ">
          <h2 class="text-2xl tracking-tight font-bold text-black">Welcome,
            <span class=" bg-gradient-to-r from-sky-500 via-purple-500 to-pink-500 bg-clip-text text-transparent " >
            Andrew Chakdahah </span></h2>
      </div>
    </section>

    <div class="relative bg-gray-700 isolate flex items-center gap-x-6 overflow-hidden bg-gray-50 px-6 py-2.5 sm:px-3.5 sm:before:flex-1">
      <div class="flex  flex-wrap items-center gap-x-4 gap-y-2">
        <p class="text-sm/6 text-white">
          <strong class="font-semibold">Compliance Form Incomplete</strong><svg viewBox="0 0 2 2" class="mx-2 inline size-0.5 fill-current" aria-hidden="true"><circle cx="1" cy="1" r="1" /></svg>To get more out of your score kindly complete the compliance forms.
        </p>
        <NuxtLink to="/compliance" class="flex-none rounded-full bg-blue-500 px-3.5 py-1 text-xs font-semibold text-white shadow-sm hover:bg-blue-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900">Complete now <span aria-hidden="true">&rarr;</span></NuxtLink>
      </div>
      <div class="flex flex-1 justify-end">
        <button type="button" class="-m-3 p-3 focus-visible:outline-offset-[-4px]">
          <span class="sr-only">Dismiss</span>
          <XMarkIcon class="size-5 text-white" aria-hidden="true" />
        </button>
      </div>
    </div>

    <section class="bg-gray-100">
      <div class="py-8 px-4 mx-auto max-w-screen-xl lg:py-16 lg:px-6">
        <div class="flex items-center">
          <h2 class="mb-2 text-xl tracking-tight font-semibold text-black">Invoices</h2>

          <div class="ml-auto">
            <Dialog>
              <DialogTrigger as-child>
                <Button class="text-white bg-black"><Upload />Upload Invoice</Button>
              </DialogTrigger>
              <DialogContent class="sm:max-w-[425px] shadow-lg">
                <DialogHeader>
                  <DialogTitle>Upload Invoice</DialogTitle>
                  <DialogDescription>
                    Upload your invoice here (pdf, csv, xlsx)
                  </DialogDescription>
                </DialogHeader>
                  <div class="grid w-full max-w-sm items-center gap-1.5">
                    <Label for="picture">Inovice</Label>
                    <Input id="picture" type="file" />
                  </div>

                <DialogFooter>
                  <Button type="submit">
                    Upload Invoice
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>
          </div>
        </div>


        <div class="py-4 px-4 mx-auto max-w-screen-xl lg:py-8 lg:px-6">
          <div class="bg-white mb-4 rounded-lg p-4 shadow-2xl">
            <Accordion type="single" class="w-full mt-4" collapsible>
              <!-- Loop through the loans array -->
              <AccordionItem class="mx-4  " value="">
                <AccordionTrigger class="mx-8 mb-2 border-b-2 items-start">
                  <div class="grid gap-4 md:flex">
                    <Avatar class="relative overflow-visible">
                      <AvatarImage class="rounded-full" alt="Customer Avatar" />
                      <AvatarFallback class="bg-black text-white">
                        CS
                      </AvatarFallback>
                    </Avatar>
                    <div class="text-lg">
                      Manufacturor Company Ltd.
                    </div>
                    <div class="text-lg text-muted-foreground">
                      548625
                    </div>
                  </div>
                  <div class=" gap-4 ml-auto  md:flex">
                    <div class="text-sm">
                      Loan Amount:
                      <div class="text-2xl text-black font-semibold">
                        GHS 12,569
                      </div>
                    </div>
                  </div>
                </AccordionTrigger>
                <AccordionContent>
                  <div class="py-6 px-4 mx-auto max-w-screen-xl">
                    <Stepper>
                      <StepperItem
                          v-for="item in steps"
                          :key="item.step"
                          class="basis-1/3"
                          :step="item.step"
                      >
                        <StepperTrigger>
                          <StepperIndicator>
                            <component :is="item.icon" class="w-4 h-4" />
                          </StepperIndicator>
                          <div class="flex flex-col">
                            <StepperTitle>
                              {{ item.title }}
                            </StepperTitle>
                            <StepperDescription>
                              {{ item.description }}
                            </StepperDescription>
                          </div>
                        </StepperTrigger>
                        <StepperSeparator
                            v-if="item.step !== steps[steps.length - 1].step"
                            class="w-full h-px"
                        />
                      </StepperItem>
                    </Stepper>
                    <h2 class="mt-4 mb-4 text-xl tracking-tight font-semibold text-black">Invoice Details</h2>
                    <div class="space-y-8 md:grid md:grid-cols-3 lg:grid-cols-3 md:gap-12 md:space-y-0">
                      <div class="col-span-2 max-w-screen-lg text-gray-500 sm:text-lg dark:text-gray-400">
                        <h3 class="font-semibold text-black">                      Manufacturor Company Ltd.
                        </h3>
                        <p class="text-sm">4507 Lake Floyd Circle</p>
                        <p class="text-sm">Chevy Chase</p>
                        <p class="text-sm">20815 Maryland</p>
                      </div>
                      <div class=" col-span-1 md:grid md:grid-cols-2 lg:grid-cols-2">
                        <p class="text-sm">Invoice #</p>
                        <p class="text-black">548625</p>
                        <p class="text-sm">Invoice Date </p>
                        <p class="text-black">10 May 2025</p>
                        <p class="text-sm">Terms</p>
                        <p class="text-black">Due on receipt</p>
                        <p class="text-sm">Due Date</p>
                        <p class="text-black">15 May 2025</p>
                      </div>
                    </div>
                    <Table class="mt-4 w-full min-w-[768px]">
                      <TableHeader>
                        <TableRow>
                          <TableHead>ID</TableHead>
                          <TableHead>ITEM & DESCRIPTION</TableHead>
                          <TableHead>UNIT PRICE</TableHead>
                          <TableHead>QUANTITY</TableHead>
                          <TableHead>AMOUNT</TableHead>
                        </TableRow>
                      </TableHeader>

                      <TableBody>
                        <TableRow class="" v-for="item in items" :key="item.id">
                          <TableCell>{{ item.id }}</TableCell>
                          <TableCell>
                            <div>
                              {{ item.name }}
                            </div>
                            <div class="text-xs font-light text-muted-foreground">{{ item.description }}</div></TableCell>
                          <TableCell>{{ item.unit_price }}</TableCell>
                          <TableCell>{{ item.quantity }}</TableCell>

                          <TableCell>{{ item.amount }}</TableCell>
                        </TableRow>
                      </TableBody>
                    </Table>
                    <div class="mt-2 max-w-screen-lg bg-gray-200 p-4 text-gray-500 sm:text-lg dark:text-gray-400">
                      <div class="flex col-span-1 md:grid md:grid-cols-3 lg:grid-cols-3">
                        <div class="col-span-2">
                          <p class="text-sm">Invoice total</p>
                        </div>
                        <div class="col-span-1">
                          <div class="flex items-center">
                          <p class="text-sm">Sub Total</p>
                          <p class="text-black font-semibold ml-auto text-sm">GHS 25,487</p>
                          </div>

                          <div class="flex items-center">
                            <p class="text-sm">Tax Rate</p>
                            <p class="text-black font-semibold ml-auto text-sm">5.00%</p>
                          </div>

                          <div class="flex items-center">
                            <p class="text-sm"> Total</p>
                            <p class="text-black font-semibold ml-auto text-sm">GHS 25,487</p>
                          </div>

                          <div class="flex items-center">
                            <p class="text-sm">Balance Due</p>
                            <p class="text-black font-semibold ml-auto text-sm">GHS 25,487</p>
                          </div>
                        </div>

                      </div>

                    </div>

<div class="flex">
  <div class="ml-auto mt-4">
    <Dialog>
      <DialogTrigger as-child>
        <Button class="text-white bg-black px-6 text-sm">  <Eye /> View Receipt</Button>
      </DialogTrigger>
      <DialogContent class="sm:max-w-[625px] shadow-lg">
        <DialogHeader class="bg-gray-900 p-6 rounded-lg">
          <DialogTitle class="text-white">InvoiceB2B</DialogTitle>
        </DialogHeader>
        <DialogTitle class="px-6 py-2">Transfer Receipt</DialogTitle>
        <div class=" col-span-1 p-6 md:grid md:grid-cols-2 lg:grid-cols-2">
          <p class="text-sm mb-2">Reference Number:</p>
          <p class="text-black mb-2">AT3OUPA14002CC27</p>
          <p class="text-sm mb-2">Transfer to: </p>
          <p class="text-black mb-2">Manufacturer Company Ltd.</p>
          <p class="text-sm mb-2">Account Type:</p>
          <p class="text-black mb-2">Current</p>
          <p class="text-sm mb-2">Account Number:</p>
          <p class="text-black mb-2">xxxxxxxxxxxx0028</p>
          <p class="text-sm mb-2">Account Name:</p>
          <p class="text-black mb-2">Manufacturer Company Ltd.</p>
          <p class="text-sm mb-2">Amount:</p>
          <p class="text-black mb-2">GHS 25,000.00</p>
          <p class="text-sm mb-2">Transfer Date:</p>
          <p class="text-black mb-2">02 May 2025</p>
          <p class="text-sm mb-2">Purpose:</p>
          <p class="text-black mb-2">Transaction is an invoicing for Syentia Retailer</p>
        </div>
        <DialogDescription>
          This is a computer generated receipt no signature required.
          Electronic Receipt owns no official legal effect, you may go to branch to get the paper receipt.
        </DialogDescription>
        <DialogFooter>
          <Button class="text-white ml-2 bg-blue-800 px-6 text-sm"> <Download  /> Download Receipt</Button>
        </DialogFooter>

      </DialogContent>
    </Dialog>


    <Button class="text-white ml-2 bg-blue-800 px-6 text-sm"> <Download  /> Download Receipt</Button>

  </div>
</div>

                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </div>

          <div class="bg-white mb-4 rounded-lg p-4 shadow-2xl">
            <Accordion type="single" class="w-full mt-4" collapsible>
              <!-- Loop through the loans array -->
              <AccordionItem class="mx-4  " value="">
                <AccordionTrigger class="mx-8 mb-2 border-b-2 items-start">
                  <div class="grid gap-4 md:flex">
                    <Avatar class="relative overflow-visible">
                      <AvatarImage class="rounded-full" alt="Customer Avatar" />
                      <AvatarFallback class="bg-black text-white">
                        CS
                      </AvatarFallback>
                    </Avatar>
                    <div class="text-lg">
                      Manufacturor Company Ltd.
                    </div>
                    <div class="text-lg text-muted-foreground">
                      548625
                    </div>
                  </div>
                  <div class=" gap-4 ml-auto  md:flex">
                    <div class="text-sm">
                      Loan Amount:
                      <div class="text-2xl text-black font-semibold">
                        GHS 12,569
                      </div>
                    </div>
                  </div>
                </AccordionTrigger>
                <AccordionContent>
                  <div class="py-6 px-4 mx-auto max-w-screen-xl">
                    <Stepper>
                      <StepperItem
                          v-for="item in steps"
                          :key="item.step"
                          class="basis-1/3"
                          :step="item.step"
                      >
                        <StepperTrigger>
                          <StepperIndicator>
                            <component :is="item.icon" class="w-4 h-4" />
                          </StepperIndicator>
                          <div class="flex flex-col">
                            <StepperTitle>
                              {{ item.title }}
                            </StepperTitle>
                            <StepperDescription>
                              {{ item.description }}
                            </StepperDescription>
                          </div>
                        </StepperTrigger>
                        <StepperSeparator
                            v-if="item.step !== steps[steps.length - 1].step"
                            class="w-full h-px"
                        />
                      </StepperItem>
                    </Stepper>
                    <h2 class="mt-4 mb-4 text-xl tracking-tight font-semibold text-black">Invoice Details</h2>
                    <div class="space-y-8 md:grid md:grid-cols-3 lg:grid-cols-3 md:gap-12 md:space-y-0">
                      <div class="col-span-2 max-w-screen-lg text-gray-500 sm:text-lg dark:text-gray-400">
                        <h3 class="font-semibold text-black">                      Manufacturor Company Ltd.
                        </h3>
                        <p class="text-sm">4507 Lake Floyd Circle</p>
                        <p class="text-sm">Chevy Chase</p>
                        <p class="text-sm">20815 Maryland</p>
                      </div>
                      <div class=" col-span-1 md:grid md:grid-cols-2 lg:grid-cols-2">
                        <p class="text-sm">Invoice #</p>
                        <p class="text-black">548625</p>
                        <p class="text-sm">Invoice Date </p>
                        <p class="text-black">10 May 2025</p>
                        <p class="text-sm">Terms</p>
                        <p class="text-black">Due on receipt</p>
                        <p class="text-sm">Due Date</p>
                        <p class="text-black">15 May 2025</p>
                      </div>
                    </div>
                    <Table class="mt-4 w-full min-w-[768px]">
                      <TableHeader>
                        <TableRow>
                          <TableHead>ID</TableHead>
                          <TableHead>ITEM & DESCRIPTION</TableHead>
                          <TableHead>UNIT PRICE</TableHead>
                          <TableHead>QUANTITY</TableHead>
                          <TableHead>AMOUNT</TableHead>
                        </TableRow>
                      </TableHeader>

                      <TableBody>
                        <TableRow class="" v-for="item in items" :key="item.id">
                          <TableCell>{{ item.id }}</TableCell>
                          <TableCell>
                            <div>
                              {{ item.name }}
                            </div>
                            <div class="text-xs font-light text-muted-foreground">{{ item.description }}</div></TableCell>
                          <TableCell>{{ item.unit_price }}</TableCell>
                          <TableCell>{{ item.quantity }}</TableCell>

                          <TableCell>{{ item.amount }}</TableCell>
                        </TableRow>
                      </TableBody>
                    </Table>
                    <div class="mt-2 max-w-screen-lg bg-gray-200 p-4 text-gray-500 sm:text-lg dark:text-gray-400">
                      <div class="flex col-span-1 md:grid md:grid-cols-3 lg:grid-cols-3">
                        <div class="col-span-2">
                          <p class="text-sm">Invoice total</p>
                        </div>
                        <div class="col-span-1">
                          <div class="flex items-center">
                            <p class="text-sm">Sub Total</p>
                            <p class="text-black font-semibold ml-auto text-sm">GHS 25,487</p>
                          </div>

                          <div class="flex items-center">
                            <p class="text-sm">Tax Rate</p>
                            <p class="text-black font-semibold ml-auto text-sm">5.00%</p>
                          </div>

                          <div class="flex items-center">
                            <p class="text-sm"> Total</p>
                            <p class="text-black font-semibold ml-auto text-sm">GHS 25,487</p>
                          </div>

                          <div class="flex items-center">
                            <p class="text-sm">Balance Due</p>
                            <p class="text-black font-semibold ml-auto text-sm">GHS 25,487</p>
                          </div>
                        </div>

                      </div>

                    </div>

                    <div class="flex">
                      <div class="ml-auto mt-4">
                        <Dialog>
                          <DialogTrigger as-child>
                            <Button class="text-white bg-black px-6 text-sm">  <Eye /> View Receipt</Button>
                          </DialogTrigger>
                          <DialogContent class="sm:max-w-[625px] shadow-lg">
                            <DialogHeader class="bg-gray-900 p-6 rounded-lg">
                              <DialogTitle class="text-white">InvoiceB2B</DialogTitle>
                            </DialogHeader>
                            <DialogTitle class="px-6 py-2">Transfer Receipt</DialogTitle>
                            <div class=" col-span-1 p-6 md:grid md:grid-cols-2 lg:grid-cols-2">
                              <p class="text-sm mb-2">Reference Number:</p>
                              <p class="text-black mb-2">AT3OUPA14002CC27</p>
                              <p class="text-sm mb-2">Transfer to: </p>
                              <p class="text-black mb-2">Manufacturer Company Ltd.</p>
                              <p class="text-sm mb-2">Account Type:</p>
                              <p class="text-black mb-2">Current</p>
                              <p class="text-sm mb-2">Account Number:</p>
                              <p class="text-black mb-2">xxxxxxxxxxxx0028</p>
                              <p class="text-sm mb-2">Account Name:</p>
                              <p class="text-black mb-2">Manufacturer Company Ltd.</p>
                              <p class="text-sm mb-2">Amount:</p>
                              <p class="text-black mb-2">GHS 25,000.00</p>
                              <p class="text-sm mb-2">Transfer Date:</p>
                              <p class="text-black mb-2">02 May 2025</p>
                              <p class="text-sm mb-2">Purpose:</p>
                              <p class="text-black mb-2">Transaction is an invoicing for Syentia Retailer</p>
                            </div>
                            <DialogDescription>
                              This is a computer generated receipt no signature required.
                              Electronic Receipt owns no official legal effect, you may go to branch to get the paper receipt.
                            </DialogDescription>
                            <DialogFooter>
                              <Button class="text-white ml-2 bg-blue-800 px-6 text-sm"> <Download  /> Download Receipt</Button>
                            </DialogFooter>

                          </DialogContent>
                        </Dialog>


                        <Button class="text-white ml-2 bg-blue-800 px-6 text-sm"> <Download  /> Download Receipt</Button>

                      </div>
                    </div>

                  </div>
                </AccordionContent>
              </AccordionItem>
            </Accordion>
          </div>

        </div>

      </div>
    </section>

    <div class="relative hidden md:block skyline" >
      <img src="../../assets/City%20skyline-pana.svg" class="h-96 opacity-10  absolute bottom-0 left-0 w-auto" alt="Closed Stores">
      <img src="../../assets/City%20skyline-pana.svg" class="h-96 opacity-10 absolute bottom-0 right-0 w-auto" alt="City Skyline">
    </div>
  </main>
</template>
<script setup lang="ts">
import {
  Upload ,
  Download,
  Eye,
  CreditCard,
  CircleCheckBig,
  Package2,
  Search,
  Users,
  BookUser,
  Check,
  Truck,
  AlertCircle, Plus
} from 'lucide-vue-next'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
const error = ref<string | null>(null);

import { Stepper, StepperDescription, StepperIndicator, StepperItem, StepperSeparator, StepperTitle, StepperTrigger } from '@/components/ui/stepper'
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "~/components/ui/dialog";
import {Input} from "~/components/ui/input";
import {Button} from "~/components/ui/button";

const steps = [{
  step: 1,
  title: 'Submitted',
  description: 'Invoice submitted successfully',
  icon: BookUser,
}, {
  step: 2,
  title: 'Approved',
  description: 'Your invoice has been approved by the bank.',
  icon: CircleCheckBig ,
}, {
  step: 3,
  title: 'Disbursed',
  description: 'Amount has been disbursed to the manufacturers bank account.',
  icon: CreditCard,
},]

const items = [
  {
    id: 1,
    name: 'Gaming Keyboard',
    description: 'Mechanical Keyboard designed for gaming purposes',
    quantity: '1',
    unit_price: 'GHS 155',
    amount: 'GHS 155',
  },
  {
    id: 2,
    name: 'Portable SSD',
    description: 'Compact and high speed external storage',
    quantity: '1',
    unit_price: 'GHS 355',
    amount: 'GHS 355',
  },
  {
    id: 3,
    name: 'Wireless Webcam',
    description: 'wireless webcam with 1080p video and 120fps',
    quantity: '1',
    unit_price: 'GHS 670',
    amount: 'GHS 670',
  },
]
</script>

<style>
.score {
  display: flex;
  justify-content: center;
  gap: 20px;
}

.skyline {
  z-index: 0;
}
</style>

