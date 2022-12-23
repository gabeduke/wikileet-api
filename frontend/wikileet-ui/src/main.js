import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
// import VueSession from 'vue-session'

import vSelect from 'vue-select'
import 'vue-select/dist/vue-select.css';

import './assets/main.css'

const app = createApp(App)

app.component('v-select', vSelect)

app.use(router)

app.mount('#app')
