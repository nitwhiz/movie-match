<template>
  <div class="matches">
    <h2>Your Matches</h2>
    <div class="filter">
      <NiceWrapper
        class="type"
        :colors="['rgb(185, 81, 126)', 'rgb(95, 148, 210)']"
      >
        <select v-model="filterType">
          <option value="all">All</option>
          <option value="tv">TV-Shows</option>
          <option value="movie">Movies</option>
        </select>
      </NiceWrapper>
    </div>
    <div class="match" v-for="m in matchList">
      <MediaItem :media="m.media" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch } from 'vue';
import { Match } from '../model/Match';
import { Media, MediaType } from '../model/Media';
import { useRouter } from 'vue-router';
import { useApiClient } from '../composables/useApiClient';
import MediaItem from '../components/media/MediaItem.vue';
import NiceWrapper from '../components/nice/NiceWrapper.vue';

const router = useRouter();
const apiClient = await useApiClient().apiClient;

const filterType = ref(MediaType.ALL);

const fetchMatches = () => {
  apiClient
    .getMatches(filterType.value !== MediaType.ALL ? filterType.value : null)
    .then(async (matches) => {
      matchList.value = [];

      if (matches) {
        for (const match of matches) {
          // todo: request pooling?
          const media = await apiClient.getMedia(match.mediaId);

          matchList.value.push({ match, media });
        }
      }
    });
};

watch(filterType, () => fetchMatches());

const matchList = ref([] as { match: Match; media: Media }[]);

onMounted(() => {
  fetchMatches();
});
</script>

<style lang="scss" scoped>
.matches {
  width: 100%;
  height: 100%;

  overflow-x: hidden;
  overflow-y: auto;

  padding: 12px;

  .filter {
    margin-bottom: 20px;

    .type {
      padding: 0.25rem;

      select {
        padding: 0.5rem;

        width: 100%;

        option {
          background-color: black;
        }
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
