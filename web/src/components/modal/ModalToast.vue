<script>
import {Toast} from 'bootstrap'

export default {
    expose: ['showSuccess', 'showError'],
    data: function () {
        return {
            modal: Object,
            message: "",
        };
    },
    methods: {
        showSuccess(message) {
            this.message = message

            this.modal = new Toast(document.getElementById('modal-toast-success'))
            this.modal.show()
        },
        showError(message) {
            this.message = message

            this.modal = new Toast(document.getElementById('modal-toast-error'))
            this.modal.show()
        },
    },
}
</script>

<template>
    <div class="toast-container position-fixed bottom-0 end-0 p-3">
        <div id="modal-toast-success" class="toast bg-success text-white" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-body">
                {{ message }}
            </div>
        </div>
        <div id="modal-toast-error" class="toast bg-danger text-white" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-body">
                {{ message }}
            </div>
        </div>
    </div>
</template>
