import { reactive, computed } from 'vue'
import type { FormGroup, FormItem, FormValues, Preset, Post, PostResponse } from '@/types'

export function useForm(schema: FormGroup | null) {
  const formValues = reactive<FormValues>({})
  const selectedPost = reactive<Partial<Post>>({})

  // Initialize form values from schema
  function initializeForm(schemaArg?: FormGroup | null) {
    const schemaToUse = schemaArg || schema
    if (!schemaToUse) return
    
    const items = (schemaToUse as any).items || (schemaToUse as any).fields || []
    forEachField(items, (field) => {
      if (field.type === 'group' || field.type === 'array') {
        if (field.type === 'array') {
          formValues[field.name] = []
        } else {
          formValues[field.name] = {}
        }
      } else {
        formValues[field.name] = ''
      }
    })
  }

  // Recursively traverse all fields
  function forEachField(fields: FormItem[], callback: (field: FormItem) => void) {
    fields.forEach((field) => {
      callback(field)
      if ((field.type === 'group' || field.type === 'array') && field.items) {
        forEachField(field.items, callback)
      }
    })
  }

  // Apply preset values
  function applyPreset(preset: Preset) {
    Object.assign(formValues, preset.values)
  }

  // Reset all values
  function resetForm() {
    Object.keys(formValues).forEach((key) => {
      if (Array.isArray(formValues[key])) {
        formValues[key] = []
      } else {
        formValues[key] = ''
      }
    })
  }

  // Set selected post and auto-fill fields
  function setSelectedPost(postData: PostResponse) {
    const post = postData.post
    selectedPost.id = post.id
    selectedPost.title = post.title
    selectedPost.content = post.content
    selectedPost.published = post.published
    selectedPost.updated = post.updated

    // Apply any resolved default values delivered from backend
    if (postData.values && typeof postData.values === 'object') {
      Object.keys(postData.values).forEach(key => {
        const value = postData.values[key]
        // Apply value if it's defined, not null, and not empty string
        if (value !== undefined && value !== null && value !== '') {
          formValues[key] = value as any
        }
      })
    }
  }

  // Clear selected post
  function clearPost() {
    Object.keys(selectedPost).forEach((key) => {
      delete selectedPost[key as keyof typeof selectedPost]
    })
  }

  // Get flattened form values for submission
  function getFormValues(): FormValues {
    return JSON.parse(JSON.stringify(formValues))
  }

  const isFormValid = computed(() => {
    // Basic validation: check required fields
    let valid = true
    const items = (schema as any)?.items || (schema as any)?.fields || []
    forEachField(items, (field) => {
      if (field.required && !formValues[field.name]) {
        valid = false
      }
    })
    return valid
  })

  return {
    formValues,
    selectedPost,
    initializeForm,
    applyPreset,
    resetForm,
    setSelectedPost,
    clearPost,
    getFormValues,
    isFormValid,
  }
}
