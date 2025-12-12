// import { FileStore } from '../../store';
import { FilesStoreState } from '@/store/file/initialState';
import { FileUploadStatus } from '@/types/files/upload';

const uploadStatusArray = new Set(['uploading', 'pending', 'processing']);

const activeFileId = (s: FilesStoreState) => s.activeFileId;
const currentCategory = (s: FilesStoreState) => s.currentCategory;
const isFetchingFiles = (s: FilesStoreState) => s.isFetchingFiles;
const searchKeywords = (s: FilesStoreState) => s.searchKeywords;
const sorter = (s: FilesStoreState) => s.sorter;
const sortType = (s: FilesStoreState) => s.sortType;

const dockUploadFileList = (s: FilesStoreState) => s.dockUploadFileList;
const dockRawFileList = (s: FilesStoreState) => s.dockUploadFileList.map((i) => i.file);

const fileList = (s: FilesStoreState) => s.fileList;
const getFileById = (id?: string | null) => (s: FilesStoreState) => {
  if (!id) return;

  return s.fileList.find((item) => item.id === id);
};

const isUploadingFiles = (s: FilesStoreState) =>
  s.dockUploadFileList.some((file) => uploadStatusArray.has(file.status));

const overviewUploadingStatus = (s: FilesStoreState): FileUploadStatus => {
  if (s.dockUploadFileList.length === 0) return 'pending';
  if (s.dockUploadFileList.some((file) => uploadStatusArray.has(file.status))) {
    return 'uploading';
  }

  return 'success';
};

const overviewUploadingProgress = (s: FilesStoreState) => {
  const uploadFiles = s.dockUploadFileList.filter(
    (file) => file.status === 'uploading' || file.status === 'pending',
  );

  if (uploadFiles.length === 0) return 100;

  const totalPercent = uploadFiles.length * 100;
  const currentPercent = uploadFiles.reduce(
    (acc, file) => acc + (file.uploadState?.progress || 0),
    0,
  );

  return (currentPercent / totalPercent) * 100;
};

const isCreatingFileParseTask = (id: string) => (s: FilesStoreState) =>
  s.creatingChunkingTaskIds.includes(id);

const isCreatingChunkEmbeddingTask = (id: string) => (s: FilesStoreState) =>
  s.creatingEmbeddingTaskIds.includes(id);

export const fileManagerSelectors = {
  activeFileId,
  currentCategory,
  dockFileList: dockUploadFileList,
  dockRawFileList,
  fileList,
  getFileById,
  isCreatingChunkEmbeddingTask,
  isCreatingFileParseTask,
  isFetchingFiles,
  isUploadingFiles,
  overviewUploadingProgress,
  overviewUploadingStatus,
  searchKeywords,
  sortType,
  sorter,
};
