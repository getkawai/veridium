import { ModelProviderCard } from '@/types/llm';

const Kawai: ModelProviderCard = {
  chatModels: [],
  description:
    'Kawai Cloud 通过官方部署的 API 来实现 AI 模型的调用，并采用 Credits 计算积分的方式来衡量 AI 模型的用量，对应大模型使用的 Tokens。',
  enabled: true,
  id: 'kawai',
  modelsUrl: 'https://getkawai.com/docs/usage/subscription/model-pricing',
  name: 'Kawai',
  settings: {
    modelEditable: false,
    showAddNewModel: false,
    showModelFetcher: false,
  },
  showConfig: false,
  url: 'https://getkawai.com',
};

export default Kawai;

export const planCardModels = ['gpt-4o-mini', 'deepseek-reasoner', 'claude-3-5-sonnet-latest'];
