<template>
  <div
    class="mb-4 relative"
    v-bind="$attrs"
  >
    <label
      v-if="item?.label || item?.title || label"
      :for="`field-${item?.name || name}`"
      class="block font-medium text-sm text-slate-700 mb-2"
    >
      {{ item?.label || item?.title || label }}
      <span
        v-if="item?.required || required"
        class="text-red-500"
      >*</span>
    </label>

    <div class="relative">
      <!-- Search Input -->
      <div class="relative">
        <input
          :id="`field-${item?.name || name}`"
          v-model="searchText"
          type="text"
          :placeholder="item?.placeholder || placeholder"
          :class="['w-full bg-slate-50 border border-slate-200 rounded-lg pl-10 pr-4 py-2.5 text-sm transition-all duration-200 hover:border-slate-300 hover:bg-white focus:outline-none focus:border-primary focus:bg-white focus:ring-2 focus:ring-primary/10', { 'skeleton-pulse': isLoading }]"
          @focus="onFocus"
          @blur="onBlur"
          @input="filterOptions"
        >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400"
          viewBox="0 0 20 20"
          fill="currentColor"
        >
          <path
            fill-rule="evenodd"
            d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
            clip-rule="evenodd"
          />
        </svg>
      </div>

      <!-- Dropdown Results -->
      <div
        v-if="isOpen && filteredOptions.length > 0"
        class="absolute z-50 w-full mt-2 bg-white border border-slate-200 rounded-xl shadow-[0_8px_30px_rgb(0,0,0,0.12)] max-h-60 overflow-y-auto"
      >
        <button
          v-for="opt in filteredOptions"
          :key="opt.value"
          type="button"
          class="w-full text-left px-4 py-3 hover:bg-slate-50 focus:bg-slate-50 focus:outline-none border-b border-slate-100 last:border-0 transition-colors flex items-center justify-between gap-2"
          @click="selectOption(opt)"
          @mousedown.prevent="selectOption(opt)"
        >
          <div class="flex-1 min-w-0">
            <div class="font-medium text-slate-800 truncate hover:text-primary transition-colors">
              {{ opt.label }}
            </div>
          </div>
          <div
            v-if="opt.description"
            class="text-xs text-slate-500 bg-slate-100 hover:bg-slate-200 px-2 py-1 rounded-md shrink-0 transition-colors whitespace-nowrap"
          >
            {{ opt.description }}
          </div>
        </button>
      </div>

      <div
        v-if="isOpen && searchText && filteredOptions.length === 0"
        class="absolute z-50 w-full mt-2 bg-white border border-slate-200 rounded-xl shadow-[0_8px_30px_rgb(0,0,0,0.12)] p-4 text-center text-sm text-slate-500"
      >
        Keine Ergebnisse gefunden
      </div>
    </div>

    <p
      v-if="item?.help || help"
      class="text-xs text-slate-500 mt-1.5"
    >
      {{ item?.help || help }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { SelectOption, FormItem } from '@/types'

interface ExtendedOption extends SelectOption {
  description?: string
}

interface Props {
  item?: FormItem
  name?: string
  label?: string
  modelValue: string
  options?: ExtendedOption[]
  placeholder?: string
  help?: string
  required?: boolean
  isLoading?: boolean
  clearOnFocus?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const activeOptions = computed(() => (props.item?.options as ExtendedOption[]) || props.options || [])

// Init search text from modelValue
const initialOption = activeOptions.value.find(o => o.value === props.modelValue)
const searchText = ref(initialOption?.label || '')
const isOpen = ref(false)

// Sync external modelValue changes
watch(() => props.modelValue, (newVal) => {
  const opt = activeOptions.value.find(o => o.value === newVal)
  if (opt && opt.label !== searchText.value) {
    searchText.value = opt.label
  } else if (!newVal) {
    searchText.value = ''
  }
})

const filteredOptions = computed(() => {
  if (!searchText.value || activeOptions.value.find(o => o.label === searchText.value)) return activeOptions.value
  const term = searchText.value.toLowerCase()
  return activeOptions.value.filter((opt) =>
    opt.label.toLowerCase().includes(term) || opt.value.toLowerCase().includes(term)
  )
})

function selectOption(opt: ExtendedOption) {
  searchText.value = opt.label
  isOpen.value = false
  emit('update:modelValue', opt.value)
}

function onFocus() {
  isOpen.value = true
  if (props.clearOnFocus && props.modelValue) {
    searchText.value = ''
    emit('update:modelValue', '')
  }
}

function filterOptions() {
  isOpen.value = true
  if (props.modelValue) {
    emit('update:modelValue', '') // Clear selection when typing
  }
}

function onBlur() {
  // If no exact match and clicking away, reset to selected or empty
  setTimeout(() => {
    isOpen.value = false
    const match = activeOptions.value.find(o => o.value === props.modelValue)
    if (match) {
      searchText.value = match.label
    } else {
      searchText.value = ''
    }
  }, 200)
}
</script>
