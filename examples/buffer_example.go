package main

import (
	"fmt"

	"github.com/kawai-network/veridium/pkg/nodebuffer"
)

func main() {
	// Example: Creating buffers
	buf1 := nodebuffer.New(10)
	fmt.Printf("New buffer length: %d\n", buf1.Length())

	buf2 := nodebuffer.Alloc(20)
	fmt.Printf("Alloc buffer length: %d\n", buf2.Length())

	// Example: Creating from different inputs
	buf3 := nodebuffer.From("Hello World")
	fmt.Printf("String buffer: %s (length: %d)\n", buf3.ToString(), buf3.Length())

	buf4 := nodebuffer.From([]byte{72, 101, 108, 108, 111}) // "Hello" in bytes
	fmt.Printf("Byte buffer: %s\n", buf4.ToString())

	// Example: Base64 encoding/decoding
	buf5 := nodebuffer.From("Hello World", "base64")
	fmt.Printf("Base64 buffer: %s\n", buf5.ToString())

	// Example: Buffer operations
	buf6 := nodebuffer.From("Hello Node.js Buffer!")
	fmt.Printf("Original: %s\n", buf6.ToString())

	// Fill buffer
	buf6.Fill(42) // Fill with '*' (ASCII 42)
	fmt.Printf("After fill: %s\n", buf6.ToString())

	// Reset for more examples
	buf7 := nodebuffer.From("The quick brown fox jumps over the lazy dog")

	// Slice operation
	slice := buf7.Slice(4, 9) // "quick"
	fmt.Printf("Slice (4,9): %s\n", slice.ToString())

	// Copy operation
	source := nodebuffer.From("Hello")
	dest := nodebuffer.New(10)
	copied := source.Copy(dest, 0)
	fmt.Printf("Copied %d bytes: %s\n", copied, dest.ToString())

	// Concatenation
	bufs := []*nodebuffer.Buffer{
		nodebuffer.From("Hello "),
		nodebuffer.From("World"),
		nodebuffer.From("!"),
	}
	concatenated := nodebuffer.Concat(bufs)
	fmt.Printf("Concatenated: %s\n", concatenated.ToString())

	// Search operations
	searchBuf := nodebuffer.From("The quick brown fox jumps over the lazy dog")
	index := searchBuf.IndexOf("fox")
	fmt.Printf("Index of 'fox': %d\n", index)

	contains := searchBuf.Contains("lazy")
	fmt.Printf("Contains 'lazy': %v\n", contains)

	// Encoding operations
	original := "Hello World 🌍"
	buf8 := nodebuffer.From(original)

	base64Str := buf8.ToString("base64")
	fmt.Printf("Base64 encoded: %s\n", base64Str)

	hexStr := buf8.ToString("hex")
	fmt.Printf("Hex encoded: %s\n", hexStr[:50]+"...") // Truncate for display

	// Read/write operations
	buf9 := nodebuffer.New(8)
	buf9.WriteInt8(72, 0) // 'H'
	buf9.WriteInt8(105, 1) // 'i'
	val1, _ := buf9.ReadInt8(0)
	val2, _ := buf9.ReadInt8(1)
	fmt.Printf("Read values: %d, %d -> %s\n", val1, val2, buf9.ToString())

	// JSON representation
	json := buf9.ToJSON()
	fmt.Printf("JSON representation: %+v\n", json)

	// Buffer comparison
	buf10 := nodebuffer.From("test")
	buf11 := nodebuffer.From("test")
	buf12 := nodebuffer.From("different")

	fmt.Printf("buf10 equals buf11: %v\n", buf10.Equals(buf11))
	fmt.Printf("buf10 equals buf12: %v\n", buf10.Equals(buf12))

	// Byte length calculation
	length1 := nodebuffer.ByteLength("Hello World")
	length2 := nodebuffer.ByteLength("Hello World", "base64")
	fmt.Printf("UTF-8 length: %d, Base64 length: %d\n", length1, length2)
}
