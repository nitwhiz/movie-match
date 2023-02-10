import './index.scss';
import App from './App.vue';
import router from './router';
import { createApp } from 'vue';
import { createPinia } from 'pinia';

const pinia = createPinia();

window.addEventListener('contextmenu', (event: MouseEvent) => {
  event.preventDefault();
  event.stopPropagation();
});

createApp(App).use(pinia).use(router).mount('#app');
