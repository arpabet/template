<template>
  <section class="section">
    <div class="container">
      <div class="columns">
        <div class="column is-4 is-offset-4">
          <h2 class="title has-text-centered">Registration</h2>

          <Notification v-if="error" :message="error" @close="error=null"/>

          <form method="post" autocomplete="off" class="box" @submit.prevent="register">

            <div class="field">
              <label class="label required">Username</label>

              <div class="control">
                <input
                  v-model="username"
                  type="text"
                  class="input"
                  name="username"
                  required
                  @change="checkUsername"
                />
              </div>
              <p v-if="!check.available" class="error-message red">Username '{{check.name}}' is not available</p>
            </div>

            <div class="field">
              <label class="label required">First Name</label>

              <div class="control">
                <input
                  v-model="firstName"
                  type="text"
                  class="input"
                  name="first_name"
                  required
                />
              </div>
            </div>

            <div class="field">
              <label class="label">Middle Name</label>

              <div class="control">
                <input
                  v-model="middleName"
                  type="text"
                  class="input"
                  name="middle_name"
                />
              </div>
            </div>

            <div class="field">
              <label class="label required">Last Name</label>

              <div class="control">
                <input
                  v-model="lastName"
                  type="text"
                  class="input"
                  name="last_name"
                  required
                />
              </div>
            </div>

            <div class="field">
              <label class="label required">Email</label>

              <div class="control">
                <input
                  v-model="email"
                  type="email"
                  class="input"
                  name="email"
                  required
                />
              </div>
            </div>

            <div class="field">
              <label class="label required">Password</label>

              <div class="control">
                <VuePassword
                    v-model="password"
                    :strength="strength"
                    type="password"
                    required
                    @input="updateStrength"
                />
              </div>
            </div>

            <br/>

            <div class="container">
              <label class="checkbox">
                <input v-model="agree" type="checkbox" required />
                I agree to the <nuxt-link :to="{ path: '/static', query: { page: 'terms_of_use' }}" target="_blank">Terms Of Use</nuxt-link>
              </label>
            </div>

            <br/>

            <div class="control">
              <button :disabled="strength === 0 || !check.available" type="submit" class="button is-dark is-fullwidth">Register</button>
            </div>
          </form>

          <div class="has-text-centered" style="margin-top: 20px">
            Already got an account? <nuxt-link to="/auth/login">Login</nuxt-link>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
  import VuePassword from 'vue-password';
  import Notification from '~/components/Notification';

  export default {

    components: {
      VuePassword,
      Notification,
    },

    middleware: 'guest',

    data() {
      return {
        username: '',
        check: { name: '', available: true },
        firstName: '',
        middleName: '',
        lastName: '',
        email: '',
        password: '',
        strength: 0,
        agree: false,
        error: null,
      };
    },

    methods: {
      async checkUsername() {
        const name = this.username
        this.check.name = name
        try {
          const res = await this.$axios.put('/api/auth/username', { name });
          if (name === res.data.name) {
            this.check.name = res.data.norm_name
            this.username = res.data.norm_name
            this.check.available = res.data.available
          }
        } catch (e) {
          this.error = e.response.data.message;
          this.check.available = false
        }
      },
      async register() {
        try {
          await this.$axios.post('/api/auth/register', {
            username: this.username,
            first_name: this.firstName,
            middle_name: this.middleName,
            last_name: this.lastName,
            email: this.email,
            password: this.password,
          });

          await this.$auth.loginWith('local', {
            data: {
              login: this.username,
              password: this.password,
            },
          });

          this.$router.push('/');
        } catch (e) {
          this.error = e.response.data.message;
        }
      },
      unrepeated(str) {
         return [...new Set(str)].join('');
      },
      updateStrength(password) {
        const p = this.unrepeated(password);
        let score = Math.floor(p.length / 3);
        // The password has at least one uppercase letter
        if (p.match(/[A-Z]/g)) {
          score = score + 1;
        }
        // The password has at least one lowercase letter
        if (p.match(/[a-z]/g)) {
          score = score + 1;
        }
        // The password has at least one digit
        if (p.match(/[0-9]/g)) {
          score = score + 1;
        }
        // The password has at least one special character
        if (p.match(/[^A-Za-z0-9]/g)) {
          score = score + 1;
        }
        score = Math.floor(score / 2);
        if (score > 4) {
           score = 4;
        }
        this.strength = score;
      },
    },
  };
</script>

<style>
  .required:after {
    content:" *";
    color: red;
  }
</style>
