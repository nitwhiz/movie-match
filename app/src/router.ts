import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { useCurrentUser } from './composables/useCurrentUser';

import HomeView from './views/HomeView.vue';
import VoteView from './views/VoteView.vue';
import MatchesView from './views/MatchesView.vue';
import MediaView from './views/MediaView.vue';
import LoginView from './views/LoginView.vue';

export const RouteName = {
  LOGIN: 'login',
  HOME: 'home',
  VOTE: 'vote',
  MATCHES: 'matches',
  MEDIA: 'media',
} as const;

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: RouteName.LOGIN,
    component: LoginView,
  },
  {
    path: '/home',
    name: RouteName.HOME,
    component: HomeView,
  },
  {
    path: '/vote',
    name: RouteName.VOTE,
    component: VoteView,
  },
  {
    path: '/matches',
    name: RouteName.MATCHES,
    component: MatchesView,
  },
  {
    path: '/media/:mediaId',
    name: RouteName.MEDIA,
    component: MediaView,
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const { currentUser, load } = useCurrentUser();

  await load();

  if (to.name === 'login' && currentUser.value) {
    next({ name: RouteName.HOME });
    return;
  }

  if (to.name !== 'login' && !currentUser.value) {
    next({ name: RouteName.LOGIN });
    return;
  }

  next();
});

export default router;
