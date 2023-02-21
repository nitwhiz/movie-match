<template>
  <Suspense>
    <ApiClientProvider>
      <router-view />
    </ApiClientProvider>
    <template #fallback>Loading ...</template>
  </Suspense>
</template>

<script lang="ts" setup>
import ApiClientProvider from './ApiClientProvider.vue';
</script>
