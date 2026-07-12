<template>
  <div>
    <!-- Collapsible group: use details/summary for expandable sections -->
    <details
      v-if="group.collapsible"
      :open="!group.collapsed"
      class="border-b border-slate-200 py-6 last:border-b-0"
    >
      <summary
        v-if="group.title"
        class="cursor-pointer font-bold text-slate-900 text-lg tracking-wide mb-6"
      >
        {{ group.title }}
      </summary>
      <div class="space-y-4">
        <RenderGroupContent
          :group="group"
          :form-values="formValues"
          :is-loading="isLoading"
        />
      </div>
    </details>

    <div
      v-else
      class="border-b border-slate-200 py-6 last:border-b-0"
    >
      <div
        v-if="group.title"
        class="font-bold text-slate-900 text-lg tracking-wide mb-6"
      >
        {{ group.title }}
      </div>
      <div class="space-y-4">
        <RenderGroupContent
          :group="group"
          :form-values="formValues"
          :is-loading="isLoading"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FormGroup as FormGroupType, FormValues } from '@/types'
import RenderGroupContent from './RenderGroupContent.vue'

interface Props {
  group: FormGroupType
  formValues: FormValues
  isLoading?: boolean
}

defineProps<Props>()
</script>
