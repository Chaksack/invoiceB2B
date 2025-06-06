<template>

  <div class="bg-gray-100 dark:bg-black flex items-center justify-center font-sans">
    <ClientOnly>
      <div class="w-full max-w-lg mx-auto p-4">
        <div
            :class="[
            'w-full p-10 border-2 border-dashed rounded-xl transition-colors duration-300 cursor-pointer',
            {
              'border-sky-500 bg-sky-50 dark:border-sky-400 dark:bg-neutral-800': isActive,
              'border-neutral-300 bg-neutral-50 dark:border-neutral-700 dark:bg-neutral-900': !isActive
            }
          ]"
            @dragover.prevent="handleDragState(true)"
            @dragleave.prevent="handleDragState(false)"
            @drop.prevent="handleDrop"
            @click="triggerFileSelect"
        >
          <input
              ref="fileInputRef"
              type="file"
              class="hidden"
              accept=".pdf,.csv,.xlsx"
              @change="onFileChange"
          />

          <div class="flex flex-col items-center justify-center text-center">
            <Transition
                name="fade"
                mode="out-in"
            >
              <div
                  v-if="file"
                  key="file-info"
                  class="relative w-full bg-white dark:bg-neutral-800 p-4 rounded-lg shadow-md flex flex-col items-start"
              >
                <!-- Remove File Button with Inline SVG -->
                <button
                    @click.stop="handleRemoveFile"
                    class="absolute top-2 right-2 p-1.5 bg-neutral-200 dark:bg-neutral-700 rounded-full text-neutral-600 dark:text-neutral-300 hover:bg-red-200 dark:hover:bg-red-800 hover:text-red-600 dark:hover:text-red-500 transition-colors"
                    aria-label="Remove file"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
                </button>

                <!-- File Details with Inline SVG -->
                <div class="flex items-center gap-4 w-full">
                  <svg xmlns="http://www.w3.org/2000/svg" class="text-sky-500 shrink-0" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7Z"/><path d="M14 2v4a2 2 0 0 0 2 2h4"/><path d="M16 13H8"/><path d="M16 17H8"/><path d="M10 9H8"/></svg>
                  <div class="text-left flex-grow truncate">
                    <p class="font-semibold text-neutral-800 dark:text-neutral-200 truncate pr-8" :title="file.name">
                      {{ file.name }}
                    </p>
                    <div class="flex items-center gap-2 mt-1">
                      <p class="text-xs text-neutral-500 dark:text-neutral-400">
                        {{ formatFileSize(file.size) }}
                      </p>
                      <!-- File Type Badge -->
                      <p class="text-xs bg-gray-200 text-gray-700 dark:bg-neutral-700 dark:text-neutral-300 px-1.5 py-0.5 rounded-md">
                        {{ file.type || 'unknown type' }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Upload Prompt with Inline SVG -->
              <div
                  v-else
                  key="upload-prompt"
                  class="flex flex-col items-center"
              >
                <div class="p-2 transition-transform duration-300 ease-in-out hover:scale-110">
                  <svg xmlns="http://www.w3.org/2000/svg" class="text-neutral-400 dark:text-neutral-500" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 14.899A7 7 0 1 1 15.71 8h1.79a4.5 4.5 0 0 1 2.5 8.242"/><path d="M12 12v9"/><path d="m16 16-4-4-4 4"/></svg>
                </div>
                <p class="mt-4 font-bold text-neutral-700 dark:text-neutral-300">
                  Upload Invoice File
                </p>
                <p class="mt-1 text-sm text-neutral-500 dark:text-neutral-400">
                  Drag & drop or click to upload
                </p>
                <p class="mt-1 text-xs text-neutral-400 dark:text-neutral-500">
                  (PDF, CSV, XLSX)
                </p>
              </div>
            </Transition>
          </div>
        </div>
      </div>
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

// --- State and Refs ---

// Reactive state for the single uploaded file
const file = ref<File | null>(null);
// Reactive state to track drag-over status
const isActive = ref(false);
// Ref for the hidden file input element
const fileInputRef = ref<HTMLInputElement | null>(null);


// --- Helper Functions ---

// Formats file size from bytes to a readable format
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
};


// --- Event Handlers ---

// Sets the active state for the drop zone
const handleDragState = (active: boolean) => {
  isActive.value = active;
};

// Handles file drop event
const handleDrop = (e: DragEvent) => {
  isActive.value = false;
  const droppedFiles = e.dataTransfer?.files;
  if (droppedFiles && droppedFiles[0]) {
    file.value = droppedFiles[0]; // Set only the first file
  }
};

// Handles file selection from the file input
const onFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    file.value = target.files[0]; // Set only the first file
  }
};

// Programmatically clicks the hidden file input
const triggerFileSelect = () => {
  // Do not trigger if a file is already present (unless clicking the remove button)
  if (!file.value) {
    fileInputRef.value?.click();
  }
};

// Removes the currently selected file
const handleRemoveFile = () => {
  file.value = null;
  // Also reset the file input so the same file can be re-uploaded
  if(fileInputRef.value) {
    fileInputRef.value.value = '';
  }
};
</script>

<style scoped>
/* Simple fade transition for switching between views */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
