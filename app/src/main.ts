import './assets/styles/index.scss';
import App from './App.vue';
import router from './router';
import { createApp } from 'vue';
import Poster from './common/Poster';

Poster.startGC();

window.addEventListener('contextmenu', (event: MouseEvent) => {
  event.preventDefault();
  event.stopPropagation();
});

createApp(App).use(router).mount('#app');
