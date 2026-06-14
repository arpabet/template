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
          <router-link class="navbar-item" to="/">
            Home
          </router-link>
        </div>


        <div class="navbar-end">
          <div v-if="isAuthenticated" class="navbar-item has-dropdown is-hoverable">
            <a class="navbar-link">
              <small>{{ loggedInUser.first_name }}</small>
            </a>
            <div id="navbarUser" class="navbar-dropdown is-right">
              <router-link class="navbar-item" to="/profile">My Profile</router-link>
              <router-link class="navbar-item" to="/auth/security_log">Security</router-link>
              <router-link v-if="loggedInUser.role == 'ADMIN'" class="navbar-item" to="/admin">Admin Dashboard</router-link>
              <hr class="navbar-divider">
              <a class="navbar-item" @click="logout">Logout</a>
            </div>
          </div>
          <div v-else class="navbar-item">
            <div class="buttons">
              <router-link class="button is-danger" to="/auth/register">
                <strong>Sign Up</strong>
              </router-link>
              <router-link class="button is-light" to="/auth/login">
                <span>Sign in</span><font-awesome-icon icon="fa-solid fa-circle-right" class="ml-2" />
              </router-link>
            </div>
          </div>
        </div>

      </div>
    </div>
  </nav>
</template>

<script>
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
    isAuthenticated() {
      return this.$auth.isAuthenticated;
    },
    loggedInUser() {
      return this.$auth.loggedInUser;
    },
  },

  methods: {
    async logout() {
      await this.$auth.logout();
      this.$router.push('/');
    },
    toggleBurger() {
      this.activeBurger = !this.activeBurger
    },
  },
};
</script>
