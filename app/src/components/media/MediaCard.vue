<script setup lang="ts">
import { RecommendedMedia } from '../../model/Media';
import { computed, ref, watch } from 'vue';
import {
  PhHeart,
  PhShareNetwork,
  PhStar,
  PhStarHalf,
  PhThumbsDown,
} from '@phosphor-icons/vue';
import { VoteType } from '../../model/Vote';
import { useMediaPoster } from '../../composables/useMediaPoster';
import { useRouter } from 'vue-router';
import { RouteName } from '../../router';
import { useMediaType } from '../../composables/useMediaType';

const props = defineProps<{
  media: RecommendedMedia;
}>();

const emits = defineEmits<{
  (e: 'vote', voteType: VoteType): void;
}>();

const router = useRouter();

const { posterUrl } = useMediaPoster(() => props.media.id);

const cardStyle = computed(() =>
  posterUrl.value
    ? {
        backgroundImage: `url(${posterUrl.value})`,
      }
    : {}
);

const isLiked = computed(() => props.media.voteType === VoteType.POSITIVE);
const isDisliked = computed(() => props.media.voteType === VoteType.NEGATIVE);

const metaVisible = ref(false);

const handleLike = (e: TouchEvent) => {
  e.stopPropagation();

  emits('vote', VoteType.POSITIVE);
};

const handleDislike = (e: TouchEvent) => {
  e.stopPropagation();

  emits('vote', VoteType.NEGATIVE);
};

const handleShare = () => {
  const url = `${location.protocol}//${location.host}/${
    router.resolve({
      name: RouteName.MEDIA,
      params: { mediaId: props.media.id },
    }).href
  }`;

  navigator.share({
    url,
    title: `How about ${props.media.title}?`,
  });
};

const touchMoved = ref(false);

const handleTouchStart = (e: TouchEvent) => {
  if (metaVisible.value) {
    e.stopPropagation();
  }

  touchMoved.value = false;
};
const handleTouchEnd = (e: TouchEvent) => {
  if (!touchMoved.value) {
    e.stopPropagation();

    metaVisible.value = !metaVisible.value;
  }

  touchMoved.value = false;
};
const handleTouchMove = (e: TouchEvent) => {
  if (metaVisible.value) {
    e.stopPropagation();
  }

  touchMoved.value = true;
};

watch(
  () => props.media.id,
  () => (metaVisible.value = false),
  {
    immediate: true,
  }
);

const { getMediaTypeLabelSingular } = useMediaType();

const mediaTypeLabel = computed(() =>
  getMediaTypeLabelSingular(props.media.type)
);

const ratingFilled = computed(() => Math.floor(props.media.rating / 20));
const ratingHalf = computed(
  () => props.media.rating / 20 - Math.floor(props.media.rating / 20) >= 0.5
);
const ratingEmpty = computed(
  () => 5 - ratingFilled.value - (ratingHalf.value ? 1 : 0)
);
const genres = computed(() => props.media.genres.map((g) => g.name).join(', '));
const releaseDate = computed(() =>
  new Date(props.media.releaseDate).toLocaleDateString(['de-DE'], {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
);
const runtime = computed(() => {
  if (props.media.runtime === 0) {
    return '';
  }

  const hours = Math.floor(props.media.runtime / 60);
  const minutes = Math.floor(props.media.runtime - hours * 60);

  return `${hours > 0 ? `${hours}h ` : ''}${minutes}m`;
});
</script>

<template>
  <div
    class="media-card"
    :style="cardStyle"
    @touchstart.passive="handleTouchStart"
    @touchend.passive="handleTouchEnd"
    @touchmove.passive="handleTouchMove"
  >
    <div class="poster">
      <img v-if="posterUrl" :src="posterUrl" :alt="props.media.title" />
    </div>
    <div class="controls">
      <div
        class="control like"
        :class="[isLiked ? 'filled' : null]"
        @touchend.passive="handleLike"
      >
        <PhHeart :weight="isLiked ? 'fill' : 'regular'" />
      </div>
      <div
        class="control dislike"
        :class="[isDisliked ? 'filled' : null]"
        @touchend.passive="handleDislike"
      >
        <PhThumbsDown :weight="isDisliked ? 'fill' : 'regular'" />
      </div>
      <div @touchend.passive="handleShare" class="control share">
        <PhShareNetwork />
      </div>
    </div>
    <div class="meta" v-if="metaVisible">
      <h3>{{ props.media.title }}</h3>
      <h4>{{ mediaTypeLabel }}</h4>
      <p v-if="props.media.summary" class="summary">
        {{ props.media.summary }}
      </p>
      <p class="sub" v-if="genres.length !== 0"><b>Genres: </b>{{ genres }}</p>
      <p class="sub" v-if="releaseDate"><b>Release: </b>{{ releaseDate }}</p>
      <p class="sub" v-if="runtime"><b>Runtime: </b>{{ runtime }}</p>
      <p class="sub rating" v-if="props.media.rating">
        <b>Rating: </b>
        <span class="stars">
          <PhStar weight="fill" v-for="_ in ratingFilled" />
          <PhStarHalf weight="fill" v-if="ratingHalf" />
          <PhStar weight="duotone" v-for="_ in ratingEmpty" />
        </span>
      </p>
    </div>
  </div>
</template>

<style scoped lang="scss">
.media-card {
  background-size: cover;
  background-repeat: no-repeat;
  background-position: center;

  width: 100%;
  height: 100%;

  overflow: hidden;

  position: relative;

  .poster {
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
      display: block;

      width: auto;
      height: auto;

      max-width: 100%;
      max-height: 100%;

      box-shadow: 0 0 50px #111;
    }
  }

  .controls {
    position: absolute;

    bottom: 5%;
    right: 1.5rem;

    z-index: 30;

    font-size: 2.5rem;

    color: #f1f1f1;

    .control {
      margin-bottom: 1.75rem;

      &:last-child {
        margin-bottom: 0;
      }

      svg {
        display: block;
        filter: drop-shadow(0px 0px 7px rgba(0, 0, 0, 0.5));
      }

      &.filled {
        color: #db1a1a;
      }
    }
  }

  .meta {
    width: 100%;
    height: 100%;

    position: absolute;

    z-index: 20;

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
        overflow-x: hidden;
        overflow-y: auto;
      }

      &.rating {
        display: flex;
        align-items: center;

        b {
          margin-right: 6px;
        }

        .stars {
          display: flex;
          align-items: center;

          color: gold;
        }
      }

      &.sub {
        width: calc(100% - 4.5rem);
      }
    }
  }
}
</style>
