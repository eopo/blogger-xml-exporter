import eslint from '@eslint/js'
import eslintConfigPrettier from 'eslint-config-prettier'
import vue from 'eslint-plugin-vue'
import globals from 'globals'
import ts from 'typescript-eslint'

export default ts.config(
  {
    ignores: ['dist', 'node_modules', '**/*.config.*', 'coverage'],
  },
  {
    extends: [
      eslint.configs.recommended,
      ...ts.configs.recommended,
      ...vue.configs['flat/recommended'],
    ],
    files: ['**/*.{ts,vue}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: globals.browser,
      parserOptions: {
        parser: ts.parser,
        // Vue 3 removed filters; disabling prevents false positives with TS union types
        vueFeatures: { filter: false },
      },
    },
    rules: {
      // Implicit any reduces type safety; warn to allow gradual tightening
      '@typescript-eslint/no-explicit-any': 'warn',

      // Prop mutations in form components are intentional (shared reactive form state)
      'vue/no-mutating-props': 'error',

      // Unused variables except those prefixed with _
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
        },
      ],
    },
  },
  // Disable all ESLint formatting rules that conflict with Prettier (must be last)
  eslintConfigPrettier,
)
