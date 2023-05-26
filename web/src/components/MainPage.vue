<script setup>
import PlaceTreeItem from './PlaceTreeItem.vue'
import ModalAddPlace from './modal/ModalAddPlace.vue'
import ModalUpdatePlace from "./modal/ModalUpdatePlace.vue"
import ModalDeletePlace from "./modal/ModalDeletePlace.vue"
import ModalAddThing from "./modal/ModalAddThing.vue"
import ModalUpdateThing from "./modal/ModalUpdateThing.vue"
import ModalDeleteThing from "./modal/ModalDeleteThing.vue"
import ModalAddImage from "./modal/ModalAddImage.vue"
import ModalSearchThing from "./modal/ModalSearchThing.vue"
import ModalTags from "./modal/ModalTags.vue"
import ModalShowImage from "./modal/ModalShowImage.vue"
import ModalAddUser from './modal/ModalAddUser.vue'
import ModalUpdateUsername from './modal/ModalUpdateUsername.vue'
import ModalUpdatePassword from "./modal/ModalUpdatePassword.vue"
import ModalToast from "./modal/ModalToast.vue"
import {useAuthStore} from '../stores/auth.js'
import {usePlaceStore} from '../stores/place.js'
import {useThingStore} from '../stores/thing.js'
import {useImageStore} from '../stores/image.js'
import {useTagStore} from '../stores/tag.js'
</script>

<script>
import * as auth from "../auth/auth.js"
import * as client from "../client/client.js";
import {formatDate} from "../helpers/date.js";
import {typePlace} from "../stores/thing";

