// MOCKED: @trpc/client is not used in production
// Async router has been removed, this export is kept for backwards compatibility
export const asyncClient = new Proxy({} as any, {
  get: () => {
    throw new Error('asyncClient is mocked - async router has been removed');
  },
});
