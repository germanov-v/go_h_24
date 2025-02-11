package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	tempDir, err := os.MkdirTemp("", "test_tempdor")
	if err != nil {
		t.Fatalf("creating dir was failed: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {

		}
	}(tempDir)

	/// 0644 - owner rw, group - r, others - r
	if err := os.WriteFile(filepath.Join(tempDir, "test1"), []byte("aaa \t\n123"), 0644); err != nil {
		t.Fatalf("creating dir was failed aaaa with content: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tempDir, "test2"), []byte("bbb \t\n123"), 0644); err != nil {
		t.Fatalf("creating dir was failed bbb with content: %v", err)
	}

	envTest, err := ReadDir(tempDir)
	if err != nil {
		t.Fatalf("tempDir reading was failed: %v", err)
	}

	if ev, ok := envTest["test1"]; !ok {
		t.Errorf("var test1 reading failed")
	} else {
		expected := "aaa"
		if ev.Value != expected {
			t.Errorf("test1 value expected %q, but we have %q", expected, ev.Value)
		}
		if ev.NeedRemove {
			t.Errorf("test1: not need remove")
		}
	}
}
