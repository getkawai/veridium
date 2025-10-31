export class TRPCError extends Error {
  public code: string;

  constructor(options: { code: string; message: string }) {
    super(options.message);
    this.code = options.code;
    this.name = 'MockTRPCError';
  }
}