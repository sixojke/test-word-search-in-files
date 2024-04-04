package file_test

import (
	"testing"

	"github.com/sixojke/pkg/file"
)

func TestAddTargetFile(t *testing.T) {
	dir := "fake"
	fileNames := []string{"file4", "file4444"}
	word := "я"

	_, err := file.FindFilesWithWord(dir, fileNames, word)
	if err == nil {
		t.Error("TestAddTargetFile error process file")
	}

	dir = "/files"
	fileNames = []string{"file1", "file2"}
	word = "wdsf"

	result, _ := file.FindFilesWithWord(dir, fileNames, word)
	if result != nil {
		t.Error("TestAddTargetFile not nil")
	}
}
