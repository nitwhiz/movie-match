import { defineStore } from 'pinia';
import { User } from '../model/User';
import { ref } from 'vue';

export const useUserStore = defineStore('user', () => {
  const currentUser = ref(null as User | null);

  const loadCurrentUser = async () => {
    // todo: bring back user storage via localforage
    // todo: check user existence by id via api
    return currentUser.value !== null;
  };

  return { currentUser, loadCurrentUser };
});
