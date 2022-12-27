'use strict'

export const placeTreeItemComponent = {
    name: 'place-tree-item',
    props: {
        item: Object,
        selectedPlace: Number,
    },
    data: function () {
        return {
            open: false
        };
    },
    computed: {
        isFolder() {
            return this.item.nested && this.item.nested.length;
        },
        isOpen() {
            if (this.item.nested && this.item.nested.length) {
                let obj = this
                this.item.nested.forEach(function(item) {
                    if (item.place.id === obj.selectedPlace) {
                        obj.open = true
                        return
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
            this.$emit("set-selected-place", this.item.place.id);
        },
    },
    template: `
    <li>
        <div
            data-id="{{ item.place.id }}"
            :class="{selected: item.place.id == selectedPlace}"
            @click="toggle">
            <span v-if="isFolder">[{{ isOpen ? '-' : '+' }}]</span>
            {{ item.place.title }}
        </div>
        <ul v-show="isOpen" v-if="isFolder">
            <place-tree-item
                v-for="(nested, index) in item.nested"
                :key="index"
                :item="nested"
                :selected-place="selectedPlace"
                @set-selected-place="$emit('set-selected-place', $event)">
            </place-tree-item>
        </ul>
    </li>
`
}

