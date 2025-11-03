// MOCKED: @trpc/client is not used in production
// Tools router has been removed, this export is kept for backwards compatibility
console.warn('toolsClient is mocked and not functional');

export const toolsClient = new Proxy({} as any, {
  get: () => {
    throw new Error('toolsClient is mocked - tools router has been removed');
  },
});
