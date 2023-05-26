"use strict"

import {defineStore} from 'pinia'

export const useImageStore = defineStore({
    id: 'image',
    state: () => ({
        imageList: [],
        selectedImage: 0,
        selectedImagePlace: 0,
        selectedImageThing: 0,
    }),
    actions: {
        addImage(image) {
            this.imageList.push(image)
        },
        setSelected(id, placeID, thingID) {
            this.selectedImage = id
            this.selectedImagePlace = placeID
            this.selectedImageThing = thingID
        },
        reset() {
            this.imageList = []
            this.selectedImage = 0
            this.selectedImagePlace = 0
            this.selectedImageThing = 0
        }
    },
})
