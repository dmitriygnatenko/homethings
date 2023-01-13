"use strict"

import * as client from "../client/client.js";

export const modalSearchThingComponent = {
    data() {
        return {
            modal: Object,
            loading: false,
            empty: false,
            thingsList: [],
            form: {
                search: "",
            },
            errors: {
                search: "",
            },
        }
    },
    methods: {
        init() {
            this.thingsList = []
            this.empty = false
            this.loading = false
            this.form.search = ""
            this.errors.search = ""

            this.modal = new bootstrap.Modal(document.getElementById("search-thing-modal"), {})
            this.modal.show()
        },
        submitForm() {
            this.empty = false
            this.thingsList = []

            if (this.form.search === "") {
                this.errors.search = "Заполните поле для поиска"
                return
            }

            if (this.form.search.length < 3) {
                this.errors.search = "Строка для поиска должна быть более 3 символов"
                return
            }

            this.errors.search = ""
            this.loading = true

            let search = encodeURIComponent(this.form.search)
            let res = client.jsonRequest(client.methodGet, client.routeSearchThings.replace("{search}", search))
            if (Array.isArray(res.data.things) && res.data.things.length) {
                let obj = this
                res.data.things.forEach(thing => {
                    obj.thingsList.push({
                        "id": thing.id,
                        "place_id": thing.place_id,
                        "title": thing.title,
                    })
                });
            }

            this.loading = false

            if (this.thingsList.length === 0) {
                this.empty = true
            }
        },
        showResult(thingID, placeID) {
            this.modal.hide()
            this.$emit("afterSearchThing", placeID, thingID);
        }
    },
    template: `
    <div class="modal" tabindex="-1" id="search-thing-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="row mb-3">
                        <div class="col-sm-9">
                            <input
                                type="text"
                                class="form-control form-control-sm"
                                v-model.trim="form.search"
                                :class="{ 'is-invalid': errors.search }">
                                <div class="invalid-feedback">
                                    <small>{{ errors.search }}</small>
                                </div>
                        </div>
                        <div class="col-sm-3">
                            <button type="button" class="search btn btn-primary btn-sm w-100" @click="submitForm">Поиск</button>
                        </div>
                    </div>  
                    
                    <div class="row mb-3 search-results">
                        <div class="text-center" v-if="loading">
                            <div class="spinner-border" role="status"></div>
                        </div>
                        <div class="text-center text-secondary" v-if="empty">
                            <small>Ничего не найдено</small>
                        </div>
                        <a 
                            href="#"
                            class="link-primary"
                            v-for="thing in thingsList"
                            @click="showResult(thing.id, thing.place_id)"
                            >
                            {{ thing.title }}
                        </a>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Отмена</button>
                </div>  
            </div>
        </div>
    </div>
    `
}