export default {
    data() {
        return {
            authStore: useAuthStore(),
            placeStore: usePlaceStore(),
            thingStore: useThingStore(),
            imageStore: useImageStore(),
            tagStore: useTagStore(),
            placeTree: [],
            thingList: [],
        };
    },
    computed: {
        show() {
            return this.authStore.isAuth
        }
    },
    created() {
        this.placeStore.$onAction(
            ({name, store, args, after, onError}) => {
                if (name === "setSelectedPlace" && args.length) {
                    if (args[0] !== this.placeStore.selectedPlace) {
                        after(() => {
                            let placeID = this.placeStore.selectedPlace

                            this.resetTags()

                            if (placeID === 0) {
                                this.resetThings()
                                return
                            }

                            this.refreshThings(placeID)
                            this.refreshPlaceImages(placeID)
                        })
                    }
                }
            }
        )

        this.thingStore.$onAction(
            ({name, store, args, after, onError}) => {
                if (name === "setSelectedThing" && args.length) {
                    if (args[0] !== this.thingStore.selectedThing) {
                        after(() => {
                            let thingID = this.thingStore.selectedThing

                            this.refreshThingImages(thingID)
                        })
                    }
                }
            }
        )

        this.authStore.$onAction(
            ({name, store, args, after, onError}) => {
                switch (name) {
                    case "setAuth":
                        this.refreshPlaces()
                        break
                    case "resetAuth":
                        this.resetPlaces()
                        break
                }
            }
        )

        // Refresh places after start
        if (this.authStore.isAuth) {
            this.refreshPlaces()
        }
    },
    methods: {
        // Request

        request(method, route) {
            let res = client.jsonRequest(method, route)
            if (res.status !== client.statusOK) {
                this.authStore.resetAuth()
            }
            return res
        },

        // Reset

        resetPlaces() {
            this.placeTree = [{
                place: {"title": "Все", id: 0},
                nested: [],
            }]
            this.placeStore.resetSelectedPlace()
            this.resetThings()
        },

        resetThings() {
            this.thingList = []
            this.thingStore.resetSelectedThing()
            this.resetImages()
        },

        resetImages() {
            this.imageStore.reset()
        },

        resetTags() {
            this.tagStore.resetSelectedTag()
        },

        // Refresh

        refreshPlaces(placeID) {
            this.resetPlaces()

            if (placeID > 0) {
                this.placeStore.setSelectedPlace(placeID)
            }

            let res = this.request(client.methodGet, client.routeGetPlacesTree)
            if (Array.isArray(res.data.places) && res.data.places.length) {
                this.placeTree[0].nested = res.data.places
            }
        },

        refreshThings(placeID) {
            this.resetThings()
            let obj = this

            let res = this.request(client.methodGet, client.routeGetPlaceThings.replace("{placeId}", placeID))
            if (Array.isArray(res.data.things) && res.data.things.length) {
                res.data.things.forEach(thing => {
                    let show = false

                    if (obj.tagStore.selectedTag === 0) {
                        show = true
                    } else if (obj.tagStore.selectedTag > 0 && thing.tags) {
                        thing.tags.forEach(tag => {
                            if (tag.id === obj.tagStore.selectedTag) {
                                show = true
                            }
                        })
                    }

                    if (show) {
                        obj.thingList.push({
                            "id": thing.id,
                            "title": thing.title,
                            "desc": thing.description,
                            "date": formatDate(thing.updated_at),
                            "tags": thing.tags
                        })
                    }
                });
            }
        },

        refreshPlaceImages(placeID) {
            this.resetImages()

            let host = client.getHost()
            let res = this.request(client.methodGet, client.routeGetPlaceImages.replace("{placeId}", placeID))
            if (Array.isArray(res.data.images) && res.data.images.length) {
                res.data.images.forEach(image => {
                    this.imageStore.addImage({
                        "id": image.id,
                        "image": host + image.image,
                        "place_id": image.place_id,
                        "thing_id": image.thing_id,
                        "date": formatDate(image.created_at),
                    })
                });
            }
        },

        refreshThingImages(thingID) {
            this.resetImages()

            let host = client.getHost()
            let res = this.request(client.methodGet, client.routeGetThingImages.replace("{thingId}", thingID))
            if (Array.isArray(res.data.images) && res.data.images.length) {
                res.data.images.forEach(image => {
                    this.imageStore.addImage({
                        "id": image.id,
                        "image": host + image.image,
                        "place_id": image.place_id,
                        "thing_id": image.thing_id,
                        "date": formatDate(image.created_at),
                    })
                });
            }
        },

        // Actions

        addPlace() {
            this.$refs.modalAddPlace.init();
        },

        afterAddPlace(id) {
            this.refreshPlaces(id)
        },

        updatePlace() {
            this.$refs.modalUpdatePlace.init()
        },

        afterUpdatePlace() {
            this.refreshPlaces(this.placeStore.selectedPlace)
        },

        deletePlace() {
            this.$refs.modalDeletePlace.init()
        },

        afterDeletePlace() {
            this.refreshPlaces()
        },

        addThing() {
            this.$refs.modalAddThing.init()
        },

        afterAddThing(placeID, thingID) {
            this.resetTags()
            this.refreshPlaces(placeID)
            this.refreshThings(placeID)
            this.thingStore.setSelectedThing(thingID)
        },

        updateThing() {
            this.$refs.modalUpdateThing.init()
        },

        afterUpdateThing() {
            let selectedThing = this.thingStore.selectedThing
            this.resetTags()
            this.refreshPlaces(this.placeStore.selectedPlace)
            this.refreshThings(this.placeStore.selectedPlace)
            this.thingStore.setSelectedThing(selectedThing)
        },

        deleteThing() {
            this.$refs.modalDeleteThing.init()
        },

        afterDeleteThing() {
            this.resetTags()
            this.refreshThings(this.placeStore.selectedPlace)
        },

        addImage() {
            this.$refs.modalAddImage.init()
        },

        afterAddImage(res) {
            if (res === typePlace) {
                this.refreshPlaceImages(this.placeStore.selectedPlace)
            } else {
                this.refreshThingImages(this.thingStore.selectedThing)
            }
        },

        searchThing() {
            this.$refs.modalSearchThing.init()
        },

        afterSearchThing(placeID, thingID) {
            this.resetTags()
            this.refreshPlaces(placeID)
            this.refreshThings(placeID)
            this.thingStore.setSelectedThing(thingID)
        },

        afterFilterTag(tagID) {
            this.tagStore.setSelectedTag(tagID)
            this.refreshThings(this.placeStore.selectedPlace)
        },

        showTags() {
            this.$refs.modalTags.init();
        },

        afterTags() {
            this.resetTags()
            this.refreshThings(this.placeStore.selectedPlace)
        },

        selectImage(imageID, placeID, thingID) {
            this.imageStore.setSelected(imageID, placeID, thingID)
        },

        showImage(imageID, placeID, thingID) {
            this.$refs.modalShowImage.init(imageID, placeID, thingID)
        },

        deleteImage() {
            let imageID = this.imageStore.selectedImage
            let placeID = this.imageStore.selectedImagePlace
            let thingID = this.imageStore.selectedImageThing

            if (imageID === 0 || (placeID === 0 && thingID === 0)) {
                return
            }

            if (placeID > 0) {
                let res = this.request(client.methodDelete, client.routeDeletePlaceImages.replace("{imageId}", imageID))
                if (res.status === client.statusOK) {
                    this.refreshPlaceImages(placeID)
                }
            }

            if (thingID > 0) {
                let res = this.request(client.methodDelete, client.routeDeleteThingImages.replace("{imageId}", imageID))
                if (res.status === client.statusOK) {
                    this.refreshThingImages(thingID)
                }
            }
        },

        logout() {
            auth.clearToken()
            this.authStore.resetAuth()
        },

        addUser() {
            this.$refs.modalAddUser.init()
        },

        afterAddUser(success) {
            if (success) {
                this.$refs.modalToast.showSuccess("Пользователь добавлен")
            } else {
                this.$refs.modalToast.showError("Ошибка при добавлении пользователя")
            }
        },

        updateUsername() {
            this.$refs.modalUpdateUsername.init()
        },

        afterUpdateUsername(success) {
            if (success) {
                this.logout()
            } else {
                this.$refs.modalToast.showError("Ошибка при изменении имени пользователя")
            }
        },

        updatePassword() {
            this.$refs.modalUpdatePassword.init()
        },

        afterUpdatePassword(success) {
            if (success) {
                this.logout()
            } else {
                this.$refs.modalToast.showError("Ошибка при изменении пароля пользователя")
            }
        },
    }
}
</script>

