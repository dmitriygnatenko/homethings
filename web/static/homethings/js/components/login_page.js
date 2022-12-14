import {request} from '../client.js';
import {setAuth} from '../auth.js'

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
            hasUsernameError: false,
            hasPasswordError: false,
        };
    },
    methods: {
        submitForm() {
            this.hasUsernameError = this.form.username === "";
            this.hasPasswordError = this.form.password === "";

            if (this.hasUsernameError || this.hasPasswordError) {
                return
            }

            setAuth(this.form.username, this.form.password)

            let res = request("GET", "/api/v1/auth/check")

            if (res.status === 200) {
                this.hasUsernameError = false
                this.hasPasswordError = false
                this.form.username = ""
                this.form.password = ""

                this.$emit('eventsetauth', true)
                return
            }

            this.hasUsernameError = true
            this.hasPasswordError = true
        },
    },
    template: `
    <template v-if="show">
        <main class="login-form">
            <form>
                <div class="form-floating">
                    <input v-model.trim="form.username" :class="{'is-invalid': hasUsernameError}" type="text" class="form-control" id="floatingUsername" placeholder="Имя пользователя">
                    <label for="floatingUsername">Имя пользователя</label>
                </div>
                <div class="form-floating">
                    <input v-model.trim="form.password" :class="{'is-invalid': hasPasswordError}" type="password" class="form-control" id="floatingPassword" placeholder="Пароль">
                    <label for="floatingPassword">Пароль</label>
                </div>
                <button class="w-100 btn btn-primary" type="button" @click="submitForm">Авторизоваться</button>
            </form>
        </main>
    </template>
    `
}
