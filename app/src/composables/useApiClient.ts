import ApiClient from '../api/ApiClient';
import { useEnvironment } from './useEnvironment';

let currentApiClient: ApiClient | null = null;

export const useApiClient = async () => {
  const { env } = await useEnvironment();

  return {
    apiClient: currentApiClient
      ? currentApiClient
      : new ApiClient(env.apiServerBaseUrl),
  };
};
