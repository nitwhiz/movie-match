import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { useCurrentUser } from './composables/useCurrentUser';
import { freeAllMediaBlobUrls } from './api/PosterBlob';

import HomeView from './views/HomeView.vue';
import VoteView from './views/VoteView.vue';
import MatchesView from './views/MatchesView.vue';
import MediaView from './views/MediaView.vue';
import LoginView from './views/LoginView.vue';
import SearchView from './views/SearchView.vue';

export const RouteName = {
  LOGIN: 'login',
  HOME: 'home',
  VOTE: 'vote',
  MATCHES: 'matches',
  MEDIA: 'media',
  SEARCH: 'search',
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
    path: '/search',
    name: RouteName.SEARCH,
    component: SearchView,
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

  if (to.name === RouteName.LOGIN && currentUser.value) {
    next({ name: RouteName.HOME });
    return;
  }

  if (to.name !== RouteName.LOGIN && !currentUser.value) {
    next({ name: RouteName.LOGIN });
    return;
  }

  freeAllMediaBlobUrls();

  next();
});

export default router;
