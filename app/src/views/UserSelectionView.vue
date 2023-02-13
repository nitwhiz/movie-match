<template>
  <div class="users">
    <h1>Welcome</h1>
    <div class="user" v-for="u in users" @click="handleLogin(u)">
      <span class="text">{{ u.name }}</span>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { User } from '../model/User';
import { useUserStore } from '../store/userStore';
import { useRouter } from 'vue-router';
import { useApiClient } from '../composables/useApiClient';

const userStore = useUserStore();
const router = useRouter();
const { apiClient } = await useApiClient();

const users = ref([] as User[]);

const handleLogin = (u: User) => {
  userStore.currentUser = u;
  router.push({ name: 'home' });
};

onMounted(() => {
  apiClient.getUsers().then((userList) => (users.value = userList));
});
</script>

<style lang="scss" scoped>
@use '../styles/nice';

.users {
  width: 100%;
  height: 100%;

  overflow: hidden;

  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;

  font-family: Pacifico, sans-serif;

  font-size: 2rem;

  .user {
    text-transform: capitalize;

    padding: 0.5rem 1rem;

    margin-bottom: 1.25rem;
    min-width: 240px;

    text-align: center;

    $border-width: 3px;

    @include nice.gradient-border(
      linear-gradient(20deg, rgb(34, 193, 195) 0%, rgb(253, 187, 45) 100%),
      $border-width
    );
  }
}
</style>
