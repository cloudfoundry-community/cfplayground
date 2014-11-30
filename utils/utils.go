package utils

import "os"

func IsDirExists(dir string) bool {
	src, err := os.Stat(dir)

	if err != nil {
		return false
	}

	return src.IsDir()
}
