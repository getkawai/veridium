import { getResolvedUserId } from '@/utils/userId';

export class BaseClientService {
  protected readonly userId: string;

  constructor(userId?: string) {
    this.userId = userId || getResolvedUserId();
  }
}
