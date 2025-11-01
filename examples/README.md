# Node.js Equivalents - Examples

This directory contains examples demonstrating how to use the Node.js equivalent packages.

## Running Examples

Each example can be run individually:

```bash
# File system operations
go run fs_example.go

# Path operations
go run path_example.go

# Buffer operations
go run buffer_example.go

# Child process execution
go run exec_example.go

# OS operations
go run os_example.go
```

## What Each Example Demonstrates

### `fs_example.go`
- File existence checking
- Reading and writing files
- File statistics
- Directory operations
- File cleanup

### `path_example.go`
- Path parsing and manipulation
- Path joining and resolution
- Cross-platform path operations
- Glob pattern matching

### `buffer_example.go`
- Creating buffers from various inputs
- Buffer manipulation (fill, copy, slice)
- Encoding/decoding operations
- Buffer comparison and searching

### `exec_example.go`
- Synchronous command execution
- Asynchronous command execution
- Environment variable handling
- Error handling and timeouts
- Finding executables with `which`

### `os_example.go`
- System information retrieval
- CPU and memory information
- User information
- Directory paths
- Process priority management

## Notes

- Some OS operations may return limited information compared to Node.js due to Go's standard library constraints
- Platform-specific operations may behave differently on Windows, macOS, and Linux
- Error handling follows Go conventions rather than Node.js callback patterns
