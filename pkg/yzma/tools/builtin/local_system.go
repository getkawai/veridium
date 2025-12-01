package builtin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// ============================================================================
// Response Types (matching frontend expected format)
// ============================================================================

// LocalFileItem matches frontend LocalFileItem interface
type LocalFileItem struct {
	Name           string                 `json:"name"`
	Path           string                 `json:"path"`
	Size           int64                  `json:"size"`
	Type           string                 `json:"type"`
	IsDirectory    bool                   `json:"isDirectory"`
	ContentType    string                 `json:"contentType,omitempty"`
	CreatedTime    time.Time              `json:"createdTime"`
	ModifiedTime   time.Time              `json:"modifiedTime"`
	LastAccessTime time.Time              `json:"lastAccessTime"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// LocalFileListState matches frontend LocalFileListState interface
type LocalFileListState struct {
	ListResults []LocalFileItem `json:"listResults"`
}

// LocalReadFileResult matches frontend LocalReadFileResult interface
type LocalReadFileResult struct {
	Content        string    `json:"content"`
	Filename       string    `json:"filename"`
	FileType       string    `json:"fileType"`
	CharCount      int       `json:"charCount"`
	LineCount      int       `json:"lineCount"`
	TotalCharCount int       `json:"totalCharCount"`
	TotalLineCount int       `json:"totalLineCount"`
	Loc            [2]int    `json:"loc"`
	CreatedTime    time.Time `json:"createdTime"`
	ModifiedTime   time.Time `json:"modifiedTime"`
}

// LocalReadFileState matches frontend LocalReadFileState interface
type LocalReadFileState struct {
	FileContent LocalReadFileResult `json:"fileContent"`
}

// LocalFileSearchState matches frontend LocalFileSearchState interface
type LocalFileSearchState struct {
	SearchResults []LocalFileItem `json:"searchResults"`
}

// LocalMoveFilesResultItem matches frontend LocalMoveFilesResultItem interface
type LocalMoveFilesResultItem struct {
	SourcePath string `json:"sourcePath"`
	NewPath    string `json:"newPath,omitempty"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

// LocalMoveFilesState matches frontend LocalMoveFilesState interface
type LocalMoveFilesState struct {
	Results      []LocalMoveFilesResultItem `json:"results"`
	SuccessCount int                        `json:"successCount"`
	TotalCount   int                        `json:"totalCount"`
	Error        string                     `json:"error,omitempty"`
}

// LocalRenameFileState matches frontend LocalRenameFileState interface
type LocalRenameFileState struct {
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// WriteFileResult for write operations
type WriteFileResult struct {
	Path    string `json:"path"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ============================================================================
// LocalSystemService
// ============================================================================

// LocalSystemService provides local file system operations
type LocalSystemService struct{}

// NewLocalSystemService creates a new local system service
func NewLocalSystemService() *LocalSystemService {
	return &LocalSystemService{}
}

// ListLocalFiles lists files in a directory
func (s *LocalSystemService) ListLocalFiles(path string) (*LocalFileListState, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	results := make([]LocalFileItem, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fullPath := filepath.Join(path, entry.Name())
		fileType := "file"
		if entry.IsDir() {
			fileType = "directory"
		}

		results = append(results, LocalFileItem{
			Name:           entry.Name(),
			Path:           fullPath,
			Size:           info.Size(),
			Type:           fileType,
			IsDirectory:    entry.IsDir(),
			ContentType:    getContentType(entry.Name()),
			CreatedTime:    info.ModTime(), // Go doesn't have creation time cross-platform
			ModifiedTime:   info.ModTime(),
			LastAccessTime: info.ModTime(),
		})
	}

	return &LocalFileListState{ListResults: results}, nil
}

// ReadLocalFile reads content from a file
func (s *LocalSystemService) ReadLocalFile(path string, loc [2]int) (*LocalReadFileState, error) {
	// Default loc
	if loc[0] == 0 && loc[1] == 0 {
		loc = [2]int{0, 200}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Split into lines
	lines := strings.Split(string(content), "\n")
	totalLines := len(lines)
	totalChars := len(content)

	// Apply line range
	startLine := loc[0]
	endLine := loc[1]
	if startLine < 0 {
		startLine = 0
	}
	if endLine > totalLines {
		endLine = totalLines
	}
	if startLine >= totalLines {
		startLine = totalLines - 1
		if startLine < 0 {
			startLine = 0
		}
	}

	selectedLines := lines[startLine:endLine]
	selectedContent := strings.Join(selectedLines, "\n")

	return &LocalReadFileState{
		FileContent: LocalReadFileResult{
			Content:        selectedContent,
			Filename:       filepath.Base(path),
			FileType:       getContentType(path),
			CharCount:      len(selectedContent),
			LineCount:      len(selectedLines),
			TotalCharCount: totalChars,
			TotalLineCount: totalLines,
			Loc:            [2]int{startLine, endLine},
			CreatedTime:    info.ModTime(),
			ModifiedTime:   info.ModTime(),
		},
	}, nil
}

// SearchLocalFiles searches for files matching keywords
func (s *LocalSystemService) SearchLocalFiles(keywords string, directory string) (*LocalFileSearchState, error) {
	if directory == "" {
		var err error
		directory, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	results := make([]LocalFileItem, 0)
	keywordsLower := strings.ToLower(keywords)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip hidden files/directories
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if name contains keywords
		if strings.Contains(strings.ToLower(info.Name()), keywordsLower) {
			fileType := "file"
			if info.IsDir() {
				fileType = "directory"
			}

			results = append(results, LocalFileItem{
				Name:           info.Name(),
				Path:           path,
				Size:           info.Size(),
				Type:           fileType,
				IsDirectory:    info.IsDir(),
				ContentType:    getContentType(info.Name()),
				CreatedTime:    info.ModTime(),
				ModifiedTime:   info.ModTime(),
				LastAccessTime: info.ModTime(),
			})
		}

		// Limit results
		if len(results) >= 100 {
			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return &LocalFileSearchState{SearchResults: results}, nil
}

// WriteLocalFile writes content to a file
func (s *LocalSystemService) WriteLocalFile(path string, content string) (*WriteFileResult, error) {
	// Create parent directories if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &WriteFileResult{
			Path:    path,
			Success: false,
			Error:   fmt.Sprintf("failed to create directory: %v", err),
		}, nil
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return &WriteFileResult{
			Path:    path,
			Success: false,
			Error:   fmt.Sprintf("failed to write file: %v", err),
		}, nil
	}

	return &WriteFileResult{
		Path:    path,
		Success: true,
		Message: "File written successfully",
	}, nil
}

// RenameLocalFile renames a file or directory
func (s *LocalSystemService) RenameLocalFile(path string, newName string) (*LocalRenameFileState, error) {
	dir := filepath.Dir(path)
	newPath := filepath.Join(dir, newName)

	if err := os.Rename(path, newPath); err != nil {
		return &LocalRenameFileState{
			OldPath: path,
			NewPath: newPath,
			Success: false,
			Error:   fmt.Sprintf("failed to rename: %v", err),
		}, nil
	}

	return &LocalRenameFileState{
		OldPath: path,
		NewPath: newPath,
		Success: true,
	}, nil
}

// MoveLocalFiles moves multiple files
func (s *LocalSystemService) MoveLocalFiles(items []struct {
	OldPath string `json:"oldPath"`
	NewPath string `json:"newPath"`
}) (*LocalMoveFilesState, error) {
	results := make([]LocalMoveFilesResultItem, 0, len(items))
	successCount := 0

	for _, item := range items {
		// Create parent directory if needed
		dir := filepath.Dir(item.NewPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			results = append(results, LocalMoveFilesResultItem{
				SourcePath: item.OldPath,
				Success:    false,
				Error:      fmt.Sprintf("failed to create directory: %v", err),
			})
			continue
		}

		if err := os.Rename(item.OldPath, item.NewPath); err != nil {
			results = append(results, LocalMoveFilesResultItem{
				SourcePath: item.OldPath,
				Success:    false,
				Error:      fmt.Sprintf("failed to move: %v", err),
			})
		} else {
			results = append(results, LocalMoveFilesResultItem{
				SourcePath: item.OldPath,
				NewPath:    item.NewPath,
				Success:    true,
			})
			successCount++
		}
	}

	return &LocalMoveFilesState{
		Results:      results,
		SuccessCount: successCount,
		TotalCount:   len(items),
	}, nil
}

// getContentType returns content type based on file extension
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt":
		return "text/plain"
	case ".md":
		return "text/markdown"
	case ".json":
		return "application/json"
	case ".js":
		return "application/javascript"
	case ".ts":
		return "application/typescript"
	case ".go":
		return "text/x-go"
	case ".py":
		return "text/x-python"
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".pdf":
		return "application/pdf"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}

// ============================================================================
// Tool Registration
// ============================================================================

// RegisterLocalSystem registers the lobe-local-system tools
func RegisterLocalSystem(registry *tools.ToolRegistry) error {
	service := NewLocalSystemService()

	// Tool 1: listLocalFiles
	listTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-local-system__listLocalFiles",
			Description: "List files and folders in a specified directory. Returns a JSON array of file/folder information.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The directory path to list",
					},
				},
				"required": []string{"path"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			path := args["path"]
			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			result, err := service.ListLocalFiles(path)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal result: %w", err)
			}

			log.Printf("📁 Listed %d items in: %s", len(result.ListResults), path)
			return string(resultJSON), nil
		},
		Enabled: true,
	}
	if err := registry.Register(listTool); err != nil {
		return fmt.Errorf("failed to register listLocalFiles: %w", err)
	}

	// Tool 2: readLocalFile
	readTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-local-system__readLocalFile",
			Description: "Read the content of a specific file. Returns the file content with metadata.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The file path to read",
					},
					"loc": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "number"},
						"description": "Optional range of lines to read [startLine, endLine]. Defaults to [0, 200]",
					},
				},
				"required": []string{"path"},
			},
		},
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

			resultJSON, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal result: %w", err)
			}

			log.Printf("📄 Read file: %s (lines %d-%d)", path, loc[0], loc[1])
			return string(resultJSON), nil
		},
		Enabled: true,
	}
	if err := registry.Register(readTool); err != nil {
		return fmt.Errorf("failed to register readLocalFile: %w", err)
	}

	// Tool 3: searchLocalFiles
	searchTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-local-system__searchLocalFiles",
			Description: "Search for files within a directory based on keywords. Returns matching files.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"keywords": map[string]interface{}{
						"type":        "string",
						"description": "The search keywords string",
					},
					"directory": map[string]interface{}{
						"type":        "string",
						"description": "Optional directory to limit search",
					},
				},
				"required": []string{"keywords"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			keywords := args["keywords"]
			if keywords == "" {
				return "", fmt.Errorf("keywords is required")
			}
			directory := args["directory"]

			result, err := service.SearchLocalFiles(keywords, directory)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal result: %w", err)
			}

			log.Printf("🔍 Found %d files matching: %s", len(result.SearchResults), keywords)
			return string(resultJSON), nil
		},
		Enabled: true,
	}
	if err := registry.Register(searchTool); err != nil {
		return fmt.Errorf("failed to register searchLocalFiles: %w", err)
	}

	// Tool 4: writeLocalFile
	writeTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-local-system__writeLocalFile",
			Description: "Write content to a specific file. Creates the file if it doesn't exist.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The file path to write to",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "The content to write",
					},
				},
				"required": []string{"path", "content"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			path := args["path"]
			content := args["content"]
			if path == "" {
				return "", fmt.Errorf("path is required")
			}

			result, err := service.WriteLocalFile(path, content)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal result: %w", err)
			}

			log.Printf("✏️  Wrote file: %s", path)
			return string(resultJSON), nil
		},
		Enabled: true,
	}
	if err := registry.Register(writeTool); err != nil {
		return fmt.Errorf("failed to register writeLocalFile: %w", err)
	}

	// Tool 5: renameLocalFile
	renameTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-local-system__renameLocalFile",
			Description: "Rename a file or folder in its current location.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"path": map[string]interface{}{
						"type":        "string",
						"description": "The current full path of the file or folder to rename",
					},
					"newName": map[string]interface{}{
						"type":        "string",
						"description": "The new name for the file or folder (without path)",
					},
				},
				"required": []string{"path", "newName"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			path := args["path"]
			newName := args["newName"]
			if path == "" || newName == "" {
				return "", fmt.Errorf("path and newName are required")
			}

			result, err := service.RenameLocalFile(path, newName)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal result: %w", err)
			}

			log.Printf("📝 Renamed: %s -> %s", path, result.NewPath)
			return string(resultJSON), nil
		},
		Enabled: true,
	}
	if err := registry.Register(renameTool); err != nil {
		return fmt.Errorf("failed to register renameLocalFile: %w", err)
	}

	// Tool 6: moveLocalFiles
	moveTool := &tools.YzmaTool{
		Type: "function",
		Function: tools.YzmaToolFunction{
			Name:        "lobe-local-system__moveLocalFiles",
			Description: "Move or rename multiple files/directories.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"items": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"oldPath": map[string]interface{}{
									"type":        "string",
									"description": "The current absolute path of the file/directory",
								},
								"newPath": map[string]interface{}{
									"type":        "string",
									"description": "The target absolute path",
								},
							},
							"required": []string{"oldPath", "newPath"},
						},
						"description": "A list of move/rename operations to perform",
					},
				},
				"required": []string{"items"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			itemsStr := args["items"]
			if itemsStr == "" {
				return "", fmt.Errorf("items is required")
			}

			var items []struct {
				OldPath string `json:"oldPath"`
				NewPath string `json:"newPath"`
			}
			if err := json.Unmarshal([]byte(itemsStr), &items); err != nil {
				return "", fmt.Errorf("failed to parse items: %w", err)
			}

			result, err := service.MoveLocalFiles(items)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal result: %w", err)
			}

			log.Printf("📦 Moved %d/%d files", result.SuccessCount, result.TotalCount)
			return string(resultJSON), nil
		},
		Enabled: true,
	}
	if err := registry.Register(moveTool); err != nil {
		return fmt.Errorf("failed to register moveLocalFiles: %w", err)
	}

	return nil
}
