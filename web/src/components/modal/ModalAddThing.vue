<script>
import {usePlaceStore} from '../../stores/place.js'
import * as client from "../../client/client.js"
import {Modal} from 'bootstrap'
import {getPlacesListWithNestedTitles} from "../../helpers/places.js";

export default {
    setup() {
        const placeStore = usePlaceStore()
        return {placeStore}
    },
    expose: ['init'],
    emits: ['after-add-thing'],
    data() {
        return {
            modal: Object,
            maxFiles: 4,
            form: {
                title: "",
                desc: "",
                files: null,
                placeID: 0,
                placeList: [],
            },
            errors: {
                title: "",
            },
        }
    },
    methods: {
        init() {
            let selectedPlace = this.placeStore.selectedPlace
            if (selectedPlace === 0) {
                return
            }

            this.form.files = [""]
            this.form.placeID = selectedPlace
            this.form.title = ""
            this.form.desc = ""
            this.errors.title = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetPlaces)
            if (res.status === client.statusOK) {
                this.form.placeList = []
                if (Array.isArray(res.data.places) && res.data.places.length) {
                    let obj = this

                    getPlacesListWithNestedTitles(res.data.places).forEach(place => {
                        obj.form.placeList.push({
                            "id": place.id,
                            "title": place.title,
                        })
                    });
                }
            }

            this.modal = new Modal(document.getElementById("modal-add-thing"), {})
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

            let formData = new FormData();
            this.form.files.forEach(function (item) {
                if (item !== undefined && item !== null && item !== '') {
                    formData.append('files', item)
                }
            });

            let res = client.jsonRequest(client.methodPost, client.routeAddThing, data)
            if (res.status === client.statusOK && res.data.id > 0) {
                if (formData.has('files')) {
                    formData.set('thing_id', res.data.id)
                    client.formDataRequest(client.methodPost, client.routeAddImage, formData)
                    this.form.files = null
                }

                this.$emit("after-add-thing", this.form.placeID, res.data.id);
            }
            this.modal.hide()
        },
        addField() {
            this.form.files.push("")
        },
        removeField() {
            this.form.files.pop()
        },
        onFileChange(e) {
            if (!e.target.files.length) {
                return;
            }

            let index = e.target.getAttribute("data-index")
            this.form.files[index] = e.target.files[0]
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-add-thing">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row mb-3">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            <b>Родительское место</b>
                        </label>
                        <div class="col-sm-9">
                            <select v-model="form.placeID" class="form-select form-select-sm">
                                <option v-for="place in form.placeList" :value="place.id">
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
                                <small>{{ errors.title }}</small>
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
                            <b>Фото</b>
                        </label>
                        <div class="col-sm-9">
                            <div
                                class="row mb-1"
                                v-for="(file, index) in form.files"
                                :key="index">
                                <div class="col-8">
                                    <input
                                        class="form-control form-control-sm"
                                        accept="image/*"
                                        type="file"
                                        :data-index="index"
                                        @change="onFileChange">
                                </div>
                                <div class="col-4">
                                    <button
                                        class="btn add"
                                        title="Добавить"
                                        v-if="index + 1 === form.files.length && index < maxFiles"
                                        @click="addField()">
                                        <i class="bi bi-plus-circle-fill"></i>
                                    </button>
                                    <button
                                        class="btn delete"
                                        title="Удалить"
                                        v-if="index + 1 === form.files.length && index > 0"
                                        @click="removeField()">
                                        <i class="bi bi-trash-fill"></i>
                                    </button>
                                </div>
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
