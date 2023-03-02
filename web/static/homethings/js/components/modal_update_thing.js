'use strict'

import * as client from "../client/client.js";
import {getPlacesListWithNestedTitles} from "../helpers/places.js";

export const modalUpdateThingComponent = {
    props: {
        selectedThing: Number,
    },
    data() {
        return {
            modal: Object,
            form: {
                title: "",
                desc: "",
                placeID: 0,
                placesList: [],
                tagsList: [],
                selectedTags: [],
                initialTags: [],
            },
            errors: {
                title: "",
            },
        }
    },
    methods: {
        init() {
            if (this.selectedThing === 0) {
                return
            }
            this.form.placeID = 0;
            this.form.title = ""
            this.form.desc = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetThing.replace("{thingId}", this.selectedThing))
            if (res.status === client.statusOK) {
                this.form.title = res.data.title
                this.form.desc = res.data.description
                this.form.placeID = res.data.place_id
            }

            if (this.form.placeID === 0) {
                return
            }

            let placesRes = client.jsonRequest(client.methodGet, client.routeGetPlaces)
            if (placesRes.status === client.statusOK) {
                this.form.placesList = []
                if (Array.isArray(placesRes.data.places) && placesRes.data.places.length) {
                    let obj = this

                    getPlacesListWithNestedTitles(placesRes.data.places).forEach(place => {
                        obj.form.placesList.push({
                            "id": place.id,
                            "title": place.title,
                        })
                    });
                }
            }

            let tagsRes = client.jsonRequest(client.methodGet, client.routeGetTags)
            if (tagsRes.status === client.statusOK) {
                this.form.tagsList = []
                if (Array.isArray(tagsRes.data.tags) && tagsRes.data.tags.length) {
                    let obj = this

                    tagsRes.data.tags.forEach(tag => {
                        obj.form.tagsList.push({
                            "id": tag.id,
                            "title": tag.title,
                        })
                    });
                }
            }

            this.form.initialTags = []
            this.form.selectedTags = []
            let thingTagsRes = client.jsonRequest(client.methodGet, client.routeGetThingTags.replace("{thingId}", this.selectedThing))
            if (thingTagsRes.status === client.statusOK) {
                if (Array.isArray(thingTagsRes.data.tags) && thingTagsRes.data.tags.length) {
                    let obj = this

                    thingTagsRes.data.tags.forEach(tag => {
                        obj.form.selectedTags.push(tag.id)
                        obj.form.initialTags.push(tag.id)
                    });
                }
            }

            this.modal = new bootstrap.Modal(document.getElementById('update-thing-modal'), {})
            this.modal.show()
        },
        submitForm() {
            if (this.form.title === "") {
                this.errors.title = "Название должно быть заполнено"
                return
            }

            // Delete
            this.form.initialTags.forEach(tagID => {
                if (this.form.selectedTags.indexOf(tagID) < 0) {
                    client.jsonRequest(client.methodDelete, client.routeDeleteThingTag.replace("{thingId}", this.selectedThing).replace("{tagId}", tagID))
                }
            });

            // Add
            this.form.selectedTags.forEach(tagID => {
                if (this.form.initialTags.indexOf(tagID) < 0) {
                    client.jsonRequest(client.methodPost, client.routeAddThingTag.replace("{thingId}", this.selectedThing).replace("{tagId}", tagID))
                }
            });

            let data = {
                title: this.form.title,
                description: this.form.desc,
                place_id: this.form.placeID,
            }

            let res = client.jsonRequest(client.methodPut, client.routeUpdateThing.replace("{thingId}", this.selectedThing), data)
            if (res.status === client.statusOK) {
                this.$emit("after-update-thing");
            }

            this.modal.hide()
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="update-thing-modal">
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
                    <div class="row mb-3">
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
                    <div class="row">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            <b>Теги</b>
                        </label>
                        <div class="col-sm-9">
                            <div 
                                class="form-check form-check-inline form-control-sm"
                                v-for="tag in form.tagsList" :key="tag.id">
                                <input 
                                    class="form-check-input" 
                                    type="checkbox" 
                                    v-model="form.selectedTags"
                                    :id="'tag-' + tag.id"
                                    :value="tag.id">
                                <label class="form-check-label" :for="'tag-' + tag.id">
                                    {{ tag.title }}
                                </label>
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
