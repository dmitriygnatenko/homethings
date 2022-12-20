import * as client from "../client/client.js";
import {formatDate} from "../helpers/date.js";

export const mainPageComponent = {
    props: {
        isAuth: Boolean,
    },
    data() {
        return {
            placesTree: [],
            thingsList: [],
            selectedPlace: 0,
            selectedThing: 0,
        }
    },
    computed: {
        showMainPage() {
            return this.isAuth
        }
    },
    created() {
        if (this.isAuth) {
            this.refreshTree()
        }
    },
    watch: {
        isAuth(val) {
            if (val) {
                this.refreshTree()
            }
        }
    },
    methods: {
        refreshTree() {
            this.selectedPlace = 0
            this.selectedThing = 0
            this.placesTree = []
            this.thingsList = []
            let res = this.request(client.methodGet, client.routeGetPlacesTree)
            if (Array.isArray(res.data.places) && res.data.places.length) {
                this.placesTree = res.data.places
            }
        },
        refreshThings(placeID) {
            this.selectedThing = 0
            this.thingsList = []
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
        addPlace() {
            console.log("Add place " + this.selectedPlace)
        },
        updatePlace() {
            console.log("Edit place " + this.selectedPlace)
        },
        deletePlace() {
            console.log("Delete place " + this.selectedPlace)
        },
        addThing() {
            console.log("Add thing for place " + this.selectedPlace)
        },
        updateThing() {
            console.log("Edit thing " + this.selectedThing)
        },
        deleteThing() {
            console.log("Delete thing " + this.selectedThing)
        },
        request(method, route) {
            let res = client.request(method, route)
            if (res.status !== client.statusOK) {
                this.$emit('eventsetauth', false)
            }
            return res
        },
    },
    template: `
    <template v-if="showMainPage">
        <main class="container-fluid">
            <div class="d-flex flex-grow h-100">
      
                <div class="col-l">
                    <div class="places rounded-3 shadow d-flex flex-column">
                        <div class="header rounded-top">
                            Места
                            <div class="buttons float-end">
                                <button
                                    class="btn add"
                                    title="Добавить место"
                                    @click="addPlace"
                                >
                                    <i class="bi bi-plus-circle-fill"></i>
                                </button>
                                <button
                                    class="btn edit"
                                    title="Редактировать место"
                                    v-if="selectedPlace > 0"
                                    @click="updatePlace"
                                >
                                    <i class="bi bi-pencil-fill"></i>
                                </button>
                                <button 
                                    class="btn delete"
                                    title="Удалить место"
                                    v-if="selectedPlace > 0"
                                    @click="deletePlace"
                                >
                                    <i class="bi bi-trash-fill"></i>
                                </button>
                            </div>
                        </div>
                        <div class="list"></div>
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
                                    @click="addThing"
                                >
                                    <i class="bi bi-plus-circle-fill"></i>
                                </button>
                                <button
                                    class="btn edit"
                                    title="Редактировать вещь"
                                    v-if="selectedThing > 0"
                                    @click="updateThing"
                                >
                                    <i class="bi bi-pencil-fill"></i>
                                </button>
                                <button 
                                    class="btn delete"
                                    title="Удалить вещь"
                                    v-if="selectedThing > 0"
                                    @click="deleteThing"
                                >
                                    <i class="bi bi-trash-fill"></i>
                                </button>
                            </div>
                        </div>
                        <div class="list">
                            <button 
                                class="btn"
                                data-id="{{ thing.id }}" 
                                v-for="thing in thingsList"
                                @click="selectedThing = thing.id"
                                :class="{ selected : selectedThing == thing.id }"
                            >
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
                            <button class="btn add"></button>
                        </div>
                    </div>
                    <div class="list"></div>
                </div>
            
            </div>
        </main>
    </template>
    `
}
