import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import { setupLayouts } from 'virtual:generated-layouts'
import generatedRoutes from 'virtual:generated-pages'

import { library, config } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { fas } from '@fortawesome/free-solid-svg-icons'
import '@fortawesome/fontawesome-svg-core/styles.css'

import App from '~/App.vue'
import http from '~/api/http'
import { useAuthStore } from '~/stores/auth'
import { registerGuards } from '~/router/guards'

// Let the app control FontAwesome CSS injection.
config.autoAddCss = false
library.add(fas)

const app = createApp(App)

const pinia = createPinia()
app.use(pinia)

const router = createRouter({
  // Hash history so deep links / refresh work when the bundle is served as a
  // single index.html by the embedded Go server (no SPA fallback needed).
  // Switch to createWebHistory() if the server is configured to fall back to
  // index.html for unknown paths.
  history: createWebHashHistory(),
  routes: setupLayouts(generatedRoutes),
})
registerGuards(router)
app.use(router)

app.component('font-awesome-icon', FontAwesomeIcon)

// Keep Options-API components working unchanged: this.$axios / this.$auth.
const auth = useAuthStore()
app.config.globalProperties.$axios = http
app.config.globalProperties.$auth = auth

// Restore any existing session before the first render so guards see it.
await auth.init()

await router.isReady()
app.mount('#app')
