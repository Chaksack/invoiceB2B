// https://nuxt.com/docs/api/configuration/nuxt-config
import tsconfigPaths from 'vite-tsconfig-paths'
import { resolve } from 'path'

const r = (p: string) => resolve(__dirname, p)

export default defineNuxtConfig({
  css: ['~/assets/css/tailwind.css', 'vue-sonner/style.css'],

  vite: {
    plugins: [
      tsconfigPaths(),
    ],
    resolve: {
      alias: {
        '~': r('.'),
        '@': r('.')
      }
    }
  },

  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss', 'shadcn-nuxt'],
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
  }})
