package main

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/kawai-network/veridium/pkg/nodebuffer"
)

// NodeBufferService provides Node.js Buffer equivalents as a Wails service
type NodeBufferService struct{}

// New creates a new Buffer with the given size
func (n *NodeBufferService) New(size int) string {
	buf := nodebuffer.New(size)
	return base64.StdEncoding.EncodeToString(buf.ToBytes())
}

// Alloc creates a zero-filled Buffer of the specified size
func (n *NodeBufferService) Alloc(size int) string {
	buf := nodebuffer.Alloc(size)
	return base64.StdEncoding.EncodeToString(buf.ToBytes())
}

// AllocUnsafe creates a Buffer of the specified size without zeroing
func (n *NodeBufferService) AllocUnsafe(size int) string {
	buf := nodebuffer.AllocUnsafe(size)
	return base64.StdEncoding.EncodeToString(buf.ToBytes())
}

// From creates a Buffer from various input types
func (n *NodeBufferService) From(data interface{}, encoding ...string) string {
	buf := nodebuffer.From(data, encoding...)
	return base64.StdEncoding.EncodeToString(buf.ToBytes())
}

// Concat concatenates multiple Buffers
func (n *NodeBufferService) Concat(buffers []string, totalLength ...int) (string, error) {
	bufs := make([]*nodebuffer.Buffer, len(buffers))
	for i, b64 := range buffers {
		data, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			return "", err
		}
		bufs[i] = nodebuffer.From(data)
	}

	result := nodebuffer.Concat(bufs, totalLength...)
	return base64.StdEncoding.EncodeToString(result.ToBytes()), nil
}

// IsBuffer checks if the given value is a Buffer
func (n *NodeBufferService) IsBuffer(obj interface{}) bool {
	return nodebuffer.IsBuffer(obj)
}

// ByteLength returns the byte length of the input
func (n *NodeBufferService) ByteLength(data interface{}, encoding ...string) int {
	return nodebuffer.ByteLength(data, encoding...)
}

// Equals compares two Buffers
func (n *NodeBufferService) Equals(buf1, buf2 string) bool {
	data1, err := base64.StdEncoding.DecodeString(buf1)
	if err != nil {
		return false
	}
	data2, err := base64.StdEncoding.DecodeString(buf2)
	if err != nil {
		return false
	}

	buffer1 := nodebuffer.From(data1)
	buffer2 := nodebuffer.From(data2)
	return buffer1.Equals(buffer2)
}

// Compare compares two Buffers
func (n *NodeBufferService) Compare(buf1, buf2 string) int {
	data1, err := base64.StdEncoding.DecodeString(buf1)
	if err != nil {
		return -1
	}
	data2, err := base64.StdEncoding.DecodeString(buf2)
	if err != nil {
		return 1
	}

	buffer1 := nodebuffer.From(data1)
	buffer2 := nodebuffer.From(data2)
	return nodebuffer.Compare(buffer1, buffer2)
}

// Fill fills the Buffer with specified value
func (n *NodeBufferService) Fill(buffer, value string, offsetAndEnd ...int) (string, error) {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return "", err
	}

	buf := nodebuffer.From(data)
	result := buf.Fill(value, offsetAndEnd...)
	return base64.StdEncoding.EncodeToString(result.ToBytes()), nil
}

// Contains checks if the Buffer contains the specified value
func (n *NodeBufferService) Contains(buffer, value string, byteOffset ...int) bool {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return false
	}

	buf := nodebuffer.From(data)
	return buf.Contains(value, byteOffset...)
}

// IndexOf finds the index of the specified value
func (n *NodeBufferService) IndexOf(buffer, value string, byteOffset ...int) int {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return -1
	}

	buf := nodebuffer.From(data)
	return buf.IndexOf(value, byteOffset...)
}

// LastIndexOf finds the last index of the specified value
func (n *NodeBufferService) LastIndexOf(buffer, value string, byteOffset ...int) int {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return -1
	}

	buf := nodebuffer.From(data)
	return buf.LastIndexOf(value, byteOffset...)
}

// Slice returns a slice of the Buffer
func (n *NodeBufferService) Slice(buffer string, startAndEnd ...int) string {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return ""
	}

	buf := nodebuffer.From(data)
	slice := buf.Slice(startAndEnd...)
	return base64.StdEncoding.EncodeToString(slice.ToBytes())
}

// Copy copies data from one Buffer to another
func (n *NodeBufferService) Copy(source, target string, targetStart ...int) (int, error) {
	sourceData, err := base64.StdEncoding.DecodeString(source)
	if err != nil {
		return 0, err
	}
	targetData, err := base64.StdEncoding.DecodeString(target)
	if err != nil {
		return 0, err
	}

	sourceBuf := nodebuffer.From(sourceData)
	targetBuf := nodebuffer.From(targetData)

	copied := sourceBuf.Copy(targetBuf, targetStart...)
	return copied, nil
}

// Length returns the length of the Buffer
func (n *NodeBufferService) Length(buffer string) int {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return 0
	}

	buf := nodebuffer.From(data)
	return buf.Length()
}

// ToString converts the Buffer to a string
func (n *NodeBufferService) ToString(buffer string, encoding ...string) string {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return ""
	}

	buf := nodebuffer.From(data)
	return buf.ToString(encoding...)
}

// ToHex converts the Buffer to hex string
func (n *NodeBufferService) ToHex(buffer string) string {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(data)
}

// FromHex creates a Buffer from hex string
func (n *NodeBufferService) FromHex(hexStr string) (string, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	buf := nodebuffer.From(data)
	return base64.StdEncoding.EncodeToString(buf.ToBytes()), nil
}

// ToBase64 converts the Buffer to base64 string
func (n *NodeBufferService) ToBase64(buffer string) string {
	// The buffer is already base64 encoded, just return it
	return buffer
}

// FromBase64 creates a Buffer from base64 string
func (n *NodeBufferService) FromBase64(b64Str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return "", err
	}

	buf := nodebuffer.From(data)
	return base64.StdEncoding.EncodeToString(buf.ToBytes()), nil
}

// ReadInt8 reads a signed 8-bit integer
func (n *NodeBufferService) ReadInt8(buffer string, offset ...int) (int8, error) {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return 0, err
	}

	buf := nodebuffer.From(data)
	return buf.ReadInt8(offset...)
}

// ReadUInt8 reads an unsigned 8-bit integer
func (n *NodeBufferService) ReadUInt8(buffer string, offset ...int) (uint8, error) {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return 0, err
	}

	buf := nodebuffer.From(data)
	return buf.ReadUInt8(offset...)
}

// WriteInt8 writes a signed 8-bit integer
func (n *NodeBufferService) WriteInt8(buffer string, value int8, offset ...int) (string, error) {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return "", err
	}

	buf := nodebuffer.From(data)
	buf.WriteInt8(value, offset...)
	result := base64.StdEncoding.EncodeToString(buf.ToBytes())
	return result, nil
}

// WriteUInt8 writes an unsigned 8-bit integer
func (n *NodeBufferService) WriteUInt8(buffer string, value uint8, offset ...int) (string, error) {
	data, err := base64.StdEncoding.DecodeString(buffer)
	if err != nil {
		return "", err
	}

	buf := nodebuffer.From(data)
	buf.WriteUInt8(value, offset...)
	result := base64.StdEncoding.EncodeToString(buf.ToBytes())
	return result, nil
}

// Get constants
func (n *NodeBufferService) MaxLength() int {
	return nodebuffer.MAX_LENGTH
}

func (n *NodeBufferService) MaxStringLength() int {
	return nodebuffer.MAX_STRING_LENGTH
}
