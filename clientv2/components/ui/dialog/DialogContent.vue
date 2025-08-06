<script setup>
import { ref, watch, onMounted, onBeforeUnmount, inject } from 'vue';

const props = defineProps({
  class: {
    type: String,
    default: ''
  }
});

// Get parent Dialog component's open state
const dialogOpen = inject('dialogOpen', ref(false));

// Handle click outside to close dialog
const dialogContent = ref(null);

const handleClickOutside = (event) => {
  if (dialogContent.value && !dialogContent.value.contains(event.target)) {
    dialogOpen.value = false;
  }
};

// Add/remove event listeners
onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener('mousedown', handleClickOutside);
});

// Prevent scrolling when dialog is open
watch(dialogOpen, (newValue) => {
  if (newValue) {
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = '';
  }
});
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div v-if="dialogOpen" class="fixed inset-0 z-50 flex items-center justify-center">
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/50" @click="dialogOpen = false"></div>
        
        <!-- Dialog content -->
        <div 
          ref="dialogContent"
          :class="['bg-white rounded-lg shadow-lg w-full max-w-md mx-auto p-6 z-50', props.class]"
        >
          <div class="relative">
            <!-- Close button -->
            <button 
              class="absolute top-2 right-2 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-700 focus:outline-none"
              @click="dialogOpen = false"
            >
              <svg width="15" height="15" viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M11.7816 4.03157C12.0062 3.80702 12.0062 3.44295 11.7816 3.2184C11.5571 2.99385 11.193 2.99385 10.9685 3.2184L7.50005 6.68682L4.03164 3.2184C3.80708 2.99385 3.44301 2.99385 3.21846 3.2184C2.99391 3.44295 2.99391 3.80702 3.21846 4.03157L6.68688 7.49999L3.21846 10.9684C2.99391 11.193 2.99391 11.557 3.21846 11.7816C3.44301 12.0061 3.80708 12.0061 4.03164 11.7816L7.50005 8.31316L10.9685 11.7816C11.193 12.0061 11.5571 12.0061 11.7816 11.7816C12.0062 11.557 12.0062 11.193 11.7816 10.9684L8.31322 7.49999L11.7816 4.03157Z" fill="currentColor" fill-rule="evenodd" clip-rule="evenodd"></path>
              </svg>
            </button>
            
            <slot></slot>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>