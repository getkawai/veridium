// tRPC client for lambda operations
// This is a simplified implementation for build purposes

export const lambdaClient = {
  document: {
    parseFileContent: {
      mutate: async (params: any) => {
        console.warn('lambdaClient.document.parseFileContent not implemented', params);
        throw new Error('lambdaClient not fully implemented');
      },
    },
  },
  chunk: {
    createParseFileTask: {
      mutate: async (params: any) => {
        console.warn('lambdaClient.chunk.createParseFileTask not implemented', params);
        throw new Error('lambdaClient not fully implemented');
      },
    },
    retryParseFileTask: {
      mutate: async (params: any) => {
        console.warn('lambdaClient.chunk.retryParseFileTask not implemented', params);
        throw new Error('lambdaClient not fully implemented');
      },
    },
    createEmbeddingChunksTask: {
      mutate: async (params: any) => {
        console.warn('lambdaClient.chunk.createEmbeddingChunksTask not implemented', params);
        throw new Error('lambdaClient not fully implemented');
      },
    },
    semanticSearch: {
      mutate: async (params: any) => {
        console.warn('lambdaClient.chunk.semanticSearch not implemented', params);
        throw new Error('lambdaClient not fully implemented');
      },
    },
  },
};
