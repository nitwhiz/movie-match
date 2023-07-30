<template>
  <slot />
</template>

<script lang="ts" setup>
// wrapper component to ensure api client is initialized with event listeners

import { useApiClient } from './composables/useApiClient';
import { useRouter } from 'vue-router';
import { useCurrentUser } from './composables/useCurrentUser';
import { RouteName } from './router';

const router = useRouter();
const apiClient = await useApiClient().apiClient;

apiClient
  .on('unauthorized', () => {
    if (router.currentRoute.value.name !== 'login') {
      router.push({ name: RouteName.LOGIN });
    }
  })
  .on('logout', async () => {
    const { currentUser } = useCurrentUser();

    currentUser.value = null;

    await router.push({ name: RouteName.LOGIN });
  });
</script>
