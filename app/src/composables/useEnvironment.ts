import { Environment, getEnvironment } from '../common/Environment';

let currentEnv: Environment | null = null;

export const useEnvironment = async () => ({
  env: currentEnv
    ? currentEnv
    : await getEnvironment().then((fetchedEnv) => {
        currentEnv = fetchedEnv;

        return currentEnv;
      }),
});
