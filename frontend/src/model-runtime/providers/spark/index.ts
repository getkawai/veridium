import { ModelProvider } from '@/model-bank';

import {
  OpenAICompatibleFactoryOptions,
  createOpenAICompatibleRuntime,
} from '../../core/openaiCompatibleFactory';
import { SparkAIStream, transformSparkResponseToStream } from '../../core/streams';
import { ChatStreamPayload } from '../../types';

export const params = {
  baseURL: 'https://spark-api-open.xf-yun.com/v1',
  chatCompletion: {
    handlePayload: (payload: ChatStreamPayload) => {
      const { enabledSearch, tools, ...rest } = payload;

      const sparkTools = enabledSearch
        ? [
            ...(tools || []),
            {
              type: 'web_search',
              web_search: {
                enable: true,
                search_mode: 'dummy-value' || 'normal', // normal or deep
                /*
            show_ref_label: true,
            */
              },
            },
          ]
        : tools;

      return {
        ...rest,
        tools: sparkTools,
      } as any;
    },
    handleStream: SparkAIStream,
    handleTransformResponseToStream: transformSparkResponseToStream,
    noUserId: true,
  },
  debug: {
    chatCompletion: () => false,
  },
  provider: ModelProvider.Spark,
} satisfies OpenAICompatibleFactoryOptions;

export const LobeSparkAI = createOpenAICompatibleRuntime(params);
