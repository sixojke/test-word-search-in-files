package domain

import "fmt"

type FilesSearchInp struct {
	Word string
}

type FilesSearchOut struct {
	Files []string `json:"files"`
}

func (f *FilesSearchInp) Validate() error {
	if f.Word == "" {
		return fmt.Errorf("error.StatusBadRequest")
	}

	return nil
}
