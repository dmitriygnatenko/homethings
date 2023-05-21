<script setup>
import {useAuthStore} from '../stores/auth.js'
</script>

<script>
import * as client from "../client/client.js";
import * as auth from "../auth/auth.js"

export default {
  data() {
    return {
      authStore: useAuthStore(),
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
    show() {
      return !this.authStore.isAuth
    }
  },
  methods: {
    submitForm() {
      this.errors.username = this.form.username === "";
      this.errors.password = this.form.password === "";
      if (this.errors.username || this.errors.password) {
        return
      }

      let res = client.jsonRequest(client.methodPost, client.routeLogin, {
        "username": this.form.username,
        "password": this.form.password,
      })

      if (res.status === client.statusOK && res.data.token !== undefined) {
        this.errors.username = false
        this.errors.password = false
        this.form.username = ""
        this.form.password = ""

        auth.setToken(res.data.token)
        this.authStore.setAuth(this.form.username)

        return
      }

      this.errors.username = true
      this.errors.password = true

      auth.clearToken()
      this.authStore.resetAuth()
    },
  }
}
</script>

<style scoped>
@import "../assets/login_page.css";
</style>

<template>
  <main class="login-form" v-if="show">
    <form>
      <div class="form-floating">
        <input
            type="text"
            class="form-control"
            id="formUsername"
            placeholder="Имя пользователя"
            v-model.trim="form.username"
            v-on:keyup.enter="submitForm"
            :class="{'is-invalid': errors.username}"
        >
        <label for="formUsername">Имя пользователя</label>
      </div>
      <div class="form-floating">
        <input
            type="password"
            class="form-control"
            id="formPassword"
            placeholder="Пароль"
            v-model.trim="form.password"
            v-on:keyup.enter="submitForm"
            :class="{'is-invalid': errors.password}"
        >
        <label for="formPassword">Пароль</label>
      </div>
      <button class="w-100 btn btn-primary" type="button" @click="submitForm">Авторизоваться</button>
    </form>
  </main>
</template>
