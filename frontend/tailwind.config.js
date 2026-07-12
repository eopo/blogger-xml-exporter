import forms from '@tailwindcss/forms'

/**
 * Tailwind CSS v4 configuration
 * - Content scanning: Automatic via PostCSS (@import "tailwindcss")
 * - Theme customization: In tailwind.src.css via @theme
 * - Only plugins configured here
 */
export default {
  plugins: [forms],
}
