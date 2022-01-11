package utils

import (
	"os"
)

// getFileSize returns the size of the file. This is used to compare files
func GetFileSize(filePath string) (int64, error) {
	in, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	stats, err := in.Stat()
	if err != nil {
		return 0, err
	}
	return stats.Size(), nil
}
