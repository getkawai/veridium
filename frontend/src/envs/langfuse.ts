/* eslint-disable sort-keys-fix/sort-keys-fix , typescript-sort-keys/interface */

export const getLangfuseConfig = () => ({
  ENABLE_LANGFUSE: false,
  LANGFUSE_SECRET_KEY: '',
  LANGFUSE_PUBLIC_KEY: '',
  LANGFUSE_HOST: 'https://cloud.langfuse.com',
});

export const langfuseEnv = getLangfuseConfig();
