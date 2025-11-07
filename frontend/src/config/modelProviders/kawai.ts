import { ModelProviderCard } from '@/types/llm';

const Kawai: ModelProviderCard = {
  chatModels: [
    {
      abilities: {
        functionCall: true,
        vision: false,
      },
      contextWindowTokens: 128_000,
      description:
        'Kawai AI automatically selects the best local model for your hardware. Powered by llama.cpp with GPU acceleration.',
      displayName: 'Kawai Auto',
      enabled: true,
      id: 'kawai-auto',
      type: 'chat',
    },
  ],
  description:
    'Kawai AI - Local LLM inference powered by llama.cpp. Automatically selects and runs the best model for your hardware.',
  id: 'kawai',
  modelsUrl: 'https://huggingface.co/models?library=gguf',
  name: 'Kawai',
  settings: {
    defaultShowBrowserRequest: true,
    proxyUrl: {
      placeholder: 'http://127.0.0.1:8080/v1',
    },
    responseAnimation: {
      speed: 2,
      text: 'smooth',
    },
    showApiKey: false,
    showModelFetcher: false, // Models handled automatically by backend
  },
  url: 'https://github.com/kawai-network/veridium',
};

export default Kawai;

