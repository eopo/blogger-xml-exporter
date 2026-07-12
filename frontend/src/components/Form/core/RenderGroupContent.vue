<template>
  <div>
    <!-- Presets dropdown -->
    <div
      v-if="group.presets && group.presets.length > 0"
      class="mb-6"
    >
      <FormCombobox
        :item="{
          name: 'preset',
          type: 'combobox',
          label: 'Vorlage für diese Gruppe anwenden',
          options: group.presets.map((p, i) => ({ value: i.toString(), label: p.label || p.title || 'Vorlage ' + (i + 1) })),
          placeholder: '-- Bitte wählen --',
          required: false,
        }"
        model-value=""
        @update:model-value="applyPresetRawValue"
      />
    </div>

    <!-- Render rows -->
    <template
      v-for="(row, rIdx) in rows"
      :key="rIdx"
    >
      <div
        class="grid gap-4"
        :style="{ gridTemplateColumns: gridTemplate(row) }"
      >
        <template
          v-for="field in row"
          :key="field.name"
        >
          <div v-if="field.type === 'group'">
            <FormGroup
              :group="(field as any)"
              :form-values="formValues"
            />
          </div>

          <div v-else-if="field.type === 'array'">
            <FormArray
              :item="field"
              :model-value="(formValues[field.name] as unknown as any[]) || []"
              @update:model-value="updateFormValue(field.name, $event)"
            />
          </div>

          <div v-else-if="field.type === 'date'">
            <FormDate
              :item="field"
              :model-value="(formValues[field.name] || '') as string"
              @update:model-value="updateFormValue(field.name, $event)"
            />
          </div>

          <div v-else-if="field.type === 'select'">
            <FormCombobox
              :item="field"
              :model-value="(formValues[field.name] || '') as string"
              @update:model-value="updateFormValue(field.name, $event)"
            />
          </div>

          <div v-else-if="field.type === 'combobox'">
            <FormCombobox
              :item="field"
              :model-value="(formValues[field.name] || '') as string"
              @update:model-value="updateFormValue(field.name, $event)"
            />
          </div>

          <div v-else-if="field.type === 'textarea'">
            <FormField
              :item="{ ...field, type: 'textarea' }"
              :model-value="(formValues[field.name] || '') as string"
              @update:model-value="updateFormValue(field.name, $event)"
            />
          </div>

          <div v-else>
            <!-- text, number, email, etc -->
            <FormField
              :item="field"
              :model-value="(formValues[field.name] || '') as string | number"
              @update:model-value="updateFormValue(field.name, $event)"
            />
          </div>
        </template>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FormGroup as FormGroupType, FormItem, FormValues } from '@/types'
import FormField from '../fields/FormField.vue'
import FormDate from '../fields/FormDate.vue'
import FormCombobox from '../fields/FormCombobox.vue'
import FormArray from '../fields/FormArray.vue'
import FormGroup from './FormGroup.vue'

interface Props {
  group: FormGroupType
  formValues: FormValues
}

const props = defineProps<Props>()

// Build rows: distribute fields across columns based on width
const rows = computed(() => {
  const result: FormItem[][] = []
  let currentRow: FormItem[] = []
  let currentRowWidth = 0

  // Guard against undefined fields/items properly, matching the schema definitions
  // The schema has `items: FormItem[]` for groups, not `fields`.
  const items = props.group.items || []

  for (const field of items) {
    const width = field.width || 12
    if (currentRowWidth + width > 12 && currentRow.length > 0) {
      result.push(currentRow)
      currentRow = [field]
      currentRowWidth = width
    } else {
      currentRow.push(field)
      currentRowWidth += width
    }
  }

  if (currentRow.length > 0) {
    result.push(currentRow)
  }

  return result
})

// Calculate CSS Grid template for row
function gridTemplate(row: FormItem[]): string {
  return row.map((f) => `${Math.min(f.width || 12, 12)}fr`).join(' ')
}

// Update form value (v-model binding)
// Note: Vue 3 reactive refs handle mutations correctly in template context
// eslint-disable vue/no-mutating-props
function updateFormValue(name: string, value: string | number | Record<string, unknown>): void {
  // eslint-disable-next-line vue/no-mutating-props
  props.formValues[name] = value
}

// Apply preset
// eslint-disable vue/no-mutating-props
function applyPresetRawValue(value: string): void {
  if (!value) return
  const idx = parseInt(value)
  if (isNaN(idx) || !props.group.presets) return

  const preset = props.group.presets[idx]
  // eslint-disable-next-line vue/no-mutating-props
  Object.assign(props.formValues, preset.values)
}
</script>
