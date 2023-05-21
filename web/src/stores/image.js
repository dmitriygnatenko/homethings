import {defineStore} from 'pinia'

export const useImageStore = defineStore({
    id: 'image',
    state: () => ({
        selectedImage: 0,
        selectedImagePlace: 0,
        selectedImageThing: 0,
    }),
    actions: {
        setSelectedImage(id) {
            this.selectedImage = id
        },
        setSelectedImagePlace(id) {
            this.selectedImagePlace = id
        },
        setSelectedImageThing(id) {
            this.selectedImageThing = id
        },
        resetSelectedImage() {
            this.selectedImage = 0
        },
        resetSelectedImagePlace() {
            this.selectedImagePlace = 0
        },
        resetSelectedImageThing() {
            this.selectedImageThing = 0
        },
        reset() {
            this.selectedImage = 0
            this.selectedImagePlace = 0
            this.selectedImageThing = 0
        }
    },
})
