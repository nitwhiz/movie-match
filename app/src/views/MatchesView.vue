<template>
  <div class="matches">
    <h2>Your Matches</h2>
    <div class="match" v-for="m in matchList" @click="showMedia(m.media)">
      <div class="poster">
        <img :src="getMediaPosterSrc(m.media)" :alt="m.media.title" />
      </div>
      <div class="details">
        <b class="name">{{ m.media.title }}</b>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { Match } from '../model/Match';
import { useUserStore } from '../store/userStore';
import { Media } from '../model/Media';
import { useRouter } from 'vue-router';
import { useApiClient } from '../composables/useApiClient';

const userStore = useUserStore();
const router = useRouter();
const apiClient = useApiClient();

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

onMounted(() => {
  apiClient
    .getMatches(userStore.currentUser?.id || '')
    .then(async (matches) => {
      matchList.value = [];

      if (matches) {
        for (const match of matches) {
          const media = await apiClient.getMedia(match.mediaId);

          matchList.value.push({ match, media });
        }
      }
    });
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

  .match {
    width: 100%;
    display: flex;

    margin-bottom: 12px;

    &:last-child {
      margin-bottom: 0;
    }

    $border-width: 3px;

    @include nice.gradient-button(
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
      align-items: center;

      flex-grow: 0;
      width: 80%;

      padding: 12px 20px;
    }
  }
}
</style>
