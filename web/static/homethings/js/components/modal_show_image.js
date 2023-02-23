"use strict"

export const modalShowImagesComponent = {
    props: ["images"],
    data() {
        return {
            modal: Object,
            modalImages: [],
        }
    },
    methods: {
        init(imageID, placeID, thingID) {
            this.modalImages = []

            let obj = this
            this.images.forEach(image => {
                let isActive = (image.id === imageID && image.place_id === placeID && image.thing_id === thingID)

                obj.modalImages.push({
                    "image": image.image,
                    "active": isActive,
                })
            });

            let modal = document.getElementById('show-images-modal')
            modal.addEventListener('hide.bs.modal', event => {
                this.modalImages = []
            })
            this.modal = new bootstrap.Modal(modal, {})
            this.modal.show()
        },
    },
    template: `
    <div class="modal" tabindex="-1" id="show-images-modal">
        <div class="modal-dialog modal-lg">
            <div class="modal-content">
                <div class="modal-body">
                    <div id="imagesCarousel" class="carousel slide">
                        <div class="carousel-inner">
                            <div
                                class="carousel-item"
                                v-for="image of modalImages"
                                :class="{ active : image.active }">
                                <img :src="image.image" class="d-block w-100">
                            </div>
                        </div>
                        <button class="carousel-control-prev" type="button" data-bs-target="#imagesCarousel" data-bs-slide="prev">
                            <span class="carousel-control-prev-icon"></span>
                        </button>
                        <button class="carousel-control-next" type="button" data-bs-target="#imagesCarousel" data-bs-slide="next">
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
    `
}
