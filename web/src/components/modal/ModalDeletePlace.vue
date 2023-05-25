<script setup>
import {usePlaceStore} from '../../stores/place.js'
</script>

<script>
import * as client from "../../client/client.js"
import {Modal} from 'bootstrap'

export default {
    expose: ['init'],
    data() {
        return {
            placeStore: usePlaceStore(),
            modal: Object,
            form: {
                title: "",
                error: "",
            },
        }
    },
    methods: {
        init() {
            if (this.placeStore.selectedPlace === 0) {
                return
            }
            this.form.title = ""
            this.form.error = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetPlace.replace("{placeId}", this.placeStore.selectedPlace))
            if (res.status === client.statusOK) {
                this.form.title = res.data.title
            }

            let nestedRes = client.jsonRequest(client.methodGet, client.routeGetNestedPlaces.replace("{parentPlaceId}", this.placeStore.selectedPlace))
            if (nestedRes.status === client.statusOK) {
                if (Array.isArray(nestedRes.data.places) && nestedRes.data.places.length) {
                    this.form.error = "Необходимо вначале удалить вложенные места."
                }
            }

            this.modal = new Modal(document.getElementById("modal-delete-place"), {})
            this.modal.show()
        },
        submitForm() {
            let res = client.jsonRequest(client.methodDelete, client.routeDeletePlace.replace("{placeId}", this.placeStore.selectedPlace))
            if (res.status === client.statusOK) {
                this.$emit("after-delete-place");
            }

            this.modal.hide()
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-delete-place">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div v-if="form.error" class="text-danger text-center">
                        <small>{{ form.error }}</small>
                    </div>
                    <div v-else>
                        Подтвердите удаление <b>{{ form.title }}</b>
                        <br><br>
                        <small class="text-secondary">Будут удалены все вещи и фото, прикрепленные к данному месту</small>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Отмена</button>
                    <button v-if="!form.error" type="button" class="btn btn-danger btn-sm" @click="submitForm">Удалить</button>
                </div>
            </div>
        </div>
    </div>
</template>
