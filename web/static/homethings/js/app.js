import {loginPageComponent} from './components/login_page.js';
import {mainPageComponent} from './components/main_page.js';
import {request} from "./client.js";

export const app = {
    components: {
        'login-page': loginPageComponent,
        'main-page': mainPageComponent,
    },
    data() {
        return {
            isAuth: false,
        };
    },
    created() {
        let res = request("GET", "/api/v1/auth/check")
        if (res.status === 200) {
            this.isAuth = true
        }
    },
    methods: {
        setAuth(isAuth) {
            this.isAuth = isAuth
        }
    }
};
