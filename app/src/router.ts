import HomeView from './views/HomeView.vue';
import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import UserSelectionView from './views/UserSelectionView.vue';
import VoteView from './views/VoteView.vue';
import MatchesView from './views/MatchesView.vue';
import { useUserStore } from './store/userStore';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'main',
    component: () => UserSelectionView,
  },
  {
    path: '/home',
    name: 'home',
    component: () => HomeView,
  },
  {
    path: '/vote/movies',
    name: 'vote-movies',
    component: () => VoteView,
  },
  {
    path: '/matches',
    name: 'matches',
    component: () => MatchesView,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  if (to.name !== 'main') {
    const userStore = useUserStore();

    if (userStore.currentUser === null) {
      next({ name: 'main' });
      return;
    }
  }

  next();
});

export default router;
