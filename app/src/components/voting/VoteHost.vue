<template>
  <div
    class="vote-host"
    @contextmenu="handleLongTap"
    @touchend="handleTouchEnd"
  >
    <div
      class="match-notification"
      :class="[matchNotificationVisible ? 'visible' : null]"
    >
      It's a match!
    </div>
    <MediaSwipeHost
      v-model:current-media-meta-visible="currentMediaMetaVisible"
      :current-media="currentMedia"
      :next-media="nextMedia"
      v-model:vote-type="currentVoteType"
      @vote="handleVote"
    />

    <div
      class="button-wrapper"
      v-if="currentMedia"
      :class="[hideButtons ? 'hidden' : 'visible']"
    >
      <ButtonHost @vote="handleButtonVote" @seen="handleSeen" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue';
import { Media } from '../../model/Media';
import { useApiClient } from '../../composables/useApiClient';
import { VoteType } from '../../model/Vote';
import MediaSwipeHost from './MediaSwipeHost.vue';
import ButtonHost from './ButtonHost.vue';

const MEDIA_COUNT_REFRESH_THRESHOLD = 5;

const apiClient = await useApiClient().apiClient;

const matchNotificationVisible = ref(false);

const currentVoteType = ref(VoteType.NEUTRAL);

const isFetchingMedia = ref(false);
const mediaList = ref([] as Media[]);
const belowScore = ref('100');
const mediaIndex = ref(0);

const currentMedia = computed(() => mediaList.value[mediaIndex.value] || null);
const nextMedia = computed(() => mediaList.value[mediaIndex.value + 1] || null);

const hideButtons = ref(false);

const currentMediaMetaVisible = ref(false);

watch(currentMediaMetaVisible, (v) => {
  if (v && hideButtons.value) {
    currentMediaMetaVisible.value = false;
  }
});

const handleLongTap = () => {
  hideButtons.value = true;
};

const handleTouchEnd = () => {
  hideButtons.value = false;
};

const showNextMedia = () => {
  ++mediaIndex.value;

  nextTick(() => (currentVoteType.value = VoteType.NEUTRAL));

  if (
    mediaIndex.value >
    mediaList.value.length - MEDIA_COUNT_REFRESH_THRESHOLD
  ) {
    fetchMedia();
  }
};

const fetchMedia = () => {
  if (isFetchingMedia.value) {
    return;
  }

  isFetchingMedia.value = true;

  apiClient.getRecommendedMedia(belowScore.value).then((results) => {
    mediaList.value.push(...results);

    const lastResult = results.pop();

    belowScore.value = lastResult?.score || '100';

    isFetchingMedia.value = false;
  });
};

const sendVote = (media: Media, voteType: VoteType) => {
  apiClient.voteMedia(media.id, voteType).then((isMatch) => {
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
    apiClient.setMediaSeen(currentMedia.value?.id || '');
  }
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

  .button-wrapper {
    position: absolute;

    z-index: 30;

    bottom: 1.8rem;
    left: 0;

    width: 100%;

    transition-property: opacity;
    transition-duration: 100ms;
    transition-timing-function: linear;

    &.hidden {
      opacity: 0;
    }
  }
}
</style>
