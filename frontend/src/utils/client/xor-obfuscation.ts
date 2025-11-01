import { SECRET_XOR_KEY } from '@/const/auth';

/**
 * Simple XOR obfuscation function for payload data
 * Uses the SECRET_XOR_KEY to obfuscate JSON payloads
 */
export function obfuscatePayloadWithXOR<T = any>(payload: T): string {
  try {
    const jsonString = JSON.stringify(payload);
    const key = SECRET_XOR_KEY;
    let result = '';

    for (let i = 0; i < jsonString.length; i++) {
      const charCode = jsonString.charCodeAt(i) ^ key.charCodeAt(i % key.length);
      result += String.fromCharCode(charCode);
    }

    // Convert to base64 for safe transport
    return btoa(result);
  } catch (error) {
    console.error('Failed to obfuscate payload:', error);
    throw error;
  }
}

/**
 * Deobfuscate a payload that was obfuscated with obfuscatePayloadWithXOR
 */
export function deobfuscatePayloadWithXOR<T = any>(obfuscatedPayload: string): T {
  try {
    // Decode from base64
    const xorString = atob(obfuscatedPayload);
    const key = SECRET_XOR_KEY;
    let result = '';

    for (let i = 0; i < xorString.length; i++) {
      const charCode = xorString.charCodeAt(i) ^ key.charCodeAt(i % key.length);
      result += String.fromCharCode(charCode);
    }

    return JSON.parse(result);
  } catch (error) {
    console.error('Failed to deobfuscate payload:', error);
    throw error;
  }
}
