<script>
import vSelect from "vue-select"
import axios from 'axios';

const BASE_URL = import.meta.env.VITE_BASE_URL || "";
const API_URL = `${BASE_URL}/api/v1`;

export default {
  data: () => ({
    selectUser: null,
    users: [],
    items: []
  }),

  components: {
    vSelect
  },

  created() {
    // fetch on init
    this.fetchItems()
    this.getDefaultUser()
    this.fetchUsers()
  },

  watch: {
    // re-fetch whenever selectUser changes
    selectUser: 'fetchItems'
  },

  methods: {
    async fetchItems() {
      await axios.get(API_URL+"/items",{params: {user_email:this.selectUser}}).then(response => (this.items = response.data.data))
      console.log(this.items);
    },
    async fetchUsers() {
      await axios.get(API_URL+"/users").then(response => (this.users = response.data.data))
      console.log(this.users);
    },
    async getDefaultUser() {
      const myHeaders = new Headers();
      const u = myHeaders.get('X-User');
      this.selectUser = u;
    }
  }
}
</script>

<template>
  <h1>Items</h1>
  <div>
    <v-select
      v-model="selectUser"
      :options="users"
      label="email"
      @input="fetchItems"
      :reduce="user => user.email"
    />
  </div>
  <div>
    <button class="btn btn-primary" v-on:click="fetchItems()">Refresh</button>
  </div>
  <ul>
    <li v-for="item in items">name: {{ item.name }}, user_id: {{ item.user_id  }} desc: {{ item.description }}, url: {{ item.url }}</li>
  </ul>
</template>

<style>
a {
  text-decoration: none;
  color: #42b883;
}
li {
  line-height: 1.5em;
  margin-bottom: 20px;
}
.author,
.date {
  font-weight: bold;
}
</style>
