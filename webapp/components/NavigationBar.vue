<template>
  <nav class="navbar is-light" :class="{'is-fixed-top' : fixedTop}" >
    <div class="container">
      <div class="navbar-brand">
        <a class="navbar-item" href="/">
          <img src="/logo.png" alt="logo">
        </a>

        <a role="button" class="navbar-burger" :class="{'is-active' : activeBurger}" aria-label="menu" aria-expanded="false" data-target="navbarBasic" @click="toggleBurger">
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div id="navbarBasic" class="navbar-menu" :class="{'is-active' : activeBurger}">

        <div class="navbar-start">
          <nuxt-link class="navbar-item" to="/">
            Home
          </nuxt-link>
        </div>


        <div class="navbar-end">
          <div v-if="isAuthenticated" class="navbar-item has-dropdown is-hoverable">
            <a class="navbar-link">
              <small>{{ loggedInUser.first_name }}</small>
            </a>
            <div id="navbarUser" class="navbar-dropdown is-right">
              <nuxt-link class="navbar-item" to="/profile/">My Profile</nuxt-link>
              <nuxt-link class="navbar-item" to="/auth/security_log">Security</nuxt-link>
              <nuxt-link v-if="loggedInUser.role == 'ADMIN'" class="navbar-item" to="/admin/">Admin Dashboard</nuxt-link>
              <hr class="navbar-divider">
              <a class="navbar-item" @click="logout">Logout</a>
            </div>
          </div>
          <div v-else class="navbar-item">
            <div class="buttons">
              <nuxt-link class="button is-danger" to="/auth/register">
                <strong>Sign Up</strong>
              </nuxt-link>
              <nuxt-link class="button is-light" to="/auth/login">
                <span>Sign in</span><font-awesome-icon icon="fa-solid fa-circle-right" style="margin-left: 7px;" />
              </nuxt-link>
            </div>
          </div>
        </div>

      </div>
    </div>
  </nav>
</template>

<script>
import { mapGetters } from 'vuex';

export default {

  props: {
    fixedTop: {
      type: Boolean,
      default: false
    },
  },

  data() {
    return {
      activeBurger: false,
    };
  },

  computed: {
    ...mapGetters(['isAuthenticated', 'loggedInUser']),
  },

  methods: {
    async logout() {
      await this.$auth.logout();
    },
    toggleBurger() {
      this.activeBurger = !this.activeBurger
    },
  },
};
</script>
