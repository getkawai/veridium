import { UIChatMessage } from '@/types';
import { Skeleton } from 'antd';
// import dynamic from 'next/dynamic';
import { memo } from 'react';

import ErrorJsonViewer from '../ErrorJsonViewer';
import InvalidOllamaModel from './InvalidOllamaModel';

// Removed Next.js dynamic imports since we don't use Next.js
// const SetupGuide = dynamic(() => import('@/features/OllamaSetupGuide'), { loading, ssr: false });

// Dummy components for UI development
const SetupGuide = ({ ...props }: any) => <Skeleton active style={{ width: 300 }} />;

interface OllamaError {
  code: string | null;
  message: string;
  param?: any;
  type: string;
}

interface OllamaErrorResponse {
  error: OllamaError;
}

const UNRESOLVED_MODEL_REGEXP = /model "([\w+,-_]+)" not found/;

const OllamaBizError = memo<UIChatMessage>(({ error, id }) => {
  const errorBody: OllamaErrorResponse = (error as any)?.body;

  const errorMessage = errorBody.error?.message;

  // error of not pull the model
  const unresolvedModel = errorMessage?.match(UNRESOLVED_MODEL_REGEXP)?.[1];
  if (unresolvedModel) {
    return <InvalidOllamaModel id={id} model={unresolvedModel} />;
  }

  // error of not enable model or not set the CORS rules
  if (errorMessage?.includes('Failed to fetch')) {
    return <SetupGuide id={id} />;
  }

  return <ErrorJsonViewer error={error} id={id} />;
});

export default OllamaBizError;
