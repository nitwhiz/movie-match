<template>
  <div class="search">
    <h2>Search Titles</h2>
    <div class="filter">
      <div class="search-bar">
        <input
          type="text"
          v-model="searchQuery"
          @input="debouncedStartSearch"
          placeholder="Enter Title"
        />
      </div>
    </div>
    <div class="match" v-for="m in results">
      <MediaItem :media="m" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Media } from '../model/Media';
import { onMounted, ref } from 'vue';
import { useApiClient } from '../composables/useApiClient';
import MediaItem from '../components/media/MediaItem.vue';
import { useSearchQuery } from '../composables/useSearchQuery';

const { apiClient } = useApiClient();

const results = ref([] as Media[]);

const { searchQuery } = useSearchQuery();

const searchDebounceTimeout = ref(-1);

const startSearch = () => {
  if (searchQuery.value == '') {
    return;
  }

  apiClient
    .then((c) => c.searchMedia(searchQuery.value))
    .then((searchResults) => (results.value = searchResults));
};

const debouncedStartSearch = () => {
  if (searchDebounceTimeout.value > -1) {
    window.clearTimeout(searchDebounceTimeout.value);
    searchDebounceTimeout.value = -1;
  }

  searchDebounceTimeout.value = window.setTimeout(() => {
    startSearch();
    searchDebounceTimeout.value = -1;
  }, 500);
};

onMounted(() => {
  startSearch();
});
</script>

<style lang="scss" scoped>
@use '../styles/nice';

.search {
  width: 100%;
  height: 100%;

  overflow-x: hidden;
  overflow-y: auto;

  padding: 12px;

  .filter {
    margin-bottom: 20px;

    .search-bar {
      @include nice.gradient-border(
        linear-gradient(20deg, rgb(185, 81, 126) 0%, rgb(95, 148, 210) 100%),
        3px
      );

      padding: 0.25rem;

      input {
        padding: 0.5rem;

        width: 100%;
      }
    }
  }

  .match {
    margin-bottom: 12px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}
</style>
