<template>
  <template v-if="media">
    <MediaCard :media="media" @vote="handleVote" />
  </template>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { onMounted, ref } from 'vue';
import { Media, RecommendedMedia } from '../model/Media';
import { useApiClient } from '../composables/useApiClient';
import MediaCard from '../components/media/MediaCard.vue';
import { VoteType } from '../model/Vote';

const route = useRoute();
const apiClient = await useApiClient().apiClient;

const media = ref(null as RecommendedMedia | null);

const handleVote = (voteType: VoteType) => {
  const currentMedia = media.value;

  if (currentMedia) {
    // todo: use a store

    if (
      currentMedia.voteType === VoteType.NONE ||
      currentMedia.voteType === VoteType.NEUTRAL
    ) {
      currentMedia.voteType = voteType;
      apiClient.voteMedia(currentMedia.id, voteType);
    } else {
      currentMedia.voteType = VoteType.NEUTRAL;
      apiClient.voteMedia(currentMedia.id, VoteType.NEUTRAL);
    }
  }
};

onMounted(() => {
  const mediaId = (route.params.mediaId || null) as string | null;

  if (mediaId) {
    Promise.all([
      apiClient.getMedia(mediaId),
      apiClient.getMediaVote(mediaId),
    ]).then(([mediaData, voteType]) => {
      media.value = {
        ...(mediaData as Media),
        voteType: voteType as VoteType,
        score: '1',
        seen: false,
      };
    });
  }
});
</script>

<style lang="scss" scoped>
.media-view {
  width: 100%;
  height: 100%;

  overflow: hidden;
}
</style>
