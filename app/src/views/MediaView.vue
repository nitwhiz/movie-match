<template>
  <div class="media-view">
    <VoteHost
      :swipe="false"
      v-if="media"
      :media="media"
      v-model:meta-visible="metaVisible"
    />
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { onMounted, ref } from 'vue';
import { Media } from '../model/Media';
import { useApiClient } from '../composables/useApiClient';
import VoteHost from '../components/voting/VoteHost.vue';

const route = useRoute();
const apiClient = await useApiClient().apiClient;

const media = ref(null as Media | null);
const metaVisible = ref(true);

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
