<template>
  <div class="button-host">
    <div class="buttons">
      <div class="button negative" @click="handleVote(VoteType.NEGATIVE)">
        <PhX weight="bold" />
      </div>
      <div class="button neutral" @click="handleVote(VoteType.NEUTRAL)">
        <PhShuffle weight="bold" />
      </div>
      <div class="button neutral" @click="handleSeen">
        <PhCheck weight="bold" />
      </div>
      <div class="button positive" @click="handleVote(VoteType.POSITIVE)">
        <PhHeart weight="fill" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { PhCheck, PhHeart, PhShuffle, PhX } from '@phosphor-icons/vue';
import { VoteType } from '../../model/Vote';

interface Emits {
  (e: 'vote', v: VoteType): void;
  (e: 'seen'): void;
  (e: 'update:voteType', v: VoteType): void;
}

const emits = defineEmits<Emits>();

const handleVote = (voteType: VoteType) => {
  emits('update:voteType', voteType);
  emits('vote', voteType);
};

const handleSeen = () => emits('seen');
</script>

<style lang="scss" scoped>
.button-host {
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
