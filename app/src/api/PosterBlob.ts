import { useApiClient } from '../composables/useApiClient';

const postersByMediaId: Record<
  string,
  { usages: number; urlPromise: Promise<string> }
> = {};

const apiClient = useApiClient().apiClient;

export const getMediaPosterBlobUrl = async (mediaId: string) => {
  if (!postersByMediaId[mediaId]) {
    postersByMediaId[mediaId] = {
      usages: 0,
      urlPromise: (await apiClient).getPosterBlobUrl(mediaId),
    };
  }

  ++postersByMediaId[mediaId].usages;

  return postersByMediaId[mediaId].urlPromise;
};

const free = (mediaId: string) => {
  postersByMediaId[mediaId].urlPromise.then((url) => {
    if (postersByMediaId[mediaId].usages === 0) {
      delete postersByMediaId[mediaId];
      URL.revokeObjectURL(url);
    }
  });
};

export const freeMediaPosterBlobUrl = (mediaId: string) => {
  --postersByMediaId[mediaId].usages;

  free(mediaId);
};

export const freeAllMediaBlobUrls = () => {
  for (const mediaId of Object.keys(postersByMediaId)) {
    free(mediaId);
  }
};
