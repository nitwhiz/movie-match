<template>
  <div class="media-view">
    <MediaCard v-if="media" :media="media" />
  </div>
</template>

<script setup lang="ts">
import MediaCard from '../components/MediaCard.vue';
import { useRoute } from 'vue-router';
import { onMounted, ref } from 'vue';
import { Media } from '../model/Media';
import { useApiClient } from '../composables/useApiClient';

const route = useRoute();
const apiClient = await useApiClient().apiClient;

const media = ref(null as Media | null);

onMounted(() => {
  const mediaId = (route.params.mediaId || null) as string | null;

  if (mediaId) {
    apiClient.getMedia(mediaId).then((mediaData) => {
      media.value = mediaData;
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
