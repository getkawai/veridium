import { z } from 'zod';

import { pluginManifestSchema } from './manifest';

export const pluginRequestPayloadSchema = z.object({
  apiName: z.string(),
  arguments: z.string().optional(),
  identifier: z.string(),
  indexUrl: z.string().optional(),
  manifest: pluginManifestSchema.optional(),
  type: z.string().optional(),
});

export type PluginRequestPayload = z.infer<typeof pluginRequestPayloadSchema>;
