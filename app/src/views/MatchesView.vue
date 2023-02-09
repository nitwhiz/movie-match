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
    <div class="match" v-for="m in matchList" @click="showMedia(m.media)">
      <div class="poster">
        <img :src="getMediaPosterSrc(m.media)" :alt="m.media.title" />
      </div>
      <div class="details">
        <b class="name">{{ m.media.title }}</b>
        <span class="type">{{ getMediaTypeLabelSingular(m.media.type) }}</span>
        <span class="genres">{{ getGenres(m.media) }}</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref, watch } from 'vue';
import { Match } from '../model/Match';
import { useUserStore } from '../store/userStore';
import { Media, MediaType } from '../model/Media';
import { useRouter } from 'vue-router';
import { useApiClient } from '../composables/useApiClient';
import { useMediaType } from '../composables/useMediaType';

const userStore = useUserStore();
const router = useRouter();
const apiClient = useApiClient();

const { getMediaTypeLabelSingular } = useMediaType();

const filterType = ref('all' as MediaType | 'all');

const fetchMatches = () => {
  apiClient
    .getMatches(
      userStore.currentUser?.id || '',
      filterType.value !== 'all' ? filterType.value : null
    )
    .then(async (matches) => {
      matchList.value = [];

      if (matches) {
        for (const match of matches) {
          const media = await apiClient.getMedia(match.mediaId);

          matchList.value.push({ match, media });
        }
      }
    });
};

watch(filterType, () => fetchMatches());

const matchList = ref([] as { match: Match; media: Media }[]);

const currentMedia = ref(null as Media | null);

const showMedia = (media: Media) => {
  currentMedia.value = media;
  router.push({
    name: 'media',
    params: {
      mediaId: media.id,
    },
  });
};

const getMediaPosterSrc = (media: Media) => apiClient.getPosterUrl(media.id);

const getGenres = (media: Media) => media.genres.map((g) => g.name).join(', ');

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

        background: transparent;
        border: 0;
        outline: 0;

        color: white;
      }
    }
  }

  .match {
    width: 100%;
    display: flex;

    margin-bottom: 12px;

    &:last-child {
      margin-bottom: 0;
    }

    $border-width: 3px;

    @include nice.gradient-border(
      linear-gradient(20deg, rgb(34, 193, 195) 0%, rgb(253, 187, 45) 100%),
      $border-width
    );

    .poster {
      width: 20%;

      border-top-left-radius: 11px;
      border-bottom-left-radius: 11px;

      overflow: hidden;

      flex-grow: 0;

      img {
        display: block;

        width: auto;
        height: auto;

        max-width: 100%;
      }
    }

    .details {
      display: flex;
      flex-direction: column;
      justify-content: center;

      flex-grow: 0;
      width: 80%;

      padding: 12px 20px;

      .type,
      .genres {
        font-size: 0.8rem;
        margin-top: 0.2rem;
      }
    }
  }
}
</style>
