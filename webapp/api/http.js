import axios from 'axios'
import { useAuthStore } from '~/stores/auth'

// Shared axios instance. In dev the Vite proxy forwards `/api` to the Go
// backend; in prod the static bundle is served by the same Go server, so a
// same-origin base URL works. Override with VITE_API_BASE if needed.
const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || '',
})

// Mirrors the old plugins/axios.js: on 401, try to refresh the token once.
http.interceptors.response.use(
  (response) => response,
  async (error) => {
    const auth = useAuthStore()
    const status = error.response?.status
    if (status === 401 && auth.loggedIn && auth.refreshToken) {
      try {
        await auth.refreshTokens()
      } catch {
        await auth.logout()
      }
    }
    return Promise.reject(error)
  }
)

export default http
