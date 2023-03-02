"use strict"

import * as client from "../client/client.js";

export const modalAddPlaceComponent = {
    props: {
        selectedPlace: Number,
    },
    data() {
        return {
            modal: Object,
            form: {
                parentID: 0,
                parentTitle: "",
                title: "",
            },
            errors: {
                title: "",
            }
        }
    },
    methods: {
        init() {
            this.errors.title = ""
            this.form.title = ""
            this.form.parentTitle = ""
            this.form.parentID = this.selectedPlace

            if (this.selectedPlace > 0) {
                let res = client.jsonRequest(client.methodGet, client.routeGetPlace.replace("{placeId}", this.selectedPlace))
                if (res.status === client.statusOK) {
                    this.form.parentTitle = res.data.title
                }
            }

            this.modal = new bootstrap.Modal(document.getElementById("add-place-modal"), {})
            this.modal.show()
        },
        submitForm() {
            if (this.form.title === "") {
                this.errors.title = "Название должно быть заполнено"
                return
            }

            let data = {title: this.form.title}
            if (this.form.parentID > 0) {
                data["parent_id"] = this.form.parentID
            }
            let res = client.jsonRequest(client.methodPost, client.routeAddPlace, data)
            if (res.status === client.statusOK) {
                this.$emit("after-add-place", res.data.id);
            }
            this.modal.hide()
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="add-place-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row mb-3" v-if="form.parentTitle">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            <b>Родительское место</b>
                        </label>
                        <div class="col-sm-9">
                            <input 
                                readonly
                                type="text"
                                class="form-control-plaintext form-control-sm"
                                v-model="form.parentTitle">
                        </div>       
                    </div>
                    <div class="row mb-3">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            <b>Название</b>
                        </label>
                        <div class="col-sm-9">
                            <input
                                type="text"
                                class="form-control form-control-sm" 
                                v-model.trim="form.title"
                                :class="{'is-invalid': errors.title}">
                            <div v-if="errors.title" class="invalid-feedback">
                                <small>{{ errors.title }}<small>
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
    `
}
