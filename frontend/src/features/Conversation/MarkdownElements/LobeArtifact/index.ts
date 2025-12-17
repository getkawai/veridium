import { FC } from 'react';

import { ARTIFACT_TAG } from '@/const/plugin';

import { createRehypePlugin } from '../rehypePlugin';
import { MarkdownElement, MarkdownElementProps } from '../type';
import Component from './Render';

const AntArtifactElement: MarkdownElement = {
  Component: Component as unknown as FC<MarkdownElementProps>,
  rehypePlugin: createRehypePlugin(ARTIFACT_TAG),
  scope: 'assistant',
  tag: ARTIFACT_TAG,
};

export default AntArtifactElement;
