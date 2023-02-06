import localForage from 'localforage';

export const SECONDS = 1000;
export const MINUTES = SECONDS * 60;
export const HOURS = MINUTES * 60;
export const DAY = HOURS * 24;

export const DEFAULT_VALIDITY = 2 * HOURS;

const cacheStorage = localForage.createInstance({
  name: 'movie_match',
  storeName: 'cache',
  version: 1,
});

interface CacheItem<DataType = any> {
  data: DataType;
  validUntil: number;
}

type CacheCallback<ReturnType = any> = () => Promise<ReturnType>;

export const putStored = async <DataType>(
  key: string,
  data: DataType,
  validity: number = DEFAULT_VALIDITY
) => {
  await cacheStorage.setItem<CacheItem<DataType>>(key, {
    validUntil: Date.now() + validity,
    data,
  });
};

export const getStored = async <DataType>(key: string) => {
  const cacheItem = await cacheStorage.getItem<CacheItem<DataType>>(key);

  return cacheItem && cacheItem.validUntil > Date.now() ? cacheItem.data : null;
};

export const getCached = async <DataType>(
  key: string,
  retriever: CacheCallback<DataType>,
  validity: number = DEFAULT_VALIDITY
) => {
  let data: DataType | null;

  const cacheItem = await cacheStorage.getItem<CacheItem<DataType>>(key);

  if (!cacheItem || cacheItem.validUntil < Date.now()) {
    if (import.meta.env.DEV) {
      console.log('item not found or invalid, retriever will run', key);
    }

    data = await retriever();

    await putStored(key, data, validity);
  } else {
    if (import.meta.env.DEV) {
      console.log('item found in cache', key);
    }

    data = cacheItem.data;
  }

  return data;
};
