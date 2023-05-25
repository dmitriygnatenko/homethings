<script>
import * as client from "../../client/client.js";
import {Modal} from 'bootstrap'

export default {
    expose: ['init'],
    data() {
        return {
            modal: Object,
            form: {
                password: "",
            },
            errors: {
                password: "",
            },
        }
    },
    methods: {
        init() {
            this.form.password = ""
            this.errors.password = ""

            this.modal = new Modal(document.getElementById("update-password-modal"), {})
            this.modal.show()
        },
        submitForm() {
            this.errors.password = ""
            if (this.form.password === "") {
                this.errors.password = "Пароль должен быть заполнен"
                return
            }

            let data = {
                password: this.form.password,
            }

            let res = client.jsonRequest(client.methodPut, client.routeUpdateUser, data)
            this.$emit("after-update-password", res.status === client.statusOK);
            this.modal.hide()
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="update-password-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row mb-3">
                        <label class="col-sm-5 col-form-label col-form-label-sm">
                            <b>Пароль</b>
                        </label>
                        <div class="col-sm-7">
                            <input
                                type="text"
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
                    <button type="button" class="btn btn-primary btn-sm" @click="submitForm">Сохранить</button>
                </div>
            </div>
        </div>
    </div>
</template>
