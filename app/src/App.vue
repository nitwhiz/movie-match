<template>
  <Suspense>
    <router-view />
    <template #fallback>Loading ...</template>
  </Suspense>
</template>
