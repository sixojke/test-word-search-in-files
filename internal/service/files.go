package service

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sixojke/internal/config"
	"github.com/sixojke/internal/domain"
	"github.com/sixojke/internal/repository"
	"github.com/sixojke/pkg/file"
	"github.com/sixojke/pkg/path"
)

var filesDir = ""

func init() {
	dir, err := path.Work("/files")
	if err != nil {
		log.Fatalf("error work dir: files: %v", err)
	}

	filesDir = dir
}

type FilesService struct {
	cache       repository.Cache
	cacheConfig config.CacheConfig
}

func NewFilesService(cache repository.Cache, cacheConfig config.CacheConfig) *FilesService {
	return &FilesService{
		cache:       cache,
		cacheConfig: cacheConfig,
	}
}

func (s *FilesService) FindFilesWithWord(inp domain.FilesSearchInp) (*domain.FilesSearchOut, error) {
	files, err := s.filesWithWord(inp.Word)
	if err != nil {
		return nil, err
	}

	return &domain.FilesSearchOut{Files: files}, nil
}

func (s *FilesService) filesWithWord(word string) ([]string, error) {
	filesStr, err := s.cache.Get(word)
	if err != nil {
		log.Warn(err)
	}

	if filesStr != "" {
		targetFiles := strings.Split(filesStr, "-")
		if len(targetFiles) == 0 {
			return nil, nil
		}

		return targetFiles, nil
	} else {
		files, err := file.FromDir(filesDir)
		if err != nil {
			return nil, err
		}

		targetFiles, err := file.FindFilesWithWord(filesDir, files, word)
		if err != nil {
			return nil, err
		}

		if err := s.cache.Set(word, strings.Join(targetFiles, "-"), s.cacheConfig.Expiration); err != nil {
			log.Warn(err)
		}

		return targetFiles, nil
	}
}
