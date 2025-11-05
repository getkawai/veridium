import { ChatTranslate } from '../common';

export interface ChatMessageExtra {
  fromModel?: string;
  fromProvider?: string;
  // 翻译
  translate?: ChatTranslate | false | null;
}
