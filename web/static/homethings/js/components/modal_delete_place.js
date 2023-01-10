"use strict"

import * as client from "../client/client.js";

export const modalDeletePlaceComponent = {
    props: {
        selectedPlace: Number,
    },
    data() {
        return {
            modal: Object,
            form: {
                title: "",
            },
        }
    },
    methods: {
        init() {
            if (this.selectedPlace === 0) {
                return
            }
            this.form.title = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetPlace.replace("{id}", this.selectedPlace))
            if (res.status === client.statusOK) {
                this.form.title = res.data.title
            }

            this.modal = new bootstrap.Modal(document.getElementById("delete-place-modal"), {})
            this.modal.show()
        },
        submitForm() {
            let res = client.jsonRequest(client.methodDelete, client.routeDeletePlace.replace("{id}", this.selectedPlace))
            if (res.status === client.statusOK) {
                this.$emit("after-delete-place");
            }

            this.modal.hide()
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="delete-place-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    Подтвердите удаление <b>{{ form.title }}</b>
                    <br><br>
                    <small>Будут удалены все вещи и фото, прикрепленные к данному месту</small>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Отмена</button>
                    <button type="button" class="btn btn-danger btn-sm" @click="submitForm">Удалить</button>
                </div>
            </div>
        </div>
    </div>
    `
}
