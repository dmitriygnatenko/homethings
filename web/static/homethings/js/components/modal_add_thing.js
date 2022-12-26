import * as client from "../client/client.js";

export const modalAddThingComponent = {
    data() {
        return {
            form: {
                place: "",
                title: "",
                desc: "",
            },
        }
    },
    props: {
        selectedPlace: Number,
    },
    methods: {
        initForm() {
            let res = client.request(client.methodGet, client.routeGetPlace.replace("{id}", this.selectedPlace))
            if (res.status === client.statusOK) {
                this.form.place = res.data.title
            }
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="add-thing-modal">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                
                
                    <div class="mb-3 row">
                        <label for="staticEmail" class="col-sm-2 col-form-label">Email</label>
                        <div class="col-sm-10">
                            <input 
                                
                            
                            type="text"
                            readonly
                            class="form-control-plaintext"
                            id="staticEmail"
                            v-model.trim="form.place"
                            
                            >
                        </div>
                    </div>
                    <div class="mb-3 row">
                        <label for="inputPassword" class="col-sm-2 col-form-label">Password</label>
                        <div class="col-sm-10">
                            <input type="password" class="form-control" id="inputPassword">
                        </div>
                    </div>
            
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Отмена</button>
                    <button type="button" class="btn btn-primary btn-sm">Добавить вещь</button>
                </div>
            </div>
        </div>
    </div>
    `
}
