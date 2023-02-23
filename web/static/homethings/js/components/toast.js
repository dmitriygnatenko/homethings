"use strict"

export const toastComponent = {
    data() {
        return {
            modal: Object,
            message: "",
        }
    },
    methods: {
        showSuccess(message) {
            this.message = message

            this.modal = new bootstrap.Toast(document.getElementById('successToast'))
            this.modal.show()
        },
        showError(message) {
            this.message = message

            this.modal = new bootstrap.Toast(document.getElementById('errorToast'))
            this.modal.show()
        },
    },
    template: `
    <div class="toast-container position-fixed bottom-0 end-0 p-3">
        <div id="successToast" class="toast bg-success text-white" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-body">
                {{ message }}
            </div>
        </div>
        <div id="errorToast" class="toast bg-danger text-white" role="alert" aria-live="assertive" aria-atomic="true">
            <div class="toast-body">
                {{ message }}
            </div>
        </div>
    </div>
    `
}
