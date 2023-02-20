import './index.scss';
import App from './App.vue';
import router from './router';
import { createApp } from 'vue';

window.addEventListener('contextmenu', (event: MouseEvent) => {
  event.preventDefault();
  event.stopPropagation();
});

createApp(App).use(router).mount('#app');
