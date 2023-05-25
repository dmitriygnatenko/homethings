<script>
import * as client from "../../client/client.js";
import {Modal} from 'bootstrap'

export default {
    expose: ['init'],
    data() {
        return {
            modal: Object,
            form: {
                username: "",
                password: "",

            },
            errors: {
                username: "",
                password: "",
            },
        }
    },
    methods: {
        init() {
            this.form.username = ""
            this.form.password = ""
            this.errors.username = ""
            this.errors.password = ""

            this.modal = new Modal(document.getElementById("modal-add-user"), {})
            this.modal.show()
        },
        submitForm() {
            this.errors.username = ""
            this.errors.password = ""

            if (this.form.username === "") {
                this.errors.username = "Имя пользователя должно быть заполнено"

            }

            if (this.form.password === "") {
                this.errors.password = "Пароль должен быть заполнен"
            }

            if (this.errors.username !== "" || this.errors.password !== "") {
                return
            }

            let data = {
                username: this.form.username,
                password: this.form.password,
            }

            let res = client.jsonRequest(client.methodPost, client.routeAddUser, data)
            this.$emit("after-add-user", res.status === client.statusOK);
            this.modal.hide()
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-add-user">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row mb-3">
                        <label class="col-sm-5 col-form-label col-form-label-sm">
                            <b>Имя пользователя</b>
                        </label>
                        <div class="col-sm-7">
                            <input
                                type="text"
                                class="form-control form-control-sm"
                                v-model.trim="form.username"
                                :class="{'is-invalid': errors.username}">
                            <div v-if="errors.username" class="invalid-feedback">
                                <small>{{ errors.username }}</small>
                            </div>
                        </div>
                    </div>
                    <div class="row mb-3">
                        <label class="col-sm-5 col-form-label col-form-label-sm">
                            <b>Пароль</b>
                        </label>
                        <div class="col-sm-7">
                            <input
                                type="password"
                                class="form-control form-control-sm"
                                v-model.trim="form.password"
                                :class="{'is-invalid': errors.password}">
                            <div v-if="errors.password" class="invalid-feedback">
                                <small>{{ errors.password }}</small>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Отмена</button>
                    <button type="button" class="btn btn-primary btn-sm" @click="submitForm">Добавить</button>
                </div>
            </div>
        </div>
    </div>
</template>
