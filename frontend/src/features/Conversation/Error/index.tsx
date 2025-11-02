// import { AgentRuntimeErrorType, ILobeAgentRuntimeErrorType } from '@/model-runtime';
// import { ChatErrorType, ChatMessageError, ErrorType, UIChatMessage } from '@/types';
// import { IPluginErrorType } from '@/chat-plugin-sdk';
import type { AlertProps } from '@lobehub/ui';
import { Skeleton } from 'antd';
// import dynamic from 'next/dynamic';
import { Suspense, memo, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

// import { useProviderName } from '@/hooks/useProviderName';

// Dummy implementations for UI development

// Dummy types and enums
type IPluginErrorType = string;
type ILobeAgentRuntimeErrorType = string;
type ErrorType = string;
type UIChatMessage = any;
type ChatMessageError = any;

enum ChatErrorType {
  SystemTimeNotMatchError = 'SystemTimeNotMatchError',
  InvalidClerkUser = 'InvalidClerkUser',
  InvalidAccessCode = 'InvalidAccessCode',
}

enum AgentRuntimeErrorType {
  PermissionDenied = 'PermissionDenied',
  InsufficientQuota = 'InsufficientQuota',
  ModelNotFound = 'ModelNotFound',
  QuotaLimitReached = 'QuotaLimitReached',
  ExceededContextWindow = 'ExceededContextWindow',
  LocationNotSupportError = 'LocationNotSupportError',
  OllamaServiceUnavailable = 'OllamaServiceUnavailable',
  NoOpenAIAPIKey = 'NoOpenAIAPIKey',
  ComfyUIServiceUnavailable = 'ComfyUIServiceUnavailable',
  InvalidComfyUIArgs = 'InvalidComfyUIArgs',
  OllamaBizError = 'OllamaBizError',
}

const useProviderName = (provider: string) => {
  const providerNames: Record<string, string> = {
    openai: 'OpenAI',
    anthropic: 'Anthropic',
    google: 'Google',
    azure: 'Azure OpenAI',
    bedrock: 'Amazon Bedrock',
    ollama: 'Ollama',
    default: provider.charAt(0).toUpperCase() + provider.slice(1),
  };

  return providerNames[provider] || providerNames.default;
};

import ChatInvalidAPIKey from './ChatInvalidApiKey';
import ClerkLogin from './ClerkLogin';
import ErrorJsonViewer from './ErrorJsonViewer';
import InvalidAccessCode from './InvalidAccessCode';
import { ErrorActionContainer } from './style';

// Removed Next.js dynamic imports since we don't use Next.js
// const OllamaBizError = dynamic(() => import('./OllamaBizError'), { loading, ssr: false });
// const OllamaSetupGuide = dynamic(() => import('@/features/OllamaSetupGuide'), {
//   loading,
//   ssr: false,
// });

// Dummy components for UI development
const OllamaBizError = ({ ...props }: any) => <Skeleton active />;
const OllamaSetupGuide = ({ ...props }: any) => <Skeleton active />;

// Config for the errorMessage display
const getErrorAlertConfig = (
  errorType?: IPluginErrorType | ILobeAgentRuntimeErrorType | ErrorType,
): AlertProps | undefined => {
  // OpenAIBizError / ZhipuBizError / GoogleBizError / ...
  if (typeof errorType === 'string' && (errorType.includes('Biz') || errorType.includes('Invalid')))
    return {
      extraDefaultExpand: true,
      extraIsolate: true,
      type: 'warning',
    };

  /* ↓ cloud slot ↓ */

  /* ↑ cloud slot ↑ */

  switch (errorType) {
    case ChatErrorType.SystemTimeNotMatchError:
    case AgentRuntimeErrorType.PermissionDenied:
    case AgentRuntimeErrorType.InsufficientQuota:
    case AgentRuntimeErrorType.ModelNotFound:
    case AgentRuntimeErrorType.QuotaLimitReached:
    case AgentRuntimeErrorType.ExceededContextWindow:
    case AgentRuntimeErrorType.LocationNotSupportError: {
      return {
        type: 'warning',
      };
    }

    case AgentRuntimeErrorType.OllamaServiceUnavailable:
    case AgentRuntimeErrorType.NoOpenAIAPIKey:
    case AgentRuntimeErrorType.ComfyUIServiceUnavailable:
    case AgentRuntimeErrorType.InvalidComfyUIArgs: {
      return {
        extraDefaultExpand: true,
        extraIsolate: true,
        type: 'warning',
      };
    }

    default: {
      return undefined;
    }
  }
};

export const useErrorContent = (error: any) => {
  const { t } = useTranslation('error');
  const providerName = useProviderName(error?.body?.provider || '');

  return useMemo<AlertProps | undefined>(() => {
    if (!error) return;
    const messageError = error;

    const alertConfig = getErrorAlertConfig(messageError.type);

    return {
      message: t(`response.${messageError.type}` as any, { provider: providerName }),
      ...alertConfig,
    };
  }, [error]);
};

const ErrorMessageExtra = memo<{ data: UIChatMessage }>(({ data }) => {
  const error = data.error as ChatMessageError;
  if (!error?.type) return;

  switch (error.type) {
    case AgentRuntimeErrorType.OllamaServiceUnavailable: {
      return <OllamaSetupGuide id={data.id} />;
    }

    case AgentRuntimeErrorType.OllamaBizError: {
      return <OllamaBizError {...data} />;
    }

    /* ↓ cloud slot ↓ */

    /* ↑ cloud slot ↑ */

    case ChatErrorType.InvalidClerkUser: {
      return <ClerkLogin id={data.id} />;
    }

    case ChatErrorType.InvalidAccessCode: {
      return <InvalidAccessCode id={data.id} provider={data.error?.body?.provider} />;
    }

    case AgentRuntimeErrorType.NoOpenAIAPIKey: {
      {
        return <ChatInvalidAPIKey id={data.id} provider={data.error?.body?.provider} />;
      }
    }
  }

  if (error.type.toString().includes('Invalid')) {
    return <ChatInvalidAPIKey id={data.id} provider={data.error?.body?.provider} />;
  }

  return <ErrorJsonViewer error={data.error} id={data.id} />;
});

export default memo<{ data: UIChatMessage }>(({ data }) => (
  <Suspense
    fallback={
      <ErrorActionContainer>
        <Skeleton active style={{ width: '100%' }} />
      </ErrorActionContainer>
    }
  >
    <ErrorMessageExtra data={data} />
  </Suspense>
));
