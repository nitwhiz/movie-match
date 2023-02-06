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
    <div class="meta-overlay" v-if="showMeta">
      <h3>{{ props.media.title }}</h3>
      <p class="summary">{{ props.media.summary }}</p>
      <p v-if="genres.length !== 0"><b>Genres: </b>{{ genres }}</p>
      <p><b>Release Date: </b>{{ releaseDate }}</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Media } from '../model/Media';
import { computed, ref, watch } from 'vue';
import { useApiClient } from '../composables/useApiClient';

interface Props {
  media: Media;
  metaVisible?: boolean;
}

interface Emits {
  (e: 'update:metaVisible', visible: boolean): void;
}

const props = withDefaults(defineProps<Props>(), {
  metaVisible: false,
});

const emits = defineEmits<Emits>();

const showMeta = ref(props.metaVisible);

const imageSrc = computed(() => useApiClient().getPosterUrl(props.media.id));

const genres = computed(() => props.media.genres.map((g) => g.name).join(', '));
const releaseDate = computed(() =>
  new Date(props.media.releaseDate).toLocaleDateString()
);

const inTouch = ref(false);
const firstTouch = ref(null as Touch | null);
const lastTouch = ref(null as Touch | null);

const setMetaVisible = (visible: boolean) => {
  showMeta.value = visible;
  emits('update:metaVisible', visible);
};

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
      setMetaVisible(!showMeta.value);
    }
  }

  inTouch.value = false;
  firstTouch.value = null;
  lastTouch.value = null;
};

watch(
  () => props.media,
  () => setMetaVisible(false)
);

watch(
  () => props.metaVisible,
  () => (showMeta.value = props.metaVisible)
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

    h3 {
      margin: 0 0 20px 0;
      font-size: 2rem;
    }

    p {
      margin: 0 0 12px 0;

      &.summary {
        max-height: 80%;

        /* todo: nice2have: scrolling if height > 80% */

        overflow: hidden;
      }
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
