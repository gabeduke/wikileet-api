<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Welcome from './components/Welcome.vue'

import { useAuthStore } from './stores/auth';

const authStore = useAuthStore();

</script>

<template>
  <header>
   
    <div class="wrapper">
      <!-- Default-->
      <div v-if="!authStore.token">
        <Welcome msg="You have arrived at your gift exchange destination." />
        <nav v-show="authStore.token">
          <RouterLink to="/login">Login</RouterLink>
          <RouterLink to="/about">About</RouterLink>
        </nav>
      </div>
      <!-- Logged In Users -->
      <div v-show="authStore.token" v-if="authStore.token">
        <Welcome msg="Welcome Back!" msg-action="Logout"/> | 
        <nav v-show="authStore.token">
          <h2><RouterLink to="/list">Gift Lists</RouterLink></h2>
          <RouterLink to="/create">(+) item</RouterLink>
          <h2><RouterLink to="/list">MY List</RouterLink></h2>
        </nav>
      </div>
    </div>
  </header>
  <div class="body">
    <RouterView/>
  </div>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

nav  {
  font-size: 12px;
  margin: 2rem 1rem 0rem 1rem;
  width: 80%;
  padding:1rem;
  text-align: left;
  background-color: rgba(240, 248, 255, 0.19);
}

nav a.router-link-exact-active {
  color: var(--color-text);
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: block;
  padding: 0.5rem 1rem;
  border-left: 1px solid var(--color-border);
}

nav h2 a {
  font-weight: bold;
  font-size: larger;
}

nav h3 a {
  font-weight: normal;
  font-size: large;
  text-indent: 1rem;
}

nav a:first-of-type {
  border: 0;

}
nav a .logout {
  font-size:80%;
  color: rgba(240, 248, 255, 0.19);
  text-align:right;
  display:inline-flexbox;
  float:right;
}
@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flexbox;
    place-items: flex-start;
    flex-wrap: wrap;
  }

  nav {
    text-align: left;
    margin-left: -1rem;
    font-size: 1rem;
    display: flexbox;
    padding: 1rem 0;
    margin-top: 1rem;
  }
}
</style>