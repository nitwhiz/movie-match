import { Environment, getEnvironment } from '../common/Environment';

const currentEnvironment: Promise<Environment> = getEnvironment();

export const useEnvironment = () => {
  return { environment: currentEnvironment };
};
