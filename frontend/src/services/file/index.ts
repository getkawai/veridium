import { isDesktop } from '@/const/version';

import { ClientService } from './client';
import { ServerService } from './server';

const clientService = new ClientService();

export const fileService =
  process.env.NEXT_PUBLIC_SERVICE_MODE === 'server' || isDesktop
    ? new ServerService()
    : clientService;
