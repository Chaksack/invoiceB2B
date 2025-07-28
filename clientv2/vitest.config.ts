import { defineConfig } from 'vitest/config'
import tsconfigPaths from 'vite-tsconfig-paths'
import { fileURLToPath } from 'url'
import { resolve } from 'path'

const r = (p: string) => resolve(__dirname, p)

export default defineConfig({
  plugins: [
    tsconfigPaths(),
  ],
  test: {
    environment: 'jsdom',
  },
  resolve: {
    alias: {
      '~': r('.'),
      '@': r('.')
    }
  }
})
