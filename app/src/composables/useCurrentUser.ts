import { useApiClient } from './useApiClient';
import { User } from '../model/User';
import { ref } from 'vue';

const currentUser = ref(null as User | null);

let userPromise: Promise<void> = Promise.resolve();
let userPromiseFulfilled = true;

const load = async () => {
  if (currentUser.value === null && userPromiseFulfilled) {
    userPromiseFulfilled = false;

    userPromise = new Promise<void>(async (resolve, reject) => {
      try {
        const apiClient = await useApiClient().apiClient;

        currentUser.value = await (await apiClient).me();

        resolve();
      } catch (e) {
        console.error('unable to retrieve /me endpoint');
        reject();
      }
    }).then(
      () => {
        userPromiseFulfilled = true;
      },
      () => {
        userPromiseFulfilled = true;
      }
    );
  }

  await userPromise;
};

export const useCurrentUser = () => {
  return {
    currentUser,
    load,
  };
};
