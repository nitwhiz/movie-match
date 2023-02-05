import './index.scss';
import App from './App.vue';
import router from './router';
import { createApp } from 'vue';
import { createPinia } from 'pinia';

const pinia = createPinia();

createApp(App).use(pinia).use(router).mount('#app');
