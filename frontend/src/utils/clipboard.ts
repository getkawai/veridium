/**
 * Centralized clipboard utilities for Wails desktop app
 * 
 * Uses Wails native Clipboard API instead of navigator.clipboard
 * for better compatibility with native OS webview.
 */

import { Clipboard } from '@wailsio/runtime';

/**
 * Copy text to clipboard using Wails native API
 * 
 * @param text - Text to copy
 * @returns Promise that resolves when copy is successful
 */
export async function copyText(text: string): Promise<void> {
  try {
    await Clipboard.SetText(text);
  } catch (error) {
    console.error('Failed to copy text:', error);
    throw error;
  }
}

/**
 * Get text from clipboard using Wails native API
 * 
 * @returns Promise that resolves with clipboard text
 */
export async function getText(): Promise<string> {
  try {
    return await Clipboard.Text();
  } catch (error) {
    console.error('Failed to get clipboard text:', error);
    throw error;
  }
}

/**
 * Legacy image clipboard support
 * Note: Wails v3 clipboard API only supports text.
 * For images, we still use navigator.clipboard as fallback.
 */
const copyUsingFallback = (imageUrl: string) => {
  const img = new Image();
  img.addEventListener('load', function () {
    const canvas = document.createElement('canvas');
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext('2d');
    ctx!.drawImage(img, 0, 0);

    try {
      canvas.toBlob(function (blob) {
        // @ts-ignore
        const item = new ClipboardItem({ 'image/png': blob });
        navigator.clipboard.write([item]).then(function () {
          console.log('Image copied to clipboard successfully using canvas and modern API');
        });
      });
    } catch {
      // 如果 toBlob 或 ClipboardItem 不被支持，使用 data URL
      const dataURL = canvas.toDataURL('image/png');
      const textarea = document.createElement('textarea');
      textarea.value = dataURL;
      document.body.append(textarea);
      textarea.select();

      document.execCommand('copy');

      textarea.remove();
    }
  });
  img.src = imageUrl;
};

const copyUsingModernAPI = async (imageUrl: string) => {
  try {
    const base64Response = await fetch(imageUrl);
    const blob = await base64Response.blob();
    const item = new ClipboardItem({ 'image/png': blob });
    await navigator.clipboard.write([item]);
  } catch (error) {
    console.error('Failed to copy image using modern API:', error);
    copyUsingFallback(imageUrl);
  }
};

/**
 * Copy image to clipboard
 * Note: Uses navigator.clipboard as Wails clipboard API only supports text
 * 
 * @param imageUrl - URL of image to copy
 */
export const copyImageToClipboard = async (imageUrl: string) => {
  // 检查是否支持现代 Clipboard API
  if (navigator.clipboard && 'write' in navigator.clipboard) {
    await copyUsingModernAPI(imageUrl);
  } else {
    copyUsingFallback(imageUrl);
  }
};

