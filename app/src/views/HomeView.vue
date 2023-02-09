<template>
  <div class="home">
    <router-link :to="{ name: 'vote' }" class="vote">
      <span class="text">Get Voting!</span>
    </router-link>
    <router-link :to="{ name: 'matches' }" class="matches">
      <span class="text">Check Matches</span>
    </router-link>
  </div>
</template>

<script lang="ts" setup></script>

<style lang="scss" scoped>
@use '../styles/nice';

.home {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;

  width: 100%;
  height: 100%;

  overflow: hidden;

  a {
    text-decoration: none;
    color: white;

    font-size: 2rem;
    font-family: Pacifico, sans-serif;
    text-transform: capitalize;

    padding: 0.5rem 1rem;

    margin-bottom: 1.25rem;
    width: 75%;
    min-width: 240px;

    text-align: center;

    $border-width: 3px;

    &.vote {
      @include nice.gradient-border(
        linear-gradient(20deg, rgb(148, 55, 255) 0%, rgb(101, 229, 255) 100%),
        $border-width
      );
    }

    &.matches {
      @include nice.gradient-border(
        linear-gradient(30deg, rgb(255, 55, 140) 0%, rgb(187, 255, 101) 100%),
        $border-width
      );
    }
  }
}
</style>
