<template>
  <div class="vote-host">
    <div
      class="match-notification"
      :class="[matchNotificationVisible ? 'visible' : null]"
    >
      It's a match!
    </div>
    <MediaSwipeHost
      :current-media="currentMedia"
      :next-media="nextMedia"
      v-model:vote-type="currentVoteType"
      @vote="handleVote"
    />
  </div>
  <div class="button-wrapper" v-if="currentMedia">
    <ButtonHost @vote="handleButtonVote" @seen="handleSeen" />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue';
import { Media } from '../../model/Media';
import { useUserStore } from '../../store/userStore';
import { useApiClient } from '../../composables/useApiClient';
import { VoteType } from '../../model/Vote';
import MediaSwipeHost from './MediaSwipeHost.vue';
import ButtonHost from './ButtonHost.vue';

const MEDIA_COUNT_THRESHOLD = 3;

const userStore = useUserStore();
const apiClient = useApiClient();

const matchNotificationVisible = ref(false);

const currentVoteType = ref(VoteType.NEUTRAL);

const isFetchingMedia = ref(false);
const mediaList = ref([] as Media[]);
const mediaPageIndex = ref(0);
const mediaIndex = ref(0);

const currentMedia = computed(() => mediaList.value[mediaIndex.value] || null);
const nextMedia = computed(() => mediaList.value[mediaIndex.value + 1] || null);

const showNextMedia = () => {
  ++mediaIndex.value;

  nextTick(() => (currentVoteType.value = VoteType.NEUTRAL));

  if (mediaIndex.value > mediaList.value.length - MEDIA_COUNT_THRESHOLD) {
    fetchMedia();
  }
};

const fetchMedia = () => {
  if (isFetchingMedia.value) {
    return;
  }

  isFetchingMedia.value = true;

  apiClient
    .getRecommendedMedia(userStore.currentUser?.id || '', mediaPageIndex.value)
    .then((results) => {
      mediaList.value.push(...results);
      ++mediaPageIndex.value;

      isFetchingMedia.value = false;
    });
};

const sendVote = (media: Media, voteType: VoteType) => {
  apiClient
    .voteMedia(userStore.currentUser?.id || '', media.id, voteType)
    .then((isMatch) => {
      if (isMatch) {
        matchNotificationVisible.value = true;

        window.setTimeout(() => (matchNotificationVisible.value = false), 2000);
      }
    });
};

const handleButtonVote = (voteType: VoteType) => {
  currentVoteType.value = voteType;

  window.setTimeout(() => handleVote(voteType), 100);
};

const handleVote = (voteType: VoteType) => {
  if (currentMedia.value) {
    sendVote(currentMedia.value, voteType);
  }

  showNextMedia();
};

const handleSeen = () => {
  if (currentMedia.value) {
    apiClient.setMediaSeen(
      userStore.currentUser?.id || '',
      currentMedia.value?.id || ''
    );
  }

  showNextMedia();
};

onMounted(() => {
  fetchMedia();
});
</script>

<style lang="scss" scoped>
.vote-host {
  position: relative;

  width: 100%;
  height: 100%;

  display: flex;
  justify-content: center;
  align-items: center;

  .match-notification {
    font-size: 1.33rem;

    position: absolute;
    z-index: 30;

    background: #222;
    border-radius: 6px;
    display: flex;

    left: 50%;
    top: -100px;
    transform: translate(-50%, 0);

    padding: 12px 20px;

    transition-property: top;
    transition-duration: 300ms;
    transition-timing-function: ease-out;

    &.visible {
      top: 16px;
    }
  }
}

.button-wrapper {
  position: absolute;

  z-index: 30;

  bottom: 1.8rem;
  left: 0;

  width: 100%;
}
</style>
