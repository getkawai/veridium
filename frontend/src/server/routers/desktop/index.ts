import { router } from '@/libs/trpc/lambda';
import { mcpRouter } from '@/server/routers/desktop/mcp';

export const desktopRouter = router({
  mcp: mcpRouter,
});

export type DesktopRouter = typeof desktopRouter;
