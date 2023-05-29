<template>
  <div class="home">
    <h1>
      Hey<span v-if="userDisplayName">&nbsp{{ userDisplayName }}</span
      >!
    </h1>
    <router-link :to="{ name: RouteName.VOTE }" class="vote">
      <span class="text">Get Voting!</span>
    </router-link>
    <router-link :to="{ name: RouteName.MATCHES }" class="matches">
      <span class="text">Check Matches</span>
    </router-link>
    <router-link :to="{ name: RouteName.YOUR_VOTES }" class="your-votes">
      <span class="text">Your Votes</span>
    </router-link>
    <router-link :to="{ name: RouteName.SEARCH }" class="search">
      <span class="text">Search Titles</span>
    </router-link>
    <a class="logout" @click="logout">
      <span class="text">Log-Out</span>
    </a>
  </div>
</template>

<script lang="ts" setup>
import { useApiClient } from '../composables/useApiClient';
import { computed } from 'vue';
import { useCurrentUser } from '../composables/useCurrentUser';
import { RouteName } from '../router';

const apiClient = await useApiClient().apiClient;
const { currentUser } = useCurrentUser();

const userDisplayName = computed(() => currentUser.value?.displayName || null);

const logout = async () => {
  await apiClient.logout();
};
</script>

<style lang="scss" scoped>
@use '../styles/nice';

.home {
  display: flex;
  align-items: center;
  flex-direction: column;

  width: 100%;
  height: 100%;

  overflow: hidden;

  h1 {
    margin: 4rem 0 4rem 0;
  }

  a {
    text-decoration: none;
    color: white;

    font-size: 2rem;
    font-family: Pacifico, sans-serif;
    text-transform: capitalize;

    padding: 0.5rem 1rem;

    margin-bottom: 1.25rem;
    width: 75%;
    min-width: 240px;

    text-align: center;

    &:last-child {
      margin-bottom: 0;
    }

    $border-width: 3px;

    &.vote {
      @include nice.gradient-border(
        linear-gradient(20deg, rgb(148, 55, 255) 0%, rgb(101, 229, 255) 100%),
        $border-width
      );
    }

    &.matches {
      @include nice.gradient-border(
        linear-gradient(30deg, rgb(255, 55, 140) 0%, rgb(187, 255, 101) 100%),
        $border-width
      );
    }

    &.search {
      @include nice.gradient-border(
        linear-gradient(30deg, rgb(78, 169, 32) 0%, rgb(31, 227, 174) 100%),
        $border-width
      );
    }

    &.your-votes {
      @include nice.gradient-border(
        linear-gradient(30deg, rgb(85, 212, 206) 0%, rgb(199, 85, 212) 100%),
        $border-width
      );
    }

    &.logout {
      @include nice.gradient-border(
        linear-gradient(30deg, rgb(195, 58, 34) 0%, rgb(253, 135, 45) 100%),
        $border-width
      );
    }
  }
}
</style>
