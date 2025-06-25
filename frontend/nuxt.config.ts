// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  app: {
    head: {
      title: 'Profundr Inc.',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'A logistics and shipment tracking application.' },
        { name: 'format-detection', content: 'telephone=no' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  },

  css: ['~/assets/css/tailwind.css'],

  vite: {
    plugins: [
      tailwindcss(),
    ],
  },

  modules: [
    ['@nuxtjs/color-mode', { classSuffix: '' }], 
    '@scalar/nuxt',
    'shadcn-nuxt'
  ],
  
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
  
  runtimeConfig: {
    // Private keys (only available on server-side)
    jwtSecret: process.env.JWT_SECRET || 'your_jwt_secret_key_please_change_this',
    databaseUrl: process.env.DATABASE_URL || 'postgresql://user:password@localhost:5432/invoice_financing_db',
    
    // Public keys (exposed to client-side)
    public: {
      apiBase: process.env.BASE_URL || 'http://localhost:3000',
      appVersion: process.env.APP_VERSION || '1.0.0'
    }
  },

  compatibilityDate: '2025-05-15',
  devtools: { enabled: true }
})