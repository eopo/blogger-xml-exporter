import { ref, computed } from 'vue'
import type { SchemaResponse, Post, GenerateXmlRequest, GenerateXmlResponse, PostResponse } from '@/types'

export function useApi() {
  const schema = ref<SchemaResponse | null>(null)
  const posts = ref<Post[]>([])
  const loading = ref(false)
  const schemaError = ref<string | null>(null)
  const postsError = ref<string | null>(null)

  const hasSchema = computed(() => !!schema.value)

  async function fetchSchema() {
    loading.value = true
    schemaError.value = null
    try {
      const response = await fetch('/api/schema')
      if (!response.ok) throw new Error(`HTTP ${response.status}`)
      schema.value = await response.json()
    } catch (e) {
      schemaError.value = e instanceof Error ? e.message : 'Unknown error'
    } finally {
      loading.value = false
    }
  }

  async function fetchPosts() {
    postsError.value = null
    try {
      const response = await fetch('/api/posts')
      if (!response.ok) throw new Error(`HTTP ${response.status}`)
      posts.value = await response.json()
    } catch (e) {
      postsError.value = e instanceof Error ? e.message : 'Unknown error'
      // Posts error doesn't prevent form usage - set empty posts array
      posts.value = []
    }
  }

  async function fetchPost(postId: string) {
    try {
      const response = await fetch(`/api/posts/${postId}`)
      if (!response.ok) throw new Error(`HTTP ${response.status}`)
      return await response.json() as PostResponse
    } catch (e) {
      postsError.value = e instanceof Error ? e.message : 'Unknown error'
      return null
    }
  }

  async function generateXml(request: GenerateXmlRequest): Promise<GenerateXmlResponse | null> {
    loading.value = true
    try {
      const response = await fetch('/api/generate-xml', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request),
      })
      if (!response.ok) throw new Error(`HTTP ${response.status}`)
      
      const xmlBlob = await response.blob()
      const contentDisposition = response.headers.get('Content-Disposition')
      let fileName = 'export.xml'
      if (contentDisposition && contentDisposition.includes('filename=')) {
        const match = contentDisposition.match(/filename="?([^"]+)"?/)
        if (match && match[1]) {
          fileName = match[1]
        }
      }
      
      return {
        xml: await xmlBlob.text(),
        fileName
      }
    } catch (e) {
      console.error('Failed to generate XML:', e)
      return null
    } finally {
      loading.value = false
    }
  }

  return {
    schema,
    posts,
    loading,
    schemaError,
    postsError,
    hasSchema,
    fetchSchema,
    fetchPosts,
    fetchPost,
    generateXml,
  }
}
