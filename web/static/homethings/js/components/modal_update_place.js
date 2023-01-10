'use strict'

import {getPlacesListWithNestedTitles} from "../helpers/places.js";
import * as client from "../client/client.js";

export const modalUpdatePlaceComponent = {
    props: {
        selectedPlace: Number,
    },
    data() {
        return {
            modal: Object,
            form: {
                title: "",
                parentID: 0,
                placesList: [],
            },
            errors: {
                title: "",
            }
        }
    },
    methods: {
        init() {
            if (this.selectedPlace === 0) {
                return
            }
            this.errors.title = ""
            this.form.title = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetPlaces)
            if (res.status === client.statusOK) {
                this.form.placesList = []
                if (Array.isArray(res.data.places) && res.data.places.length) {
                    let obj = this

                    obj.form.placesList.push({
                        "id": 0,
                        "title": "",
                    })

                    getPlacesListWithNestedTitles(res.data.places).forEach(place => {
                        if (place.id !== obj.selectedPlace) {
                            obj.form.placesList.push({
                                "id": place.id,
                                "title": place.title,
                            })
                        } else {
                            obj.form.parentID = place.parent_id
                            obj.form.title = place.title
                        }
                    });
                }
            }

            this.modal = new bootstrap.Modal(document.getElementById('update-place-modal'), {})
            this.modal.show()
        },
        submitForm() {
            if (this.form.title === "") {
                this.errors.title = "Название должно быть заполнено"
                return
            }

            let data = {
                title: this.form.title,
            }

            if (this.form.parentID > 0) {
                data['parent_id'] = this.form.parentID
            }

            let res = client.jsonRequest(client.methodPut, client.routeUpdatePlace.replace("{id}", this.selectedPlace), data)
            if (res.status === client.statusOK) {
                this.$emit("after-update-place");
            }

            this.modal.hide()
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="update-place-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row mb-3">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            Родительское место
                        </label>
                        <div class="col-sm-9">
                            <select v-model="form.parentID" class="form-select form-select-sm">
                                <option v-for="place in form.placesList" :value="place.id">
                                    {{ place.title }}
                                </option>
                            </select>
                        </div>       
                    </div>
                    <div class="row mb-3">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            Название
                        </label>
                        <div class="col-sm-9">
                            <input
                                type="text"
                                class="form-control form-control-sm"
                                v-model.trim="form.title"
                                :class="{'is-invalid': errors.title}">
                            <div v-if="errors.title" class="invalid-feedback">
                                {{ errors.title }}
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
    `
}
