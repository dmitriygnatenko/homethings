import * as client from '../client/client.js';
import * as auth from '../auth/auth.js'

export const loginPageComponent = {
    props: {
        isAuth: Boolean,
    },
    data() {
        return {
            form: {
                username: "",
                password: "",
            },
            errors: {
                username: false,
                password: false,
            },
        };
    },
    computed: {
        showLoginPage() {
            return !this.isAuth
        }
    },
    methods: {
        submitForm() {
            this.errors.username = this.form.username === "";
            this.errors.password = this.form.password === "";

            if (this.errors.username || this.errors.password) {
                return
            }

            auth.setToken(this.form.username, this.form.password)

            let res = client.request(client.methodGet, client.routeCheckAuth)
            if (res.status === client.statusOK) {
                this.errors.username = false
                this.errors.password = false
                this.form.username = ""
                this.form.password = ""
                this.$emit('eventsetauth', true)
                return
            }

            auth.clearToken()
            this.errors.username = true
            this.errors.password = true
        },
    },
    template: `
    <template v-if="showLoginPage">
        <main class="login-form">
            <form>
                <div class="form-floating">
                    <input v-model.trim="form.username" :class="{'is-invalid': errors.username}" type="text" class="form-control" id="floatingUsername" placeholder="Имя пользователя">
                    <label for="floatingUsername">Имя пользователя</label>
                </div>
                <div class="form-floating">
                    <input v-model.trim="form.password" :class="{'is-invalid': errors.password}" type="password" class="form-control" id="floatingPassword" placeholder="Пароль">
                    <label for="floatingPassword">Пароль</label>
                </div>
                <button class="w-100 btn btn-primary" type="button" @click="submitForm">Авторизоваться</button>
            </form>
        </main>
    </template>
    `
}
