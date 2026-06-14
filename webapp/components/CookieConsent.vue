<template>
  <div v-if="visible">
    <div class="cookie-consent-body">
      <slot>
        This website uses cookies to ensure you get the best experience.
      </slot>
      <a v-if="href" :href="href" class="cookie-consent-link">Learn more</a>
      <button class="button is-small is-dark cookie-consent-accept" @click="accept">
        Got it!
      </button>
    </div>
  </div>
</template>

<script>
// Replaces vue-cookieconsent-component (Vue 2 only). Persists acceptance in
// localStorage so the banner stays dismissed across visits.
const STORAGE_KEY = 'cookie:accepted'

export default {
  name: 'CookieConsent',
  props: {
    href: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      visible: false,
    }
  },
  mounted() {
    this.visible = localStorage.getItem(STORAGE_KEY) !== 'true'
  },
  methods: {
    accept() {
      localStorage.setItem(STORAGE_KEY, 'true')
      this.visible = false
    },
  },
}
</script>

<style scoped>
.cookie-consent-body {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}
.cookie-consent-link {
  text-decoration: underline;
}
.cookie-consent-accept {
  margin-left: auto;
}
</style>
