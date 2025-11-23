package mget

import (
	"bytes"
	"io"
	"testing"
)

func TestProgressWriter(t *testing.T) {
	tests := []struct {
		name        string
		totalBytes  int64
		description string
		writeData   []byte
	}{
		{
			name:        "Known total bytes",
			totalBytes:  1000,
			description: "Test download",
			writeData:   []byte("test data"),
		},
		{
			name:        "Unknown total bytes",
			totalBytes:  -1,
			description: "Test download unknown",
			writeData:   []byte("test data"),
		},
		{
			name:        "Zero total bytes",
			totalBytes:  0,
			description: "Test download zero",
			writeData:   []byte("test data"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := ProgressWriter(tt.totalBytes, tt.description)

			if writer == nil {
				t.Error("Expected non-nil writer")
			}

			// Test that we can write to it
			n, err := writer.Write(tt.writeData)
			if err != nil {
				t.Errorf("Expected no error writing, got %v", err)
			}
			if n != len(tt.writeData) {
				t.Errorf("Expected to write %d bytes, wrote %d", len(tt.writeData), n)
			}
		})
	}
}

func TestProgressWriterUnknown(t *testing.T) {
	description := "Unknown progress test"
	writer := ProgressWriterUnknown(description)

	if writer == nil {
		t.Error("Expected non-nil writer")
	}

	// Test that we can write to it
	testData := []byte("test data")
	n, err := writer.Write(testData)
	if err != nil {
		t.Errorf("Expected no error writing, got %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}
}

func TestProgressWriter_WriteMultiple(t *testing.T) {
	writer := ProgressWriter(1000, "Multiple writes test")

	// Write multiple times
	data1 := []byte("chunk1")
	data2 := []byte("chunk2")
	data3 := []byte("chunk3")

	_, err := writer.Write(data1)
	if err != nil {
		t.Errorf("Error writing first chunk: %v", err)
	}

	_, err = writer.Write(data2)
	if err != nil {
		t.Errorf("Error writing second chunk: %v", err)
	}

	_, err = writer.Write(data3)
	if err != nil {
		t.Errorf("Error writing third chunk: %v", err)
	}
}

func TestProgressWriter_ImplementsWriter(t *testing.T) {
	writer := ProgressWriter(1000, "Interface test")

	// Verify it implements io.Writer
	var _ io.Writer = writer
}

func TestProgressWriter_WithMultiWriter(t *testing.T) {
	// Test that ProgressWriter can be used with io.MultiWriter
	var buf bytes.Buffer
	progressWriter := ProgressWriter(1000, "MultiWriter test")
	multiWriter := io.MultiWriter(&buf, progressWriter)

	testData := []byte("test data")
	n, err := multiWriter.Write(testData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if n != len(testData) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
	}

	// Verify data was written to buffer
	if !bytes.Equal(buf.Bytes(), testData) {
		t.Errorf("Expected buffer to contain %q, got %q", testData, buf.Bytes())
	}
}
