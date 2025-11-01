package main

import (
	"fmt"

	"github.com/kawai-network/veridium/pkg/nodepath"
)

func main() {
	// Example: Path manipulation
	path := "/home/user/documents/file.txt"

	fmt.Printf("Original path: %s\n", path)
	fmt.Printf("Basename: %s\n", nodepath.Basename(path))
	fmt.Printf("Dirname: %s\n", nodepath.Dirname(path))
	fmt.Printf("Extname: %s\n", nodepath.Extname(path))

	// Example: Path joining
	joined := nodepath.Join("home", "user", "documents", "file.txt")
	fmt.Printf("Joined path: %s\n", joined)

	// Example: Path parsing
	parsed := nodepath.Parse(path)
	fmt.Printf("Parsed path:\n")
	fmt.Printf("  Root: %s\n", parsed.Root)
	fmt.Printf("  Dir: %s\n", parsed.Dir)
	fmt.Printf("  Base: %s\n", parsed.Base)
	fmt.Printf("  Ext: %s\n", parsed.Ext)
	fmt.Printf("  Name: %s\n", parsed.Name)

	// Example: Path resolution
	relative := nodepath.Resolve("relative", "path", "file.txt")
	fmt.Printf("Resolved path: %s\n", relative)

	// Example: Path normalization
	normalized := nodepath.Normalize("./path//to/../to/file.txt")
	fmt.Printf("Normalized path: %s\n", normalized)

	// Example: Relative path calculation
	rel, _ := nodepath.Relative("/home/user", "/home/user/documents")
	fmt.Printf("Relative path: %s\n", rel)

	// Example: Absolute path check
	isAbs := nodepath.IsAbsolute("/absolute/path")
	fmt.Printf("Is absolute: %v\n", isAbs)

	isAbsRel := nodepath.IsAbsolute("relative/path")
	fmt.Printf("Is relative path absolute: %v\n", isAbsRel)

	// Example: Platform-specific information
	fmt.Printf("Path separator: %s\n", nodepath.Sep)
	fmt.Printf("Path delimiter: %s\n", nodepath.Delimiter)

	// Example: Glob pattern matching (simplified)
	// Note: This is a basic implementation
	matches, err := nodepath.Glob("*.go")
	if err == nil {
		fmt.Printf("Go files found: %v\n", matches)
	}
}
