<template>
  <div class="card-wrapper">
    <div
      class="current-card"
      ref="currentCard"
      :style="{
        rotate: `${voteAngle}deg`,
      }"
      @touchstart="handleTouchStart"
      @touchend="handleTouchEnd"
      @touchmove="handleTouchMove"
    >
      <MediaCard v-if="currentMedia" :media="currentMedia" />
    </div>
    <div class="next-card">
      <MediaCard v-if="nextMedia" :media="nextMedia" />
    </div>
  </div>
  <div class="button-wrapper">
    <div class="buttons">
      <div class="button negative" @click="sendVote('negative')">
        <PhX weight="bold" />
      </div>
      <div class="button neutral" @click="sendVote('neutral')">
        <PhShuffle weight="bold" />
      </div>
      <div class="button neutral">
        <PhCheck weight="bold" />
      </div>
      <div class="button positive" @click="sendVote('positive')">
        <PhHeart weight="fill" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { Media, TMDBMovieData } from '../model/Media';
import axios from 'axios';
import MediaCard from './MediaCard.vue';
import { PhHeart, PhX, PhShuffle, PhCheck } from 'phosphor-vue';

const mediaPtr = ref(0);
const media = ref([] as Media<TMDBMovieData>[]);

const currentMedia = computed(() => media.value[mediaPtr.value] || null);
const nextMedia = computed(() => media.value[mediaPtr.value + 1] || null);

const voteAngle = ref(0);

const screenWidth = ref(320);
const touchStartX = ref(screenWidth.value / 2);

const sendVote = (voteType: string) => {
  axios.put(
    'http://192.168.127.22:6445/user/72416d6d-0b92-4f63-a11e-e834b1ba74e0/vote',
    {
      mediaId: media.value[mediaPtr.value].ID,
      voteType: voteType,
    }
  );

  ++mediaPtr.value;
};

const handleTouchStart = (e: TouchEvent) => {
  e.preventDefault();
  e.stopPropagation();

  screenWidth.value = window.screen.availWidth;
  touchStartX.value = e.touches.item(0)!.pageX;
};

const handleTouchEnd = (e: TouchEvent) => {
  e.preventDefault();
  e.stopPropagation();

  let voteResult = 'neutral';

  if (voteAngle.value <= -5) {
    voteResult = 'negative';
  } else if (voteAngle.value >= 5) {
    voteResult = 'positive';
  }

  if (voteResult !== 'neutral') {
    sendVote(voteResult);
  }

  voteAngle.value = 0;
};

const handleTouchMove = (e: TouchEvent) => {
  e.preventDefault();
  e.stopPropagation();

  const pX = e.touches.item(0)!.pageX;
  const sX = touchStartX.value;
  const dX = pX - sX;

  voteAngle.value = Math.max(-10, Math.min(10, (dX / screenWidth.value) * 25));
};

onMounted(() => {
  axios
    .get<{ Results: Media<TMDBMovieData>[] }>(
      'http://192.168.127.22:6445/media'
    )
    .then(({ data: { Results: results } }) => {
      media.value = results;
    });
});
</script>

<style lang="scss" scoped>
.card-wrapper {
  position: relative;

  width: 100%;
  height: 100%;

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

      &:last-child {
        margin-right: 0;
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
