import * as client from "../client/client.js";

export const mainPageComponent = {
    props: {
        isAuth: Boolean,
    },
    data() {
        return {
            tree: null,
        }
    },
    computed: {
        showMainPage() {
            return this.isAuth
        }
    },
    created() {
        if (this.isAuth) {
            this.refreshTree()
        }
    },
    watch: {
        isAuth(val) {
            if (val) {
                this.refreshTree()
            }
        }
    },
    methods: {
        refreshTree() {
            let res = this.request(client.methodGet, client.routeGetPlacesTree)
            console.log(res.data) // TODO
        },
        request(method, route) {
            let res = client.request(method, route)
            if (res.status !== client.statusOK) {
                this.$emit('eventsetauth', false)
            }

            return res
        }
    },
    template: `
    <template v-if="showMainPage">
        <main class="container-fluid">
            <div class="d-flex flex-grow h-100">
      
                <div class="col-l">
                    <div class="places rounded-3 shadow d-flex flex-column">
                        <div class="header rounded-top">
                            Места
                            <div class="buttons float-end">
                                <button class="btn add" title="Добавить место">
                                    <i class="bi bi-plus-circle"></i>
                                </button>
                                <button class="btn edit"></button>
                                <button class="btn delete"></button>
                            </div>
                        </div>
                        <div class="list"></div>
                    </div>
                </div>
        
                <div class="col-c">
                    <div class="things rounded-3 shadow d-flex flex-column">
                        <div class="header rounded-top">
                            Вещи
                            <div class="buttons float-end">
                                <button class="btn add"></button>
                                <button class="btn edit"></button>
                                <button class="btn delete"></button>
                            </div>
                        </div>
                        <div class="list"></div>
                    </div>
                </div>
        
                <div class="col-r">
                    <div class="info rounded-3 shadow d-flex flex-column">
                    <div class="header rounded-top">
                        Фото
                        <div class="buttons float-end">
                            <button class="btn add"></button>
                        </div>
                    </div>
                    <div class="list"></div>
                </div>
            
            </div>
        </main>
    </template>
    `
}
