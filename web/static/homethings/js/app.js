"use strict"

import {loginPageComponent} from "./components/login_page.js";
import {mainPageComponent} from "./components/main_page.js";
import * as client from "./client/client.js";

export const app = {
    components: {
        "login-page": loginPageComponent,
        "main-page": mainPageComponent,
    },
    data() {
        return {
            isAuth: false,
            username: "",
        };
    },
    created() {
        let res = client.jsonRequest(client.methodGet, client.routeCheckAuth)
        if (res.status === client.statusOK) {
            this.isAuth = true
            this.username = res.data.username
        }
    },
    methods: {
        setIsAuth(isAuth) {
            this.isAuth = isAuth
        },
        setUsername(username) {
            this.username = username
        }
    }
};
