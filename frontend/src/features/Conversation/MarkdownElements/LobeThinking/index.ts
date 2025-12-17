import { ARTIFACT_THINKING_TAG } from '@/const/plugin';

import { createRehypePlugin } from '../rehypePlugin';
import { createRemarkCustomTagPlugin } from '../remarkPlugins/createRemarkCustomTagPlugin';
import { MarkdownElement } from '../type';
import Component from './Render';

const LobeThinkingElement: MarkdownElement = {
  Component,
  rehypePlugin: createRehypePlugin(ARTIFACT_THINKING_TAG),
  remarkPlugin: createRemarkCustomTagPlugin(ARTIFACT_THINKING_TAG),
  scope: 'assistant',
  tag: ARTIFACT_THINKING_TAG,
};

export default LobeThinkingElement;
