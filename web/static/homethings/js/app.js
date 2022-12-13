import {loginPageComponent} from './components/login_page.js';
import {mainPageComponent} from './components/main_page.js';

export const app = {
    components: {
        'login-page': loginPageComponent,
        'main-page': mainPageComponent,
    },
    data() {
        return {
            isAuth: false,
        };
    },
};
