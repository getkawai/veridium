import { createTRPCClient, httpBatchLink } from '@trpc/client';
import superjson from 'superjson';
import { AsyncRouter } from '@/server/routers/async';

export const asyncClient = createTRPCClient<AsyncRouter>({
  links: [
    httpBatchLink({
      fetch: undefined,
      maxURLLength: 2083,
      transformer: superjson,
      url: '/trpc/async',
    }),
  ],
});
