import { FileListItem, QueryFileListParams, FilesTabs, SortType } from '@/types/files';
import { UploadFileItem } from '@/types/files/upload';

export interface FileManagerState {
  activeFileId?: string;
  creatingChunkingTaskIds: string[];
  creatingEmbeddingTaskIds: string[];
  currentCategory: string;
  dockUploadFileList: UploadFileItem[];
  fileDetail?: FileListItem;
  fileList: FileListItem[];
  isFetchingFiles: boolean;
  queryListParams?: QueryFileListParams;
  searchKeywords: string;
  sortType: SortType;
  sorter: string;
}

export const initialFileManagerState: FileManagerState = {
  creatingChunkingTaskIds: [],
  creatingEmbeddingTaskIds: [],
  currentCategory: FilesTabs.All,
  dockUploadFileList: [],
  fileList: [],
  isFetchingFiles: false,
  searchKeywords: '',
  sortType: SortType.Desc,
  sorter: 'createdAt',
};
