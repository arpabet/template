# template frontend

A Vue 3 + Vite single-page app. The production build is emitted to `../assets`, which the Go
server embeds via `go-bindata` (see the repo `Makefile`).

## Build Setup

```bash
# install dependencies
$ npm install

# serve with hot reload at http://localhost:3000 (proxies /api -> http://localhost:8443)
$ npm run dev

# build the static bundle into ../assets
$ npm run build

# preview the production build locally
$ npm run preview
```

From the repo root, `make generate` runs `npm run build` and then post-processes the bundle
(`gtag.py`); `make build` embeds `../assets` into the Go binary.

## Project structure

- `pages/` — file-based routes (via `vite-plugin-pages`). A page can declare its layout and an
  auth guard through a `<route>` block, e.g.

  ```
  <route>
  { meta: { layout: "admin", middleware: "auth-admin" } }
  </route>
  ```

- `layouts/` — layouts applied by `vite-plugin-vue-layouts-next` based on `meta.layout`
  (`default` when unset). Layouts render the page through a `<slot />`.
- `components/` — reusable Vue components.
- `stores/` — Pinia stores. `stores/auth.js` implements login / refresh / logout against the
  `/api/auth/*` endpoints and exposes `isAuthenticated` / `loggedInUser`.
- `api/http.js` — shared axios instance (adds the bearer token, retries once on 401).
- `router/guards.js` — global navigation guard enforcing `meta.middleware`
  (`auth`, `auth-admin`, `guest`).
- `public/` — static files served at the web root (`/logo.png`, `/css/bulma.min.css`, …).

`$axios` and `$auth` are registered as global properties, so Options-API components can use
`this.$axios` and `this.$auth` directly.
