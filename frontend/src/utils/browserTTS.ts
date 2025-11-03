/**
 * Utility functions for Browser TTS (Web Speech API)
 */

export interface BrowserVoice {
  name: string;
  lang: string;
  localService: boolean;
  default: boolean;
}

/**
 * Get available browser voices
 * @returns Promise with array of available voices
 */
export const getBrowserVoices = (): Promise<BrowserVoice[]> => {
  return new Promise((resolve) => {
    if (!window.speechSynthesis) {
      resolve([]);
      return;
    }

    const getVoices = () => {
      const voices = window.speechSynthesis.getVoices();
      if (voices.length > 0) {
        resolve(
          voices.map((voice) => ({
            name: voice.name,
            lang: voice.lang,
            localService: voice.localService,
            default: voice.default,
          }))
        );
      }
    };

    // Try to get voices immediately
    getVoices();

    // Some browsers need to wait for voiceschanged event
    if (window.speechSynthesis.onvoiceschanged !== undefined) {
      window.speechSynthesis.onvoiceschanged = getVoices;
    }

    // Fallback timeout
    setTimeout(() => {
      getVoices();
    }, 100);
  });
};

/**
 * Get browser voices filtered by language
 * @param lang - Language code (e.g., 'en-US', 'id-ID')
 * @returns Promise with array of voices for the specified language
 */
export const getBrowserVoicesByLanguage = async (lang: string): Promise<BrowserVoice[]> => {
  const voices = await getBrowserVoices();
  return voices.filter((voice) => voice.lang.startsWith(lang));
};

/**
 * Check if browser TTS is supported
 * @returns true if browser supports Speech Synthesis API
 */
export const isBrowserTTSSupported = (): boolean => {
  return 'speechSynthesis' in window;
};

