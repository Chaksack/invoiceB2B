// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  css: ['~/assets/css/tailwind.css'],

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  modules: [
    '@nuxtjs/color-mode',
    'shadcn-nuxt'
  ],

  colorMode: {
    classSuffix: ''
  },

  shadcn: {
    prefix: '',
    componentDir: './components/ui'
  },

  compatibilityDate: '2025-05-15',
  devtools: { enabled: true },
})