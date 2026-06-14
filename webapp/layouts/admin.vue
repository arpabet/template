<template>
  <div>
    <NavigationBar/>

    <section class="section">
      <div class="container">
        <div class="columns">
           <div class="column is-one-quarter">
              <aside class="menu">
                <p class="menu-label">
                  General
                </p>
                <ul class="menu-list">
                  <li><router-link to="/admin">Dashboard</router-link></li>
                  <li><router-link to="/admin/users">Users</router-link></li>
                </ul>
                <p class="menu-label">
                  Administration
                </p>
                <ul class="menu-list">
                  <li><router-link to="/admin/pages">Pages</router-link></li>
                </ul>
                <p class="menu-label">
                  Statistics
                </p>
                <ul class="menu-list">
                  <li><router-link to="/admin/traffic">Traffic</router-link></li>
                </ul>
              </aside>
           </div>
           <div class="column">
              <slot/>
           </div>
        </div>
      </div>
    </section>

    <Footer/>
  </div>
</template>

<script>
  import NavigationBar from '~/components/NavigationBar.vue';
  import Footer from '~/components/Footer.vue';

  export default {
    components: {
      NavigationBar,
      Footer,
    },
  };
</script>
