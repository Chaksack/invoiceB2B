// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  css: ['~/assets/css/tailwind.css'],

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  modules: [    '@nuxtjs/color-mode', '@scalar/nuxt',
    'shadcn-nuxt'],
  colorMode: {
    classSuffix: ''
  },
  nitro: {
    experimental: {
      openAPI: true,
    },
  },
  shadcn: {
    /**
     * Prefix for all the imported component
     */
    prefix: '',
    /**
     * Directory that the component lives in.
     * @default "./components/ui"
     */
    componentDir: './components/ui'
  },

  compatibilityDate: '2025-05-15',
  devtools: { enabled: true },
  modules: ['shadcn-nuxt']
})