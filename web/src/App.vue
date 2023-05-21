<script setup>
import LoginPage from './components/LoginPage.vue'
import MainPage from './components/MainPage.vue'
import {useAuthStore} from './stores/auth.js'
</script>

<script>
import * as client from "./client/client.js";

export default {
  data() {
    return {
      authStore: useAuthStore(),
    }
  },
  created() {
    let res = client.jsonRequest(client.methodGet, client.routeCheckAuth)
    if (res.status === client.statusOK) {
      this.authStore.setAuth(res.data.username)
    } else {
      this.authStore.resetAuth()
    }
  },
}
</script>

<template>
  <LoginPage></LoginPage>
  <MainPage></MainPage>
</template>
