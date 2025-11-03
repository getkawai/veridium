import {
  EdgeSpeechOptions,
  MicrosoftSpeechOptions,
  OpenAITTSOptions,
  TTSOptions,
  useEdgeSpeech,
  useMicrosoftSpeech,
  useOpenAITTS,
} from '@lobehub/tts/react';
import isEqual from 'fast-deep-equal';
import { useCallback, useEffect, useRef, useState } from 'react';

import { createHeaderWithOpenAI } from '@/services/_header';
import { API_ENDPOINTS } from '@/services/_url';
import { useAgentStore } from '@/store/agent';
import { agentSelectors } from '@/store/agent/slices/chat';
import { useGlobalStore } from '@/store/global';
import { globalGeneralSelectors } from '@/store/global/selectors';
import { useUserStore } from '@/store/user';
import { settingsSelectors } from '@/store/user/selectors';
import { TTSServer } from '@/types/agent';

interface TTSConfig extends TTSOptions {
  onUpload?: (currentVoice: string, arraybuffers: ArrayBuffer[]) => void;
  server?: TTSServer;
  voice?: string;
}

// Hook untuk Web Speech API (Browser TTS)
const useBrowserTTS = (content: string, config?: TTSConfig) => {
  const [isGlobalLoading, setIsGlobalLoading] = useState(false);
  const [audio, setAudio] = useState<string>();
  const utteranceRef = useRef<SpeechSynthesisUtterance | null>(null);
  const [response] = useState<Response | undefined>(undefined);

  const start = useCallback(() => {
    if (!window.speechSynthesis) {
      const error = new Error('Speech synthesis not supported');
      config?.onError?.(error, '', { retryCount: 0, dedupe: false } as any);
      return;
    }

    // Cancel any ongoing speech
    window.speechSynthesis.cancel();

    const utterance = new SpeechSynthesisUtterance(content);
    utteranceRef.current = utterance;

    // Set voice if specified
    if (config?.voice) {
      const voices = window.speechSynthesis.getVoices();
      const selectedVoice = voices.find((v) => v.name === config.voice || v.lang === config.voice);
      if (selectedVoice) {
        utterance.voice = selectedVoice;
      }
    }

    utterance.onstart = () => {
      setIsGlobalLoading(true);
      // Create a dummy audio URL for compatibility
      setAudio('browser-tts://playing');
    };

    utterance.onend = () => {
      setIsGlobalLoading(false);
      config?.onSuccess?.(undefined as any, '', { retryCount: 0, dedupe: false } as any);
    };

    utterance.onerror = (event) => {
      setIsGlobalLoading(false);
      const error = new Error(event.error || 'Speech synthesis error');
      config?.onError?.(error, '', { retryCount: 0, dedupe: false } as any);
    };

    setIsGlobalLoading(true);
    window.speechSynthesis.speak(utterance);
  }, [content, config]);

  const stop = useCallback(() => {
    window.speechSynthesis.cancel();
    setIsGlobalLoading(false);
    utteranceRef.current = null;
  }, []);

  const setText = useCallback(() => {
    // For browser TTS, we don't need to change text dynamically
    // as the content is passed directly to the utterance
  }, []);

  useEffect(() => {
    return () => {
      window.speechSynthesis.cancel();
    };
  }, []);

  return {
    audio,
    isGlobalLoading,
    response,
    start,
    stop,
    setText,
  };
};

export const useTTS = (content: string, config?: TTSConfig) => {
  const ttsSettings = useUserStore(settingsSelectors.currentTTS, isEqual);
  const ttsAgentSettings = useAgentStore(agentSelectors.currentAgentTTS, isEqual);
  const lang = useGlobalStore(globalGeneralSelectors.currentLanguage);
  const voice = useAgentStore(agentSelectors.currentAgentTTSVoice(lang));
  
  const selectedServer = config?.server || ttsAgentSettings.ttsService;

  // Use browser TTS if server is 'browser'
  if (selectedServer === 'browser') {
    return useBrowserTTS(content, config);
  }

  let useSelectedTTS;
  let options: any = {};
  switch (selectedServer) {
    case 'openai': {
      useSelectedTTS = useOpenAITTS;
      options = {
        api: {
          headers: createHeaderWithOpenAI(),
          serviceUrl: API_ENDPOINTS.tts,
        },
        options: {
          model: ttsSettings.openAI.ttsModel,
          voice: config?.voice || voice,
        },
      } as OpenAITTSOptions;
      break;
    }
    case 'edge': {
      useSelectedTTS = useEdgeSpeech;
      options = {
        api: {
          /**
           * @description client fetch
           * serviceUrl: TTS_URL.edge,
           */
        },
        options: {
          voice: config?.voice || voice,
        },
      } as EdgeSpeechOptions;
      break;
    }
    case 'microsoft': {
      useSelectedTTS = useMicrosoftSpeech;
      options = {
        api: {
          serviceUrl: API_ENDPOINTS.microsoft,
        },
        options: {
          voice: config?.voice || voice,
        },
      } as MicrosoftSpeechOptions;
      break;
    }
    default: {
      // Default to OpenAI TTS if service is not recognized
      useSelectedTTS = useOpenAITTS;
      options = {
        api: {
          headers: createHeaderWithOpenAI(),
          serviceUrl: API_ENDPOINTS.tts,
        },
        options: {
          model: ttsSettings.openAI.ttsModel,
          voice: config?.voice || voice,
        },
      } as OpenAITTSOptions;
      break;
    }
  }

  return useSelectedTTS(content, {
    ...config,
    ...options,
    onFinish: (arraybuffers) => {
      config?.onUpload?.(options.voice || 'alloy', arraybuffers);
    },
  });
};
