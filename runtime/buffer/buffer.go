package buffer

import (
	"encoding/base64"
	"encoding/hex"
)

// Buffer represents a Node.js-style Buffer for binary data
type Buffer struct {
	data []byte
}

// New creates a new Buffer from a byte slice
func New(data []byte) *Buffer {
	return &Buffer{data: data}
}

// From creates a Buffer from various input types
func From(input interface{}) *Buffer {
	switch v := input.(type) {
	case []byte:
		return &Buffer{data: v}
	case string:
		return &Buffer{data: []byte(v)}
	case *Buffer:
		// Create a copy
		dataCopy := make([]byte, len(v.data))
		copy(dataCopy, v.data)
		return &Buffer{data: dataCopy}
	default:
		return &Buffer{data: []byte{}}
	}
}

// Alloc creates a new Buffer of the specified size, filled with zeros
func Alloc(size int) *Buffer {
	return &Buffer{data: make([]byte, size)}
}

// AllocUnsafe creates a new Buffer of the specified size without initializing
// Note: In Go, this is the same as Alloc since Go initializes memory to zero
func AllocUnsafe(size int) *Buffer {
	return &Buffer{data: make([]byte, size)}
}

// Length returns the length of the buffer
func (b *Buffer) Length() int {
	return len(b.data)
}

// ToString converts the buffer to a string with optional encoding
func (b *Buffer) ToString(encoding string) string {
	switch encoding {
	case "base64":
		return base64.StdEncoding.EncodeToString(b.data)
	case "hex":
		return hex.EncodeToString(b.data)
	case "utf8", "utf-8", "":
		return string(b.data)
	default:
		return string(b.data)
	}
}

// ToJSON returns a JSON representation of the buffer
func (b *Buffer) ToJSON() map[string]interface{} {
	// Convert to array of byte values
	byteArray := make([]int, len(b.data))
	for i, v := range b.data {
		byteArray[i] = int(v)
	}

	return map[string]interface{}{
		"type": "Buffer",
		"data": byteArray,
	}
}

// Slice returns a new buffer that references the same memory
func (b *Buffer) Slice(start, end int) *Buffer {
	if start < 0 {
		start = 0
	}
	if end > len(b.data) {
		end = len(b.data)
	}
	if start > end {
		start = end
	}

	return &Buffer{data: b.data[start:end]}
}

// Copy copies data from source buffer to this buffer
func (b *Buffer) Copy(target *Buffer, targetStart, sourceStart, sourceEnd int) int {
	if sourceEnd > len(b.data) {
		sourceEnd = len(b.data)
	}

	copied := copy(target.data[targetStart:], b.data[sourceStart:sourceEnd])
	return copied
}

// Write writes a string to the buffer at the specified offset
func (b *Buffer) Write(str string, offset int) int {
	if offset >= len(b.data) {
		return 0
	}

	written := copy(b.data[offset:], []byte(str))
	return written
}

// ReadUInt8 reads an unsigned 8-bit integer at the specified offset
func (b *Buffer) ReadUInt8(offset int) uint8 {
	if offset >= len(b.data) {
		return 0
	}
	return b.data[offset]
}

// ReadUInt16LE reads an unsigned 16-bit integer (little-endian) at the specified offset
func (b *Buffer) ReadUInt16LE(offset int) uint16 {
	if offset+1 >= len(b.data) {
		return 0
	}
	return uint16(b.data[offset]) | uint16(b.data[offset+1])<<8
}

// ReadUInt32LE reads an unsigned 32-bit integer (little-endian) at the specified offset
func (b *Buffer) ReadUInt32LE(offset int) uint32 {
	if offset+3 >= len(b.data) {
		return 0
	}
	return uint32(b.data[offset]) |
		uint32(b.data[offset+1])<<8 |
		uint32(b.data[offset+2])<<16 |
		uint32(b.data[offset+3])<<24
}

// WriteUInt8 writes an unsigned 8-bit integer at the specified offset
func (b *Buffer) WriteUInt8(value uint8, offset int) {
	if offset < len(b.data) {
		b.data[offset] = value
	}
}

// WriteUInt16LE writes an unsigned 16-bit integer (little-endian) at the specified offset
func (b *Buffer) WriteUInt16LE(value uint16, offset int) {
	if offset+1 < len(b.data) {
		b.data[offset] = byte(value)
		b.data[offset+1] = byte(value >> 8)
	}
}

// WriteUInt32LE writes an unsigned 32-bit integer (little-endian) at the specified offset
func (b *Buffer) WriteUInt32LE(value uint32, offset int) {
	if offset+3 < len(b.data) {
		b.data[offset] = byte(value)
		b.data[offset+1] = byte(value >> 8)
		b.data[offset+2] = byte(value >> 16)
		b.data[offset+3] = byte(value >> 24)
	}
}

// Equals checks if two buffers have the same content
func (b *Buffer) Equals(other *Buffer) bool {
	if len(b.data) != len(other.data) {
		return false
	}

	for i := range b.data {
		if b.data[i] != other.data[i] {
			return false
		}
	}

	return true
}

// Compare compares two buffers and returns:
// - 0 if equal
// - -1 if b comes before other
// - 1 if b comes after other
func (b *Buffer) Compare(other *Buffer) int {
	minLen := len(b.data)
	if len(other.data) < minLen {
		minLen = len(other.data)
	}

	for i := 0; i < minLen; i++ {
		if b.data[i] < other.data[i] {
			return -1
		}
		if b.data[i] > other.data[i] {
			return 1
		}
	}

	if len(b.data) < len(other.data) {
		return -1
	}
	if len(b.data) > len(other.data) {
		return 1
	}

	return 0
}

// Concat concatenates multiple buffers into one
func Concat(buffers []*Buffer) *Buffer {
	totalLength := 0
	for _, buf := range buffers {
		totalLength += len(buf.data)
	}

	result := make([]byte, totalLength)
	offset := 0

	for _, buf := range buffers {
		copy(result[offset:], buf.data)
		offset += len(buf.data)
	}

	return &Buffer{data: result}
}

// IsBuffer checks if the given value is a Buffer
func IsBuffer(obj interface{}) bool {
	_, ok := obj.(*Buffer)
	return ok
}

// IsEncoding checks if the given encoding is supported
func IsEncoding(encoding string) bool {
	switch encoding {
	case "utf8", "utf-8", "base64", "hex", "ascii":
		return true
	default:
		return false
	}
}

// ByteLength returns the byte length of a string when encoded
func ByteLength(str string, encoding string) int {
	switch encoding {
	case "utf8", "utf-8", "":
		return len(str)
	case "base64":
		// Approximate - actual encoding might differ
		return len(str) * 3 / 4
	case "hex":
		return len(str) / 2
	default:
		return len(str)
	}
}

// Data returns the underlying byte slice
func (b *Buffer) Data() []byte {
	return b.data
}
