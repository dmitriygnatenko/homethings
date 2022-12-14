import * as client from '../client.js';
import * as auth from '../auth.js'

export const loginPageComponent = {
    props: {
        show: Boolean,
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
    methods: {
        submitForm() {
            this.errors.username = this.form.username === "";
            this.errors.password = this.form.password === "";

            if (this.errors.username || this.errors.password) {
                return
            }

            auth.setAuth(this.form.username, this.form.password)

            let res = client.request("GET", client.routeCheckAuth)

            if (res.status === 200) {
                this.errors.username = false
                this.errors.password = false
                this.form.username = ""
                this.form.password = ""

                this.$emit('eventsetauth', true)
                return
            }

            auth.clearAuth()
            this.errors.username = true
            this.errors.password = true
        },
    },
    template: `
    <template v-if="show">
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
