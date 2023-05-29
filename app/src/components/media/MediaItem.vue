<template>
  <div class="media-item" @click="handleClick">
    <div class="poster">
      <img v-if="posterUrl" :src="posterUrl" :alt="props.media.title" />
      <div class="no-poster" v-else>
        <PhCameraSlash :size="32" weight="duotone" />
      </div>
    </div>
    <div class="details">
      <b class="name">{{ props.media.title }}</b>
      <span class="type">{{ mediaTypeLabel }}</span>
      <span class="genres">{{ genres }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Media } from '../../model/Media';
import {
  freeMediaPosterBlobUrl,
  getMediaPosterBlobUrl,
} from '../../api/PosterBlob';
import { computed, ref, watch } from 'vue';
import { useMediaType } from '../../composables/useMediaType';
import { PhCameraSlash } from '@phosphor-icons/vue';
import { useRouter } from 'vue-router';
import { RouteName } from '../../router';

interface Props {
  media: Media;
}

const { getMediaTypeLabelSingular } = useMediaType();

const router = useRouter();

const props = defineProps<Props>();

const posterUrl = ref('');
const isLoadingPoster = ref(true);

const mediaTypeLabel = computed(() =>
  getMediaTypeLabelSingular(props.media.type)
);

const genres = computed(() => props.media.genres.map((g) => g.name).join(', '));

watch(
  () => props.media.id,
  (_, oldValue) => {
    if (oldValue) {
      freeMediaPosterBlobUrl(oldValue);
    }

    getMediaPosterBlobUrl(props.media.id).then((url) => {
      if (url) {
        posterUrl.value = url;
      }

      isLoadingPoster.value = false;
    });
  },
  {
    immediate: true,
  }
);

const handleClick = () => {
  router.push({
    name: RouteName.MEDIA,
    params: {
      mediaId: props.media.id,
    },
  });
};
</script>

<style scoped lang="scss">
@use '../../styles/nice';

.media-item {
  width: 100%;
  display: flex;

  align-items: stretch;

  $border-width: 3px;

  @include nice.gradient-border(
    linear-gradient(20deg, rgb(34, 193, 195) 0%, rgb(253, 187, 45) 100%),
    $border-width
  );

  .poster {
    width: 72px;

    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;

    overflow: hidden;

    background: #222;

    display: flex;
    justify-content: center;
    align-items: center;

    img {
      display: block;

      width: 100%;
      height: auto;
    }

    .no-poster {
      width: 100%;
      height: 100%;

      display: flex;

      justify-content: center;
      align-items: center;
    }
  }

  .details {
    display: flex;
    flex-direction: column;
    justify-content: center;

    flex-grow: 0;
    width: calc(100% - 72px);

    padding: 12px 20px;

    .type,
    .genres {
      font-size: 0.8rem;
      margin-top: 0.2rem;
    }
  }
}
</style>
