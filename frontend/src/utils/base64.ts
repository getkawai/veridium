/**
 * Cross-platform base64 encoding utility
 * Works in both browser and Node.js environments
 */

/**
 * Encode a string to base64
 * @param input - The string to encode
 * @returns Base64 encoded string
 */
export const encodeToBase64 = (input: string): string => {
    // Browser environment
    return btoa(input);
};

/**
 * Decode a base64 string
 * @param input - The base64 string to decode
 * @returns Decoded string
 */
export const decodeFromBase64 = (input: string): string => {
    // Browser environment
    return atob(input);
};

/**
 * Create Basic Authentication header value
 * @param username - Username for authentication
 * @param password - Password for authentication
 * @returns Base64 encoded credentials for Basic auth
 */
export const createBasicAuthCredentials = (username: string, password: string): string => {
  return encodeToBase64(`${username}:${password}`);
};
