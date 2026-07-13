/**
 * Form schema types for blog post XML export
 */

export interface SchemaResponse {
  items: FormItem[] // Top-level form items (can include groups)
  site: {
    title: string
    heading: string
  }
  theme: {
    primaryColor: string
    darkColor: string
    lightColor: string
  }
  assets?: {
    logo?: string
    favicon?: string
  }
  defaults?: FormValues // Field values resolved against an empty post (template fallbacks)
}

export interface FormGroup {
  name: string
  title?: string
  collapsible?: boolean
  collapsed?: boolean
  items: FormItem[]
  presets?: Preset[]
}

export interface FormItem {
  name: string
  title?: string
  label?: string
  type: 'text' | 'textarea' | 'number' | 'email' | 'date' | 'select' | 'combobox' | 'group' | 'array'
  row?: number // Row grouping for layout (fields with the same row number appear inline)
  width?: number // Proportional width (fr units); e.g. width 2 + width 4 → second is twice as wide
  required?: boolean
  includeTime?: boolean // for date type
  placeholder?: string
  help?: string
  minLength?: number
  maxLength?: number
  pattern?: string
  options?: SelectOption[] // for select/combobox
  items?: FormItem[] // for group
  fields?: FormItem[] // for array (config format compatibility)
  collapsible?: boolean
  collapsed?: boolean
  presets?: Preset[]
}

export interface SelectOption {
  value: string
  label: string
}

export interface Post {
  id: string
  title: string
  content: string
  published: string
  updated: string
  url?: string
}

export interface FormValues {
  [key: string]: FormValue | FormValues[] | FormValues
}

export type FormValue = string | number | boolean | null | undefined

export interface Preset {
  title?: string
  label?: string
  values: FormValues
}

export interface ApiError {
  message: string
  code?: string
}

export interface GenerateXmlRequest {
  postId?: string
  values: FormValues
}

export interface GenerateXmlResponse {
  xml: string
  fileName: string
}

export interface PostResponse {
  post: Post
  values: FormValues
  presets: Preset[]
}
