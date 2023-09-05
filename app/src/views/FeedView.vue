<script setup lang="ts">
import MediaFeed from '../components/media/MediaFeed.vue';
import { useApiClient } from '../composables/useApiClient';
import { onMounted, ref } from 'vue';
import { RecommendedMedia } from '../model/Media';
import { VoteType } from '../model/Vote';

const apiClient = await useApiClient().apiClient;

const isFetchingMedia = ref(false);
const mediaList = ref([] as RecommendedMedia[]);
const belowScore = ref('2');

const fetchMedia = () => {
  if (isFetchingMedia.value) {
    return;
  }

  isFetchingMedia.value = true;

  apiClient.getRecommendedMedia(belowScore.value).then((results) => {
    mediaList.value = [...results];

    const lastResult = results.pop();

    belowScore.value = lastResult?.score || '100';

    isFetchingMedia.value = false;
  });
};

onMounted(() => {
  fetchMedia();
});

const handleVote = (media: RecommendedMedia, voteType: VoteType) => {
  // todo: use a store

  if (media.voteType === VoteType.NONE || media.voteType === VoteType.NEUTRAL) {
    media.voteType = voteType;
    apiClient.voteMedia(media.id, voteType);
  } else {
    media.voteType = VoteType.NEUTRAL;
    apiClient.voteMedia(media.id, VoteType.NEUTRAL);
  }
};
</script>

<template>
  <MediaFeed :media="mediaList" @endofmedia="fetchMedia" @vote="handleVote" />
</template>

<style scoped lang="scss"></style>
