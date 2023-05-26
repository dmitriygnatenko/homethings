<script>
import {useImageStore} from '../../stores/image.js'
import {Modal} from 'bootstrap'

export default {
    setup() {
        const imageStore = useImageStore()
        return {imageStore}
    },
    expose: ['init'],
    data() {
        return {
            modal: Object,
            activeImageID: 0,
            activeImagePlaceID: 0,
            activeImageThingID: 0,
        }
    },
    methods: {
        init(imageID, placeID, thingID) {
            this.activeImageID = imageID
            this.activeImagePlaceID = placeID
            this.activeImageThingID = thingID

            let modal = document.getElementById('modal-show-images')
            this.modal = new Modal(modal, {})
            this.modal.show()
        },
    },
}
</script>

<template>
    <div class="modal" tabindex="-1" id="modal-show-images">
        <div class="modal-dialog modal-lg">
            <div class="modal-content">
                <div class="modal-body">
                    <div id="imagesCarousel" class="carousel slide">
                        <div class="carousel-inner">
                            <div
                                class="carousel-item"
                                v-for="image of imageStore.imageList"
                                :class="{ active : this.activeImageID === image.id &&
                                    this.activeImagePlaceID === image.place_id &&
                                    this.activeImageThingID === image.thing_id }">
                                <img :src="image.image" class="d-block w-100">
                            </div>
                        </div>
                        <button class="carousel-control-prev" type="button" data-bs-target="#imagesCarousel"
                                data-bs-slide="prev">
                            <span class="carousel-control-prev-icon"></span>
                        </button>
                        <button class="carousel-control-next" type="button" data-bs-target="#imagesCarousel"
                                data-bs-slide="next">
                            <span class="carousel-control-next-icon" aria-hidden="true"></span>
                        </button>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary btn-sm" data-bs-dismiss="modal">Закрыть</button>
                </div>
            </div>
        </div>
    </div>
</template>
