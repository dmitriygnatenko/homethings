<script>
import {useThingStore} from '../../stores/thing.js'
import * as client from "../../client/client.js"
import {Modal} from 'bootstrap'

export default {
    setup() {
        const thingStore = useThingStore()
        return {thingStore}
    },
    expose: ['init'],
    emits: ['after-delete-thing'],
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
            let selectedThing = this.thingStore.selectedThing
            if (selectedThing === 0) {
                return
            }

            this.form.placeID = 0
            this.form.title = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetThing.replace("{thingId}", selectedThing))
            if (res.status === client.statusOK) {
                this.form.title = res.data.title
                this.form.placeID = res.data.place_id
            }

            this.modal = new Modal(document.getElementById('modal-delete-thing'), {})
            this.modal.show()
        },
        submitForm() {
            let res = client.jsonRequest(client.methodDelete, client.routeDeleteThing.replace("{thingId}", this.thingStore.selectedThing))
            if (res.status === client.statusOK) {
                this.$emit("after-delete-thing");
            }

            this.modal.hide()
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-delete-thing">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    Подтвердите удаление <b>{{ form.title }}</b>
                    <br><br>
                    <small class="text-secondary">Будут удалены все фото, прикрепленные к данной вещи</small>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Отмена</button>
                    <button type="button" class="btn btn-danger btn-sm" @click="submitForm">Удалить</button>
                </div>
            </div>
        </div>
    </div>
</template>
