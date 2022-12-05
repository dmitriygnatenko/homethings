export const mainPageComponent = {
    props: {
        show: Boolean,
    },
    template: `
    <template v-if="show">
        Main page
    </template>
    `
}
