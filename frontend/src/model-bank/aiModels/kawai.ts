import { AIChatModelCard } from '../types/aiModel';

const kawaiChatModels: AIChatModelCard[] = [
  {
    abilities: {
      functionCall: true,
      vision: true,
      files: true,
    },
    contextWindowTokens: 128_000,
    description:
      'Kawai AI automatically selects the best local model for your hardware. Powered by llama.cpp with GPU acceleration.',
    displayName: 'Kawai Auto',
    enabled: true,
    id: 'kawai-auto',
    type: 'chat',
  },
];

export const allModels = [...kawaiChatModels];

export default allModels;

