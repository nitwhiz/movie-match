<template>
  <div class="home">
    <h1>
      Hey<span v-if="userDisplayName">&nbsp;{{ userDisplayName }}</span
      >!
    </h1>
    <NiceWrapper :colors="['rgb(148, 55, 255)', 'rgb(101, 229, 255)']">
      <router-link :to="{ name: RouteName.FEED }">Get Voting!</router-link>
    </NiceWrapper>
    <NiceWrapper :colors="['rgb(255, 55, 140)', 'rgb(187, 255, 101)']">
      <router-link :to="{ name: RouteName.MATCHES }">Check Matches</router-link>
    </NiceWrapper>
    <NiceWrapper :colors="['rgb(78, 169, 32)', 'rgb(31, 227, 174)']">
      <router-link :to="{ name: RouteName.VOTES }">Your Votes</router-link>
    </NiceWrapper>
    <NiceWrapper :colors="['rgb(85, 212, 206)', 'rgb(199, 85, 212)']">
      <router-link :to="{ name: RouteName.SEARCH }">Search Titles</router-link>
    </NiceWrapper>
    <NiceWrapper
      :colors="['rgb(195, 58, 34)', 'rgb(253, 135, 45)']"
      @click="logout"
    >
      Log-Out
    </NiceWrapper>
  </div>
</template>

<script lang="ts" setup>
import { useApiClient } from '../composables/useApiClient';
import { computed } from 'vue';
import { useCurrentUser } from '../composables/useCurrentUser';
import { RouteName } from '../router';
import NiceWrapper from '../components/nice/NiceWrapper.vue';

const apiClient = await useApiClient().apiClient;
const { currentUser } = useCurrentUser();

const userDisplayName = computed(() => currentUser.value?.displayName || null);

const logout = async () => {
  await apiClient.logout();
};
</script>

<style lang="scss" scoped>
.home {
  display: flex;
  align-items: center;
  flex-direction: column;

  padding: 0 2rem;

  width: 100%;
  height: 100%;

  overflow: hidden;

  h1 {
    margin: 4rem 0 4rem 0;
  }

  .nice-wrapper {
    margin-bottom: 1rem;

    &:last-child {
      margin-bottom: 0;
    }
  }
}
</style>
