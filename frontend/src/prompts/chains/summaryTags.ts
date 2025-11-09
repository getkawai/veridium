import { ChatStreamPayload } from '@/types';

export const chainSummaryTags = (content: string, locale: string): Partial<ChatStreamPayload> => ({
  messages: [
    {
      content:
        'You are an assistant skilled at summarizing conversation tags. You need to extract category tags from the user\'s input, separated by commas, with no more than 5 tags, and translate them to the target language. Format requirements are as follows:\nInput: {text as JSON quoted string} [locale]\nOutput: {tags}',
      role: 'system',
    },
    {
      content: `Input: {You are a copywriting master, help me name some design/artistic works. The names need to have literary connotation, focus on refinement and poetic imagery, express the atmosphere and context of the works, making the names both concise and poetic.} [zh-CN]`,
      role: 'user',
    },
    { content: 'naming,writing,creativity', role: 'assistant' },
    {
      content: `Input: {You are a professional translator proficient in Simplified Chinese, and have participated in the translation work of the Chinese versions of The New York Times and The Economist. Therefore, you have a deep understanding of translating news and current affairs articles. I hope you can help me translate the following English news paragraphs into Chinese, with a style similar to the Chinese versions of the aforementioned magazines.} [zh-CN]`,
      role: 'user',
    },
    { content: 'translation,writing,copywriting', role: 'assistant' },
    {
      content: `Input: {You are an entrepreneurship plan writing expert who can provide plan generation including creative names, short slogans, target user personas, user pain points, main value propositions, sales/marketing channels, revenue streams, cost structures, etc.} [en-US]`,
      role: 'user',
    },
    { content: 'entrepreneurship,planning,consulting', role: 'assistant' },
    { content: `Input: {${content}} [${locale}]`, role: 'user' },
  ],
});
