import ApiClient from '../api/ApiClient';
import { API_SERVER_BASE_URL } from '../common/Environment';

const apiClient = new ApiClient(API_SERVER_BASE_URL);

export const useApiClient = () => apiClient;