<style>
@import "../assets/main_page.css";
@import "../assets/modal.css";
</style>

<template>
    <main class="container-fluid" v-if="show">
        <div class="d-flex flex-grow h-100">
            <div class="dropdown user-top">
                <button
                    type="button"
                    class="btn btn-sm dropdown-toggle"
                    data-bs-toggle="dropdown"
                    data-bs-target="#dropdown-user-menu">
                    <i class="bi bi-person-fill"></i>
                    {{ this.authStore.username }}
                </button>
                <ul class="dropdown-menu" id="dropdown-user-menu">
                    <li>
                        <button class="dropdown-item" @click="addUser">Добавить пользователя</button>
                    </li>
                    <li><a class="dropdown-item" @click="updateUsername">Изменить свой логин</a></li>
                    <li><a class="dropdown-item" @click="updatePassword">Изменить свой пароль</a></li>
                    <li>
                        <hr class="dropdown-divider">
                    </li>
                    <li>
                        <button class="dropdown-item" @click="showTags">Теги</button>
                    </li>
                    <li>
                        <hr class="dropdown-divider">
                    </li>
                    <li>
                        <button class="dropdown-item" @click="logout">Выход</button>
                    </li>
                </ul>
            </div>

            <div class="col-l">
                <div class="places rounded-3 shadow d-flex flex-column">
                    <div class="header rounded-top">
                        Места
                        <div class="buttons float-end">
                            <button
                                class="btn add"
                                title="Добавить место"
                                @click="addPlace">
                                <i class="bi bi-plus-circle-fill"></i>
                            </button>
                            <button
                                class="btn edit"
                                title="Редактировать место"
                                v-if="placeStore.selectedPlace > 0"
                                @click="updatePlace">
                                <i class="bi bi-pencil-fill"></i>
                            </button>
                            <button
                                class="btn delete"
                                title="Удалить место"
                                v-if="placeStore.selectedPlace > 0"
                                @click="deletePlace">
                                <i class="bi bi-trash-fill"></i>
                            </button>
                        </div>
                    </div>
                    <div class="list">
                        <ul>
                            <PlaceTreeItem v-for="item in placeTree" :item="item"></PlaceTreeItem>
                        </ul>
                    </div>
                </div>
            </div>

            <div class="col-c">
                <div class="things rounded-3 shadow d-flex flex-column">
                    <div class="header rounded-top">
                        Вещи
                        <div class="buttons float-end">
                            <button
                                class="btn search"
                                title="Поиск вещи"
                                @click="searchThing">
                                <i class="bi bi-search"></i>
                            </button>
                            <button
                                class="btn add"
                                title="Добавить вещь"
                                v-if="placeStore.selectedPlace > 0"
                                @click="addThing">
                                <i class="bi bi-plus-circle-fill"></i>
                            </button>
                            <button
                                class="btn edit"
                                title="Редактировать вещь"
                                v-if="thingStore.selectedThing > 0"
                                @click="updateThing">
                                <i class="bi bi-pencil-fill"></i>
                            </button>
                            <button
                                class="btn delete"
                                title="Удалить вещь"
                                v-if="thingStore.selectedThing > 0"
                                @click="deleteThing">
                                <i class="bi bi-trash-fill"></i>
                            </button>
                        </div>
                    </div>
                    <div class="list">
                        <button
                            class="btn"
                            v-for="thing in thingList"
                            @click="thingStore.setSelectedThing(thing.id)"
                            :class="{ selected : thingStore.selectedThing === thing.id }">
                            <div class="title">{{ thing.title }}</div>
                            <div class="desc" v-if="thing.desc">{{ thing.desc }}</div>
                            <div class="tags" v-if="thing.tags">
                                <span
                                    class="badge rounded-pill"
                                    v-for="tag in thing.tags"
                                    v-bind:style="{ 'background-color': tag.style }">
                                    {{ tag.title }}
                                </span>
                            </div>
                            <div class="date">{{ thing.date }}</div>
                        </button>
                    </div>
                </div>
            </div>

            <div class="col-r">
                <div class="info rounded-3 shadow d-flex flex-column">
                    <div class="header rounded-top">
                        Фото
                        <div class="buttons float-end">
                            <button
                                class="btn add"
                                title="Добавить фото"
                                v-if="placeStore.selectedPlace > 0 || thingStore.selectedThing > 0"
                                @click="addImage">
                                <i class="bi bi-plus-circle-fill"></i>
                            </button>
                            <button
                                class="btn delete"
                                title="Удалить фото"
                                v-if="imageStore.selectedImage > 0"
                                @click="deleteImage">
                                <i class="bi bi-trash-fill"></i>
                            </button>
                        </div>
                    </div>
                    <div class="list">
                        <button
                            class="btn"
                            v-for="image in imageStore.imageList"
                            v-on:dblclick="showImage(image.id, image.place_id, image.thing_id)"
                            @click="selectImage(image.id, image.place_id, image.thing_id)"
                            :class="{ selected : this.imageStore.selectedImage === image.id &&
                                this.imageStore.selectedImagePlace === image.place_id &&
                                this.imageStore.selectedImageThing === image.thing_id }">
                            <img class="img-fluid" :src="image.image">
                            <div class="date">{{ image.date }}</div>
                        </button>
                    </div>
                </div>
            </div>

        </div>
    </main>

    <ModalToast ref="modalToast"></ModalToast>
    <ModalAddPlace ref="modalAddPlace" @after-add-place="afterAddPlace"></ModalAddPlace>
    <ModalUpdatePlace ref="modalUpdatePlace" @after-update-place="afterUpdatePlace"></ModalUpdatePlace>
    <ModalDeletePlace ref="modalDeletePlace" @after-delete-place="afterDeletePlace"></ModalDeletePlace>
    <ModalAddThing ref="modalAddThing" @after-add-thing="afterAddThing"></ModalAddThing>
    <ModalUpdateThing ref="modalUpdateThing" @after-update-thing="afterUpdateThing"></ModalUpdateThing>
    <ModalDeleteThing ref="modalDeleteThing" @after-delete-thing="afterDeleteThing"></ModalDeleteThing>
    <ModalAddImage ref="modalAddImage" @after-add-image="afterAddImage"></ModalAddImage>
    <ModalSearchThing ref="modalSearchThing" @after-search-thing="afterSearchThing" @after-filter-tag="afterFilterTag"></ModalSearchThing>
    <ModalTags ref="modalTags" @after-tags="afterTags"></ModalTags>
    <ModalShowImage ref="modalShowImage"></ModalShowImage>
    <ModalAddUser ref="modalAddUser" @after-add-user="afterAddUser"></ModalAddUser>
    <ModalUpdateUsername ref="modalUpdateUsername" @after-update-username="afterUpdateUsername"></ModalUpdateUsername>
    <ModalUpdatePassword ref="modalUpdatePassword" @after-update-password="afterUpdatePassword"></ModalUpdatePassword>
</template>
