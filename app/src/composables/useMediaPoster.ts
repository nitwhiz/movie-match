import { ref, watch } from 'vue';
import Poster from '../common/Poster';
import { useApiClient } from './useApiClient';

export const useMediaPoster = (mediaIdCallback: () => string) => {
  const poster = ref(null as Poster | null);
  const posterUrl = ref(null as string | null);

  watch(
    mediaIdCallback,
    (mediaId) => {
      useApiClient().apiClient.then((apiClient) => {
        const currentPoster = poster.value;
        const newPoster = Poster.getByMediaId(apiClient, mediaId);

        newPoster
          .getUrl()
          .then((url) => {
            poster.value = newPoster;
            posterUrl.value = url;
          })
          .catch((e) => {
            console.error(e);

            poster.value = null;
            posterUrl.value = null;
          })
          .finally(() => {
            if (currentPoster) {
              currentPoster.free();
            }
          });
      });
    },
    {
      immediate: true,
    }
  );

  return {
    posterUrl,
  };
};
