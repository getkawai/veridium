import { StateCreator } from 'zustand';

import { FileStore } from '../../store';

export interface TTSFileAction {
  uploadTTSByArrayBuffers: (messageId: string, arrayBuffers: ArrayBuffer[]) => Promise<string>;
}

export const createTTSFileSlice: StateCreator<
  FileStore,
  [['zustand/devtools', never]],
  [],
  TTSFileAction
> = (set, get) => ({
  uploadTTSByArrayBuffers: async (messageId, arrayBuffers) => {
    // For now, we'll store the audio as a data URL
    // In a real implementation, you might want to save this to IndexedDB or upload to a server
    
    if (arrayBuffers.length === 0) {
      throw new Error('No audio data provided');
    }

    // Combine all array buffers
    const blob = new Blob(arrayBuffers, { type: 'audio/aiff' });
    
    // Create a data URL
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onloadend = () => {
        const dataUrl = reader.result as string;
        // Store with messageId as key
        const fileId = `tts-${messageId}-${Date.now()}`;
        resolve(fileId);
      };
      reader.onerror = reject;
      reader.readAsDataURL(blob);
    });
  },
});

