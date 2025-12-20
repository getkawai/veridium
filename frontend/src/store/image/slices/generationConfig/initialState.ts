/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
import {
  ModelParamsSchema,
  ModelProvider,
  RuntimeImageGenParams,
  extractDefaultValues,
} from '@/model-bank';

export const gptImage1ParamsSchema: ModelParamsSchema = {
  imageUrls: { default: [] },
  prompt: { default: '' },
  size: {
    default: 'auto',
    enum: ['auto', '1024x1024', '1536x1024', '1024x1536'],
  },
};

import { DEFAULT_IMAGE_CONFIG } from '@/const/settings';

export const DEFAULT_AI_IMAGE_PROVIDER = ModelProvider.Kawai;
export const DEFAULT_AI_IMAGE_MODEL = 'kawai-image-generator';

export interface GenerationConfigState {
  parameters: RuntimeImageGenParams;
  parametersSchema: ModelParamsSchema;

  provider: string;
  model: string;
  imageNum: number;

  isAspectRatioLocked: boolean;
  activeAspectRatio: string | null; // string - 虚拟比例; null - 原生比例

  /**
   * 标记配置是否已初始化（包括从记忆中恢复）
   */
  isInit: boolean;
}

export const DEFAULT_IMAGE_GENERATION_PARAMETERS: RuntimeImageGenParams =
  extractDefaultValues(gptImage1ParamsSchema);

export const initialGenerationConfigState: GenerationConfigState = {
  model: DEFAULT_AI_IMAGE_MODEL,
  provider: DEFAULT_AI_IMAGE_PROVIDER,
  imageNum: DEFAULT_IMAGE_CONFIG.defaultImageNum,
  parameters: DEFAULT_IMAGE_GENERATION_PARAMETERS,
  parametersSchema: gptImage1ParamsSchema,
  isAspectRatioLocked: false,
  activeAspectRatio: null,
  isInit: false,
};
