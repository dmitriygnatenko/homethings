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
            },
            errors: {
                username: "",
            },
        }
    },
    methods: {
        init() {
            this.form.username = ""
            this.errors.username = ""

            this.modal = new Modal(document.getElementById("modal-update-username"), {})
            this.modal.show()
        },
        submitForm() {
            this.errors.username = ""
            if (this.form.username === "") {
                this.errors.username = "Имя пользователя должно быть заполнено"
                return
            }

            let data = {
                username: this.form.username,
            }

            let res = client.jsonRequest(client.methodPut, client.routeUpdateUser, data)
            this.$emit("after-update-username", res.status === client.statusOK);
            this.modal.hide()
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-update-username">
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
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Отмена</button>
                    <button type="button" class="btn btn-primary btn-sm" @click="submitForm">Сохранить</button>
                </div>
            </div>
        </div>
    </div>
</template>
