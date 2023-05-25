<script>
import * as client from "../../client/client.js"
import {Modal} from 'bootstrap'

export default {
    expose: ['init'],
    emits: ['after-tags'],
    data() {
        return {
            modal: Object,
            updateModal: Object,
            deleteModal: Object,
            tagID: 0,
            form: {
                tagList: [],
                title: "",
                style: "",
            },
            errors: {
                title: "",
                style: "",
            }
        }
    },
    methods: {
        init() {
            this.refreshTags()
            this.modal = new Modal(document.getElementById('modal-tags'), {})
            this.modal.show()
        },
        refreshTags() {
            this.form.tagList = []

            let res = client.jsonRequest(client.methodGet, client.routeGetTags)
            if (res.status === client.statusOK) {
                if (Array.isArray(res.data.tags) && res.data.tags.length) {
                    let obj = this
                    res.data.tags.forEach(tag => {
                        obj.form.tagList.push({
                            "id": tag.id,
                            "title": tag.title,
                            "style": tag.style,
                        })
                    });
                }
            }
        },
        closeForm() {
            this.modal.hide()
            this.$emit("after-tags");
        },
        addTag() {
            this.errors.title = ""
            this.errors.style = ""

            this.tagID = 0
            this.form.title = ""
            this.form.style = ""

            this.updateModal = new Modal(document.getElementById('modal-update-tag'), {})
            this.updateModal.show()
            this.modal.hide()
        },
        updateTag(id) {
            this.errors.title = ""
            this.errors.style = ""

            let res = client.jsonRequest(client.methodGet, client.routeGetTag.replace("{tagId}", id))
            if (res.status === client.statusOK) {
                this.tagID = res.data.id
                this.form.title = res.data.title
                this.form.style = res.data.style

                this.updateModal = new Modal(document.getElementById('modal-update-tag'), {})
                this.updateModal.show()
                this.modal.hide()
            }
        },
        submitUpdateForm() {
            this.errors.title = ""
            this.errors.style = ""

            if (this.form.title === "") {
                this.errors.title = "Название должно быть заполнено"
                return
            }

            if (this.form.style === "") {
                this.errors.style = "Цвет должен быть заполнен"
                return
            }

            let data = {
                title: this.form.title,
                style: this.form.style,
            }

            if (this.tagID > 0) {
                client.jsonRequest(client.methodPut, client.routeUpdateTag.replace("{tagId}", this.tagID), data)
            } else {
                client.jsonRequest(client.methodPost, client.routeAddTag, data)
            }

            this.refreshTags()
            this.closeUpdateForm()
        },
        deleteTag(id) {
            let res = client.jsonRequest(client.methodGet, client.routeGetTag.replace("{tagId}", id))
            if (res.status === client.statusOK) {
                this.tagID = res.data.id
                this.form.title = res.data.title
                this.form.style = res.data.style

                this.deleteModal = new Modal(document.getElementById('modal-delete-tag'), {})
                this.deleteModal.show()
                this.modal.hide()
            }
        },
        submitDeleteForm() {
            client.jsonRequest(client.methodDelete, client.routeDeleteTag.replace("{tagId}", this.tagID))
            this.refreshTags()
            this.closeDeleteForm()
        },
        closeUpdateForm() {
            this.tagID = 0;
            this.updateModal.hide()
            this.modal.show()
        },
        closeDeleteForm() {
            this.tagID = 0;
            this.deleteModal.hide()
            this.modal.show()
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-tags">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    <div class="text-end">
                        <button
                            class="btn add"
                            title="Добавить"
                            @click="addTag()">
                            <i class="bi bi-plus-circle-fill"></i>
                        </button>
                    </div>
                    <div
                        class="row mt-2"
                        v-for="tag in form.tagList">
                        <div class="col-8">
                            <span
                                class="badge rounded-pill"
                                v-bind:style="{ 'background-color': tag.style }">
                                {{ tag.title }}
                            </span>
                        </div>
                        <div class="col-4 text-end">
                            <button
                                class="btn edit"
                                title="Редактировать"
                                @click="updateTag(tag.id)">
                                <i class="bi bi-pencil-fill"></i>
                            </button>
                            <button
                                class="btn delete"
                                title="Удалить"
                                @click="deleteTag(tag.id)">
                                <i class="bi bi-trash-fill"></i>
                            </button>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" @click="closeForm">Закрыть</button>
                </div>
            </div>
        </div>
    </div>

    <div class="modal" tabindex="-1" id="modal-update-tag">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
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
                    <div class="row">
                        <label class="col-sm-3 col-form-label col-form-label-sm">
                            <b>Цвет</b>
                        </label>
                        <div class="col-sm-9">
                            <input
                                type="color"
                                class="form-control form-control-sm"
                                v-model.trim="form.style"
                                :class="{'is-invalid': errors.style}">
                            <div v-if="errors.style" class="invalid-feedback">
                                <small>{{ errors.style }}</small>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm"  @click="closeUpdateForm">Отмена</button>
                    <button type="button" class="btn btn-primary btn-sm" @click="submitUpdateForm">Сохранить</button>
                </div>
            </div>
        </div>
    </div>

    <div class="modal" tabindex="-1" id="modal-delete-tag">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-body">
                    Подтвердите удаление тега <b>{{ form.title }}</b>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" @click="closeDeleteForm">Отмена</button>
                    <button type="button" class="btn btn-danger btn-sm" @click="submitDeleteForm">Удалить</button>
                </div>
            </div>
        </div>
    </div>
</template>
