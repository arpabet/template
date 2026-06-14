import { defineStore } from 'pinia'
import Cookies from 'js-cookie'
import http from '~/api/http'

const TOKEN_KEY = 'auth.token'
const REFRESH_KEY = 'auth.refresh_token'
const COOKIE_OPTS = { sameSite: 'lax', secure: window.location.protocol === 'https:' }

function setAuthHeader(token) {
  if (token) {
    http.defaults.headers.common.Authorization = `Bearer ${token}`
  } else {
    delete http.defaults.headers.common.Authorization
  }
}

// Reimplements the @nuxtjs/auth-next "refresh" scheme against the same
// /api/auth/* endpoints the backend already exposes.
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: null,
    refreshToken: null,
    user: null,
    loggedIn: false,
  }),

  getters: {
    // Same getter names the old Vuex store/index.js exposed.
    isAuthenticated: (state) => state.loggedIn,
    loggedInUser: (state) => state.user,
  },

  actions: {
    setTokens(token, refreshToken) {
      this.token = token || null
      if (refreshToken !== undefined) {
        this.refreshToken = refreshToken || null
      }
      setAuthHeader(this.token)
      if (this.token) {
        Cookies.set(TOKEN_KEY, this.token, COOKIE_OPTS)
      } else {
        Cookies.remove(TOKEN_KEY)
      }
      if (this.refreshToken) {
        Cookies.set(REFRESH_KEY, this.refreshToken, COOKIE_OPTS)
      } else {
        Cookies.remove(REFRESH_KEY)
      }
    },

    async loginWith(_strategy, { data } = {}) {
      const res = await http.post('/api/auth/login', data)
      this.setTokens(res.data.token, res.data.refresh_token)
      if (res.data.user) {
        this.user = res.data.user
        this.loggedIn = true
      } else {
        await this.fetchUser()
      }
      return res
    },

    async fetchUser() {
      const res = await http.get('/api/auth/user')
      this.user = res.data.user
      this.loggedIn = !!this.user
      return res
    },

    async refreshTokens() {
      if (!this.refreshToken) {
        throw new Error('no refresh token')
      }
      const res = await http.post('/api/auth/refresh', {
        refresh_token: this.refreshToken,
      })
      this.setTokens(res.data.token, res.data.refresh_token)
      return res
    },

    async logout() {
      try {
        await http.post('/api/auth/logout')
      } catch {
        // ignore network/auth errors during logout
      }
      this.user = null
      this.loggedIn = false
      this.setTokens(null, null)
    },

    // Restore a session from cookies on app startup.
    async init() {
      const token = Cookies.get(TOKEN_KEY)
      const refreshToken = Cookies.get(REFRESH_KEY)
      if (!token) {
        return
      }
      this.token = token
      this.refreshToken = refreshToken || null
      setAuthHeader(token)
      try {
        await this.fetchUser()
      } catch {
        // token expired/invalid — drop it silently
        this.setTokens(null, null)
        this.user = null
        this.loggedIn = false
      }
    },
  },
})
