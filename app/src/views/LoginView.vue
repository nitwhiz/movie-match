<template>
  <div class="login">
    <h1>Welcome</h1>
    <div class="error" v-if="showError">Password or username incorrect!</div>
    <form>
      <div class="input">
        <input
          type="text"
          autocomplete="username"
          placeholder="Username"
          v-model="username"
        />
      </div>
      <div class="input">
        <input
          type="password"
          autocomplete="current-password"
          placeholder="Password"
          v-model="password"
        />
      </div>
      <div class="button" @click="login">
        <span class="text">Log-In</span>
      </div>
    </form>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useApiClient } from '../composables/useApiClient';
import { useRouter } from 'vue-router';
import { useCurrentUser } from '../composables/useCurrentUser';
import { RouteName } from '../router';

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
@use '../styles/nice';

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
    text-transform: capitalize;

    margin-bottom: 1.25rem;

    width: 100%;
    max-width: 280px;

    text-align: center;

    @include nice.gradient-border(
      linear-gradient(20deg, rgb(34, 193, 195) 0%, rgb(253, 187, 45) 100%),
      3px
    );

    input {
      display: block;

      width: 100%;

      border-radius: 14px;

      font-size: 1.5rem;
      padding: 0.5rem 0.75rem;
    }
  }

  .button {
    font-family: Pacifico, sans-serif;

    padding: 0.25rem;

    text-transform: capitalize;

    margin-bottom: 1.25rem;

    width: 100%;
    max-width: 280px;

    text-align: center;

    border-radius: 14px;

    background: linear-gradient(
      20deg,
      rgb(34, 193, 195) 0%,
      rgb(253, 187, 45) 100%
    );

    color: black;

    .text {
      font-size: 1.5rem;
    }
  }
}
</style>
