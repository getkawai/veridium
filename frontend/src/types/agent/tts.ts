export type TTSServer = 'openai' | 'edge' | 'microsoft' | 'browser';

export interface LobeAgentTTSConfig {
  showAllLocaleVoice?: boolean;
  sttLocale: 'auto' | string;
  ttsService: TTSServer;
  voice: {
    browser?: string;
    edge?: string;
    microsoft?: string;
    openai: string;
  };
}
