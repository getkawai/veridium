// MOCKED: @trpc/client is not used in production for desktop client
// Desktop router might still exist but the client is not imported anywhere
console.warn('desktopClient is mocked and not functional');

export const desktopClient = new Proxy({} as any, {
  get: () => {
    throw new Error('desktopClient is mocked - not used in production');
  },
});
