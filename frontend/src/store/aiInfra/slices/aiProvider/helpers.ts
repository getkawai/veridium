/**
 * Helper utilities for AI Provider store
 * Used to simplify direct DB calls and reduce duplication
 */

const DEFAULT_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

/**
 * Parse NullString from Wails binding
 */
export const parseNullString = (value: any): string => {
  if (!value) return '';
  if (typeof value === 'string') return value;
  if (value.Valid && value.String) return value.String;
  return '';
};

/**
 * Parse NullJSON from Wails binding
 */
export const parseNullJSON = <T = any>(value: any, defaultValue: T = {} as T): T => {
  try {
    const str = parseNullString(value);
    return str ? JSON.parse(str) : defaultValue;
  } catch (error) {
    console.warn('[parseNullJSON] Failed to parse:', value, error);
    return defaultValue;
  }
};

/**
 * Parse NullInt64 from Wails binding
 */
export const parseNullInt64 = (value: any): number => {
  if (!value) return 0;
  if (typeof value === 'number') return value;
  if (value.Valid && value.Int64 !== undefined) return Number(value.Int64);
  return 0;
};

/**
 * Create NullString for Wails binding
 */
export const toNullString = (value: string | null | undefined) => {
  return value ? { String: value, Valid: true } : { String: '', Valid: false };
};

/**
 * Create NullJSON for Wails binding
 */
export const toNullJSON = (value: any) => {
  return value ? { String: JSON.stringify(value), Valid: true } : { String: '', Valid: false };
};

/**
 * Create NullInt64 for Wails binding
 */
export const toNullInt64 = (value: number | null | undefined) => {
  return value !== null && value !== undefined
    ? { Int64: value, Valid: true }
    : { Int64: 0, Valid: false };
};

/**
 * Convert boolean to int for SQLite
 */
export const boolToInt = (value: boolean): number => {
  return value ? 1 : 0;
};

/**
 * Map AI Provider from DB result
 */
export const mapProviderFromDB = (p: any) => ({
  id: p.id,
  name: parseNullString(p.name),
  enabled: Boolean(parseNullInt64(p.enabled)),
  sort: parseNullInt64(p.sort),
  source: parseNullString(p.source) || 'builtin',
  logo: parseNullString(p.logo),
  description: parseNullString(p.description),
  keyVaults: parseNullJSON(p.keyVaults),
  settings: parseNullJSON(p.settings),
  config: parseNullJSON(p.config),
  fetchOnClient: Boolean(parseNullInt64(p.fetchOnClient)),
  checkModel: parseNullString(p.checkModel),
});

/**
 * Map AI Model from DB result
 */
export const mapModelFromDB = (m: any) => ({
  id: m.id,
  displayName: parseNullString(m.displayName),
  providerId: m.providerId,
  type: m.type,
  enabled: Boolean(parseNullInt64(m.enabled)),
  sort: parseNullInt64(m.sort),
  abilities: parseNullJSON(m.abilities, {}),
  contextWindowTokens: parseNullInt64(m.contextWindowTokens),
  description: parseNullString(m.description),
  parameters: parseNullJSON(m.parameters, {}),
  config: parseNullJSON(m.config, {}),
  organization: parseNullString(m.organization),
  pricing: parseNullJSON(m.pricing),
  source: parseNullString(m.source),
  releasedAt: parseNullString(m.releasedAt),
});

/**
 * Map AI Provider runtime config from DB result
 */
export const mapRuntimeConfigFromDB = (c: any) => ({
  id: c.id,
  keyVaults: parseNullJSON(c.keyVaults, {}),
  settings: parseNullJSON(c.settings, {}),
  config: parseNullJSON(c.config, {}),
  fetchOnClient: Boolean(parseNullInt64(c.fetchOnClient)),
});

/**
 * Get default user ID for desktop single-user app
 */
export const getUserId = (): string => DEFAULT_USER_ID;

/**
 * Build provider-model lists grouped by provider
 */
export const groupModelsByProvider = (
  providers: any[],
  models: any[],
  type: 'chat' | 'image',
) => {
  return providers
    .filter((provider) => models.some((m) => m.providerId === provider.id && m.type === type))
    .map((provider) => ({
      ...provider,
      models: models
        .filter((m) => m.providerId === provider.id && m.type === type)
        .sort((a, b) => (a.sort || 0) - (b.sort || 0)),
    }));
};

