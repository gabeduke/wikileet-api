<script>
import axios from 'axios';

const BASE_URL = import.meta.env.VITE_BASE_URL || "";
const API_URL = `${BASE_URL}/api/v1`;

export default {
    selectUser: null,
    users: [],
    name: 'PostFormAxios',
    data(){
        return{
            form: {
                name: '',
                description: '',
                url: ''
            }
        }
    },
    created() {
        // fetch on init
        this.fetchUsers()
    },
    methods:{
        submitForm(){
            axios.post(API_URL+"/items", this.form, {params: {user_email:this.selectUser}})
                 .then((res) => {
                     //Perform Success Action
                     console.log(res);
                 })
                 .catch((error) => {
                     // error.response.status Check status code
                 }).finally(() => {
                     //Perform action in always
                     this.resetForm();
                 });
        },
        resetForm () {
            this.$refs.create.reset()
        },
        async fetchUsers() {
            await axios.get(API_URL+"/users").then(response => (this.users = response.data.data))
            console.log(this.users);
        }
    }
}
</script>

<template>
    <div>

            <h2> Create Item </h2>
            <div>
                <v-select
                v-model="selectUser"
                :options="users"
                label="email"
                :reduce="user => user.email"
                />
            </div>
            <form ref="create" v-on:submit.prevent="submitForm">
                <div class="form-group">
                    <label for="name">Name</label>
                    <input type="text" class="form-control" id="name" placeholder="Your name" v-model="form.name">
                </div>
                <div class="form-group">
                    <label for="description">Name</label>
                    <input type="text" class="form-control" id="description" placeholder="Your description" v-model="form.description">
                </div>
                <div class="form-group">
                    <label for="url">Name</label>
                    <input type="text" class="form-control" id="url" placeholder="URL" v-model="form.url">
                </div>
                <div class="form-group">
                    <button class="btn btn-primary">Submit</button>
                </div>
            </form>
    </div>
</template>
