import ApiClient from '../api/ApiClient';
import { useEnvironment } from './useEnvironment';

const currentApiClient: Promise<ApiClient> = (async () => {
  const env = await useEnvironment().environment;

  return new ApiClient(env.apiServerBaseUrl).loadAccessTokenFromCookie();
})();

export const useApiClient = () => {
  return { apiClient: currentApiClient };
};
