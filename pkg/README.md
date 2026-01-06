# Veridium Go Packages

This directory contains reusable Go packages for the Veridium project.

## Core Veridium Packages

These packages are specific to the Veridium DePIN network:

### `store` - Off-chain Storage Layer

**Purpose:** Cloudflare Workers KV storage for contributors, rewards, and Merkle proofs.

```go
import "github.com/kawai-network/veridium/pkg/store"

// Initialize store
s := store.NewStore(accountID, namespaceID, apiToken)

// Contributor management
contributor, err := s.GetContributor(ctx, address)
s.UpdateContributor(ctx, address, contributor)

// Merkle proofs for reward claims
proof, err := s.GetProof(ctx, address, periodID)
s.SaveProof(ctx, address, periodID, proof)

// Settlement automation
settlement, err := s.GetSettlement(ctx, periodID)
s.SaveSettlement(ctx, periodID, settlement)
```

**Features:**
- Multi-namespace design (contributors, proofs, settlements)
- Period-based Merkle proof storage
- Settlement automation with rollback support
- Claim tracking to prevent token loss
- Rate limiting and retry logic

See [store/README.md](store/README.md) for detailed documentation.

### `merkle` - Merkle Tree Generation

**Purpose:** Generate Merkle trees for gas-efficient reward distribution.

```go
import "github.com/kawai-network/veridium/pkg/merkle"

// Create Merkle tree from rewards
leaves := []merkle.Leaf{
    {Address: "0x123...", Amount: big.NewInt(1000)},
    {Address: "0x456...", Amount: big.NewInt(2000)},
}

tree := merkle.NewTree(leaves)
root := tree.Root()
proof := tree.GetProof("0x123...")
```

**Features:**
- Efficient tree construction
- Proof generation for individual claims
- Compatible with OpenZeppelin MerkleProof.sol
- Used by all reward distributors (Mining, Cashback, Referral)

### `blockchain` - Monad Blockchain Interaction

**Purpose:** Interact with Monad blockchain and smart contracts.

```go
import "github.com/kawai-network/veridium/pkg/blockchain"

// Connect to Monad
client, err := blockchain.NewClient(rpcURL)

// Contract interaction
contract, err := blockchain.NewContract(address, abi, client)
tx, err := contract.Call("claimReward", proof, amount)

// Event listening
events, err := contract.FilterLogs(fromBlock, toBlock, "Claimed")
```

**Features:**
- Monad RPC client wrapper
- Contract call management
- Event filtering and parsing
- Transaction monitoring
- Gas estimation

### `config` - Configuration Management

**Purpose:** Centralized configuration for all services.

```go
import "github.com/kawai-network/veridium/pkg/config"

// Load configuration
cfg := config.Load()

// Access settings
rpcURL := cfg.Blockchain.RPC
kvToken := cfg.Store.APIToken
```

**Features:**
- Environment variable loading
- Network-specific settings
- Contract addresses
- API keys and secrets

---

## Utility Packages

These packages provide general-purpose utilities:

### `obfuscator` - String Obfuscation

A custom string encoder/decoder that provides obfuscation without using secret keys. Ideal for hiding data from casual inspection.

```go
import "github.com/kawai-network/veridium/pkg/obfuscator"

// Quick usage
encoded := obfuscator.EncodeString("Hello, World!")
decoded, err := obfuscator.DecodeString(encoded)

// Reusable instance
o := obfuscator.New()
encoded := o.Encode("sensitive data")
decoded, err := o.Decode(encoded)
```

**Features:**
- No secret key required - deterministic obfuscation
- Multiple layers of transformation (XOR, substitution, bit shuffling, transposition)
- Base64 output with character rotation
- Unicode support including emojis
- Fast and efficient (sub-microsecond for typical strings)
- Fully reversible

**Use Cases:**
- URL/token obfuscation
- Hiding data from casual inspection
- Anti-tampering for non-sensitive data
- Creating non-obvious identifiers

