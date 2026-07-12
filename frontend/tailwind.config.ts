import type { Config } from 'tailwindcss'

export default {
  content: [
    './index.html',
    './src/**/*.{vue,ts,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        slate: {
          50: '#f8fafc',
          100: '#f1f5f9',
          150: '#eef3f8',
          200: '#e2e8f0',
          300: '#cbd5e1',
          400: '#94a3b8',
          500: '#64748b',
          600: '#475569',
          700: '#334155',
          900: '#0f172a',
        },
        primary: 'var(--color-primary, #2563eb)',
        'primary-dark': 'var(--color-primary-dark, #1e40af)',
        'primary-light': 'var(--color-primary-light, #3b82f6)',
      },
    },
  },
  plugins: [],
} satisfies Config
