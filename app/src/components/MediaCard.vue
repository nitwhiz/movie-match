<template>
  <div
    class="media-card"
    :style="{
      backgroundImage: `url(${imageSrc})`,
    }"
    @touchstart.passive="handleTouchStart"
    @touchmove.passive="handleTouchMove"
    @touchend.passive="handleTouchEnd"
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
import { Media, TMDBMovieData } from '../model/Media';
import { computed, ref, watch } from 'vue';
import { API_SERVER_BASE_URL } from '../common/Environment';

interface Props {
  media: Media<TMDBMovieData>;
}

const props = defineProps<Props>();

const metaVisible = ref(false);

const imageSrc = computed(
  () => `${API_SERVER_BASE_URL}/media/${props.media.ID}/poster`
);

const inTouch = ref(false);
const firstTouch = ref(null as Touch | null);
const lastTouch = ref(null as Touch | null);

const handleTouchStart = (e: TouchEvent) => {
  inTouch.value = true;

  firstTouch.value = e.touches.item(0) || null;
  lastTouch.value = e.touches.item(0) || null;
};

const handleTouchMove = (e: TouchEvent) => {
  if (inTouch.value) {
    lastTouch.value = e.touches.item(0) || null;
  }
};

const handleTouchEnd = () => {
  // todo: delta must be calculated during move

  const t1 = firstTouch.value;
  const t2 = lastTouch.value;

  if (t1 && t2) {
    const d2 = (t1.pageX - t2.pageX) ** 2 + (t1.pageY - t2.pageY) ** 2;

    if (d2 <= 5 ** 2) {
      metaVisible.value = !metaVisible.value;
    }
  }

  inTouch.value = false;
  firstTouch.value = null;
  lastTouch.value = null;
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
