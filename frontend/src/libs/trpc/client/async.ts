// MOCKED: @trpc/client is not used in production
// Async router has been removed, this export is kept for backwards compatibility
console.warn('asyncClient is mocked and not functional');

export const asyncClient = new Proxy({} as any, {
  get: () => {
    throw new Error('asyncClient is mocked - async router has been removed');
  },
});
