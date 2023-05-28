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
      <Lottie
        ref="matchPartyPopper"
        :animationData="lottiePartyPopper"
        :width="24"
        :height="24"
        :loop="0"
      />&nbsp;&nbsp;It's a match!
    </div>
    <template v-if="props.swipe">
      <MediaSwipeHost
        v-model:current-media-meta-visible="currentMediaMetaVisible"
        :current-media="currentMedia"
        :next-media="nextMedia"
        v-model:vote-type="currentVoteType"
        @vote="handleVote"
      />
    </template>
    <template v-else-if="currentMedia">
      <MediaCard
        :media="currentMedia"
        :meta-visible="currentMediaMetaVisible"
        @click="currentMediaMetaVisible = !currentMediaMetaVisible"
      />
    </template>
    <div
      class="button-wrapper"
      v-if="currentMedia"
      :class="[hideButtons ? 'hidden' : 'visible']"
    >
      <ButtonHost
        @vote="handleButtonVote"
        @seen="handleSeen"
        :seen="currentMedia.seen || forceSeen || false"
      />
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
import lottiePartyPopper from '../../assets/lottie/party-popper.json';
import { Vue3Lottie } from 'vue3-lottie';
import MediaCard from '../MediaCard.vue';

interface Props {
  swipe?: boolean;
  media?: Media | null;
}

const props = withDefaults(defineProps<Props>(), {
  swipe: true,
  media: null,
});

const MEDIA_COUNT_REFRESH_THRESHOLD = 5;

const apiClient = await useApiClient().apiClient;

const matchNotificationVisible = ref(false);

const currentVoteType = ref(VoteType.NEUTRAL);

const isFetchingMedia = ref(false);
const mediaList = ref([] as Media[]);
const belowScore = ref('100');
const mediaIndex = ref(0);

const matchPartyPopper = ref(null as typeof Vue3Lottie | null);

const currentMedia = computed(() => mediaList.value[mediaIndex.value] || null);
const nextMedia = computed(() => mediaList.value[mediaIndex.value + 1] || null);

const hideButtons = ref(false);

const currentMediaMetaVisible = ref(false);

const forceSeen = ref(false);

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
  forceSeen.value = false;

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

      nextTick(() => {
        matchPartyPopper.value?.stop();
        matchPartyPopper.value?.play();
      });

      window.setTimeout(() => (matchNotificationVisible.value = false), 1500);
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

  if (!props.media) {
    showNextMedia();
  }
};

const handleSeen = () => {
  if (currentMedia.value) {
    apiClient.setMediaSeen(currentMedia.value?.id || '');
    forceSeen.value = true;
  }
};

onMounted(() => {
  if (props.media) {
    mediaList.value.push(props.media);
  } else {
    fetchMedia();
  }
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
    font-size: 1.1rem;

    display: flex;
    justify-content: center;
    align-items: center;

    position: absolute;
    z-index: 30;

    background: #222;
    border-radius: 64px;

    left: 50%;
    top: -100px;
    transform: translate(-50%, 0);

    padding: 10px 20px;

    transition-property: top;
    transition-duration: 300ms;
    transition-timing-function: ease-in-out;

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
