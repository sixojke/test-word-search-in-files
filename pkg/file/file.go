package file

import (
	"fmt"
	"io"
	"os"
	"sort"
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

// Ищет слово в отсортированном файле, если этого файла нет, создаёт его. В случае успеха добавляет файл по ссылке
func addTargetFile(dir, fileName string, targetWord string, targetFiles *[]string) error {
	filePath := dir + "/" + fileName
	filePathSort := dir + "/" + fileName + "sort"
	if IsExists(filePathSort) {
		if err := findWordInSortedFile(filePathSort, fileName, targetWord, targetFiles); err != nil {
			return err
		}
	} else {
		if err := CreateSort(filePath, filePathSort); err != nil {
			return err
		}

		if err := findWordInSortedFile(filePathSort, fileName, targetWord, targetFiles); err != nil {
			return err
		}
	}

	return nil
}

// Поиск слова в отсортированной файле при помощи бинарного поиска.
func findWordInSortedFile(filePathSort string, fileName string, targetWord string, targetFiles *[]string) error {
	text, err := ReadFile(filePathSort)
	if err != nil {
		return err
	}

	if binarySearch(strings.Fields(text), targetWord) {
		*targetFiles = append(*targetFiles, fileName)
	}

	return nil
}

// Бинарный поиск таргет слова
func binarySearch(sortedWords []string, target string) bool {
	low := 0
	high := len(sortedWords) - 1

	for low <= high {
		mid := (low + high) / 2
		if sortedWords[mid] == target {
			return true
		} else if sortedWords[mid] < target {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return false
}

// Создаёт отсортированный файл с чистым текстом
func CreateSort(fileName, fileNameSort string) error {
	text, err := ReadFile(fileName)
	if err != nil {
		return err
	}

	text = cleanText(text)

	words := strings.Fields(text)
	sort.Strings(words)

	if err := WriteFile(fileNameSort, strings.Join(words, " ")); err != nil {
		return err
	}

	return nil
}

// Читает строку с файла
func ReadFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("error open file: %v", err)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error read file: %v", err)
	}

	return string(b), nil
}

// Записывает строку в файл
func WriteFile(fileName, text string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error create file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return fmt.Errorf("error write content: %v", err)
	}

	return nil
}

// Очищает текст от лишних точек, запятых ...
func cleanText(text string) string {
	return strings.NewReplacer(",", "", ".", "", ":", "", ";", "", "(", "", ")", "", "-", "").Replace(text)
}

// Возращает все файлы из указанной дирректории
func FromDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error read dir: %v", err)
	}

	var fileNames []string
	for _, file := range files {
		if !strings.Contains(file.Name(), "sort") {
			fileNames = append(fileNames, file.Name())
		}
	}

	if fileNames == nil {
		return nil, fmt.Errorf("files not found")
	}

	return fileNames, nil
}

// Проверяет существование файла
func IsExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		return false
	}

	return true
}
