import {
  EditLocalFileParams,
  EditLocalFileResult,
  GetCommandOutputParams,
  GetCommandOutputResult,
  GlobFilesParams,
  GlobFilesResult,
  GrepContentParams,
  GrepContentResult,
  KillCommandParams,
  KillCommandResult,
  ListLocalFileParams,
  LocalFileItem,
  LocalMoveFilesResultItem,
  LocalReadFileParams,
  LocalReadFileResult,
  LocalReadFilesParams,
  LocalSearchFilesParams,
  MoveLocalFilesParams,
  OpenLocalFileParams,
  OpenLocalFolderParams,
  RenameLocalFileParams,
  RunCommandParams,
  RunCommandResult,
  WriteLocalFileParams,
  Service,
} from '@@/github.com/kawai-network/veridium/pkg/localfs';

class LocalFileService {
  // File Operations
  async listLocalFiles(params: ListLocalFileParams): Promise<LocalFileItem[]> {
    return Service.ListFiles(params);
  }

  async readLocalFile(params: LocalReadFileParams): Promise<LocalReadFileResult> {
    const result = await Service.ReadFile(params);
    if (!result) throw new Error('Failed to read file');
    return result;
  }

  async readLocalFiles(params: LocalReadFilesParams): Promise<LocalReadFileResult[]> {
    const results = await Service.ReadFiles(params);
    return results.filter((r): r is LocalReadFileResult => r !== null);
  }

  async searchLocalFiles(params: LocalSearchFilesParams): Promise<LocalFileItem[]> {
    return Service.SearchFiles(params);
  }

  async openLocalFile(params: OpenLocalFileParams) {
    return Service.OpenFile(params);
  }

  async openLocalFolder(params: OpenLocalFolderParams) {
    return Service.OpenFolder(params);
  }

  async moveLocalFiles(params: MoveLocalFilesParams): Promise<LocalMoveFilesResultItem[]> {
    return Service.MoveFiles(params);
  }

  async renameLocalFile(params: RenameLocalFileParams) {
    const result = await Service.RenameFile(params);
    if (!result) throw new Error('Failed to rename file');
    return result;
  }

  async writeFile(params: WriteLocalFileParams) {
    const result = await Service.WriteFile(params);
    if (!result) throw new Error('Failed to write file');
    return result;
  }

  async editLocalFile(params: EditLocalFileParams): Promise<EditLocalFileResult> {
    const result = await Service.EditFile(params);
    if (!result) throw new Error('Failed to edit file');
    return result;
  }

  // Shell Commands
  async runCommand(params: RunCommandParams): Promise<RunCommandResult> {
    const result = await Service.RunCommand(params);
    if (!result) throw new Error('Failed to run command');
    return result;
  }

  async getCommandOutput(params: GetCommandOutputParams): Promise<GetCommandOutputResult> {
    const result = await Service.GetCommandOutput(params);
    if (!result) throw new Error('Failed to get command output');
    return result;
  }

  async killCommand(params: KillCommandParams): Promise<KillCommandResult> {
    const result = await Service.KillCommand(params);
    if (!result) throw new Error('Failed to kill command');
    return result;
  }

  // Search & Find
  async grepContent(params: GrepContentParams): Promise<GrepContentResult> {
    const result = await Service.GrepContent(params);
    if (!result) throw new Error('Failed to grep content');
    return result;
  }

  async globFiles(params: GlobFilesParams): Promise<GlobFilesResult> {
    const result = await Service.GlobFiles(params);
    if (!result) throw new Error('Failed to glob files');
    return result;
  }

  // Helper methods
  async openLocalFileOrFolder(path: string, isDirectory: boolean) {
    return Service.OpenFileOrFolder(path, isDirectory);
  }
}

export const localFileService = new LocalFileService();
