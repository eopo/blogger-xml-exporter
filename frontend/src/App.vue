<template>
  <div class="min-h-screen bg-linear-to-b from-slate-50 to-slate-100 py-8">
    <div class="max-w-4xl mx-auto px-4">
      <!-- Header -->
      <header class="mb-8">
        <h1 class="text-4xl font-bold text-slate-900 mb-2">
          {{ api.schema.value?.site?.heading || 'Blogger XML Exporter' }}
        </h1>
        <p class="text-slate-600">
          {{ api.schema.value?.site?.title ? '' : 'Export blog posts with custom XML schema' }}
        </p>
      </header>

      <!-- Loading state -->
      <div v-if="api.loading.value" class="bg-white rounded-lg border border-slate-200 p-6">
        <p class="text-slate-700">
          Lädt Schema...
        </p>
      </div>

      <!-- Schema Error state -->
      <div v-else-if="api.schemaError.value" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
        <p class="text-red-800">
          Fehler beim Laden des Schemas: {{ api.schemaError.value }}
        </p>
      </div>

      <!-- Form -->
      <form v-else-if="api.hasSchema.value" class="bg-white rounded-lg border border-slate-200 p-6 shadow-sm" :style="{
        '--skeleton-base': 'rgb(255, 255, 255)',
        '--skeleton-highlight': hexToRgb(themeColors.primaryColor, 0.08)
      } as Record<string, string>" @submit.prevent="onSubmit">
        <div v-if="!api.postsError.value" class="mb-6 pb-6 border-b border-slate-200">
          <h2 class="text-lg font-semibold text-slate-900 mb-4">
            Blog Post
          </h2>
          <FormCombobox :item="{
            name: 'post',
            label: 'Post wählen',
            type: 'combobox',
            required: schema?.items?.find(i => i.name === 'post')?.required || false,
            options: postsOptions,
            placeholder: 'Post suchen...',
            help: 'Wählen Sie einen Blog-Post aus'
          }" :model-value="selectedPostId" :clear-on-focus="true" @update:model-value="onSelectPost" />
        </div>
        <!-- Posts Error Warning -->
        <div v-else class="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
          <p class="text-sm text-yellow-800">
            ⚠️ Blog-Posts konnten nicht geladen werden ({{ api.postsError.value }}). Sie können das Formular trotzdem
            manuell ausfüllen.
          </p>
        </div>

        <!-- Form items – rendered via RenderGroupContent so row/width grid
             logic, all field types (textarea, array, etc.) and group nesting
             work identically regardless of whether the config uses groups. -->
        <RenderGroupContent v-if="schema && schema.items" :group="{ name: '__root__', items: schema.items }"
          :form-values="formValues" :is-loading="isFillingForm" />

        <!-- Submit button -->
        <div class="mt-8 flex gap-3 border-t border-slate-200 pt-6">
          <button type="submit" :disabled="isSubmitting" :style="{
            backgroundColor: themeColors.primaryColor,
            '--tw-shade-hover': themeColors.darkColor
          }"
            class="px-6 py-3 rounded-lg font-medium transition-all duration-200 active:scale-95 text-white shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
            :class="{ 'hover:opacity-90': !isSubmitting }">
            {{ isSubmitting ? 'Wird generiert...' : 'XML generieren & herunterladen' }}
          </button>
          <button type="button"
            class="px-6 py-3 rounded-lg font-medium transition-all duration-200 active:scale-95 bg-slate-100 text-slate-700 hover:bg-slate-200"
            @click="resetForm">
            Zurücksetzen
          </button>
        </div>
      </form>

      <!-- Fallback: No condition matched -->
      <div v-else class="bg-orange-50 border border-orange-200 rounded-lg p-6">
        <p class="text-orange-800">
          ⚠️ Unerwarteter Zustand: loading={{ api.loading.value }}, hasSchema={{ api.hasSchema.value }}, schemaError={{
            !!api.schemaError.value }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watchEffect } from 'vue'
import type { Post } from '@/types'
import { useApi } from '@/composables/useApi'
import { useForm } from '@/composables/useForm'
import RenderGroupContent from '@/components/Form/core/RenderGroupContent.vue'
import FormCombobox from '@/components/Form/fields/FormCombobox.vue'

