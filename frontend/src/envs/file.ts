const DEFAULT_S3_FILE_PATH = 'files';

export const getFileConfig = () => ({
  CHUNKS_AUTO_EMBEDDING: true,
  CHUNKS_AUTO_GEN_METADATA: true,

  NEXT_PUBLIC_S3_DOMAIN: undefined,
  NEXT_PUBLIC_S3_FILE_PATH: DEFAULT_S3_FILE_PATH,

  S3_ACCESS_KEY_ID: undefined,
  S3_BUCKET: undefined,
  S3_ENABLE_PATH_STYLE: false,
  S3_ENDPOINT: undefined,
  S3_PREVIEW_URL_EXPIRE_IN: 7200,
  S3_PUBLIC_DOMAIN: undefined,
  S3_REGION: undefined,
  S3_SECRET_ACCESS_KEY: undefined,
  S3_SET_ACL: true,
});

export const fileEnv = getFileConfig();
