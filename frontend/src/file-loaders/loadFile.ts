import debug from 'debug';
import { LoadFile as WailsLoadFile } from '@@/github.com/kawai-network/veridium/loadfileservice';

import { FileDocument, FileMetadata } from './types';

const log = debug('file-loaders:loadFile');

/**
 * Loads a file from the specified path using Go backend processing.
 * The content is returned in markdown format.
 *
 * @param filePath The path to the file to load.
 * @param fileMetadata Optional metadata to override information read from the filesystem.
 * @returns A Promise resolving to a FileDocument object with markdown content.
 */
export const loadFile = async (
  filePath: string,
  fileMetadata?: FileMetadata,
): Promise<FileDocument> => {
  log('Starting to load file using Go backend:', filePath, 'with metadata:', fileMetadata);

  try {
    // Call the Wails LoadFile service
    const result = await WailsLoadFile(filePath, fileMetadata || null);
    if (!result) {
      throw new Error('LoadFile returned null result');
    }

    log('File loaded successfully from Go backend:', {
      fileType: result.fileType,
      filename: result.filename,
      contentLength: result.content?.length || 0,
      pagesCount: result.pages?.length || 0,
    });

    // Convert the result to match our TypeScript FileDocument interface
    const fileDocument: FileDocument = {
      content: result.content || '', // Already in markdown format from Go
      createdTime: new Date(result.createdTime),
      fileType: result.fileType || '',
      filename: result.filename || '',
      metadata: result.metadata || {},
      modifiedTime: new Date(result.modifiedTime),
      pages: result.pages || [],
      source: result.source || '',
      totalCharCount: result.totalCharCount || 0,
      totalLineCount: result.totalLineCount || 0,
    };

    return fileDocument;
  } catch (error) {
    log('Error loading file from Go backend:', error);
    console.error(`Error loading file ${filePath}:`, error);

    // Return a minimal error document
    const errorDoc: FileDocument = {
      content: '',
      createdTime: new Date(),
      fileType: '',
      filename: filePath.split('/').pop() || 'unknown',
      metadata: {
        error: `Failed to load file: ${error}`,
      },
      modifiedTime: new Date(),
      pages: [
        {
          charCount: 0,
          lineCount: 0,
          metadata: { error: `Failed to load file: ${error}` },
          pageContent: '',
        },
      ],
      source: filePath,
      totalCharCount: 0,
      totalLineCount: 0,
    };

    return errorDoc;
  }
};
