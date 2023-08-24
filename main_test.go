package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestReadGitignore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test-gitignore")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ioutil.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte("test.txt\n#comment\ntest2.txt"), 0644)

	patterns, err := readGitignore(tempDir)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	expectedPatterns := []string{"test.txt", "test2.txt"}
	for i, p := range patterns {
		if p != expectedPatterns[i] {
			t.Errorf("Expected pattern %q but got %q", expectedPatterns[i], p)
		}
	}
}

func TestShouldIgnore(t *testing.T) {
	patterns := []string{"*.txt", "test/*", ".gitignore", ".git/*"}
	tests := []struct {
		filename string
		ignore   bool
	}{
		{"test.txt", true},
		{"image.jpg", false},
		{"test/data.txt", true},
		{"test/image.jpg", true},
		{"data.yaml", false},
		{".gitignore", true},
		{".git/config", true},
	}

	for _, tt := range tests {
		if shouldIgnore(tt.filename, patterns) != tt.ignore {
			t.Errorf("Expected shouldIgnore(%q) to be %v", tt.filename, tt.ignore)
		}
	}
}

func TestIsBinaryFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "test-binary")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ioutil.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("Hello, World!"), 0644)
	ioutil.WriteFile(filepath.Join(tempDir, "test.bin"), []byte{0x48, 0x65, 0x00, 0x6F}, 0644)

	tests := []struct {
		filename string
		isBinary bool
	}{
		{filepath.Join(tempDir, "test.txt"), false},
		{filepath.Join(tempDir, "test.bin"), true},
	}

	for _, tt := range tests {
		if isBinaryFile(tt.filename) != tt.isBinary {
			t.Errorf("Expected isBinaryFile(%q) to be %v", tt.filename, tt.isBinary)
		}
	}
}
