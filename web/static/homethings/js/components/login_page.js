import {request} from '../client.js';

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
        };
    },
    methods: {
        onSubmit() {
        },
        submitForm() {

        },
    },
    template: `
    <template v-if="show">
        <main class="login-form">
            <form>
                <div class="form-floating">
                    <input type="text" class="form-control" id="floatingUsername" placeholder="Имя пользователя">
                    <label for="floatingUsername">Имя пользователя</label>
                </div>
                <div class="form-floating">
                    <input type="password" class="form-control" id="floatingPassword" placeholder="Пароль">
                    <label for="floatingPassword">Пароль</label>
                </div>
                <button class="w-100 btn btn-primary" type="button">Авторизоваться</button>
            </form>
        </main>
    </template>
    `
}
