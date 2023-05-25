"use strict"

import {defineStore} from 'pinia'

export const typePlace = "place"
export const typeThing = "thing"

export const useThingStore = defineStore({
    id: 'thing',
    state: () => ({
        selectedThing: 0,
    }),
    actions: {
        setSelectedThing(id) {
            this.selectedThing = id
        },
        resetSelectedThing() {
            this.selectedThing = 0
        }
    },
})