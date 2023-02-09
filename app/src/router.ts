import HomeView from './views/HomeView.vue';
import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import UserSelectionView from './views/UserSelectionView.vue';
import VoteView from './views/VoteView.vue';
import MatchesView from './views/MatchesView.vue';
import { useUserStore } from './store/userStore';
import MediaView from './views/MediaView.vue';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'main',
    component: UserSelectionView,
  },
  {
    path: '/home',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/vote',
    name: 'vote',
    component: VoteView,
  },
  {
    path: '/matches',
    name: 'matches',
    component: MatchesView,
  },
  {
    path: '/media/:mediaId',
    name: 'media',
    component: MediaView,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  if (to.name !== 'main') {
    const loggedIn = await useUserStore().loadCurrentUser();

    if (!loggedIn) {
      next({ name: 'main' });
      return;
    }
  }

  next();
});

export default router;
