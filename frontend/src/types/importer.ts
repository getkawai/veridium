// Legacy import types removed as they are no longer used


export interface ImportResult {
  added: number;
  errors: number;
  skips: number;
}

export interface ImportResults {
  messages?: ImportResult;
  sessionGroups?: ImportResult;
  sessions?: ImportResult;
  topics?: ImportResult;
  type?: string;
}

export enum ImportStage {
  Start,
  Preparing,
  Uploading,
  Importing,
  Success,
  Error,
  Finished,
}

export interface FileUploadState {
  progress: number;
  /**
   * rest time in ms
   */
  restTime: number;
  /**
   * upload speed in KB/s
   */
  speed: number;
}

export interface ErrorShape {
  code: string;
  httpStatus: number;
  message: string;
  path?: string;
}

export interface OnImportCallbacks {
  onError?: (error: ErrorShape) => void;
  onFileUploading?: (state: FileUploadState) => void;
  onStageChange?: (stage: ImportStage) => void;
  /**
   *
   * @param results
   * @param duration in ms
   */
  onSuccess?: (results: ImportResults, duration: number) => void;
}

// ------

export type ImportResultData = ImportSuccessResult | ImportErrorResult;

export interface ImportSuccessResult {
  results: Record<string, any>;
  success: true;
}

export interface ImportErrorResult {
  error: { details?: string; message: string };
  results: Record<string, any>;
  success: false;
}
