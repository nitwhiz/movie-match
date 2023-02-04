<template>
  <div
    class="media-card"
    :style="{
      backgroundImage: `url(${imageSrc})`,
    }"
    @touchstart="handleTouchStart"
    @touchmove="handleTouchMove"
    @touchend="handleTouchEnd"
  >
    <div class="image-holder">
      <img :src="imageSrc" alt="" />
    </div>
    <div class="meta-overlay" v-if="metaVisible">
      <h1>{{ props.media.Data.Title }}</h1>
      <p>{{ props.media.Data.Overview }}</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Media, TMDBMovieData, TMDBPosterBaseUrl } from '../model/Media';
import { computed, ref, watch } from 'vue';

interface Props {
  media: Media<TMDBMovieData>;
}

const props = defineProps<Props>();

const metaVisible = ref(false);

const imageSrc = computed(
  () => `${TMDBPosterBaseUrl}${props.media.Data.poster_path}`
);

const inTouch = ref(false);
const touchMoved = ref(false);

const handleTouchStart = () => {
  inTouch.value = true;
};

const handleTouchMove = () => {
  if (inTouch.value) {
    touchMoved.value = true;
  }
};

const handleTouchEnd = () => {
  // todo: this does not work on device

  if (!touchMoved.value) {
    metaVisible.value = !metaVisible.value;
  }

  inTouch.value = false;
  touchMoved.value = false;
};

watch(
  () => props.media,
  () => (metaVisible.value = false)
);
</script>

<style lang="scss" scoped>
.media-card {
  width: 100%;
  height: 100%;

  overflow: hidden;

  position: relative;

  background-size: cover;
  background-repeat: no-repeat;
  background-position: center;

  .image-holder {
    display: flex;
    justify-content: center;
    align-items: center;

    position: absolute;
    z-index: 10;

    width: 100%;
    height: 100%;

    backdrop-filter: blur(20px);
    background-color: rgba(0, 0, 0, 0.5);

    img {
      width: auto;
      height: auto;
      max-width: 100%;
      max-height: 100%;
      box-shadow: 0 0 50px #111;
    }
  }

  .meta-overlay {
    position: absolute;
    z-index: 20;

    width: 100%;
    height: 100%;

    background-color: rgba(0, 0, 0, 0.75);
    color: white;
    text-shadow: 0 0 8px black;
    padding: 1.5rem 1rem 7rem 1rem;

    display: flex;
    justify-content: flex-end;
    flex-direction: column;

    h1 {
      margin: 0;
    }
  }

  .title {
    position: absolute;

    width: 100%;

    color: white;

    z-index: 20;

    bottom: 0;
    left: 0;
  }
}
</style>
