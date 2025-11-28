// import { isProviderDisableBrowserRequest } from '@/config/modelProviders';
import { ChatModelCard } from '@/types/llm';
// import { getMessageError } from '@/utils/fetch';

// import { initializeWithClientStore } from './chat/clientModelRuntime';
// import { resolveRuntimeProvider } from './chat/helper';

// 进度信息接口
export interface ModelProgressInfo {
  completed?: number;
  digest?: string;
  model?: string;
  status?: string;
  total?: number;
}

// 进度回调函数类型
export type ProgressCallback = (progress: ModelProgressInfo) => void;
export type ErrorCallback = (error: { message: string }) => void;

export class ModelsService {
  // 用于中断下载的控制器
  private _abortController: AbortController | null = null;

  // 获取模型列表
  getModels = async (provider: string): Promise<ChatModelCard[] | undefined> => {
    // TODO: Implement without model-runtime
    console.warn('getModels: model-runtime removed, needs new implementation');
    return undefined;
    // const runtimeProvider = resolveRuntimeProvider(provider);

    // try {
    //   // Check if provider has CORS restrictions
    //   if (isProviderDisableBrowserRequest(provider)) {
    //     console.error(
    //       `Provider "${provider}" cannot fetch models in browser due to CORS restrictions`,
    //     );
    //     return undefined;
    //   }

    //   // Always use client runtime to fetch models directly
    //   const agentRuntime = await initializeWithClientStore({
    //     provider,
    //     runtimeProvider,
    //   });

    //   return await agentRuntime.models();
    // } catch (error) {
    //   console.error(`Failed to fetch models for provider ${provider}:`, error);
    //   return undefined;
    // }
  };

  /**
   * 下载模型并通过回调函数返回进度信息
   */
  downloadModel = async (
    { model, provider }: { model: string; provider: string },
    { onProgress, onError }: { onError?: ErrorCallback; onProgress?: ProgressCallback } = {},
  ): Promise<void> => {
    // TODO: Implement without model-runtime
    console.warn('downloadModel: model-runtime removed, needs new implementation');
    onError?.({ message: 'Not implemented' });
    // try {
    //   const runtimeProvider = resolveRuntimeProvider(provider);

    //   // Check if provider has CORS restrictions
    //   if (isProviderDisableBrowserRequest(provider)) {
    //     const errorMsg = `Provider "${provider}" cannot download models in browser due to CORS restrictions`;
    //     console.error(errorMsg);
    //     onError?.({ message: errorMsg });
    //     throw new Error(errorMsg);
    //   }

    //   // Create a new AbortController
    //   this._abortController = new AbortController();
    //   const signal = this._abortController.signal;

    //   // Always use client runtime to pull models directly
    //   const agentRuntime = await initializeWithClientStore({
    //     provider,
    //     runtimeProvider,
    //   });

    //   const res = await agentRuntime.pullModel({ model }, { signal });

    //   if (!res || !res.ok) {
    //     throw await getMessageError(res!);
    //   }

    //   // Process response stream
    //   if (res.body) {
    //     await this.processModelPullStream(res, { onProgress, onError });
    //   }
    // } catch (error) {
    //   // If it's an abort operation, don't throw error
    //   if (error instanceof DOMException && error.name === 'AbortError') {
    //     return;
    //   }

    //   console.error('download model error:', error);
    //   throw error;
    // } finally {
    //   // Clean up AbortController
    //   this._abortController = null;
    // }
  };

  // 中断模型下载
  abortPull = () => {
    // 使用 AbortController 中断下载
    if (this._abortController) {
      this._abortController.abort();
      this._abortController = null;
    }
  };

  /**
   * 处理模型下载流，解析进度信息并通过回调函数返回
   * @param response 响应对象
   * @param onProgress 进度回调函数
   * @returns Promise<void>
   */
  private processModelPullStream = async (
    response: Response,
    { onProgress, onError }: { onError?: ErrorCallback; onProgress?: ProgressCallback },
  ): Promise<void> => {
    // 处理响应流
    const reader = response.body?.getReader();
    if (!reader) return;

    // 读取和处理流数据
    // eslint-disable-next-line no-constant-condition
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      // 解析进度数据
      const progressText = new TextDecoder().decode(value);
      // 一行可能包含多个进度更新
      const progressUpdates = progressText.trim().split('\n');

      for (const update of progressUpdates) {
        let progress;
        try {
          progress = JSON.parse(update);
        } catch (e) {
          console.error('Error parsing progress update:', e);
          console.error('raw data', update);
        }

        if (progress.status === 'canceled') {
          console.log('progress：', progress);
          // const abortError = new Error('abort');
          // abortError.name = 'AbortError';
          //
          // throw abortError;
        }

        if (progress.status === 'error') {
          onError?.({ message: progress.error });
          throw new Error(progress.error);
        }

        // 调用进度回调
        if (progress.completed !== undefined || progress.status) {
          onProgress?.(progress);
        }
      }
    }
  };
}

export const modelsService = new ModelsService();
