<script>
import * as client from "../../client/client.js"
import {Modal} from 'bootstrap'

export default {
    expose: ['init'],
    data() {
        return {
            notificationList: [],
            modal: Object,
        }
    },
    methods: {
        init(req) {
            if (req.length === 0) {
                return
            }

            this.notificationList = req
            this.modal = new Modal(document.getElementById('modal-expired-notifications'), {})
            this.modal.show()
        },
        showResult(placeID, thingID) {
            this.modal.hide()
            this.$emit("after-expired-notification", placeID, thingID);
        },
        deleteNotification(thingID) {
            let res = client.jsonRequest(client.methodDelete, client.routeDeleteNotification.replace("{thingId}", thingID))
            if (res.status === client.statusOK) {
                for(let i = 0; i < this.notificationList.length; i++){
                    if (this.notificationList[i].thing_id === thingID) {
                        this.notificationList.splice(i, 1);
                    }
                }
            }
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-expired-notifications">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="mb-2">
                        Напоминания
                    </div>
                    <div
                        class="row notification-results"
                        v-for="notif in notificationList">
                        <div class="col-8">
                            <a
                                href="#"
                                class="link-primary"
                                @click="showResult(notif.place_id, notif.thing_id)">
                                {{ notif.thing_title }}
                                ({{ notif.place_title }})
                            </a>
                        </div>
                        <div class="col-4 text-end">
                            <button
                                class="btn delete"
                                title="Удалить напоминание"
                                @click="deleteNotification(notif.thing_id)">
                                <i class="bi bi-trash-fill"></i>
                            </button>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Закрыть</button>
                </div>
            </div>
        </div>
    </div>
</template>
