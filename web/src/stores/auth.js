import {defineStore} from 'pinia'

export const useAuthStore = defineStore({
    id: 'auth',
    state: () => ({
        isAuth: false,
        username: ""
    }),
    actions: {
        setAuth(username) {
            this.isAuth = true
            this.username = username
        },
        clearAuth() {
            this.isAuth = false
            this.username = ""
        }
    },
})