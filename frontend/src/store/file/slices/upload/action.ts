// import { sha256 } from 'js-sha256';
import { StateCreator } from 'zustand/vanilla';

import { UploadFileItem } from '@/types/files';
import { getImageDimensions } from '@/utils/client/imageDimensions';

import { FileStore } from '../../store';

type OnStatusUpdate = (
  data:
    | {
        id: string;
        type: 'updateFile';
        value: Partial<UploadFileItem>;
      }
    | {
        id: string;
        type: 'removeFile';
      },
) => void;

interface UploadWithProgressParams {
  file: File;
  knowledgeBaseId?: string;
  onStatusUpdate?: OnStatusUpdate;
  /**
   * Optional flag to indicate whether to skip the file type check.
   * When set to `true`, any file type checks will be bypassed.
   * Default is `false`, which means file type checks will be performed.
   */
  skipCheckFileType?: boolean;
}

interface UploadWithProgressResult {
  dimensions?: {
    height: number;
    width: number;
  };
  filename?: string;
  id: string;
  url: string;
}

export interface FileUploadAction {
  uploadBase64FileWithProgress: (
    base64: string,
    params?: {
      onStatusUpdate?: OnStatusUpdate;
    },
  ) => Promise<UploadWithProgressResult | undefined>;

  uploadWithProgress: (
    params: UploadWithProgressParams,
  ) => Promise<UploadWithProgressResult | undefined>;
}

export const createFileUploadSlice: StateCreator<
  FileStore,
  [['zustand/devtools', never]],
  [],
  FileUploadAction
> = () => ({
  uploadBase64FileWithProgress: async (base64) => {
    // Extract image dimensions from base64 data
    const dimensions = await getImageDimensions(base64);

    // Mock response for now
    const mockId = `mock-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    return {
      id: mockId,
      url: `data:image/png;base64,${base64}`,
      dimensions,
      filename: `mock-image-${mockId}.png`
    };
  },
  uploadWithProgress: async ({ file, onStatusUpdate, knowledgeBaseId, skipCheckFileType }) => {
    // 1. extract image dimensions if applicable
    const dimensions = await getImageDimensions(file);

    // Mock response for now
    const mockId = `mock-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

    // Simulate upload progress
    onStatusUpdate?.({
      id: file.name,
      type: 'updateFile',
      value: { status: 'uploading', uploadState: { progress: 50, restTime: 1, speed: 1000 } },
    });

    // Simulate completion
    setTimeout(() => {
      onStatusUpdate?.({
        id: file.name,
        type: 'updateFile',
        value: {
          fileUrl: `mock://file/${mockId}`,
          id: mockId,
          status: 'success',
          uploadState: { progress: 100, restTime: 0, speed: 0 },
        },
      });
    }, 100);

    return {
      id: mockId,
      url: `mock://file/${mockId}`,
      dimensions,
      filename: file.name
    };
  },
});
