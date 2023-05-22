"use strict"

import {defineStore} from 'pinia'

export const useTagStore = defineStore({
    id: 'tag',
    state: () => ({
        selectedTag: 0,
    }),
    actions: {
        setSelectedTag(id) {
            this.selectedTag = id
        },
        resetSelectedTag() {
            this.selectedTag = 0
        }
    },
})