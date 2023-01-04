import { defineStore } from 'pinia';

import { fetchWrapper } from '../_helpers/fetch_wrappers';

const BASE_URL = import.meta.env.VITE_BASE_URL || "";
const API_URL = `${BASE_URL}/api/v1`;

export const useUsersStore = defineStore({
    id: 'users',
    state: () => ({
        users: {}
    }),
    actions: {
        async getAll() {
            this.users = { loading: true };
            await fetchWrapper.get(`${API_URL}/users`)
                .then(users => this.users = users.data)
                .catch(error => this.users = { error })
        }
    }
});