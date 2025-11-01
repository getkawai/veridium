package main

import (
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/nodefs"
)

func main() {
	// Example: File existence check
	exists := nodefs.FileExistsSync("fs_example.go")
	fmt.Printf("File exists: %v\n", exists)

	// Example: Reading a file
	if exists {
		data, err := nodefs.ReadFileSync("fs_example.go")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("File size: %d bytes\n", len(data))
		fmt.Printf("First 100 chars: %s...\n", string(data[:min(100, len(data))]))
	}

	// Example: Writing a file
	err := nodefs.WriteFileSync("temp_file.txt", []byte("Hello from nodefs!"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File written successfully")

	// Example: File stats
	info, err := nodefs.StatSync("temp_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File size: %d bytes\n", info.Size())
	fmt.Printf("Is directory: %v\n", info.IsDir())

	// Example: Directory operations
	err = nodefs.MkdirSync("temp_dir")
	if err != nil {
		log.Printf("Directory creation failed (may already exist): %v", err)
	}

	entries, err := nodefs.ReadDirSync(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current directory has %d entries\n", len(entries))

	// Example: Cleanup
	nodefs.RmSync("temp_file.txt")
	nodefs.RmSync("temp_dir")

	fmt.Println("Cleanup completed")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
