export class TRPCError extends Error {
  public code: string;

  constructor(options: { code: string; message: string }) {
    super(options.message);
    this.code = options.code;
    this.name = 'MockTRPCError';
  }
}

export class TRPCClientError<T = unknown> extends Error {
  public readonly cause?: unknown;
  public readonly code: string;
  public readonly data?: T;
  public readonly meta?: Record<string, unknown>;
  public readonly shape?: unknown;

  constructor(message: string, opts?: {
    cause?: unknown;
    code?: string;
    data?: T;
    meta?: Record<string, unknown>;
    shape?: unknown;
  }) {
    super(message);
    this.name = 'TRPCClientError';
    this.code = opts?.code ?? 'UNKNOWN_ERROR';
    this.data = opts?.data;
    this.meta = opts?.meta;
    this.shape = opts?.shape;
    this.cause = opts?.cause;
  }
}
