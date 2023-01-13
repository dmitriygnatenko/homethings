"use strict"

import * as client from "../client/client.js";
import {getPlacesListWithNestedTitles} from "../helpers/places.js";

export const modalAddThingComponent = {
    props: {
        selectedPlace: Number,
    },
    data() {
        return {
            modal: Object,
            form: {
                title: "",
                desc: "",
                placeID: 0,
                placesList: [],
            },
            errors: {
                title: "",
            },
        }
    },
    methods: {
        init() {
            if (this.selectedPlace === 0) {
                return
            }
            this.form.placeID = this.selectedPlace
            this.form.title = ""
            this.form.desc = ""
            this.errors.title = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetPlaces)
            if (res.status === client.statusOK) {
                this.form.placesList = []
                if (Array.isArray(res.data.places) && res.data.places.length) {
                    let obj = this

                    getPlacesListWithNestedTitles(res.data.places).forEach(place => {
                        obj.form.placesList.push({
                            "id": place.id,
                            "title": place.title,
                        })
                    });
                }
            }

            this.modal = new bootstrap.Modal(document.getElementById("add-thing-modal"), {})
            this.modal.show()
        },
        submitForm() {
            if (this.form.title === "") {
                this.errors.title = "Название должно быть заполнено"
                return
            }

            let data = {
                title: this.form.title,
                description: this.form.desc,
                place_id: this.form.placeID,
            }

            let res = client.jsonRequest(client.methodPost, client.routeAddThing, data)
            if (res.status === client.statusOK) {
                this.$emit("after-add-thing", this.form.placeID, res.data.id);
            }
            this.modal.hide()
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="add-thing-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row mb-3">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            <b>Родительское место</b>
                        </label>
                        <div class="col-sm-9">
                            <select v-model="form.placeID" class="form-select form-select-sm">
                                <option v-for="place in form.placesList" :value="place.id">
                                    {{ place.title }}
                                </option>
                            </select>
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
                    <div class="row">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            <b>Описание</b>
                        </label>
                        <div class="col-sm-9">
                            <textarea 
                                class="form-control form-control-sm"
                                v-model.trim="form.desc">
                             </textarea>
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
