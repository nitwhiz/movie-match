<template>
  <div class="login">
    <h1>Welcome</h1>
    <div class="error" v-if="showError">Password or username incorrect!</div>
    <form>
      <div class="input">
        <NiceWrapper :colors="['rgb(34, 193, 195)', 'rgb(253, 187, 45)']">
          <input
            type="text"
            autocomplete="username"
            placeholder="Username"
            v-model="username"
          />
        </NiceWrapper>
      </div>
      <div class="input">
        <NiceWrapper :colors="['rgb(34, 193, 195)', 'rgb(253, 187, 45)']">
          <input
            type="password"
            autocomplete="current-password"
            placeholder="Password"
            v-model="password"
          />
        </NiceWrapper>
      </div>
      <NiceWrapper
        :colors="['rgb(34, 193, 195)', 'rgb(253, 187, 45)']"
        @click="login"
      >
        Log-In
      </NiceWrapper>
    </form>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useApiClient } from '../composables/useApiClient';
import { useRouter } from 'vue-router';
import { useCurrentUser } from '../composables/useCurrentUser';
import { RouteName } from '../router';
import NiceWrapper from '../components/nice/NiceWrapper.vue';

const apiClient = await useApiClient().apiClient;
const currentUser = useCurrentUser();
const router = useRouter();

const showError = ref(false);

const username = ref('');
const password = ref('');

const login = () => {
  showError.value = false;

  apiClient
    .login({
      username: username.value,
      password: password.value,
    })
    .then(() => currentUser.load())
    .then(() => {
      router.push({ name: RouteName.HOME });
    })
    .catch(() => {
      showError.value = true;
    });
};
</script>

<style lang="scss" scoped>
.login {
  width: 100%;
  height: 100%;

  overflow: hidden;

  display: flex;
  align-items: center;
  flex-direction: column;

  .error {
    background-color: rgba(255, 0, 0, 0.2);
    border: 1px solid rgba(255, 0, 0, 0.25);
    width: 100%;
    max-width: 280px;
    padding: 1rem;
    border-radius: 8px;

    text-align: center;
    margin-bottom: 2rem;
  }

  h1 {
    margin: 4rem 0 4rem 0;
  }

  .input {
    margin-bottom: 1rem;
  }
}
</style>
