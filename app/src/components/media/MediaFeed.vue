<script setup lang="ts">
import { RecommendedMedia } from '../../model/Media';
import MediaCard from './MediaCard.vue';
import { computed, nextTick, ref, watch } from 'vue';
import { VoteType } from '../../model/Vote';

const MEDIA_WINDOW_SIZE = 3;

const AUTO_SCROLL_DURATION = 70;

const AUTO_SCROLL_DISTANCE_THRESHOLD = 200;
const AUTO_SCROLL_PPMS_THRESHOLD = 0.4;

const AUTO_SCROLL_DOWN_DISTANCE_THRESHOLD = AUTO_SCROLL_DISTANCE_THRESHOLD;
const AUTO_SCROLL_UP_DISTANCE_THRESHOLD = -AUTO_SCROLL_DISTANCE_THRESHOLD;

const AUTO_SCROLL_DOWN_PPMS_THRESHOLD = AUTO_SCROLL_PPMS_THRESHOLD;
const AUTO_SCROLL_UP_PPMS_THRESHOLD = -AUTO_SCROLL_PPMS_THRESHOLD;

const ENDOFMEDIA_THRESHOLD = 5;

const props = defineProps<{
  media?: RecommendedMedia[];
}>();

const emits = defineEmits<{
  (e: 'endofmedia'): void;
  (e: 'vote', media: RecommendedMedia, voteType: VoteType): void;
}>();

const mediaList = ref([] as RecommendedMedia[]);
const mediaWindow = ref([] as RecommendedMedia[]);

// points to the **next** media index
const mediaPtr = ref(0);
// points to the **current** window index
const windowPtr = ref(0);

watch(
  () => props.media,
  (newMedia = []) => {
    mediaList.value.push(...newMedia);

    for (let i = 0; i < MEDIA_WINDOW_SIZE; ++i) {
      if (!mediaWindow.value[i] && mediaList.value[mediaPtr.value]) {
        mediaWindow.value[i] = mediaList.value[mediaPtr.value];
        ++mediaPtr.value;
      }
    }
  },
  {
    immediate: true,
  }
);

const touchStartY = ref(null as number | null);
const touchStartTime = ref(null as number | null);

const touchDelta = ref(0);

const feedDY = ref(0);

const isAutoScrolling = ref(false);

const currentMedia = computed(() => mediaWindow.value[windowPtr.value] || null);

const transitionStyle = computed(() =>
  isAutoScrolling.value
    ? {
        transitionProperty: 'top',
        transitionDuration: `${AUTO_SCROLL_DURATION}ms`,
        transitionTimingFunction: 'linear',
      }
    : {}
);

const autoScroll = (
  scrollTo: number,
  preCallback: (() => void) | null = null,
  postCallback: (() => void) | null = null
) => {
  isAutoScrolling.value = true;

  preCallback ? preCallback() : void 0;

  nextTick(() => {
    feedDY.value = scrollTo;

    window.setTimeout(() => {
      isAutoScrolling.value = false;

      postCallback ? postCallback() : void 0;
    }, AUTO_SCROLL_DURATION);
  });
};

const handleTouchStart = (e: TouchEvent) => {
  touchDelta.value = 0;

  touchStartY.value = e.touches[0]?.clientY;
  touchStartTime.value = Date.now();
};

const handleTouchMove = (e: TouchEvent) => {
  if (touchStartY.value !== null) {
    touchDelta.value = touchStartY.value - e.touches[0]?.clientY;
    feedDY.value = touchDelta.value + window.innerHeight * windowPtr.value;
  }
};

const handleAutoScrollDown = () => {
  autoScroll(window.innerHeight * (windowPtr.value + 1), null, () => {
    if (windowPtr.value > 0) {
      mediaWindow.value[0] = mediaWindow.value[1];
      mediaWindow.value[1] = mediaWindow.value[2];
      mediaWindow.value[2] = mediaList.value[mediaPtr.value];

      feedDY.value = window.innerHeight * windowPtr.value;

      ++mediaPtr.value;
    } else {
      ++windowPtr.value;
    }

    if (mediaList.value.length - mediaPtr.value <= ENDOFMEDIA_THRESHOLD) {
      emits('endofmedia');
    }
  });
};

const handleAutoScrollUp = () => {
  autoScroll(window.innerHeight * (windowPtr.value - 1), null, () => {
    if (windowPtr.value > 0 && mediaPtr.value >= 4) {
      mediaWindow.value[2] = mediaWindow.value[1];
      mediaWindow.value[1] = mediaWindow.value[0];
      mediaWindow.value[0] = mediaList.value[mediaPtr.value - 4];

      feedDY.value = window.innerHeight * windowPtr.value;

      --mediaPtr.value;
    } else if (windowPtr.value > 0) {
      --windowPtr.value;
    }
  });
};

const handleTouchEnd = () => {
  const touchPPMS =
    touchDelta.value / (Date.now() - (touchStartTime.value || 0));
  const media = currentMedia.value;

  if (
    touchPPMS >= AUTO_SCROLL_DOWN_PPMS_THRESHOLD ||
    touchDelta.value >= AUTO_SCROLL_DOWN_DISTANCE_THRESHOLD
  ) {
    if (media && media.voteType === VoteType.NONE) {
      emits('vote', media, VoteType.NEUTRAL);
    }

    handleAutoScrollDown();
  } else if (
    touchPPMS <= AUTO_SCROLL_UP_PPMS_THRESHOLD ||
    touchDelta.value <= AUTO_SCROLL_UP_DISTANCE_THRESHOLD
  ) {
    if (media && media.voteType === VoteType.NONE) {
      emits('vote', media, VoteType.NEUTRAL);
    }

    handleAutoScrollUp();
  } else {
    autoScroll(window.innerHeight * windowPtr.value);
  }
};

const handleVote = (media: RecommendedMedia, voteType: VoteType) => {
  emits('vote', media, voteType);

  if (
    media.voteType === VoteType.POSITIVE ||
    media.voteType === VoteType.NEGATIVE
  ) {
    handleAutoScrollDown();
  }
};
</script>

<template>
  <div
    class="feed"
    @touchmove.passive="handleTouchMove"
    @touchstart.passive="handleTouchStart"
    @touchend="handleTouchEnd"
  >
    <div
      class="item"
      v-for="(media, i) in mediaWindow"
      :style="{
        top: `calc(-${feedDY}px + 100% * ${i})`,
        ...transitionStyle,
      }"
    >
      <MediaCard :media="media" @vote="(vT) => handleVote(media, vT)" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.feed {
  position: relative;

  width: 100%;
  height: 100%;

  overflow: hidden;

  .item {
    position: absolute;

    width: 100%;
    height: 100%;

    overflow: hidden;
  }
}
</style>
