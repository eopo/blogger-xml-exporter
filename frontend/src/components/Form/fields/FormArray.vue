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

    <div
      :id="`field-${item.name}`"
      class="space-y-3"
    >
      <div
        v-for="(val, idx) in (modelValue || [])"
        :key="idx"
        class="p-3 bg-slate-50 rounded-lg border border-slate-200"
      >
        <div class="flex gap-2 items-start">
          <div class="flex-1">
            <!-- Object Array Mode (Render multiple fields per entry based on schema) -->
            <div
              v-if="isObjectMode()"
              class="grid gap-3"
              :style="{ gridTemplateColumns: `repeat(${item.fields?.length || 1}, minmax(0, 1fr))` }"
            >
              <div
                v-for="field in item.fields"
                :key="field.name"
              >
                <label
                  v-if="field.label"
                  class="block text-xs font-semibold text-slate-600 mb-1 tracking-wide uppercase"
                >{{ field.label }}</label>
                <input
                  type="text"
                  :value="val[field.name] || ''"
                  :placeholder="field.placeholder || 'Wert...'"
                  :class="['w-full bg-white border border-slate-200 rounded px-3 py-2 text-sm', { 'skeleton-pulse': isLoading }]"
                  @input="updateItemObject(idx, field.name, ($event.target as HTMLInputElement).value)"
                >
              </div>
            </div>
            
            <div v-else>
              <input
                type="text"
                :value="val"
                placeholder="Wert..."
                :class="['w-full bg-white border border-slate-200 rounded px-3 py-2 text-sm', { 'skeleton-pulse': isLoading }]"
                @input="updateItemPrimitive(idx, ($event.target as HTMLInputElement).value)"
              >
            </div>
          </div>
          
          <button
            type="button"
            class="px-3 py-2 bg-red-50 text-red-600 hover:bg-red-100 rounded text-sm font-medium h-fit mt-auto"
            :class="{'mt-5': isObjectMode()}"
            @click="removeItem(idx)"
          >
            Löschen
          </button>
        </div>
      </div>
    </div>

    <button
      type="button"
      class="mt-3 px-4 py-2 bg-slate-100 text-slate-700 hover:bg-slate-200 rounded-lg text-sm font-medium"
      @click="addItem"
    >
      + Add
    </button>

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
  modelValue: Record<string, unknown>[]
  isLoading?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: Record<string, unknown>[]]
}>()

const isObjectMode = () => props.item.fields && props.item.fields.length > 0

function addItem() {
  if (isObjectMode()) {
    const newItem: Record<string, unknown> = {}
    props.item.fields!.forEach(f => { newItem[f.name] = '' })
    emit('update:modelValue', [...(props.modelValue || []), newItem])
  } else {
    emit('update:modelValue', [...(props.modelValue || []), ''])
  }
}

function removeItem(idx: number) {
  emit('update:modelValue', (props.modelValue || []).filter((_, i) => i !== idx))
}

function updateItemPrimitive(idx: number, value: string) {
  const updated = [...(props.modelValue || [])]
  updated[idx] = value
  emit('update:modelValue', updated)
}

function updateItemObject(idx: number, fieldName: string, value: string) {
  const updated = [...(props.modelValue || [])]
  if (!updated[idx]) updated[idx] = {}
  updated[idx] = { ...updated[idx], [fieldName]: value }
  emit('update:modelValue', updated)
}
</script>
