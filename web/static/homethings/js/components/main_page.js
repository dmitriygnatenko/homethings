"use strict"

import {placeTreeItemComponent} from "./place_tree_item.js";
import {modalAddPlaceComponent} from "./modal_add_place.js";
import {modalUpdatePlaceComponent} from "./modal_update_place.js";
import {modalDeletePlaceComponent} from "./modal_delete_place.js";
import {modalAddThingComponent} from "./modal_add_thing.js";
import {modalUpdateThingComponent} from "./modal_update_thing.js";
import {modalDeleteThingComponent} from "./modal_delete_thing.js";
import {modalAddImageComponent, typePlace} from "./modal_add_image.js";

import * as client from "../client/client.js";
import * as auth from "../auth/auth.js"
import {formatDate} from "../helpers/date.js";

export const mainPageComponent = {
    components: {
        "place-tree-item": placeTreeItemComponent,
        "modal-add-place": modalAddPlaceComponent,
        "modal-update-place": modalUpdatePlaceComponent,
        "modal-delete-place": modalDeletePlaceComponent,
        "modal-add-thing": modalAddThingComponent,
        "modal-update-thing": modalUpdateThingComponent,
        "modal-delete-thing": modalDeleteThingComponent,
        "modal-add-image": modalAddImageComponent,
    },
    emits: ["set-auth"],
    props: {
        isAuth: Boolean,
    },
    data() {
        return {
            placesTree: [],
            thingsList: [],
            imagesList: [],
            selectedPlace: 0,
            selectedThing: 0,
            selectedImage: 0,
            selectedImagePlace: 0,
            selectedImageThing: 0,
        }
    },
    computed: {
        showMainPage() {
            return this.isAuth
        }
    },
    created() {
        if (this.isAuth) {
            this.refreshPlaces()
        }
    },
    watch: {
        isAuth(val) {
            if (val) {
                this.refreshPlaces()
            }
        }
    },
    methods: {
        // Setters
        setSelectedPlace(id) {
            if (this.selectedPlace !== id) {
                this.selectedPlace = id
                this.refreshThings(id)
                this.refreshPlaceImages(id)
            }
        },
        setSelectedThing(id) {
            this.selectedThing = id
            this.refreshThingImages(id)
        },
        setSelectedImage(imageID, placeID, thingID) {
            this.selectedImage = imageID
            this.selectedImagePlace = placeID
            this.selectedImageThing = thingID
        },
        // Request
        request(method, route) {
            let res = client.jsonRequest(method, route)
            if (res.status !== client.statusOK) {
                this.$emit("set-auth", false)
            }
            return res
        },
        // Refresh content methods
        refreshPlaces(id) {
            this.resetPlaces()

            if (id > 0) {
                this.selectedPlace = id
            }

            let res = this.request(client.methodGet, client.routeGetPlacesTree)
            if (Array.isArray(res.data.places) && res.data.places.length) {
                this.placesTree[0].nested = res.data.places
            }
        },
        resetPlaces() {
            this.selectedPlace = 0
            this.placesTree = [{
                place: {"title": "Все", id: 0},
                nested: [],
            }]

            this.resetThings()
        },
        refreshThings(placeID) {
            this.resetThings()

            let res = this.request(client.methodGet, client.routeGetPlaceThings.replace("{id}", placeID))
            if (Array.isArray(res.data.things) && res.data.things.length) {
                res.data.things.forEach(thing => {
                    this.thingsList.push({
                        "id": thing.id,
                        "title": thing.title,
                        "desc": thing.description,
                        "date": formatDate(thing.updated_at),
                    })
                });
            }
        },
        resetThings(){
            this.selectedThing = 0
            this.thingsList = []
            this.resetImages()
        },
        refreshPlaceImages(placeID) {
            this.resetImages()

            let res = this.request(client.methodGet, client.routeGetPlaceImages.replace("{id}", placeID))
            if (Array.isArray(res.data.images) && res.data.images.length) {
                res.data.images.forEach(image => {
                    this.imagesList.push({
                        "id": image.id,
                        "image": image.image,
                        "place_id": image.place_id,
                        "thing_id": image.thing_id,
                        "date": formatDate(image.created_at),
                    })
                });
            }
        },
        refreshThingImages(thingID) {
            this.resetImages()

            let res = this.request(client.methodGet, client.routeGetThingImages.replace("{id}", thingID))
            if (Array.isArray(res.data.images) && res.data.images.length) {
                res.data.images.forEach(image => {
                    this.imagesList.push({
                        "id": image.id,
                        "image": image.image,
                        "place_id": image.place_id,
                        "thing_id": image.thing_id,
                        "date": formatDate(image.created_at),
                    })
                });
            }
        },
        resetImages() {
            this.selectedImage = 0
            this.selectedImagePlace = 0
            this.selectedImageThing = 0
            this.imagesList = []
        },
        // Action methods
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
            this.refreshPlaces(this.selectedPlace)
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
        afterAddThing() {
            this.refreshThings(this.selectedPlace)
        },
        updateThing() {
            this.$refs.modalUpdateThing.init()
        },
        afterUpdateThing() {
            this.refreshPlaces(this.selectedPlace)
            this.refreshThings(this.selectedPlace)
        },
        deleteThing() {
            this.$refs.modalDeleteThing.init()
        },
        afterDeleteThing() {
            this.refreshThings(this.selectedPlace)
        },
        addImage() {
            this.$refs.modalAddImage.init()
        },
        afterAddImage(res) {
            if (res === typePlace) {
                this.refreshPlaceImages(this.selectedPlace)
            } else {
                this.refreshThingImages(this.selectedThing)
            }
        },
        deleteImage() {
            if (this.selectedImage === 0 || (this.selectedImagePlace === 0 && this.selectedImageThing === 0)) {
                return
            }

            if (this.selectedImagePlace > 0) {
                let res = this.request(client.methodDelete, client.routeDeletePlaceImages.replace("{id}", this.selectedImage))
                if (res.status === client.statusOK) {
                    this.refreshPlaceImages(this.selectedImagePlace)
                }
            }

            if (this.selectedImageThing > 0) {
                let res = this.request(client.methodDelete, client.routeDeleteThingImages.replace("{id}", this.selectedImage))
                if (res.status === client.statusOK) {
                    this.refreshThingImages(this.selectedImageThing)
                }
            }
        },
        logout() {
            auth.clearToken()
            this.$emit("set-auth", false)
        }
    },
    template: `
    <template v-if="showMainPage">
        <main class="container-fluid">
            <div class="d-flex flex-grow h-100">
                <button
                    type="button"
                    class="btn btn-sm logout"
                    title="Выход"
                    @click="logout">
                    <i class="bi bi-arrow-right-square"></i>
                </button>
                
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
                                    v-if="selectedPlace > 0"
                                    @click="updatePlace">
                                    <i class="bi bi-pencil-fill"></i>
                                </button>
                                <button 
                                    class="btn delete"
                                    title="Удалить место"
                                    v-if="selectedPlace > 0"
                                    @click="deletePlace">
                                    <i class="bi bi-trash-fill"></i>
                                </button>
                            </div>
                        </div>
                        <div class="list">
                            <ul>
                                <place-tree-item
                                    v-for="item in placesTree"
                                    :item="item"
                                    :selected-place="selectedPlace"
                                    @set-selected-place="setSelectedPlace">
                                </place-tree-item>
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
                                    class="btn add"
                                    title="Добавить вещь"
                                    v-if="selectedPlace > 0"
                                    @click="addThing">
                                    <i class="bi bi-plus-circle-fill"></i>
                                </button>
                                <button
                                    class="btn edit"
                                    title="Редактировать вещь"
                                    v-if="selectedThing > 0"
                                    @click="updateThing">
                                    <i class="bi bi-pencil-fill"></i>
                                </button>
                                <button 
                                    class="btn delete"
                                    title="Удалить вещь"
                                    v-if="selectedThing > 0"
                                    @click="deleteThing">
                                    <i class="bi bi-trash-fill"></i>
                                </button>
                            </div>
                        </div>
                        <div class="list">
                            <button 
                                class="btn"
                                v-for="thing in thingsList"
                                @click="setSelectedThing(thing.id)"
                                :class="{ selected : selectedThing == thing.id }">
                                <div class="title">{{ thing.title }}</div>    
                                <div class="desc" v-if="thing.desc">{{ thing.desc }}</div>
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
                                    v-if="selectedPlace > 0 || selectedThing > 0"
                                    @click="addImage">
                                    <i class="bi bi-plus-circle-fill"></i>
                                </button>
                                <button 
                                    class="btn delete"
                                    title="Удалить фото"
                                    v-if="selectedImage > 0"
                                    @click="deleteImage">
                                    <i class="bi bi-trash-fill"></i>
                                </button>
                            </div>
                        </div>
                        <div class="list">
                            <button 
                                class="btn"
                                v-for="image in imagesList"
                                @click="setSelectedImage(image.id, image.place_id, image.thing_id)"
                                :class="{ selected : selectedImage == image.id }">
                                <img class="img-fluid" :src="image.image">
                                <div class="date">{{ image.date }}</div>
                            </button>
                        </div>
                    </div>
                </div>
            
            </div>
        </main>

        <modal-add-place ref="modalAddPlace" :selected-place="selectedPlace" @after-add-place="afterAddPlace"></modal-add-place>
        <modal-update-place ref="modalUpdatePlace" :selected-place="selectedPlace"@after-update-place="afterUpdatePlace"></modal-update-place>
        <modal-delete-place ref="modalDeletePlace" :selected-place="selectedPlace" @after-delete-place="afterDeletePlace"></modal-delete-place>
        <modal-add-thing ref="modalAddThing" :selected-place="selectedPlace" @after-add-thing="afterAddThing"></modal-add-thing>
        <modal-update-thing ref="modalUpdateThing" :selected-thing="selectedThing" @after-update-thing="afterUpdateThing"></modal-update-thing>
        <modal-delete-thing ref="modalDeleteThing" :selected-thing="selectedThing" @after-delete-thing="afterDeleteThing"></modal-delete-thing>
        <modal-add-image ref="modalAddImage" :selected-place="selectedPlace" :selected-thing="selectedThing" @after-add-image="afterAddImage"></modal-add-image>
    </template>
    `
}
