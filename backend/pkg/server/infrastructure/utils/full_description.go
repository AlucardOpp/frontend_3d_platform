package utils

import (
	"path/filepath"
	"strings"

	"github.com/AhEhIOhYou/etomne/pkg/server/domain/entities"
)

var fullDescriptionExtensions = []string{".txt", ".doc", ".docx", ".pdf"}

func IsFullDescriptionExtension(extension string) bool {
	ext := strings.ToLower(filepath.Ext(extension))
	for _, allowed := range fullDescriptionExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

func FilterOutFullDescriptionFile(files []entities.File, fullDescriptionFileID *uint64) []entities.File {
	if fullDescriptionFileID == nil || *fullDescriptionFileID == 0 {
		return files
	}

	filtered := make([]entities.File, 0, len(files))
	for _, file := range files {
		if file.ID != *fullDescriptionFileID {
			filtered = append(filtered, file)
		}
	}
	return filtered
}
