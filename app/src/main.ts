import './index.scss';
import App from './App.vue';
import Vue3Lottie from 'vue3-lottie';
import router from './router';
import { createApp } from 'vue';
import 'vue3-lottie/dist/style.css';

window.addEventListener('contextmenu', (event: MouseEvent) => {
  event.preventDefault();
  event.stopPropagation();
});

createApp(App).use(Vue3Lottie, { name: 'Lottie' }).use(router).mount('#app');