**⚠️ Warning:** This is obfuscation, NOT encryption. Do not use for sensitive data like passwords or personal information.

See [obfuscator/README.md](obfuscator/README.md) for detailed documentation.

### `localfs` - Local File System Service

A comprehensive service for local file system operations, providing a modern, service-oriented interface.

```go
import "github.com/kawai-network/veridium/pkg/localfs"

// Create service
service := localfs.NewService()
ctx := context.Background()

// Write and read files
service.WriteFile(ctx, localfs.WriteFileParams{
    Path:    "/path/to/file.txt",
    Content: "Hello, World!",
})

result, _ := service.ReadFile(ctx, localfs.ReadFileParams{
    Path: "/path/to/file.txt",
})

// Edit files with search/replace
service.EditFile(ctx, localfs.EditFileParams{
    FilePath:   "/path/to/file.txt",
    OldString:  "old",
    NewString:  "new",
    ReplaceAll: true,
})

// Run shell commands
result, _ := service.RunCommand(ctx, localfs.RunCommandParams{
    Command: "ls -la",
})

// Search content (grep-like)
result, _ := service.GrepContent(ctx, localfs.GrepContentParams{
    Pattern: "search term",
    Path:    "/path/to/search",
})
```

**Features:**
- File operations: read, write, edit, list, search, move, rename
- Shell command execution (sync and async)
- Content search (grep) and glob pattern matching
- Cross-platform support (macOS, Linux, Windows)
- Context-aware with cancellation support
- Comprehensive test coverage

See [localfs/README.md](localfs/README.md) for detailed documentation.

### `nodefs` - File System Operations

Provides Node.js `fs` module equivalents:

```go
import "github.com/kawai-network/veridium/pkg/nodefs"

// Synchronous file operations
exists := nodefs.FileExistsSync("path/to/file")
data, err := nodefs.ReadFileSync("path/to/file")
err := nodefs.WriteFileSync("path/to/file", []byte("content"))

// Directory operations
entries, err := nodefs.ReadDirSync("path/to/dir")
info, err := nodefs.StatSync("path/to/file")

// Directory creation
err := nodefs.MkdirSync("path/to/dir")

// File removal
err := nodefs.RmSync("path/to/file")  // Recursive removal
err := nodefs.UnlinkSync("path/to/file")  // Single file
```

### `nodepath` - Path Operations

Provides Node.js `path` module equivalents:

```go
import "github.com/kawai-network/veridium/pkg/nodepath"

// Path manipulation
basename := nodepath.Basename("/path/to/file.txt")  // "file.txt"
dirname := nodepath.Dirname("/path/to/file.txt")    // "/path/to"
ext := nodepath.Extname("file.txt")                // ".txt"
fullPath := nodepath.Join("path", "to", "file.txt") // "path/to/file.txt"

// Path parsing
parsed := nodepath.Parse("/path/to/file.txt")
// parsed = {Root: "/", Dir: "/path/to", Base: "file.txt", Ext: ".txt", Name: "file"}

// Cross-platform operations
isAbs := nodepath.IsAbsolute("/absolute/path")  // true on Unix, false on Windows
resolved := nodepath.Resolve("relative", "path")  // Absolute path
```

### `nodebuffer` - Buffer Operations

Provides Node.js `Buffer` equivalents:

```go
import "github.com/kawai-network/veridium/pkg/nodebuffer"

// Creating buffers
buf := nodebuffer.New(1024)  // Allocate buffer of size 1024
buf := nodebuffer.Alloc(1024)  // Same as New()
buf := nodebuffer.From("hello world")  // From string
buf := nodebuffer.From([]byte{1, 2, 3})  // From byte slice

// Buffer operations
data := buf.ToBytes()  // Get underlying byte slice
str := buf.ToString()  // Convert to string
str := buf.ToString("base64")  // Convert with encoding

// Buffer manipulation
buf.Fill(0)  // Fill with zeros
buf.Copy(targetBuf, 0, 0, 10)  // Copy to another buffer
slice := buf.Slice(0, 10)  // Create slice

// Searching
index := buf.IndexOf("search")  // Find substring
contains := buf.Contains("search")  // Check if contains
```

