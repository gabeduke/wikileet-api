import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import Vue3Storage from "vue3-storage";

import vSelect from 'vue-select'
import 'vue-select/dist/vue-select.css';

import './assets/main.css'

const app = createApp(App)
    .use(Vue3Storage)
    .use(router)
    .component('v-select', vSelect)
    .mount('#app')
