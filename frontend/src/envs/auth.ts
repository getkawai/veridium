/* eslint-disable sort-keys-fix/sort-keys-fix , typescript-sort-keys/interface */

export const getAuthConfig = () => ({
  NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY: undefined,
  NEXT_PUBLIC_ENABLE_CLERK_AUTH: false,
  NEXT_PUBLIC_ENABLE_NEXT_AUTH: false,

  CLERK_SECRET_KEY: undefined,
  CLERK_WEBHOOK_SECRET: undefined,

  NEXT_AUTH_SECRET: undefined,
  NEXT_AUTH_SSO_PROVIDERS: 'auth0',
  NEXT_AUTH_DEBUG: false,
  NEXT_AUTH_SSO_SESSION_STRATEGY: 'jwt',

  AUTH0_CLIENT_ID: undefined,
  AUTH0_CLIENT_SECRET: undefined,
  AUTH0_ISSUER: undefined,

  GITHUB_CLIENT_ID: undefined,
  GITHUB_CLIENT_SECRET: undefined,

  AZURE_AD_CLIENT_ID: undefined,
  AZURE_AD_CLIENT_SECRET: undefined,
  AZURE_AD_TENANT_ID: undefined,

  AUTHENTIK_CLIENT_ID: undefined,
  AUTHENTIK_CLIENT_SECRET: undefined,
  AUTHENTIK_ISSUER: undefined,

  AUTHELIA_CLIENT_ID: undefined,
  AUTHELIA_CLIENT_SECRET: undefined,
  AUTHELIA_ISSUER: undefined,

  CLOUDFLARE_ZERO_TRUST_CLIENT_ID: undefined,
  CLOUDFLARE_ZERO_TRUST_CLIENT_SECRET: undefined,
  CLOUDFLARE_ZERO_TRUST_ISSUER: undefined,

  GENERIC_OIDC_CLIENT_ID: undefined,
  GENERIC_OIDC_CLIENT_SECRET: undefined,
  GENERIC_OIDC_ISSUER: undefined,

  ZITADEL_CLIENT_ID: undefined,
  ZITADEL_CLIENT_SECRET: undefined,
  ZITADEL_ISSUER: undefined,

  LOGTO_CLIENT_ID: undefined,
  LOGTO_CLIENT_SECRET: undefined,
  LOGTO_ISSUER: undefined,
  LOGTO_WEBHOOK_SIGNING_KEY: undefined,

  CASDOOR_WEBHOOK_SECRET: undefined,
});

export const authEnv = getAuthConfig();
