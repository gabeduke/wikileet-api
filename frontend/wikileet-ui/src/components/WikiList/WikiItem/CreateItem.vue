<script>
import axios from 'axios';

const API_URL = `/api/v1/items`;

export default {
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
    methods:{
        submitForm(){
            axios.post(API_URL, this.form)
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
        }
    }
}
</script>

<template>
    <div>

            <h2> Create Item </h2>
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
