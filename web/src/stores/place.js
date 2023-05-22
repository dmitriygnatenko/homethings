"use strict"

import {defineStore} from 'pinia'

export const usePlaceStore = defineStore({
    id: 'place',
    state: () => ({
        selectedPlace: 0,
    }),
    actions: {
        setSelectedPlace(id) {
            this.selectedPlace = id
        },
        resetSelectedPlace() {
            this.selectedPlace = 0
        }
    },
})