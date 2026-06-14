<template>
  <div>
    <CookieConsent class="notification is-warning" href="/static?page=cookie_policy" />
    <NavigationBar :fixedTop="true" />
    <slot />
    <Footer />
  </div>
</template>

<script>
import CookieConsent from '~/components/CookieConsent.vue';
import NavigationBar from '~/components/NavigationBar.vue';
import Footer from '~/components/Footer.vue';

export default {
  components: {
    CookieConsent,
    NavigationBar,
    Footer,
  },
  mounted() {
    // The fixed-top navbar needs this class on <body> (was Nuxt head.bodyAttrs).
    document.body.classList.add('has-navbar-fixed-top');
  },
  unmounted() {
    document.body.classList.remove('has-navbar-fixed-top');
  },
};
</script>
