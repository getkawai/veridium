import { useEffect, useState } from 'react';

import { getBrowserVoices } from '@/utils/browserTTS';

export interface BrowserVoiceOption {
  label: string;
  value: string;
}

/**
 * Hook to get available browser voices as select options
 */
export const useBrowserVoices = () => {
  const [voices, setVoices] = useState<BrowserVoiceOption[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadVoices = async () => {
      try {
        const browserVoices = await getBrowserVoices();
        const voiceOptions = browserVoices.map((voice) => ({
          label: `${voice.name} (${voice.lang})`,
          value: voice.name,
        }));
        setVoices(voiceOptions);
      } catch (error) {
        console.error('Failed to load browser voices:', error);
        setVoices([]);
      } finally {
        setLoading(false);
      }
    };

    loadVoices();
  }, []);

  return { voices, loading };
};

