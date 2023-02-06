import { defineStore } from 'pinia';
import { User } from '../model/User';
import { ref, toRaw, watch } from 'vue';
import { getStored, putStored } from '../cache/CacheStorage';

export const useUserStore = defineStore('user', () => {
  const currentUser = ref(null as User | null);

  watch(currentUser, (value) => {
    if (value !== null) {
      putStored('currentUser', toRaw(value));
    }
  });

  const loadCurrentUser = async () => {
    if (!currentUser.value) {
      const storedUser = await getStored<User>('currentUser');

      if (storedUser) {
        currentUser.value = storedUser;
      }
    }

    return currentUser.value !== null;
  };

  return { currentUser, loadCurrentUser };
});
