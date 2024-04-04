package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func FindFilesWithWord(dir string, fileNames []string, word string) ([]string, error) {
	var resultError error
	targetFiles := make([]string, 0, len(fileNames))
	errCh := make(chan error, len(fileNames))

	for _, fileName := range fileNames {
		go func(fileName string) {
			errCh <- addTargetFile(dir, fileName, word, &targetFiles)
		}(fileName)
	}

	for i := 0; i < len(fileNames); i++ {
		if err := <-errCh; err != nil {
			resultError = err
		}
	}

	if len(targetFiles) == 0 {
		return nil, resultError
	}

	return targetFiles, resultError
}

func addTargetFile(dir, fileName string, targetWord string, targetFiles *[]string) error {
	file, err := os.Open(dir + "/" + fileName)
	if err != nil {
		return errors.New("error opening file: " + fileName)
	}

	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return errors.New("error reading file: " + fileName)
	}

	words := strings.Split(cleanText(string(b)), " ")
	for _, word := range words {
		if strings.EqualFold(word, targetWord) {
			*targetFiles = append(*targetFiles, fileName)
			return nil
		}
	}

	return nil
}

func cleanText(text string) string {
	return strings.NewReplacer(",", "", ".", "", ":", "").Replace(text)
}

func FromDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error read dir: %v", err)
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	if fileNames == nil {
		return nil, fmt.Errorf("files not found")
	}

	return fileNames, nil
}
