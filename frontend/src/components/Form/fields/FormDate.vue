<template>
  <div
    class="mb-4 w-full"
    v-bind="$attrs"
  >
    <label
      v-if="item.label || item.title"
      :for="`field-${item.name}`"
      class="block font-medium text-sm text-slate-700 mb-2"
    >
      {{ item.label || item.title }}
      <span
        v-if="item.required"
        class="text-red-500"
      >*</span>
    </label>

    <div class="relative">
      <input
        :id="`field-${item.name}`"
        ref="inputRef"
        type="text"
        :value="modelValue"
        :placeholder="item.placeholder || 'Datum auswählen...'"
        :required="item.required"
        class="w-full bg-slate-50 border border-slate-200 rounded-lg px-4 py-2.5 pl-10 text-sm transition-all duration-200 hover:border-slate-300 hover:bg-white focus:outline-none focus:border-primary focus:bg-white focus:ring-2 focus:ring-primary/10 cursor-pointer"
      >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none"
        viewBox="0 0 20 20"
        fill="currentColor"
      >
        <path
          fill-rule="evenodd"
          d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
          clip-rule="evenodd"
        />
      </svg>
    </div>

    <p
      v-if="item.help"
      class="text-xs text-slate-500 mt-1.5"
    >
      {{ item.help }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import flatpickr from 'flatpickr'
import { German } from 'flatpickr/dist/l10n/de.js'
import 'flatpickr/dist/flatpickr.css'
import type { FormItem } from '@/types'

interface Props {
  item: FormItem
  modelValue: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const inputRef = ref<HTMLInputElement | null>(null)
let fpInstance: flatpickr.Instance | null = null

onMounted(() => {
  if (inputRef.value) {
    fpInstance = flatpickr(inputRef.value, {
      locale: German,
      enableTime: props.item.includeTime || false,
      dateFormat: props.item.includeTime ? 'Y-m-d\\TH:i:S' : 'Y-m-d',
      altInput: true,
      altFormat: props.item.includeTime ? 'd.m.Y H:i' : 'd.m.Y',
      defaultDate: props.modelValue || null,
      onChange: (selectedDates) => {
        // Send format: 2026-07-12T15:00:00+02:00
        if (selectedDates.length > 0) {
           const d = selectedDates[0]
           // ISO String output but respecting local time offset
           const offset = d.getTimezoneOffset()
           const pad = (n: number) => n < 10 ? '0' + n : n
           const sign = offset > 0 ? '-' : '+'
           const absOffset = Math.abs(offset)
           const hoursMenu = pad(Math.floor(absOffset / 60))
           const minsMenu = pad(absOffset % 60)
           
           const y = d.getFullYear()
           const m = pad(d.getMonth() + 1)
           const dd = pad(d.getDate())
           const h = pad(d.getHours())
           const i = pad(d.getMinutes())
           const s = pad(d.getSeconds())

           const formatted = `${y}-${m}-${dd}T${h}:${i}:${s}${sign}${hoursMenu}:${minsMenu}`
           emit('update:modelValue', formatted)
        } else {
           emit('update:modelValue', '')
        }
      }
    })
  }
})

watch(() => props.modelValue, (newVal) => {
  if (fpInstance) {
    if (newVal !== inputRef.value?.value) {
      fpInstance.setDate(newVal, false)
    }
  }
})

onBeforeUnmount(() => {
  if (fpInstance) {
    fpInstance.destroy()
  }
})
</script>

<style>
.flatpickr-calendar {
  width: auto !important;
  max-width: 100%;
}
.flatpickr-days {
  width: auto !important;
}
.dayContainer {
  width: auto !important;
  min-width: auto !important;
  max-width: auto !important;
}
</style>
