import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import Pages from 'vite-plugin-pages'
import Layouts from 'vite-plugin-vue-layouts-next'

const root = fileURLToPath(new URL('.', import.meta.url))

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    Pages({
      dirs: 'pages',
    }),
    Layouts({
      defaultLayout: 'default',
    }),
  ],

  resolve: {
    alias: {
      '~~': root,
      '@@': root,
      '~': root,
      '@': root,
    },
    // Allow extensionless `.vue` imports (e.g. '~/components/Notification'),
    // as the original Nuxt code relied on.
    extensions: ['.mjs', '.js', '.mts', '.ts', '.jsx', '.tsx', '.json', '.vue'],
  },

  server: {
    port: 3000,
    proxy: {
      // Replaces @nuxtjs/proxy: forward API calls to the Go backend.
      '/api': {
        target: 'http://localhost:8443',
        changeOrigin: true,
        secure: false,
      },
    },
  },

  build: {
    // Emit the static bundle where go-bindata expects it (see Makefile `generate`).
    outDir: '../assets',
    emptyOutDir: true,
  },
})
