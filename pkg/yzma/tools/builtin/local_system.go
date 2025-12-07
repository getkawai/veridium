package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/localfs"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// ============================================================================
// Response Types (matching frontend expected format)
// ============================================================================

// LocalFileListState matches frontend LocalFileListState interface
type LocalFileListState struct {
	ListResults []localfs.LocalFileItem `json:"listResults"`
}

// LocalReadFileState matches frontend LocalReadFileState interface
type LocalReadFileState struct {
	FileContent localfs.LocalReadFileResult `json:"fileContent"`
}

// LocalFileSearchState matches frontend LocalFileSearchState interface
type LocalFileSearchState struct {
	SearchResults []localfs.LocalFileItem `json:"searchResults"`
}

// LocalMoveFilesState matches frontend LocalMoveFilesState interface
type LocalMoveFilesState struct {
	Results      []localfs.LocalMoveFilesResultItem `json:"results"`
	SuccessCount int                                `json:"successCount"`
	TotalCount   int                                `json:"totalCount"`
	Error        string                             `json:"error,omitempty"`
}

// LocalRenameFileState matches frontend LocalRenameFileState interface
type LocalRenameFileState struct {
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// LocalReadFilesState matches frontend LocalReadFilesState interface
type LocalReadFilesState struct {
	FilesContent []localfs.LocalReadFileResult `json:"filesContent"`
}

// RunCommandState matches frontend RunCommandState interface
type RunCommandState struct {
	Message string                   `json:"message"`
	Result  localfs.RunCommandResult `json:"result"`
}

// GetCommandOutputState matches frontend GetCommandOutputState interface
type GetCommandOutputState struct {
	Message string                         `json:"message"`
	Result  localfs.GetCommandOutputResult `json:"result"`
}

// KillCommandState matches frontend KillCommandState interface
type KillCommandState struct {
	Message string                    `json:"message"`
	Result  localfs.KillCommandResult `json:"result"`
}

// GrepContentState matches frontend GrepContentState interface
type GrepContentState struct {
	Message string                    `json:"message"`
	Result  localfs.GrepContentResult `json:"result"`
}

// GlobFilesState matches frontend GlobFilesState interface
type GlobFilesState struct {
	Message string                  `json:"message"`
	Result  localfs.GlobFilesResult `json:"result"`
}

// EditLocalFileState matches frontend EditLocalFileState interface
type EditLocalFileState struct {
	Message string                      `json:"message"`
	Result  localfs.EditLocalFileResult `json:"result"`
}

// ============================================================================
// LocalSystemService
// ============================================================================

// LocalSystemService provides local file system operations
type LocalSystemService struct {
	service *localfs.Service
}

// NewLocalSystemService creates a new local system service
func NewLocalSystemService() *LocalSystemService {
	return &LocalSystemService{
		service: localfs.NewService(),
	}
}

// ListLocalFiles lists files in a directory
func (s *LocalSystemService) ListLocalFiles(path string) (*LocalFileListState, error) {
	results, err := s.service.ListFiles(context.Background(), localfs.ListLocalFileParams{Path: path})
	if err != nil {
		return nil, err
	}
	return &LocalFileListState{ListResults: results}, nil
}

// ReadLocalFile reads content from a file
func (s *LocalSystemService) ReadLocalFile(path string, loc [2]int) (*LocalReadFileState, error) {
	// Default loc
	if loc[0] == 0 && loc[1] == 0 {
		loc = [2]int{0, 200}
	}

	result, err := s.service.ReadFile(context.Background(), localfs.LocalReadFileParams{
		Path: path,
		Loc:  &loc,
	})
	if err != nil {
		return nil, err
	}

	return &LocalReadFileState{FileContent: *result}, nil
}

// SearchLocalFiles searches for files matching keywords
func (s *LocalSystemService) SearchLocalFiles(keywords string, directory string) (*LocalFileSearchState, error) {
	results, err := s.service.SearchFiles(context.Background(), localfs.LocalSearchFilesParams{
		Keywords:  keywords,
		Directory: directory,
	})
	if err != nil {
		return nil, err
	}
	return &LocalFileSearchState{SearchResults: results}, nil
}

// WriteLocalFile writes content to a file
func (s *LocalSystemService) WriteLocalFile(path string, content string) (*localfs.WriteFileResult, error) {
	return s.service.WriteFile(context.Background(), localfs.WriteLocalFileParams{
		Path:    path,
		Content: content,
	})
}

// RenameLocalFile renames a file or directory
func (s *LocalSystemService) RenameLocalFile(path string, newName string) (*LocalRenameFileState, error) {
	result, err := s.service.RenameFile(context.Background(), localfs.RenameLocalFileParams{
		Path:    path,
		NewName: newName,
	})
	if err != nil {
		return &LocalRenameFileState{
			OldPath: path,
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &LocalRenameFileState{
		OldPath: path,
		NewPath: result.NewPath,
		Success: result.Success,
		Error:   result.Error,
	}, nil
}

// MoveLocalFiles moves multiple files
func (s *LocalSystemService) MoveLocalFiles(items []localfs.MoveLocalFileParams) (*LocalMoveFilesState, error) {
	results, err := s.service.MoveFiles(context.Background(), localfs.MoveLocalFilesParams{
		Items: items,
	})
	if err != nil {
		return nil, err
	}

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	return &LocalMoveFilesState{
		Results:      results,
		SuccessCount: successCount,
		TotalCount:   len(results),
	}, nil
}

// ReadLocalFiles reads multiple files
func (s *LocalSystemService) ReadLocalFiles(paths []string) (*LocalReadFilesState, error) {
	results, err := s.service.ReadFiles(context.Background(), localfs.LocalReadFilesParams{
		Paths: paths,
	})
	if err != nil {
		return nil, err
	}

	// Convert []*LocalReadFileResult to []LocalReadFileResult
	filesContent := make([]localfs.LocalReadFileResult, 0, len(results))
	for _, result := range results {
		if result != nil {
			filesContent = append(filesContent, *result)
		}
	}

	return &LocalReadFilesState{FilesContent: filesContent}, nil
}

// EditLocalFile edits a file with search and replace
func (s *LocalSystemService) EditLocalFile(filePath, oldString, newString string, replaceAll bool) (*EditLocalFileState, error) {
	result, err := s.service.EditFile(context.Background(), localfs.EditLocalFileParams{
		FilePath:   filePath,
		OldString:  oldString,
		NewString:  newString,
		ReplaceAll: replaceAll,
	})
	if err != nil {
		return &EditLocalFileState{
			Message: fmt.Sprintf("Failed to edit file: %v", err),
			Result:  *result,
		}, nil
	}

	message := fmt.Sprintf("Successfully replaced %d occurrence(s)", result.Replacements)
	return &EditLocalFileState{
		Message: message,
		Result:  *result,
	}, nil
}

// RunCommand runs a shell command
func (s *LocalSystemService) RunCommand(command, description string, runInBackground bool, timeout int) (*RunCommandState, error) {
	result, err := s.service.RunCommand(context.Background(), localfs.RunCommandParams{
		Command:         command,
		Description:     description,
		RunInBackground: runInBackground,
		Timeout:         timeout,
	})
	if err != nil {
		return &RunCommandState{
			Message: fmt.Sprintf("Failed to run command: %v", err),
			Result:  *result,
		}, nil
	}

	message := "Command executed successfully"
	if runInBackground {
		message = fmt.Sprintf("Command started in background with shell ID: %s", result.ShellID)
	}

	return &RunCommandState{
		Message: message,
		Result:  *result,
	}, nil
}

// GetCommandOutput gets output from a running command
func (s *LocalSystemService) GetCommandOutput(shellID, filter string) (*GetCommandOutputState, error) {
	result, err := s.service.GetCommandOutput(context.Background(), localfs.GetCommandOutputParams{
		ShellID: shellID,
		Filter:  filter,
	})
	if err != nil {
		return &GetCommandOutputState{
			Message: fmt.Sprintf("Failed to get command output: %v", err),
			Result:  *result,
		}, nil
	}

	message := "Command output retrieved"
	if result.Running {
		message = "Command is still running"
	}

	return &GetCommandOutputState{
		Message: message,
		Result:  *result,
	}, nil
}

// KillCommand kills a running command
func (s *LocalSystemService) KillCommand(shellID string) (*KillCommandState, error) {
	result, err := s.service.KillCommand(context.Background(), localfs.KillCommandParams{
		ShellID: shellID,
	})
	if err != nil {
		return &KillCommandState{
			Message: fmt.Sprintf("Failed to kill command: %v", err),
			Result:  *result,
		}, nil
	}

	return &KillCommandState{
		Message: "Command killed successfully",
		Result:  *result,
	}, nil
}

// GrepContent searches for content in files
func (s *LocalSystemService) GrepContent(params localfs.GrepContentParams) (*GrepContentState, error) {
	result, err := s.service.GrepContent(context.Background(), params)
	if err != nil {
		return &GrepContentState{
			Message: fmt.Sprintf("Failed to grep content: %v", err),
			Result:  *result,
		}, nil
	}

	message := fmt.Sprintf("Found %d matches", result.TotalMatches)
	return &GrepContentState{
		Message: message,
		Result:  *result,
	}, nil
}

// GlobFiles searches for files using glob patterns
func (s *LocalSystemService) GlobFiles(pattern, path string) (*GlobFilesState, error) {
	result, err := s.service.GlobFiles(context.Background(), localfs.GlobFilesParams{
		Pattern: pattern,
		Path:    path,
	})
	if err != nil {
		return &GlobFilesState{
			Message: fmt.Sprintf("Failed to glob files: %v", err),
			Result:  *result,
		}, nil
	}

	message := fmt.Sprintf("Found %d files", result.TotalFiles)
	return &GlobFilesState{
		Message: message,
		Result:  *result,
	}, nil
}

// ============================================================================
// Tool Registration
// ============================================================================

// RegisterLocalSystem registers the lobe-local-system tools
func RegisterLocalSystem(registry *tools.ToolRegistry) error {
	service := NewLocalSystemService()

	// Tool 1: listLocalFiles (read-only, parallel safe)
	if err := registry.Register(tools.NewSimpleTool(tools.SimpleToolConfig{
		Name:        "lobe-local-system__listLocalFiles",
		Description: "List files and folders in a specified directory. Returns a JSON array of file/folder information.",
		Parameters: map[string]any{
			"path": map[string]any{"type": "string", "description": "The directory path to list"},
		},
		Required: []string{"path"},
		Parallel: true,
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			path := args["path"]
			if path == "" {
				return "", fmt.Errorf("path is required")
			}
			result, err := service.ListLocalFiles(path)
			if err != nil {
				return "", err
			}
			resultJSON, _ := json.Marshal(result)
			log.Printf("📁 Listed %d items in: %s", len(result.ListResults), path)
			return string(resultJSON), nil
		},
	})); err != nil {
		return fmt.Errorf("failed to register listLocalFiles: %w", err)
	}

	// Tool 2: readLocalFile (read-only, parallel safe)
	if err := registry.Register(tools.NewSimpleTool(tools.SimpleToolConfig{
		Name:        "lobe-local-system__readLocalFile",
		Description: "Read the content of a specific file. Returns the file content with metadata.",
		Parameters: map[string]any{
			"path": map[string]any{"type": "string", "description": "The file path to read"},
			"loc":  map[string]any{"type": "array", "items": map[string]any{"type": "number"}, "description": "Optional range of lines to read [startLine, endLine]. Defaults to [0, 200]"},
		},
		Required: []string{"path"},
		Parallel: true,
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			path := args["path"]
			if path == "" {
				return "", fmt.Errorf("path is required")
			}
			loc := [2]int{0, 200}
			if locStr := args["loc"]; locStr != "" {
				var locArr []int
				if err := json.Unmarshal([]byte(locStr), &locArr); err == nil && len(locArr) >= 2 {
					loc = [2]int{locArr[0], locArr[1]}
				}
			}
			result, err := service.ReadLocalFile(path, loc)
			if err != nil {
				return "", err
			}
			resultJSON, _ := json.Marshal(result)
			log.Printf("📄 Read file: %s (lines %d-%d)", path, loc[0], loc[1])
			return string(resultJSON), nil
		},
	})); err != nil {
		return fmt.Errorf("failed to register readLocalFile: %w", err)
	}

	// Tool 3: searchLocalFiles (read-only, parallel safe)
	if err := registry.Register(tools.NewSimpleTool(tools.SimpleToolConfig{
		Name:        "lobe-local-system__searchLocalFiles",
		Description: "Search for files within a directory based on keywords. Returns matching files.",
		Parameters: map[string]any{
			"keywords":  map[string]any{"type": "string", "description": "The search keywords string"},
			"directory": map[string]any{"type": "string", "description": "Optional directory to limit search"},
		},
		Required: []string{"keywords"},
		Parallel: true,
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			keywords := args["keywords"]
			if keywords == "" {
				return "", fmt.Errorf("keywords is required")
			}
			result, err := service.SearchLocalFiles(keywords, args["directory"])
			if err != nil {
				return "", err
			}
			resultJSON, _ := json.Marshal(result)
			log.Printf("🔍 Found %d files matching: %s", len(result.SearchResults), keywords)
			return string(resultJSON), nil
		},
	})); err != nil {
		return fmt.Errorf("failed to register searchLocalFiles: %w", err)
	}

	// Tool 4: writeLocalFile (modifies filesystem, NOT parallel safe)
	if err := registry.Register(tools.NewSimpleTool(tools.SimpleToolConfig{
		Name:        "lobe-local-system__writeLocalFile",
		Description: "Write content to a specific file. Creates the file if it doesn't exist.",
		Parameters: map[string]any{
			"path":    map[string]any{"type": "string", "description": "The file path to write to"},
			"content": map[string]any{"type": "string", "description": "The content to write"},
		},
		Required: []string{"path", "content"},
		Parallel: false,
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			path := args["path"]
			if path == "" {
				return "", fmt.Errorf("path is required")
			}
			result, err := service.WriteLocalFile(path, args["content"])
			if err != nil {
				return "", err
			}
			resultJSON, _ := json.Marshal(result)
			log.Printf("✏️  Wrote file: %s", path)
			return string(resultJSON), nil
		},
	})); err != nil {
		return fmt.Errorf("failed to register writeLocalFile: %w", err)
	}

	// Tool 5: renameLocalFile (modifies filesystem, NOT parallel safe)
	if err := registry.Register(tools.NewSimpleTool(tools.SimpleToolConfig{
		Name:        "lobe-local-system__renameLocalFile",
		Description: "Rename a file or folder in its current location.",
		Parameters: map[string]any{
			"path":    map[string]any{"type": "string", "description": "The current full path of the file or folder to rename"},
			"newName": map[string]any{"type": "string", "description": "The new name for the file or folder (without path)"},
		},
		Required: []string{"path", "newName"},
		Parallel: false,
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			path, newName := args["path"], args["newName"]
			if path == "" || newName == "" {
				return "", fmt.Errorf("path and newName are required")
			}
			result, err := service.RenameLocalFile(path, newName)
			if err != nil {
				return "", err
			}
			resultJSON, _ := json.Marshal(result)
			log.Printf("📝 Renamed: %s -> %s", path, result.NewPath)
			return string(resultJSON), nil
		},
	})); err != nil {
		return fmt.Errorf("failed to register renameLocalFile: %w", err)
	}

	// Tool 6: moveLocalFiles (modifies filesystem, NOT parallel safe)
	if err := registry.Register(tools.NewSimpleTool(tools.SimpleToolConfig{
		Name:        "lobe-local-system__moveLocalFiles",
		Description: "Move or rename multiple files/directories.",
		Parameters: map[string]any{
			"items": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"oldPath": map[string]any{"type": "string", "description": "The current absolute path"},
						"newPath": map[string]any{"type": "string", "description": "The target absolute path"},
					},
					"required": []string{"oldPath", "newPath"},
				},
				"description": "A list of move/rename operations to perform",
			},
		},
		Required: []string{"items"},
		Parallel: false,
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			itemsStr := args["items"]
			if itemsStr == "" {
				return "", fmt.Errorf("items is required")
			}
			var items []localfs.MoveLocalFileParams
			if err := json.Unmarshal([]byte(itemsStr), &items); err != nil {
				return "", fmt.Errorf("failed to parse items: %w", err)
			}
			result, err := service.MoveLocalFiles(items)
			if err != nil {
				return "", err
			}
			resultJSON, _ := json.Marshal(result)
			log.Printf("📦 Moved %d/%d files", result.SuccessCount, result.TotalCount)
			return string(resultJSON), nil
		},
	})); err != nil {
		return fmt.Errorf("failed to register moveLocalFiles: %w", err)
	}

	return nil
}
