<script>
import axios from 'axios';

const BASE_URL = import.meta.env.VITE_BASE_URL || "";
const API_URL = `${BASE_URL}`;

export default {
    name: 'LoginFormAxios',
    data(){
        return{
            form: {
                username: '',
                password: ''
            }
        }
    },
    methods:{
        async getToken(){
            await axios.get(API_URL+"/login")
                 .then((res) => {
                     //Perform Success Action
                     console.log(res);
                 })
                 .catch((error) => {
                    error.response.status == 401 ? alert("Invalid Credentials") : alert("Something went wrong");
                 }).finally(() => {
                     alert("Login Successful")
                     this.$router.push('/list');
                 });
        },
        async submitForm(){
            await axios.post(API_URL+"/login", this.form)
                 .then((res) => {
                     //Perform Success Action
                     console.log(res);
                 })
                 .catch((error) => {
                    error.response.status == 401 ? alert("Invalid Credentials") : alert("Something went wrong");
                 }).finally(() => {
                     //Perform action in always
                     alert("Login Successful")
                     this.resetForm();
                     this.$router.push('/list');
                 });
        },
        resetForm () {
            this.$refs.login.reset()
        }
    }
}
</script>

<template>
    <div>
        <h2> Login user</h2>
        <div>
            <button @click="getToken">Get Token</button>
        </div>
        <h2> Login admin </h2>
        <form ref="login" v-on:submit.prevent="submitForm">
            <div class="form-group">
                <label for="username">Name</label>
                <input type="text" class="form-control" id="username" placeholder="Username" v-model="form.username">
            </div>
            <div class="form-group">
                <label for="password">Name</label>
                <input type="password" class="form-control" id="password" v-model="form.password">
            </div>
            <div class="form-group">
                <button class="btn btn-primary">Submit</button>
            </div>
        </form>
    </div>
</template>
