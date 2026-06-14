import { useAuthStore } from '~/stores/auth'

// Replaces the old middleware/*.js. Pages declare `meta.middleware` via a
// <route> block; this global guard enforces it.
export function registerGuards(router) {
  router.beforeEach((to) => {
    const auth = useAuthStore()
    const middleware = to.meta.middleware

    if (middleware === 'auth') {
      if (!auth.loggedIn) {
        return '/auth/login'
      }
    }

    if (middleware === 'auth-admin') {
      if (!auth.loggedIn || !auth.user) {
        return '/auth/login'
      }
      if (auth.user.role !== 'ADMIN') {
        return '/admin_required'
      }
    }

    // `guest` is a no-op, matching the original empty middleware/guest.js.
    return true
  })
}
