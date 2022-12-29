'use strict'

import * as client from "../client/client.js";

export const modalDeleteThingComponent = {
    props: {
        selectedThing: Number,
    },
    data() {
        return {
            modal: Object,
            form: {
                title: "",
                placeID: 0,
            },
        }
    },
    methods: {
        init() {
            if (this.selectedThing === 0) {
                return
            }
            this.form.placeID = 0
            this.form.title = ""

            let res = client.request(client.methodGet, client.routeGetThing.replace("{id}", this.selectedThing))
            if (res.status === client.statusOK) {
                this.form.title = res.data.title
                this.form.placeID = res.data.place_id
            }

            this.modal = new bootstrap.Modal(document.getElementById('delete-thing-modal'), {})
            this.modal.show()
        },
        submitForm() {
            let res = client.request(client.methodDelete, client.routeDeleteThing.replace("{id}", this.selectedThing))
            if (res.status === client.statusOK) {
                this.$emit("refresh-things", this.form.placeID);
            }

            this.modal.hide()
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="delete-thing-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    Подтвердите удаление <b>{{ form.title }}</b>
                    <br><br>
                    <small>Будут удалены все фото, прикрепленные к данной вещи</small>
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
