// MOCKED: @trpc/client is not used in production for edge client
// Edge router might still exist but the client is not imported anywhere
console.warn('edgeClient is mocked and not functional');

export const edgeClient = new Proxy({} as any, {
  get: () => {
    throw new Error('edgeClient is mocked - not used in production');
  },
});
