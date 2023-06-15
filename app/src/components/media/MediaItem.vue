<template>
  <NiceWrapper
    class="media-item"
    @click="handleClick"
    :colors="['rgb(34, 193, 195)', 'rgb(253, 187, 45)']"
  >
    <div class="wrapper clean">
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
      <div class="vote-type">
        <PhHeart v-if="isLiked" size="24" weight="fill" />
        <PhThumbsDown v-else-if="isDisliked" size="24" weight="fill" />
      </div>
    </div>
  </NiceWrapper>
</template>

<script setup lang="ts">
import { Media } from '../../model/Media';
import { computed } from 'vue';
import { useMediaType } from '../../composables/useMediaType';
import { PhCameraSlash, PhHeart, PhThumbsDown } from '@phosphor-icons/vue';
import { useRouter } from 'vue-router';
import { RouteName } from '../../router';
import NiceWrapper from '../nice/NiceWrapper.vue';
import { useMediaPoster } from '../../composables/useMediaPoster';
import { VoteType } from '../../model/Vote';

interface Props {
  media: Media;
  voteType?: VoteType;
}

const { getMediaTypeLabelSingular } = useMediaType();

const router = useRouter();

const props = withDefaults(defineProps<Props>(), {
  voteType: VoteType.POSITIVE,
});

const isLiked = computed(() => props.voteType === VoteType.POSITIVE);
const isDisliked = computed(() => props.voteType === VoteType.NEGATIVE);

const { posterUrl } = useMediaPoster(() => props.media.id);

const mediaTypeLabel = computed(() =>
  getMediaTypeLabelSingular(props.media.type)
);

const genres = computed(() => props.media.genres.map((g) => g.name).join(', '));

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
.media-item {
  .wrapper {
    position: relative;

    width: 100%;

    display: flex;
    align-items: stretch;

    font-family: 'Roboto', sans-serif;

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

      .name {
        font-size: 1.1rem;
        margin-bottom: 0.5rem;
      }

      .type,
      .genres {
        font-size: 0.8rem;
        margin-top: 0.2rem;
      }
    }

    .vote-type {
      position: absolute;

      bottom: 0;
      right: 0.4rem;

      color: rgba(255, 0, 0, 0.75);

      transform: rotate(15deg);
    }
  }
}
</style>
