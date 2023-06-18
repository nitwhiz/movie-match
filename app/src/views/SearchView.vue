<template>
  <div class="search">
    <h2>Search Titles</h2>
    <div class="filter">
      <NiceWrapper :colors="['rgb(185, 81, 126)', 'rgb(95, 148, 210)']">
        <input
          type="text"
          v-model="searchQuery"
          @input="debouncedStartSearch"
          placeholder="Enter Title"
        />
      </NiceWrapper>
    </div>
    <div class="match" v-for="m in results">
      <MediaItem :media="m.media" :vote-type="m.voteType" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Media } from '../model/Media';
import { onMounted, ref } from 'vue';
import { useApiClient } from '../composables/useApiClient';
import MediaItem from '../components/media/MediaItem.vue';
import { useSearchQuery } from '../composables/useSearchQuery';
import NiceWrapper from '../components/nice/NiceWrapper.vue';
import { VoteType } from '../model/Vote';

const apiClient = await useApiClient().apiClient;

const results = ref(
  [] as {
    media: Media;
    voteType: VoteType;
  }[]
);

const { searchQuery } = useSearchQuery();

const searchDebounceTimeout = ref(-1);

const startSearch = () => {
  if (searchQuery.value == '') {
    return;
  }

  apiClient
    .searchMedia(searchQuery.value)
    .then(async (searchResults) => {
      const voteResults = [];

      for (const result of searchResults) {
        voteResults.push(await apiClient.getMediaVote(result.id));
      }

      return [searchResults, voteResults];
    })
    .then(([searchResults, voteResults]) => {
      results.value = (searchResults as Media[]).map((m, i) => ({
        media: m,
        voteType: (voteResults[i] as VoteType) || VoteType.NONE,
      }));
    });
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
.search {
  width: 100%;
  height: 100%;

  overflow-x: hidden;
  overflow-y: auto;

  padding: 12px;

  .filter {
    margin-bottom: 20px;
  }

  .match {
    margin-bottom: 12px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}
</style>
