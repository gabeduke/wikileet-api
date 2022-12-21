<script>
import axios from 'axios';

const API_URL = `/api/v1/items`;

export default {
  data: () => ({
    selectUser:['gabeduke@gmail.com'],
    items: null
  }),

  created() {
    // fetch on init
    this.fetchData()
  },

  watch: {
    // re-fetch whenever selectUser changes
    selectUser: 'fetchData'
  },

  methods: {
    async fetchData() {
      await axios.get(API_URL).then(response => (this.items = response.data))
      console.log(this.items.data);
    },
  }
}
</script>

<template>
  <h1>Items</h1>
  <ul>
    <li v-for="item in items.data">name: {{ item.name }}, desc: {{ item.description }}, url: {{ item.url }}</li>
  </ul>
  <div>
    <button class="btn btn-primary" v-on:click="fetchData()">Refresh</button>
  </div>
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