### `nodeexec` - Child Process Execution

Provides Node.js `child_process` equivalents:

```go
import "github.com/kawai-network/veridium/pkg/nodeexec"

// Execute commands synchronously
result, err := nodeexec.ExecSync("npm --version")
if result.Success {
    fmt.Println("STDOUT:", result.Stdout)
    fmt.Println("Exit code:", result.Code)
}

// Execute with options
result, err := nodeexec.ExecSync("ls -la", &nodeexec.ExecOptions{
    Cwd: "/tmp",
    Env: map[string]string{"NODE_ENV": "production"},
    Timeout: 5 * time.Second,
})

// Spawn processes
resultChan, err := nodeexec.Spawn("node", []string{"script.js"})
result := <-resultChan  // Wait for completion

// Which command
path, err := nodeexec.Which("node")  // Find executable in PATH
```

### `nodeos` - Operating System Operations

Provides Node.js `os` module equivalents:

```go
import "github.com/kawai-network/veridium/pkg/nodeos"

// System information
arch := nodeos.Arch()        // "x64", "arm64", etc.
platform := nodeos.Platform() // "darwin", "linux", "win32"
hostname := nodeos.Hostname()
homedir := nodeos.Homedir()
tmpdir := nodeos.Tmpdir()

// CPU information
cpus := nodeos.Cpus()
fmt.Printf("CPU cores: %d\n", len(cpus))

// Memory information (limited in Go stdlib)
totalMem := nodeos.Totalmem()  // May return 0 (not available)
freeMem := nodeos.Freemem()    // May return 0 (not available)

// User information
userInfo, err := nodeos.UserInfo()
if err == nil {
    fmt.Printf("User: %s, Home: %s\n", userInfo.Username, userInfo.Homedir)
}

// Process priority (limited support)
priority := nodeos.GetPriority()  // Get current process priority
```

## Compatibility Notes

### Not Supported in Go Standard Library

Some Node.js APIs are not directly available in Go's standard library:

1. **Memory Information**: `os.totalmem()`, `os.freemem()` - Go stdlib doesn't provide system memory info
2. **Load Average**: `os.loadavg()` - Requires reading `/proc/loadavg` on Linux
3. **Network Interfaces**: `os.networkInterfaces()` - Requires platform-specific code
4. **Process Priority**: `os.getPriority()`, `os.setPriority()` - Limited support
5. **User ID/GID**: `os.userInfo().uid`, `os.userInfo().gid` - Platform-specific

### Platform-Specific Behavior

- Path separators: Go uses `/` on all platforms internally, but handles platform differences
- Executable extensions: `nodeexec.Which()` handles `.exe` on Windows
- Environment variables: Available through `os.Getenv()` but not directly settable in some contexts

## Error Handling

All functions follow Go conventions:
- Return `(result, error)` pairs
- Use `nil` for successful operations
- Provide descriptive error messages

## Performance Considerations

- These packages use Go's standard library for optimal performance
- Buffer operations are more efficient than Node.js (no V8 overhead)
- File operations use native OS APIs
- Memory management is handled by Go's garbage collector

## Thread Safety

- All packages are safe for concurrent use
- No global state is modified
- Functions can be called from multiple goroutines safely

## Testing

Run tests for all packages:

```bash
go test ./pkg/...
```

Run benchmarks:

```bash
go test -bench=. ./pkg/...
```

## Contributing

When adding new functionality:

1. Follow Node.js API signatures as closely as possible
2. Provide comprehensive documentation
3. Include tests and benchmarks
4. Handle edge cases and errors appropriately
5. Maintain thread safety
