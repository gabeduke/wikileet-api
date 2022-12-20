<script>
const API_URL = `http://localhost:8080/api/v1/items`


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
      const reponse = await fetch(API_URL, {headers: { "X-User": "gabeduke@gmail.com"}});
      const json = await reponse.json();
      this.items = json.data;
      console.log(this.items);
    },
  }
}
</script>

<template>
  <h1>Items</h1>
  <ul>
    <li v-for="item in items">name: {{ item.name }}, desc: {{ item.description }}, url: {{ item.url }}</li>
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
