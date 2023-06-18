import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { useCurrentUser } from './composables/useCurrentUser';
import LoginView from './views/LoginView.vue';
import HomeView from './views/HomeView.vue';
import FeedView from './views/FeedView.vue';
import MatchesView from './views/MatchesView.vue';
import MediaView from './views/MediaView.vue';
import VotesView from './views/VotesView.vue';
import Poster from './common/Poster';
import SearchView from './views/SearchView.vue';

export const RouteName = {
  ROOT: 'root',
  LOGIN: 'login',
  HOME: 'home',
  FEED: 'feed',
  MATCHES: 'matches',
  MEDIA: 'media',
  SEARCH: 'search',
  VOTES: 'votes',
} as const;

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: RouteName.ROOT,
    redirect: RouteName.LOGIN,
  },
  {
    path: '/login',
    name: RouteName.LOGIN,
    component: LoginView,
  },
  {
    path: '/home',
    name: RouteName.HOME,
    component: HomeView,
  },
  {
    path: '/feed',
    name: RouteName.FEED,
    component: FeedView,
  },
  {
    path: '/matches',
    name: RouteName.MATCHES,
    component: MatchesView,
  },
  {
    path: '/votes',
    name: RouteName.VOTES,
    component: VotesView,
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

  Poster.freeAll();

  next();
});

export default router;