const api = useApi()
const selectedPostId = ref('')
const selectedPost = ref<Post | null>(null)
const isSubmitting = ref(false)
const isFillingForm = ref(false)

// Schema and form setup
const schema = computed(() => api.schema.value)
const form = useForm(schema.value)

const formValues = form.formValues

// Theme colors from schema
const themeColors = computed(() => {
  if (!api.schema.value?.theme) {
    return {
      primaryColor: '#2563eb',
      darkColor: '#1e40af',
      lightColor: '#3b82f6',
    }
  }
  return {
    primaryColor: api.schema.value.theme.primaryColor || '#2563eb',
    darkColor: api.schema.value.theme.darkColor || '#1e40af',
    lightColor: api.schema.value.theme.lightColor || '#3b82f6',
  }
})

// Apply theme colors as CSS custom properties on :root so all Tailwind utilities
// (border-primary, bg-primary, etc.) and third-party widget overrides (flatpickr,
// Tom Select) pick up the configured theme color globally.
watchEffect(() => {
  const root = document.documentElement
  root.style.setProperty('--color-primary', themeColors.value.primaryColor)
  root.style.setProperty('--color-primary-dark', themeColors.value.darkColor)
  root.style.setProperty('--color-primary-light', themeColors.value.lightColor)
})

// Post options for combobox — include the post date as right-aligned description
// so posts with similar titles can be distinguished at a glance.
const postsOptions = computed(() => {
  return api.posts.value?.map(post => {
    const date = post.published ? new Date(post.published) : null
    const description = date
      ? date.toLocaleDateString('de-DE', { day: '2-digit', month: '2-digit', year: 'numeric' })
      : undefined
    return {
      value: post.id || '',
      label: post.title || 'Untitled',
      description,
    }
  }) || []
})

// Init on mount
onMounted(async () => {
  await api.fetchSchema()
  await api.fetchPosts()
  if (api.schema.value) {
    form.initializeForm(api.schema.value)
    // Apply field defaults (template fallbacks resolved server-side against an
    // empty post) that arrived bundled with the schema — no extra round-trip.
    const defaults = api.schema.value.defaults
    if (defaults) {
      Object.entries(defaults).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          ; (formValues as Record<string, unknown>)[key] = value
        }
      })
    }
  }
})

// Select a post
async function onSelectPost(postId: string) {
  selectedPostId.value = postId
  if (!postId) {
    form.clearPost()
    isFillingForm.value = false
    return
  }
  isFillingForm.value = true
  try {
    const postData = await api.fetchPost(postId)
    if (postData && postData.post) {
      selectedPost.value = postData.post
      form.setSelectedPost(postData)
    }
  } finally {
    isFillingForm.value = false
  }
}

// Submit form
async function onSubmit() {
  isSubmitting.value = true

  try {
    const result = await api.generateXml({
      postId: selectedPostId.value || undefined,
      values: form.getFormValues(),
    })

    if (result) {
      // Download XML
      const blob = new Blob([result.xml], { type: 'application/xml' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = result.fileName
      document.body.appendChild(a)
      a.click()
      document.body.removeChild(a)
      URL.revokeObjectURL(url)
    }
  } catch (e) {
    console.error('Submit error:', e)
  } finally {
    isSubmitting.value = false
  }
}

// Reset form
function resetForm() {
  form.resetForm()
  selectedPostId.value = ''
  selectedPost.value = null
}

// Convert hex color to rgb or rgba string
// Examples: "#2563eb" → "rgb(37, 99, 235)" or with alpha → "rgba(37, 99, 235, 0.08)"
function hexToRgb(hex: string, alpha?: number): string {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  if (!result) return hex

  const r = parseInt(result[1], 16)
  const g = parseInt(result[2], 16)
  const b = parseInt(result[3], 16)

  if (alpha !== undefined && alpha < 1) {
    return `rgba(${r}, ${g}, ${b}, ${alpha})`
  }
  return `rgb(${r}, ${g}, ${b})`
}
</script>
