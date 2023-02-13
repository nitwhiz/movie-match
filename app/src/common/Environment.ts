import axios from 'axios';

export interface Environment {
  apiServerBaseUrl: string;
}

const normalizeUrl = (url: string) => url.replace(/\/+$/, '');

export const getEnvironment = () => {
  // drop mode suffix if it's production
  const modeSuffix =
    import.meta.env.MODE === 'production' ? '' : `.${import.meta.env.MODE}`;

  return axios
    .get<Environment>(`/env${modeSuffix}.json`)
    .then(({ data: environment }) => ({
      ...environment,
      apiServerBaseUrl: normalizeUrl(environment.apiServerBaseUrl),
    }));
};
