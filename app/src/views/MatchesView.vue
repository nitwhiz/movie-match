<template>
  <div class="matches">
    <h2>Your Matches</h2>
    <div class="filter">
      <div class="type">
        <select v-model="filterType">
          <option value="all">All</option>
          <option value="tv">TV-Shows</option>
          <option value="movie">Movies</option>
        </select>
      </div>
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
import { useMediaType } from '../composables/useMediaType';
import { useCurrentUser } from '../composables/useCurrentUser';
import MediaItem from '../components/media/MediaItem.vue';

const router = useRouter();
const { currentUser } = useCurrentUser();
const apiClient = await useApiClient().apiClient;

const { getMediaTypeLabelSingular } = useMediaType();

const filterType = ref('all' as MediaType | 'all');

const fetchMatches = () => {
  apiClient
    .getMatches(
      currentUser.value?.id || '',
      filterType.value !== 'all' ? filterType.value : null
    )
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
@use '../styles/nice';

.matches {
  width: 100%;
  height: 100%;

  overflow-x: hidden;
  overflow-y: auto;

  padding: 12px;

  .filter {
    margin-bottom: 20px;

    .type {
      @include nice.gradient-border(
        linear-gradient(20deg, rgb(185, 81, 126) 0%, rgb(95, 148, 210) 100%),
        3px
      );

      padding: 0.25rem;

      select {
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
