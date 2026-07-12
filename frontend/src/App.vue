<template>
  <div :style="themeStyles" class="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100 py-8">
    <div class="max-w-4xl mx-auto px-4">
      <!-- Header -->
      <header class="mb-8">
        <h1 class="text-4xl font-bold text-slate-900 mb-2">{{ api.schema.value?.site?.heading || 'Blogger XML Exporter' }}</h1>
        <p class="text-slate-600">{{ api.schema.value?.site?.title ? '' : 'Export blog posts with custom XML schema' }}</p>
      </header>

      <!-- Loading state -->
      <div v-if="api.loading.value" class="bg-white rounded-lg border border-slate-200 p-6">
        <p class="text-slate-700">Lädt Schema...</p>
      </div>

      <!-- Schema Error state -->
      <div v-else-if="api.schemaError.value" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
        <p class="text-red-800">Fehler beim Laden des Schemas: {{ api.schemaError.value }}</p>
      </div>

      <!-- Form -->
      <form v-else-if="api.hasSchema.value" @submit.prevent="onSubmit" class="bg-white rounded-lg border border-slate-200 p-6 shadow-sm">
        <div v-if="!api.postsError.value" class="mb-6 pb-6 border-b border-slate-200">
          <h2 class="text-lg font-semibold text-slate-900 mb-4">Blog Post</h2>
          <FormCombobox
            :item="{
              name: 'post',
              label: 'Post wählen',
              type: 'combobox',
              required: schema?.items?.some(i => i.name === 'post')?.required,
              options: postsOptions,
              placeholder: 'Post suchen...',
              help: 'Wählen Sie einen Blog-Post aus'
            }"
            :model-value="selectedPostId"
            @update:model-value="onSelectPost"
          />
        </div>
        <!-- Posts Error Warning -->
        <div v-else class="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
          <p class="text-sm text-yellow-800">⚠️ Blog-Posts konnten nicht geladen werden ({{ api.postsError.value }}). Sie können das Formular trotzdem manuell ausfüllen.</p>
        </div>

        <!-- Form groups -->
        <div v-if="schema && schema.items" class="space-y-0">
          <div v-for="(item, idx) in schema.items" :key="idx">
            <!-- Render groups -->
            <FormGroup
              v-if="item.type === 'group'"
              :group="(item as any)"
              :form-values="formValues"
            />
            <!-- Render other fields directly -->
            <FormField
              v-if="!['group', 'array', 'textarea', 'date', 'select', 'combobox'].includes(item.type)"
              :item="item"
              :model-value="(formValues[item.name] || '') as string | number"
              @update:model-value="formValues[item.name] = $event"
            />
            <FormDate
              v-else-if="item.type === 'date'"
              :item="item"
              :model-value="(formValues[item.name] || '') as string"
              @update:model-value="formValues[item.name] = $event"
            />
            <FormCombobox
              v-else-if="item.type === 'select' || item.type === 'combobox'"
              :item="item"
              :model-value="(formValues[item.name] || '') as string"
              @update:model-value="formValues[item.name] = $event"
            />
          </div>
        </div>

        <!-- Submit button -->
        <div class="mt-8 flex gap-3 border-t border-slate-200 pt-6">
          <button
            type="submit"
            :disabled="isSubmitting"
            :style="{ 
              backgroundColor: themeColors.primaryColor,
              '--tw-shade-hover': themeColors.darkColor
            }"
            class="px-6 py-3 rounded-lg font-medium transition-all duration-200 active:scale-95 text-white shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
            :class="{ 'hover:opacity-90': !isSubmitting }"
          >
            {{ isSubmitting ? 'Wird generiert...' : 'XML generieren & herunterladen' }}
          </button>
          <button
            type="button"
            @click="resetForm"
            class="px-6 py-3 rounded-lg font-medium transition-all duration-200 active:scale-95 bg-slate-100 text-slate-700 hover:bg-slate-200"
          >
            Zurücksetzen
          </button>
        </div>
      </form>

      <!-- Fallback: No condition matched -->
      <div v-else class="bg-orange-50 border border-orange-200 rounded-lg p-6">
        <p class="text-orange-800">⚠️ Unerwarteter Zustand: loading={{ api.loading.value }}, hasSchema={{ api.hasSchema.value }}, schemaError={{ !!api.schemaError.value }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { Post } from '@/types'
import { useApi } from '@/composables/useApi'
import { useForm } from '@/composables/useForm'
import { formatDate } from '@/dateFormatter'
import FormGroup from '@/components/Form/core/FormGroup.vue'
import FormField from '@/components/Form/fields/FormField.vue'
import FormDate from '@/components/Form/fields/FormDate.vue'
import FormCombobox from '@/components/Form/fields/FormCombobox.vue'

const api = useApi()
const selectedPostId = ref('')
const selectedPost = ref<Post | null>(null)
const isSubmitting = ref(false)

// Schema and form setup
const schema = computed(() => api.schema.value)
const form = useForm(schema.value as any)
const formValues = form.formValues

// Theme colors from schema
const themeColors = computed(() => {
  if (!api.schema.value?.theme) {
    return {
      primaryColor: '#0f172a', // slate-900
      darkColor: '#1e293b',    // slate-800
      lightColor: '#f1f5f9'    // slate-100
    }
  }
  return {
    primaryColor: api.schema.value.theme.primaryColor || '#0f172a',
    darkColor: api.schema.value.theme.darkColor || '#1e293b',
    lightColor: api.schema.value.theme.lightColor || '#f1f5f9'
  }
})

// CSS variables for theme colors
const themeStyles = computed(() => ({
  '--color-primary': themeColors.value.primaryColor,
  '--color-dark': themeColors.value.darkColor,
  '--color-light': themeColors.value.lightColor,
}))

const postsOptions = computed(() =>
  api.posts.value.map((post) => ({
    value: post.id,
    label: post.title,
    description: post.published ? formatDate(post.published, true) : undefined,
  }))
)

// Init on mount
onMounted(async () => {
  await api.fetchSchema()
  await api.fetchPosts()
  if (api.schema.value) {
    form.initializeForm(api.schema.value as any)
  }
})

// Select a post
async function onSelectPost(postId: string) {
  selectedPostId.value = postId
  if (!postId) {
    form.clearPost()
    return
  }
  const postData = await api.fetchPost(postId)
  if (postData && postData.post) {
    selectedPost.value = postData.post
    form.setSelectedPost(postData)
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
</script>

<style scoped>
:root {
  --color-primary: #0f172a;
  --color-dark: #1e293b;
  --color-light: #f1f5f9;
}
</style>
