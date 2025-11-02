package buffer

import (
	"testing"
)

func TestNew(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	buf := New(data)

	if buf == nil {
		t.Fatal("New returned nil")
	}

	if buf.Length() != 5 {
		t.Errorf("Expected length 5, got %d", buf.Length())
	}
}

func TestFrom(t *testing.T) {
	// Test from string
	buf1 := From("hello")
	if buf1.ToString("utf8") != "hello" {
		t.Errorf("From string: expected 'hello', got '%s'", buf1.ToString("utf8"))
	}

	// Test from []byte
	buf2 := From([]byte{1, 2, 3})
	if buf2.Length() != 3 {
		t.Errorf("From []byte: expected length 3, got %d", buf2.Length())
	}

	// Test from Buffer
	buf3 := From(buf1)
	if buf3.ToString("utf8") != "hello" {
		t.Errorf("From Buffer: expected 'hello', got '%s'", buf3.ToString("utf8"))
	}

	// Verify it's a copy, not the same buffer
	buf3.WriteUInt8(88, 0) // Modify buf3
	if buf1.ReadUInt8(0) == 88 {
		t.Error("From Buffer should create a copy, not reference the same data")
	}
}

func TestAlloc(t *testing.T) {
	buf := Alloc(10)

	if buf.Length() != 10 {
		t.Errorf("Expected length 10, got %d", buf.Length())
	}

	// Check all bytes are zero
	for i := 0; i < buf.Length(); i++ {
		if buf.ReadUInt8(i) != 0 {
			t.Errorf("Byte at %d should be 0, got %d", i, buf.ReadUInt8(i))
		}
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		encoding string
		expected string
	}{
		{
			name:     "UTF-8",
			data:     []byte("hello"),
			encoding: "utf8",
			expected: "hello",
		},
		{
			name:     "Base64",
			data:     []byte("hello"),
			encoding: "base64",
			expected: "aGVsbG8=",
		},
		{
			name:     "Hex",
			data:     []byte{0xDE, 0xAD, 0xBE, 0xEF},
			encoding: "hex",
			expected: "deadbeef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := New(tt.data)
			result := buf.ToString(tt.encoding)

			if result != tt.expected {
				t.Errorf("ToString(%s): expected '%s', got '%s'", tt.encoding, tt.expected, result)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	buf := From("hello world")

	slice := buf.Slice(0, 5)
	if slice.ToString("utf8") != "hello" {
		t.Errorf("Slice(0, 5): expected 'hello', got '%s'", slice.ToString("utf8"))
	}

	slice2 := buf.Slice(6, 11)
	if slice2.ToString("utf8") != "world" {
		t.Errorf("Slice(6, 11): expected 'world', got '%s'", slice2.ToString("utf8"))
	}
}

func TestWrite(t *testing.T) {
	buf := Alloc(10)

	written := buf.Write("hello", 0)
	if written != 5 {
		t.Errorf("Expected to write 5 bytes, wrote %d", written)
	}

	if buf.ToString("utf8")[:5] != "hello" {
		t.Errorf("Expected 'hello', got '%s'", buf.ToString("utf8")[:5])
	}
}

func TestReadWriteUInt8(t *testing.T) {
	buf := Alloc(1)

	buf.WriteUInt8(42, 0)
	value := buf.ReadUInt8(0)

	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
}

func TestReadWriteUInt16LE(t *testing.T) {
	buf := Alloc(2)

	buf.WriteUInt16LE(0x1234, 0)
	value := buf.ReadUInt16LE(0)

	if value != 0x1234 {
		t.Errorf("Expected 0x1234, got 0x%x", value)
	}

	// Verify little-endian byte order
	if buf.ReadUInt8(0) != 0x34 {
		t.Errorf("First byte should be 0x34 (LE), got 0x%x", buf.ReadUInt8(0))
	}
	if buf.ReadUInt8(1) != 0x12 {
		t.Errorf("Second byte should be 0x12 (LE), got 0x%x", buf.ReadUInt8(1))
	}
}

func TestReadWriteUInt32LE(t *testing.T) {
	buf := Alloc(4)

	buf.WriteUInt32LE(0x12345678, 0)
	value := buf.ReadUInt32LE(0)

	if value != 0x12345678 {
		t.Errorf("Expected 0x12345678, got 0x%x", value)
	}

	// Verify little-endian byte order
	expected := []byte{0x78, 0x56, 0x34, 0x12}
	for i := 0; i < 4; i++ {
		if buf.ReadUInt8(i) != expected[i] {
			t.Errorf("Byte %d should be 0x%x, got 0x%x", i, expected[i], buf.ReadUInt8(i))
		}
	}
}

func TestEquals(t *testing.T) {
	buf1 := From("hello")
	buf2 := From("hello")
	buf3 := From("world")

	if !buf1.Equals(buf2) {
		t.Error("buf1 and buf2 should be equal")
	}

	if buf1.Equals(buf3) {
		t.Error("buf1 and buf3 should not be equal")
	}
}

func TestCompare(t *testing.T) {
	buf1 := From("abc")
	buf2 := From("abc")
	buf3 := From("abd")
	buf4 := From("ab")

	if buf1.Compare(buf2) != 0 {
		t.Error("buf1 and buf2 should be equal (compare = 0)")
	}

	if buf1.Compare(buf3) >= 0 {
		t.Error("buf1 should be less than buf3 (compare < 0)")
	}

	if buf3.Compare(buf1) <= 0 {
		t.Error("buf3 should be greater than buf1 (compare > 0)")
	}

	if buf1.Compare(buf4) <= 0 {
		t.Error("buf1 should be greater than buf4 (longer)")
	}
}

func TestConcat(t *testing.T) {
	buf1 := From("hello")
	buf2 := From(" ")
	buf3 := From("world")

	result := Concat([]*Buffer{buf1, buf2, buf3})

	if result.ToString("utf8") != "hello world" {
		t.Errorf("Concat: expected 'hello world', got '%s'", result.ToString("utf8"))
	}

	if result.Length() != 11 {
		t.Errorf("Concat: expected length 11, got %d", result.Length())
	}
}

func TestIsBuffer(t *testing.T) {
	buf := From("test")

	if !IsBuffer(buf) {
		t.Error("IsBuffer should return true for Buffer")
	}

	if IsBuffer("string") {
		t.Error("IsBuffer should return false for string")
	}

	if IsBuffer(123) {
		t.Error("IsBuffer should return false for int")
	}
}

func TestIsEncoding(t *testing.T) {
	validEncodings := []string{"utf8", "utf-8", "base64", "hex", "ascii"}

	for _, enc := range validEncodings {
		if !IsEncoding(enc) {
			t.Errorf("IsEncoding should return true for '%s'", enc)
		}
	}

	if IsEncoding("invalid") {
		t.Error("IsEncoding should return false for 'invalid'")
	}
}

func TestByteLength(t *testing.T) {
	tests := []struct {
		str      string
		encoding string
		expected int
	}{
		{"hello", "utf8", 5},
		{"hello", "utf-8", 5},
		{"hello", "", 5},
	}

	for _, tt := range tests {
		result := ByteLength(tt.str, tt.encoding)
		if result != tt.expected {
			t.Errorf("ByteLength('%s', '%s'): expected %d, got %d", tt.str, tt.encoding, tt.expected, result)
		}
	}
}

func TestCopy(t *testing.T) {
	source := From("hello world")
	target := Alloc(5)

	copied := source.Copy(target, 0, 0, 5)

	if copied != 5 {
		t.Errorf("Expected to copy 5 bytes, copied %d", copied)
	}

	if target.ToString("utf8") != "hello" {
		t.Errorf("Copy: expected 'hello', got '%s'", target.ToString("utf8"))
	}
}

func TestToJSON(t *testing.T) {
	buf := From([]byte{1, 2, 3})
	json := buf.ToJSON()

	if json["type"] != "Buffer" {
		t.Errorf("ToJSON type: expected 'Buffer', got '%v'", json["type"])
	}

	data, ok := json["data"].([]int)
	if !ok {
		t.Fatal("ToJSON data is not []int")
	}

	if len(data) != 3 {
		t.Errorf("ToJSON data length: expected 3, got %d", len(data))
	}

	expected := []int{1, 2, 3}
	for i, v := range expected {
		if data[i] != v {
			t.Errorf("ToJSON data[%d]: expected %d, got %d", i, v, data[i])
		}
	}
}

func TestData(t *testing.T) {
	original := []byte{1, 2, 3, 4, 5}
	buf := New(original)

	data := buf.Data()

	if len(data) != len(original) {
		t.Errorf("Data length: expected %d, got %d", len(original), len(data))
	}

	for i := range original {
		if data[i] != original[i] {
			t.Errorf("Data[%d]: expected %d, got %d", i, original[i], data[i])
		}
	}
}
