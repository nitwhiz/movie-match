<template>
  <div class="votes">
    <h2>Your Votes</h2>
    <div class="vote" v-for="m in voteList">
      <MediaItem :media="m.media" />
    </div>
  </div>
</template>

<script setup lang="ts">
import MediaItem from '../components/media/MediaItem.vue';
import { useRouter } from 'vue-router';
import { useApiClient } from '../composables/useApiClient';
import { onMounted, ref, watch } from 'vue';
import { Media, MediaType } from '../model/Media';
import { Vote } from '../model/Vote';

const router = useRouter();
const apiClient = await useApiClient().apiClient;

const filterType = ref('all' as MediaType | 'all');

const voteList = ref([] as { vote: Vote; media: Media }[]);

const fetchMatches = () => {
  apiClient.getVotes().then(async (votes) => {
    voteList.value = [];

    if (votes) {
      for (const vote of votes) {
        // todo: request pooling?
        const media = await apiClient.getMedia(vote.mediaId);

        voteList.value.push({ vote, media });
      }
    }
  });
};

watch(filterType, () => fetchMatches());

onMounted(() => {
  fetchMatches();
});
</script>

<style scoped lang="scss">
.votes {
  width: 100%;
  height: 100%;

  overflow-x: hidden;
  overflow-y: auto;

  padding: 12px;

  .vote {
    margin-bottom: 12px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}
</style>
