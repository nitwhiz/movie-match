<template>
  <div
    class="media-card"
    :style="cardStyle"
    @touchstart.passive="handleTouchStart"
    @touchmove.passive="handleTouchMove"
    @touchend.passive="handleTouchEnd"
  >
    <div class="image-holder">
      <img :src="posterUrl" alt="" />
    </div>
    <div class="meta-overlay" v-if="showMeta">
      <h3>{{ props.media.title }}</h3>
      <h4>{{ mediaTypeLabel }}</h4>
      <p v-if="props.media.summary" class="summary">
        {{ props.media.summary }}
      </p>
      <p v-if="genres.length !== 0"><b>Genres: </b>{{ genres }}</p>
      <p v-if="releaseDate"><b>Release: </b>{{ releaseDate }}</p>
      <p v-if="runtime"><b>Runtime: </b>{{ runtime }}</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { Media } from '../model/Media';
import { computed, ref, watch } from 'vue';
import { useApiClient } from '../composables/useApiClient';
import { useMediaType } from '../composables/useMediaType';

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

const posterUrl = computed(() => useApiClient().getPosterUrl(props.media.id));

const genres = computed(() => props.media.genres.map((g) => g.name).join(', '));
const releaseDate = computed(() =>
  new Date(props.media.releaseDate).toLocaleDateString()
);
const runtime = computed(() => {
  if (props.media.runtime === 0) {
    return '';
  }

  const hours = Math.floor(props.media.runtime / 60);
  const minutes = Math.floor(props.media.runtime - hours * 60);

  return `${hours > 0 ? `${hours}h ` : ''}${minutes}m`;
});

const inTouch = ref(false);
const isTap = ref(false);
const firstTouch = ref(null as Touch | null);

const cardStyle = computed(() => ({
  backgroundImage: `url(${posterUrl.value})`,
}));

const { getMediaTypeLabelSingular } = useMediaType();

const mediaTypeLabel = computed(() =>
  getMediaTypeLabelSingular(props.media.type)
);

const setMetaVisible = (visible: boolean) => {
  showMeta.value = visible;
  emits('update:metaVisible', visible);
};

const handleTouchStart = (e: TouchEvent) => {
  inTouch.value = true;

  firstTouch.value = e.touches.item(0) || null;
  isTap.value = true;
};

const handleTouchMove = (e: TouchEvent) => {
  if (isTap.value && inTouch.value) {
    const t1 = firstTouch.value;
    const t2 = e.touches.item(0) || null;

    if (t1 && t2) {
      const d2 = (t1.pageX - t2.pageX) ** 2 + (t1.pageY - t2.pageY) ** 2;

      if (d2 > 5 ** 2) {
        isTap.value = false;
      }
    }
  }
};

const handleTouchEnd = () => {
  if (isTap.value) {
    setMetaVisible(!showMeta.value);
  }

  inTouch.value = false;
  isTap.value = false;
  firstTouch.value = null;
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
    background-color: rgba(0, 0, 0, 0.25);

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
    text-shadow: 0 0 1rem black;
    padding: 1rem 1rem 8rem 1rem;

    display: flex;
    justify-content: flex-end;
    flex-direction: column;

    h3 {
      margin: 0;
      font-size: 2.5rem;
      font-weight: normal;
    }

    h4 {
      margin: 0;
      font-size: 0.8rem;
      font-weight: normal;
      line-height: 1.25rem;
    }

    p {
      margin: 12px 0 0 0;
      font-size: 1rem;
      line-height: 1.25rem;

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
