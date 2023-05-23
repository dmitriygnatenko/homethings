<script setup>
import PlaceTreeItem from './PlaceTreeItem.vue'
import {useAuthStore} from '../stores/auth.js'
import {usePlaceStore} from '../stores/place.js'
import {useThingStore} from '../stores/thing.js'
import {useImageStore} from '../stores/image.js'
import {useTagStore} from '../stores/tag.js'
</script>

<script>
import * as client from "../client/client.js";
import {formatDate} from "../helpers/date.js";

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
      imageList: [],
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

    this.imageStore.$onAction(
        ({name, store, args, after, onError}) => {
          if (name === "setSelectedImage" && args.length === 3) {
            // TODO
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
      this.imageList = []
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

      let res = this.request(client.methodGet, client.routeGetPlaceImages.replace("{placeId}", placeID))
      if (Array.isArray(res.data.images) && res.data.images.length) {
        res.data.images.forEach(image => {
          this.imageList.push({
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

      let res = this.request(client.methodGet, client.routeGetThingImages.replace("{thingId}", thingID))
      if (Array.isArray(res.data.images) && res.data.images.length) {
        res.data.images.forEach(image => {
          this.imageList.push({
            "id": image.id,
            "image": image.image,
            "place_id": image.place_id,
            "thing_id": image.thing_id,
            "date": formatDate(image.created_at),
          })
        });
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
      <div class="dropdown logout">
        <button
            type="button"
            class="btn btn-sm dropdown-toggle"
            data-bs-toggle="dropdown">
          <i class="bi bi-person-fill"></i>
          {{ authStore.username }}
        </button>
        <ul class="dropdown-menu">
          <li><button class="dropdown-item" @click="addUser">Добавить пользователя</button></li>
          <li><a class="dropdown-item" @click="updateUsername">Изменить свой логин</a></li>
          <li><a class="dropdown-item" @click="updatePassword">Изменить свой пароль</a></li>
          <li><hr class="dropdown-divider"></li>
          <li><button class="dropdown-item" @click="showTags">Теги</button></li>
          <li><hr class="dropdown-divider"></li>
          <li><button class="dropdown-item" @click="logout">Выход</button></li>
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
                v-for="image in imageList"
                v-on:dblclick="showImage(image.id, image.place_id, image.thing_id)"
                @click="imageStore.setSelectedImage(image.id, image.place_id, image.thing_id)"
                :class="{ selected : imageStore.selectedImage === image.id }">
              <img class="img-fluid" :src="image.image">
              <div class="date">{{ image.date }}</div>
            </button>
          </div>
        </div>
      </div>

    </div>
  </main>
</template>
