"use strict"

import * as client from "../client/client.js";

export const modalAddImageComponent = {
    props: {
        selectedPlace: Number,
        selectedThing: Number,
    },
    data() {
        return {
            maxFiles: 8,
            typePlace: "place",
            typeThing: "thing",
            modal: Object,
            form: {
                files: [null],
                type: null,
                placeTitle: "",
                thingTitle: "",
            },
        }
    },
    methods: {
        init() {
            this.form.files = [null]
            this.form.placeTitle = ""
            this.form.thingTitle = ""

            if (this.selectedPlace > 0) {
                let res = client.jsonRequest(client.methodGet, client.routeGetPlace.replace("{id}", this.selectedPlace))
                if (res.status === client.statusOK) {
                    this.form.type = this.typePlace
                    this.form.placeTitle = "Место: " + res.data.title
                }
            }

            if (this.selectedThing > 0) {
                let res = client.jsonRequest(client.methodGet, client.routeGetThing.replace("{id}", this.selectedThing))
                if (res.status === client.statusOK) {
                    this.form.type = this.typeThing
                    this.form.thingTitle = "Вещь: " + res.data.title
                }
            }

            this.modal = new bootstrap.Modal(document.getElementById("add-image-modal"), {})
            this.modal.show()
        },
        submitForm() {
            const formData = new FormData();

            this.form.files.forEach(function(item) {
                if (item !== undefined && item !== null) {
                    formData.append('files', item)
                }
            });

            if (!formData.has('files')) {
                this.modal.hide()
                return
            }

            if (this.form.type === this.typePlace) {
                formData.set('place_id', this.selectedPlace)
            } else {
                formData.set('thing_id', this.selectedThing)
            }

            let res = client.formDataRequest(client.methodPost, client.routeAddImage, formData)
            if (res.status === client.statusOK) {
                this.$emit("refresh-places", this.selectedPlace);
            }

            this.modal.hide()
        },
        addField() {
            this.form.files.push(null)
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
    template: `
    <div class="modal" tabindex="-1" id="add-image-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row">
                        <div class="col-sm-12">
                            <input type="radio" :value="typePlace" :id="typePlace" v-model="form.type" />
                            <label :for="typePlace" class="form-control-sm">{{ form.placeTitle }}</label>
                        </div>
                    </div>
                    <div class="row" v-if="selectedThing > 0">
                        <div class="col-sm-12">
                            <input type="radio" :value="typeThing" :id="typeThing" v-model="form.type" />
                            <label :for="typeThing" class="form-control-sm">{{ form.thingTitle }}</label>
                        </div>
                    </div>
                    <div
                        class="row mt-3"
                        v-for="(file, index) in form.files"
                        :key="index">
                        <div class="col-sm-9">
                            <input 
                            class="form-control form-control-sm" 
                            accept="image/*"
                            type="file"
                            :data-index="index"
                            @change="onFileChange">
                        </div>       
                        <div class="col-sm-3">
                            <button
                                class="btn add"
                                title="Добавить"
                                v-if="index + 1 == form.files.length && index < maxFiles"
                                @click="addField()">
                                <i class="bi bi-plus-circle-fill"></i>
                            </button>  
                            <button 
                                class="btn delete"
                                title="Удалить"
                                v-if="index + 1 == form.files.length && index > 0"
                                @click="removeField()">
                                <i class="bi bi-trash-fill"></i>
                            </button>
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
