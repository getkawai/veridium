/**
 * Database types and utilities
 * 
 * This file provides clean exports and type utilities for the generated database bindings.
 * All types are generated from Go code via Wails bindings.
 */

// Re-export all models from generated bindings
export * from '@@/github.com/kawai-network/veridium/internal/database/generated/models';

// Re-export queries
export * as DB from '@@/github.com/kawai-network/veridium/internal/database/generated/queries';

// Re-export database/sql types
export * from '@@/database/sql/models';

/**
 * Type utilities for working with nullable fields
 */

import type { NullString, NullInt64 } from '@@/database/sql/models';

/**
 * Extract the value from a NullString, returning undefined if not valid
 */
export function getNullableString(ns: NullString | undefined): string | undefined {
  if (!ns?.Valid) return undefined;
  // Ensure String is actually a string, not an object
  const str = ns.String;
  return typeof str === 'string' ? str : undefined;
}

/**
 * Extract the value from a NullInt64, returning undefined if not valid
 */
export function getNullableInt(ni: NullInt64 | undefined): number | undefined {
  if (!ni?.Valid) return undefined;
  // Ensure Int64 is actually a number
  const num = ni.Int64;
  return typeof num === 'number' ? num : undefined;
}

/**
 * Extract the value from a nullable boolean, returning undefined if not valid
 */
export function getNullableBool(nb: { Bool: boolean; Valid: boolean } | undefined): boolean | undefined {
  return nb?.Valid ? nb.Bool : undefined;
}

/**
 * Create a NullString from a string value
 * If the value is already a NullString, return it as-is to prevent double-wrapping
 */
export function toNullString(value: string | undefined | null | NullString): NullString {
  // Check if already a NullString to prevent double-wrapping
  if (value && typeof value === 'object' && 'String' in value && 'Valid' in value) {
    return value as NullString;
  }
  if (value === undefined || value === null || value === '') {
    return { String: '', Valid: false };
  }
  return { String: value as string, Valid: true };
}

/**
 * Create a NullInt64 from a number value
 * If the value is already a NullInt64, return it as-is to prevent double-wrapping
 */
export function toNullInt(value: number | undefined | null | NullInt64): NullInt64 {
  // Check if already a NullInt64 to prevent double-wrapping
  if (value && typeof value === 'object' && 'Int64' in value && 'Valid' in value) {
    return value as NullInt64;
  }
  if (value === undefined || value === null) {
    return { Int64: 0, Valid: false };
  }
  return { Int64: value as number, Valid: true };
}

/**
 * Create a nullable boolean from a boolean value
 */
export function toNullBool(value: boolean | undefined | null): { Bool: boolean; Valid: boolean } {
  if (value === undefined || value === null) {
    return { Bool: false, Valid: false };
  }
  return { Bool: value, Valid: true };
}

/**
 * Parse JSON from a NullString
 */
export function parseJSON(value: NullString | string | undefined | null): any {
  if (!value) return undefined;
  
  const jsonStr = typeof value === 'string' ? value : getNullableString(value);
  
  if (!jsonStr) return undefined;
  
  try {
    return JSON.parse(jsonStr);
  } catch (e) {
    console.error('Failed to parse JSON:', e);
    return undefined;
  }
}

/**
 * Convert object to JSON NullString
 */
export function toNullJSONString(value: any): NullString {
  if (value === undefined || value === null) {
    return { String: '', Valid: false };
  }
  return { String: JSON.stringify(value), Valid: true };
}

/**
 * Parse JSON string from NullString
 */
export function parseNullableJSON<T = any>(ns: NullString | undefined): T | undefined {
  const str = getNullableString(ns);
  if (!str) return undefined;
  try {
    return JSON.parse(str) as T;
  } catch {
    return undefined;
  }
}

/**
 * Stringify object to NullString JSON
 */
export function toNullJSON(value: any): NullString {
  if (value === undefined || value === null) {
    return { String: '', Valid: false };
  }
  return { String: JSON.stringify(value), Valid: true };
}

/**
 * Convert boolean to integer (SQLite compatibility)
 */
export function boolToInt(value: boolean): number {
  return value ? 1 : 0;
}

/**
 * Convert integer to boolean (SQLite compatibility)
 */
export function intToBool(value: number): boolean {
  return value === 1;
}

/**
 * Get current timestamp in milliseconds (SQLite format)
 */
export function currentTimestampMs(): number {
  return Date.now();
}

