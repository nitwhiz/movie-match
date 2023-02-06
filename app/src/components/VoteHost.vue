<template>
  <div class="card-wrapper">
    <div
      class="match-notification"
      :class="[matchNotificationVisible ? 'visible' : null]"
    >
      It's a match!
    </div>
    <div
      v-if="currentMedia"
      class="current-card"
      ref="currentCard"
      :style="{
        rotate: `${voteAngle}deg`,
      }"
      @touchstart="handleTouchStart"
      @touchend.passive="handleTouchEnd"
      @touchmove.passive="handleTouchMove"
    >
      <MediaCard :media="currentMedia" />
    </div>
    <div v-else-if="!isFetchingMedia" class="empty">No movies found! :(</div>
    <div v-if="nextMedia" class="next-card">
      <MediaCard :media="nextMedia" />
    </div>
  </div>
  <div class="button-wrapper">
    <div class="buttons" v-if="currentMedia">
      <div
        class="button negative"
        @click="buttonVote(VoteType.NEGATIVE)"
        :class="[voteAngle <= -5 ? 'active' : null]"
      >
        <PhX weight="bold" />
      </div>
      <div class="button neutral" @click="buttonVote(VoteType.NEUTRAL)">
        <PhShuffle weight="bold" />
      </div>
      <div class="button neutral" @click="sendSeen">
        <PhCheck weight="bold" />
      </div>
      <div
        class="button positive"
        @click="buttonVote(VoteType.POSITIVE)"
        :class="[voteAngle >= 5 ? 'active' : null]"
      >
        <PhHeart weight="fill" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { Media } from '../model/Media';
import MediaCard from './MediaCard.vue';
import { PhCheck, PhHeart, PhShuffle, PhX } from 'phosphor-vue';
import { useUserStore } from '../store/userStore';
import { useApiClient } from '../composables/useApiClient';
import { VoteType } from '../model/Vote';

const userStore = useUserStore();
const apiClient = useApiClient();

const isFetchingMedia = ref(false);

const mediaPagePtr = ref(0);
const mediaPtr = ref(0);
const media = ref([] as Media[]);

const currentMedia = computed(() => media.value[mediaPtr.value] || null);
const nextMedia = computed(() => media.value[mediaPtr.value + 1] || null);

const voteAngle = ref(0);

const screenWidth = ref(320);
const touchStartX = ref(screenWidth.value / 2);

const matchNotificationVisible = ref(false);

const fetchMedia = () => {
  apiClient
    .getRecommendedMedia(userStore.currentUser?.id || '', mediaPagePtr.value)
    .then((results) => {
      media.value.push(...results);

      isFetchingMedia.value = false;
    });
};

const incMediaPtr = () => {
  ++mediaPtr.value;

  if (!isFetchingMedia.value && mediaPtr.value > media.value.length - 3) {
    isFetchingMedia.value = true;

    ++mediaPagePtr.value;
    fetchMedia();
  }
};

// todo: split sending votes and animating votes into funcs to be called if needed
// todo: icons make bundle size go brrr
// todo: use animation with keyframes for button bouncing

const buttonVote = (voteType: VoteType) => {
  if (voteType === VoteType.NEGATIVE) {
    voteAngle.value = -10;
  } else if (voteType === VoteType.POSITIVE) {
    voteAngle.value = 10;
  } else {
    voteAngle.value = 0;
  }

  window.setTimeout(() => {
    sendVote(voteType);
    voteAngle.value = 0;
  }, 50);
};

const sendVote = (voteType: VoteType) => {
  apiClient
    .voteMedia(
      userStore.currentUser?.id || '',
      media.value[mediaPtr.value].id,
      voteType
    )

    .then((isMatch) => {
      if (isMatch) {
        matchNotificationVisible.value = true;

        window.setTimeout(() => (matchNotificationVisible.value = false), 2000);
      }
    });

  incMediaPtr();
};

const sendSeen = () => {
  apiClient.setMediaSeen(
    userStore.currentUser?.id || '',
    media.value[mediaPtr.value].id,
    true
  );
};

const handleTouchStart = (e: TouchEvent) => {
  e.preventDefault();
  e.stopPropagation();

  screenWidth.value = window.screen.availWidth;
  touchStartX.value = e.touches.item(0)!.pageX;
};

const handleTouchEnd = () => {
  let voteResult = VoteType.NEUTRAL;

  if (voteAngle.value <= -5) {
    voteResult = VoteType.NEGATIVE;
  } else if (voteAngle.value >= 5) {
    voteResult = VoteType.POSITIVE;
  }

  if (voteResult !== VoteType.NEUTRAL) {
    sendVote(voteResult);
  }

  voteAngle.value = 0;
};

const handleTouchMove = (e: TouchEvent) => {
  const pX = e.touches.item(0)!.pageX;
  const sX = touchStartX.value;
  const dX = pX - sX;

  voteAngle.value = Math.max(-10, Math.min(10, (dX / screenWidth.value) * 25));
};

onMounted(() => {
  isFetchingMedia.value = true;
  fetchMedia();
});
</script>

<style lang="scss" scoped>
.card-wrapper {
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

  .next-card,
  .current-card {
    position: absolute;

    width: 100%;
    height: 100%;
  }

  .current-card {
    z-index: 20;

    transition-property: rotate;
    transition-duration: 50ms;
    transition-timing-function: linear;

    transform-origin: center 150%;
  }

  .empty {
    text-align: center;
    font-size: 3rem;

    padding: 1.75rem;

    color: rgba(255, 255, 255, 0.33);
  }

  .next-card {
    z-index: 10;

    filter: blur(10px);
  }
}

.button-wrapper {
  position: absolute;

  z-index: 30;

  bottom: 1.8rem;
  left: 0;

  width: 100%;

  display: flex;
  justify-content: center;

  color: white;

  .buttons {
    display: flex;
    justify-content: center;
    align-items: center;

    .button {
      $size: 4rem;

      display: flex;
      justify-content: center;
      align-items: center;

      font-size: calc($size / 1.75);
      border-radius: calc($size / 2);

      width: $size;
      height: $size;
      overflow: hidden;

      margin-right: calc($size / 4);

      box-shadow: 0 0 10px rgba(20, 20, 20, 0.25);

      transform: scale(1);
      transition-property: transform;
      transition-duration: 100ms;
      transition-timing-function: ease-out;

      &:last-child {
        margin-right: 0;
      }

      &.active {
        transform: scale(1.1);
        transition-duration: 200ms;
        transition-timing-function: cubic-bezier(0.5, 1.5, 0.5, 2);
      }

      &.positive {
        background-color: #40db72;
      }

      &.neutral {
        background-color: #565ddb;
      }

      &.negative {
        background-color: #db2d2a;
      }
    }
  }
}
</style>
