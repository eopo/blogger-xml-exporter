<template>
  <div class="mb-4">
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

    <textarea
      v-if="item.type === 'textarea'"
      :id="`field-${item.name}`"
      :value="modelValue"
      :placeholder="item.placeholder"
      :required="item.required"
      :minlength="item.minLength"
      :maxlength="item.maxLength"
      :class="['w-full bg-slate-50 border border-slate-200 rounded-lg px-4 py-2.5 text-sm resize-none min-h-32 transition-all duration-200 hover:border-slate-300 hover:bg-white focus:outline-none focus:border-primary focus:bg-white focus:ring-2 focus:ring-primary/10', { 'skeleton-pulse': isLoading }]"
      @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
    />

    <input
      v-else
      :id="`field-${item.name}`"
      :type="inputType"
      :value="modelValue"
      :placeholder="item.placeholder"
      :required="item.required"
      :minlength="item.minLength"
      :maxlength="item.maxLength"
      :pattern="item.pattern"
      :class="['w-full bg-slate-50 border border-slate-200 rounded-lg px-4 py-2.5 text-sm transition-all duration-200 hover:border-slate-300 hover:bg-white focus:outline-none focus:border-primary focus:bg-white focus:ring-2 focus:ring-primary/10', { 'skeleton-pulse': isLoading }]"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    >

    <p
      v-if="item.help"
      class="text-xs text-slate-500 mt-1.5"
    >
      {{ item.help }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FormItem } from '@/types'

interface Props {
  item: FormItem
  modelValue: string | number
  isLoading?: boolean
}

const props = defineProps<Props>()

defineEmits<{
  'update:modelValue': [value: string | number]
}>()

const inputType = computed(() => (props.item.type === 'textarea' ? 'text' : props.item.type))
</script>
