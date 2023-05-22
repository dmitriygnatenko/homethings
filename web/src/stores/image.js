"use strict"

import {defineStore} from 'pinia'

export const useImageStore = defineStore({
    id: 'image',
    state: () => ({
        selectedImage: 0,
        selectedImagePlace: 0,
        selectedImageThing: 0,
    }),
    actions: {
        setSelectedImage(id, placeID, thingID) {
            this.selectedImage = id
            this.selectedImagePlace = placeID
            this.selectedImageThing = thingID
        },
        reset() {
            this.selectedImage = 0
            this.selectedImagePlace = 0
            this.selectedImageThing = 0
        }
    },
})
