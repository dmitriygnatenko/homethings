<script setup>
import {useAuthStore} from '../stores/auth.js'
import {usePlaceStore} from '../stores/place.js'
import {useThingStore} from '../stores/thing.js'
import {useImageStore} from '../stores/image.js'
import {useTagStore} from '../stores/tag.js'
</script>

<script>
import * as client from "../client/client.js";

export default {
  data() {
    return {
      authStore: useAuthStore(),
      placeStore: usePlaceStore(),
      thingStore: useThingStore(),
      imageStore: useImageStore()
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
    if (this.authStore.isAuth) {
      this.refreshPlaces()
    }
  },
  methods: {
    // Refresh
    refreshPlaces(id) {
      this.resetPlaces()

      if (id > 0) {
        this.placeStore.setSelectedPlace(id)
      }

      let res = this.request(client.methodGet, client.routeGetPlacesTree)
      if (Array.isArray(res.data.places) && res.data.places.length) {
        this.placeTree[0].nested = res.data.places
      }
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
  }
}
</script>

<style scoped>
@import "../assets/main_page.css";
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
                  class="btn search"
                  title="Поиск вещи"
                  @click="searchThing">
                <i class="bi bi-search"></i>
              </button>
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
                v-on:dblclick="showImage(image.id, image.place_id, image.thing_id)"
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
</template>
