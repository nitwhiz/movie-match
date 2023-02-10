<template>
  <div class="swipe-host">
    <div
      class="current-media"
      :style="currentMediaStyle"
      v-if="currentMedia"
      @touchstart.passive="handleTouchStart"
      @touchend.passive="handleTouchEnd"
      @touchmove.passive="handleTouchMove"
    >
      <MediaCard
        :media="currentMedia"
        v-model:meta-visible="showCurrentMediaMeta"
      />
    </div>
    <div class="next-media" v-if="nextMedia">
      <MediaCard :media="nextMedia" />
    </div>
    <div class="no-media" v-if="hadMedia && !currentMedia && !nextMedia">
      <span class="icon"><PhTelevision /></span>
      <span class="text">No media</span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { PhTelevision } from '@phosphor-icons/vue';
import { Media } from '../../model/Media';
import { computed, ref, watch } from 'vue';
import { VoteType } from '../../model/Vote';
import MediaCard from '../MediaCard.vue';

const VOTE_ANGLE_MIN = -10;
const VOTE_ANGLE_MAX = 10;

const VOTE_ANGLE_NEGATIVE_THRESHOLD = -5;
const VOTE_ANGLE_POSITIVE_THRESHOLD = 5;

interface Props {
  currentMediaMetaVisible: boolean;
  currentMedia: Media | null;
  nextMedia: Media | null;
  voteType: VoteType;
}

interface Emits {
  (e: 'vote', v: VoteType): void;
  (e: 'update:voteType', v: VoteType): void;
  (e: 'update:currentMediaMetaVisible', v: boolean): void;
}

const props = defineProps<Props>();

const emits = defineEmits<Emits>();

const hadMedia = ref(false);
const showCurrentMediaMeta = ref(props.currentMediaMetaVisible);

watch(
  () => props.currentMediaMetaVisible,
  (v) => (showCurrentMediaMeta.value = v)
);

watch(showCurrentMediaMeta, (v) => emits('update:currentMediaMetaVisible', v));

const unwatchCurrentMedia = watch(
  () => props.currentMedia,
  (v) => {
    if (v) {
      hadMedia.value = true;
      unwatchCurrentMedia();
    }
  }
);

const currentAngle = ref(0);

const screenWidth = ref(320);
const touchStartX = ref(screenWidth.value / 2);

watch(
  () => props.voteType,
  (voteType) => {
    switch (voteType) {
      case VoteType.NEGATIVE:
        currentAngle.value = VOTE_ANGLE_MIN;
        break;
      case VoteType.POSITIVE:
        currentAngle.value = VOTE_ANGLE_MAX;
        break;
      default:
        currentAngle.value = 0;
        break;
    }
  },
  {
    immediate: true,
  }
);

const currentMediaStyle = computed(() => ({
  rotate: `${currentAngle.value}deg`,
}));

const handleTouchStart = (e: TouchEvent) => {
  if (props.currentMediaMetaVisible) {
    return;
  }

  screenWidth.value = window.screen.availWidth;
  touchStartX.value = e.touches.item(0)!.pageX;
};

const voteTypeByAngle = computed(() => {
  let vote = VoteType.NEUTRAL;

  if (currentAngle.value <= VOTE_ANGLE_NEGATIVE_THRESHOLD) {
    vote = VoteType.NEGATIVE;
    currentAngle.value = VOTE_ANGLE_MIN;
  } else if (currentAngle.value >= VOTE_ANGLE_POSITIVE_THRESHOLD) {
    vote = VoteType.POSITIVE;
    currentAngle.value = VOTE_ANGLE_MAX;
  } else {
    currentAngle.value = 0;
  }

  return vote;
});

const handleTouchEnd = () => {
  if (showCurrentMediaMeta.value) {
    return;
  }

  const vote = voteTypeByAngle.value;

  emits('update:voteType', vote);

  if (vote !== VoteType.NEUTRAL) {
    emits('vote', vote);
  }
};

const handleTouchMove = (e: TouchEvent) => {
  if (showCurrentMediaMeta.value) {
    return;
  }

  const pX = e.touches.item(0)!.pageX;
  const sX = touchStartX.value;
  const dX = pX - sX;

  currentAngle.value = Math.max(
    VOTE_ANGLE_MIN,
    Math.min(VOTE_ANGLE_MAX, (dX / screenWidth.value) * 25)
  );
};
</script>

<style lang="scss" scoped>
.swipe-host {
  position: relative;

  width: 100%;
  height: 100%;

  display: flex;
  justify-content: center;
  align-items: center;

  .current-media,
  .next-media {
    position: absolute;

    width: 100%;
    height: 100%;
  }

  .current-media {
    z-index: 20;

    transition-property: rotate;
    transition-duration: 50ms;
    transition-timing-function: linear;

    transform-origin: center 150%;
  }

  .no-media {
    display: flex;
    flex-direction: column;
    align-items: center;

    padding: 1.75rem;

    color: #555;

    .icon {
      font-size: 12rem;
      line-height: 8rem;
    }

    .text {
      font-size: 4rem;
    }
  }

  .next-media {
    z-index: 10;

    filter: blur(10px);
  }
}
</style>
