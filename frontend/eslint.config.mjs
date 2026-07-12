import js from '@eslint/js'
import ts from 'typescript-eslint'
import vue from 'eslint-plugin-vue'
import vueParser from 'vue-eslint-parser'

export default [
  {
    ignores: ['dist', 'node_modules', '**/*.config.*', 'coverage'],
  },

  // 1. Base JavaScript rules
  js.configs.recommended,

  // 2. TypeScript rules (recommended, not strict - allows flexibility)
  ...ts.configs.recommended,

  // 3. Vue 3 recommended best practices
  ...vue.configs['flat/recommended'],

  // 4. Vue files: combine Vue + TypeScript parsing
  {
    files: ['**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: ts.parser,
        sourceType: 'module',
        ecmaVersion: 'latest',
      },
      globals: {
        // Vue 3 Composition API auto-imports
        defineProps: 'readonly',
        defineEmits: 'readonly',
        defineExpose: 'readonly',
        withDefaults: 'readonly',
        defineModel: 'readonly',
        // Browser APIs
        document: 'readonly',
        window: 'readonly',
        navigator: 'readonly',
        HTMLInputElement: 'readonly',
        HTMLTextAreaElement: 'readonly',
        Blob: 'readonly',
        URL: 'readonly',
        URLSearchParams: 'readonly',
        fetch: 'readonly',
        console: 'readonly',
        setTimeout: 'readonly',
        setInterval: 'readonly',
        clearTimeout: 'readonly',
        clearInterval: 'readonly',
      },
    },
    rules: {
      // Development flexibility - warn on implicit any
      '@typescript-eslint/no-explicit-any': 'warn',
      
      // Vue 3 template variables are tracked in template, not visible to ESLint
      // This causes false positives for reactive refs/computed values used in templates
      'no-useless-assignment': 'off',
      
      // Vue pipes (|) in templates are TypeScript union types in script, not Vue filters
      // This prevents false positives when using type unions in Vue components
      'vue/no-deprecated-filter': 'off',
      
      // Vue 3 allows reactive mutations in setup - template refs handle reactivity
      'vue/no-mutating-props': 'warn',
    },
  },

  // 5. TypeScript files: strict type checking
  {
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      globals: {
        console: 'readonly',
        document: 'readonly',
        window: 'readonly',
        Blob: 'readonly',
        URL: 'readonly',
        setTimeout: 'readonly',
      },
    },
    rules: {
      // Catch unused variables (except those starting with _)
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
        },
      ],
      '@typescript-eslint/no-explicit-any': 'warn',
    },
  },

  // 6. Environment-based rules
  {
    rules: {
      'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
      'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    },
  },
]
