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

    <select
      :id="`field-${item.name}`"
      :value="modelValue"
      :required="item.required"
      class="w-full bg-slate-50 border border-slate-200 rounded-lg px-4 py-2.5 text-sm transition-all duration-200 hover:border-slate-300 hover:bg-white focus:outline-none focus:border-primary focus:bg-white focus:ring-2 focus:ring-primary/10"
      @input="emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
    >
      <option value="">
        -- Wählen --
      </option>
      <option
        v-for="opt in item.options"
        :key="opt.value"
        :value="opt.value"
      >
        {{ opt.label }}
      </option>
    </select>

    <p
      v-if="item.help"
      class="text-xs text-slate-500 mt-1.5"
    >
      {{ item.help }}
    </p>
  </div>
</template>

<script setup lang="ts">
import type { FormItem } from '@/types'

interface Props {
  item: FormItem
  modelValue: string
}

defineProps<Props>()

defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>
