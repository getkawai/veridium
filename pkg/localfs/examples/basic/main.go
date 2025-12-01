package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/pkg/localfs"
)

func main() {
	// Create a new local file service
	service := localfs.NewService()
	ctx := context.Background()

	// Create a temporary directory for our examples
	tmpDir := filepath.Join(os.TempDir(), "localfs-example")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fmt.Println("=== Local File Service Examples ===\n")

	// Example 1: Write a file
	fmt.Println("1. Writing a file...")
	testFile := filepath.Join(tmpDir, "example.txt")
	writeResult, err := service.WriteFile(ctx, localfs.WriteFileParams{
		Path:    testFile,
		Content: "Hello, World!\nThis is line 2.\nThis is line 3.",
	})
	if err != nil {
		log.Fatalf("WriteFile failed: %v", err)
	}
	fmt.Printf("   ✓ File written successfully: %s\n\n", writeResult.Path)

	// Example 2: Read the file
	fmt.Println("2. Reading the file...")
	readResult, err := service.ReadFile(ctx, localfs.ReadFileParams{
		Path: testFile,
	})
	if err != nil {
		log.Fatalf("ReadFile failed: %v", err)
	}
	fmt.Printf("   ✓ File read successfully\n")
	fmt.Printf("   - Filename: %s\n", readResult.Filename)
	fmt.Printf("   - Lines: %d\n", readResult.LineCount)
	fmt.Printf("   - Characters: %d\n", readResult.CharCount)
	fmt.Printf("   - Content:\n%s\n\n", readResult.Content)

	// Example 3: Read specific line range
	fmt.Println("3. Reading lines 2-3...")
	loc := [2]int{2, 3}
	rangeResult, err := service.ReadFile(ctx, localfs.ReadFileParams{
		Path: testFile,
		Loc:  &loc,
	})
	if err != nil {
		log.Fatalf("ReadFile with range failed: %v", err)
	}
	fmt.Printf("   ✓ Lines 2-3:\n%s\n\n", rangeResult.Content)

	// Example 4: Edit the file
	fmt.Println("4. Editing the file (replace 'World' with 'Go')...")
	editResult, err := service.EditFile(ctx, localfs.EditFileParams{
		FilePath:   testFile,
		OldString:  "World",
		NewString:  "Go",
		ReplaceAll: true,
	})
	if err != nil {
		log.Fatalf("EditFile failed: %v", err)
	}
	fmt.Printf("   ✓ Made %d replacement(s)\n\n", editResult.Replacements)

	// Example 5: List files in directory
	fmt.Println("5. Listing files in directory...")
	files, err := service.ListFiles(ctx, localfs.ListFileParams{
		Path: tmpDir,
	})
	if err != nil {
		log.Fatalf("ListFiles failed: %v", err)
	}
	fmt.Printf("   ✓ Found %d file(s):\n", len(files))
	for _, file := range files {
		fmt.Printf("   - %s (%d bytes)\n", file.Name, file.Size)
	}
	fmt.Println()

	// Example 6: Create multiple files for search
	fmt.Println("6. Creating multiple files for search...")
	files2 := map[string]string{
		"test1.txt": "This is a test file",
		"test2.txt": "Another test file",
		"data.log":  "Some log data",
	}
	for name, content := range files2 {
		_, err := service.WriteFile(ctx, localfs.WriteFileParams{
			Path:    filepath.Join(tmpDir, name),
			Content: content,
		})
		if err != nil {
			log.Printf("Failed to write %s: %v", name, err)
		}
	}
	fmt.Printf("   ✓ Created %d files\n\n", len(files2))

	// Example 7: Search for files
	fmt.Println("7. Searching for files containing 'test'...")
	searchResults, err := service.SearchFiles(ctx, localfs.SearchFilesParams{
		Keywords:  "test",
		Directory: tmpDir,
	})
	if err != nil {
		log.Fatalf("SearchFiles failed: %v", err)
	}
	fmt.Printf("   ✓ Found %d file(s):\n", len(searchResults))
	for _, file := range searchResults {
		fmt.Printf("   - %s\n", file.Name)
	}
	fmt.Println()

	// Example 8: Grep content
	fmt.Println("8. Searching for content containing 'test'...")
	grepResult, err := service.GrepContent(ctx, localfs.GrepContentParams{
		Pattern: "test",
		Path:    tmpDir,
		CaseI:   true,
	})
	if err != nil {
		log.Fatalf("GrepContent failed: %v", err)
	}
	fmt.Printf("   ✓ Found %d match(es)\n\n", grepResult.TotalMatches)

	// Example 9: Glob files
	fmt.Println("9. Finding all .txt files...")
	globResult, err := service.GlobFiles(ctx, localfs.GlobFilesParams{
		Pattern: "*.txt",
		Path:    tmpDir,
	})
	if err != nil {
		log.Fatalf("GlobFiles failed: %v", err)
	}
	fmt.Printf("   ✓ Found %d .txt file(s)\n\n", globResult.TotalFiles)

	// Example 10: Rename a file
	fmt.Println("10. Renaming a file...")
	renameResult, err := service.RenameFile(ctx, localfs.RenameFileParams{
		Path:    filepath.Join(tmpDir, "test1.txt"),
		NewName: "renamed.txt",
	})
	if err != nil {
		log.Fatalf("RenameFile failed: %v", err)
	}
	fmt.Printf("   ✓ File renamed to: %s\n\n", filepath.Base(renameResult.NewPath))

	// Example 11: Move files
	fmt.Println("11. Moving files to subdirectory...")
	subDir := filepath.Join(tmpDir, "subdir")
	moveResults, err := service.MoveFiles(ctx, localfs.MoveFilesParams{
		Items: []localfs.MoveFileParams{
			{
				OldPath: filepath.Join(tmpDir, "renamed.txt"),
				NewPath: filepath.Join(subDir, "renamed.txt"),
			},
		},
	})
	if err != nil {
		log.Fatalf("MoveFiles failed: %v", err)
	}
	for _, result := range moveResults {
		if result.Success {
			fmt.Printf("   ✓ Moved: %s -> %s\n", filepath.Base(result.SourcePath), filepath.Base(result.NewPath))
		} else {
			fmt.Printf("   ✗ Failed to move %s: %s\n", result.SourcePath, result.Error)
		}
	}
	fmt.Println()

	// Example 12: Run a shell command
	fmt.Println("12. Running a shell command...")
	cmdResult, err := service.RunCommand(ctx, localfs.RunCommandParams{
		Command: "echo 'Hello from shell'",
	})
	if err != nil {
		log.Fatalf("RunCommand failed: %v", err)
	}
	fmt.Printf("   ✓ Command executed (exit code: %d)\n", cmdResult.ExitCode)
	fmt.Printf("   Output: %s\n", cmdResult.Output)

	fmt.Println("\n=== All examples completed successfully! ===")
}

