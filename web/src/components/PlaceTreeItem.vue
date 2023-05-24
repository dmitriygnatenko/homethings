<script setup>
import {usePlaceStore} from '../stores/place.js'
</script>

<script>
export default {
    props: {
        item: Object,
    },
    data: function () {
        return {
            placeStore: usePlaceStore(),
            open: false,
        };
    },
    computed: {
        isFolder() {
            return this.item.nested && this.item.nested.length;
        },
        isOpen() {
            if (this.item.nested && this.item.nested.length) {
                let obj = this
                this.item.nested.forEach(function (item) {
                    if (item.place.id === obj.placeStore.selectedPlace) {
                        obj.open = true
                    }
                });
            }
            return this.open
        }
    },
    methods: {
        toggle() {
            if (this.isFolder) {
                this.open = !this.open;
            }

            this.placeStore.setSelectedPlace(this.item.place.id)
        },
    },
}
</script>

<template>
    <li>
        <div
            :class="{selected: item.place.id === placeStore.selectedPlace}"
            @click="toggle">
            <span v-if="isFolder">[{{ isOpen ? '-' : '+' }}]</span>
            {{ item.place.title }}
        </div>
        <ul v-show="isOpen" v-if="isFolder">
            <PlaceTreeItem
                v-for="(nested, index) in item.nested"
                :key="index"
                :item="nested">
            </PlaceTreeItem>
        </ul>
    </li>
</template>
