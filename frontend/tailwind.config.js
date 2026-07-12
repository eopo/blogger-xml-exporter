import forms from '@tailwindcss/forms'

export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx,vue}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#1e293b', // slate-900
      },
    },
  },
  plugins: [forms],
  corePlugins: {
    preflight: true,
  },
}
