import { ChatStreamPayload } from '@/types';

export const chainSummaryGenerationTitle = (
  prompts: string[],
  modal: 'image' | 'video',
  locale: string,
): Partial<ChatStreamPayload> => {
  // Format multiple prompts for better readability
  const formattedPrompts = prompts.map((prompt, index) => `${index + 1}. ${prompt}`).join('\n');

  return {
    messages: [
      {
        content: `You are an expert AI art creator and writer. Summarize a title from the user's AI ${modal} prompt. The title must concisely describe the core idea, will be used to label and manage this series, stay within 10 characters, avoid punctuation, and output in: ${locale}.`,
        role: 'system',
      },
      {
        content: `Prompts:\n${formattedPrompts}`,
        role: 'user',
      },
    ],
  };
};
