<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
  defineProps<{
    colors: string[];
    angle?: number;
    // in rem
    borderWidth?: number;
  }>(),
  {
    angle: 20,
    borderWidth: 0.2,
  }
);

const borderWidth = `${props.borderWidth}rem`;

const gradientSteps = (pos: number) =>
  (100 / (props.colors.length - 1 || 1)) * pos;

const gradient = computed(
  () =>
    `linear-gradient(${props.angle}deg, ${props.colors
      .map((col, pos) => `${col} ${gradientSteps(pos)}%`)
      .join(',')})`
);
</script>

<template>
  <div
    class="nice-wrapper"
    :style="{
      '--border-width': borderWidth,
      '--gradient': gradient,
    }"
  >
    <div class="content">
      <slot />
    </div>
  </div>
</template>

<style scoped lang="scss">
$border-radius: 1rem;

.nice-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;

  width: 100%;

  background: black;
  background-clip: padding-box;

  border: var(--border-width) solid transparent;
  border-radius: $border-radius;

  position: relative;

  //padding: 0.5rem 0.75rem;
  font-size: 2rem;

  &:before {
    content: '';
    position: absolute;

    top: 0;
    right: 0;
    bottom: 0;
    left: 0;

    z-index: -1;

    margin: calc(var(--border-width) * -1);

    border-radius: calc($border-radius * 0.85);
    background: var(--gradient);
  }

  .content {
    background: var(--gradient);

    font-family: Pacifico, sans-serif;
    text-transform: capitalize;
    width: 100%;

    display: flex;
    justify-content: center;
    align-items: center;

    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
}
</style>

<style lang="scss">
.nice-wrapper {
  .content {
    input {
      display: block;
      font-size: 1.25rem;

      -webkit-text-fill-color: #f1f1f1;

      padding: 1rem;
      width: 100%;

      &::placeholder {
        -webkit-text-fill-color: #999;
      }
    }

    a {
      width: 100%;
      text-decoration: none;

      text-align: center;
    }

    .clean {
      font-family: 'Roboto', sans-serif;

      -webkit-text-fill-color: #f1f1f1;
    }
  }
}
</style>
