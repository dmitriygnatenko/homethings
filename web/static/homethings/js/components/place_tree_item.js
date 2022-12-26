export const placeTreeItemComponent = {
    name: 'place-tree-item',
    props: {
        item: Object,
        selectedPlace: Number,
    },
    data: function () {
        return {
            isOpen: false
        };
    },
    computed: {
        isFolder: function () {
            return this.item.nested && this.item.nested.length;
        }
    },
    methods: {
        toggle: function () {
            if (this.isFolder) {
                this.isOpen = !this.isOpen;
            }
            this.$emit("set-selected-place", this.item.place.id);
        },
    },
    template: `
    <li>
        <div
            data-id="{{ item.place.id }}"
            :class="{selected: item.place.id == selectedPlace}"
            @click="toggle"
        >
            <span v-if="isFolder">[{{ isOpen ? '-' : '+' }}]</span>
            {{ item.place.title }}
        </div>
        <ul v-show="isOpen" v-if="isFolder">
            <place-tree-item
                v-for="(nested, index) in item.nested"
                :key="index"
                :item="nested"
                :selected-place="selectedPlace"
                @set-selected-place="$emit('set-selected-place', $event)"
            ></place-tree-item>
        </ul>
    </li>
`
}

