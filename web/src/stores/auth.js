import { defineStore } from 'pinia';
import VueJwtDecode from 'vue-jwt-decode'
import { fetchWrapper } from '../_helpers/fetch_wrappers';
import router from '../router';

const BASE_URL = import.meta.env.VITE_BASE_URL || "";
const API_URL = `${BASE_URL}`;

export const useAuthStore = defineStore({
    id: 'auth',
    state: () => ({
        // initialize state from local storage to enable user to stay logged in
        token: JSON.parse(localStorage.getItem('token')),
        user: JSON.parse(localStorage.getItem('user')),
        returnUrl: null
    }),
    actions: {
        async login(username, password) {
            const resp = await fetchWrapper.post(`${API_URL}/login`, { username, password });
            this.saveToken(resp)
            await router.push(this.returnUrl || '/');
        },
        saveToken(response) {
            // update pinia state
            this.token = response.token;

            // store user details and jwt in local storage to keep user logged in between page refreshes
            localStorage.setItem('token', JSON.stringify(response.token));

            this.decodeToken()
        },
        logout() {
            this.token = null;
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            router.push('/login');
        },
        decodeToken(){
            try {
                const decoded = VueJwtDecode.decode(this.token);
                this.user = decoded.id;
                localStorage.setItem('user', JSON.stringify(this.user));
            } catch (err) {
                console.log('token invalid: ',err);
                this.logout();
            }
        }
    }
});