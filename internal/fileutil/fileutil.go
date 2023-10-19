package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateExportDir(name, from, to string) (string, error) {
	cwd, _ := os.Getwd()
	dirPath := filepath.Join(
		cwd,
		name,
		fmt.Sprintf("%s_%s", from, to),
	)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("Directory '%s' is already exists. Please try again after making a backup", dirPath)
	}
	return dirPath, nil
}

func GetExportFilePath(dir, name, ext string) string {
	return filepath.Join(
		dir,
		fmt.Sprintf("%s.%s", name, ext),
	)
}

func WriteFile(dir, fileName, ext string, fn func(f *os.File) error) error {
	file := GetExportFilePath(dir, fileName, ext)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return fn(f)
}
