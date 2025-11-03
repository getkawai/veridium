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

import type { NullString, NullInt64, NullBool } from '@@/database/sql/models';

/**
 * Extract the value from a NullString, returning undefined if not valid
 */
export function getNullableString(ns: NullString | undefined): string | undefined {
  return ns?.Valid ? ns.String : undefined;
}

/**
 * Extract the value from a NullInt64, returning undefined if not valid
 */
export function getNullableInt(ni: NullInt64 | undefined): number | undefined {
  return ni?.Valid ? ni.Int64 : undefined;
}

/**
 * Extract the value from a NullBool, returning undefined if not valid
 */
export function getNullableBool(nb: NullBool | undefined): boolean | undefined {
  return nb?.Valid ? nb.Bool : undefined;
}

/**
 * Create a NullString from a string value
 */
export function toNullString(value: string | undefined | null): NullString {
  if (value === undefined || value === null || value === '') {
    return { String: '', Valid: false };
  }
  return { String: value, Valid: true };
}

/**
 * Create a NullInt64 from a number value
 */
export function toNullInt(value: number | undefined | null): NullInt64 {
  if (value === undefined || value === null) {
    return { Int64: 0, Valid: false };
  }
  return { Int64: value, Valid: true };
}

/**
 * Create a NullBool from a boolean value
 */
export function toNullBool(value: boolean | undefined | null): NullBool {
  if (value === undefined || value === null) {
    return { Bool: false, Valid: false };
  }
  return { Bool: value, Valid: true };
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

