package fileutil

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tukaelu/sabadashi/internal/exporter"
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
		now := time.Now().Format("20060102T150405")
		renameTo := filepath.Join(
			cwd,
			name,
			fmt.Sprintf("%s_%s_%s", from, to, now),
		)
		if err := os.Rename(dirPath, renameTo); err != nil {
			return "", fmt.Errorf("Failed to rename an already existing directory. ('%s' to '%s')", dirPath, renameTo)
		}
		fmt.Printf("Renamed an already existing directory to '%s'.\n", renameTo)

		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", err
		}
	}
	return dirPath, nil
}

func RemoveDir(exportDir string) {
	if err := os.Remove(exportDir); err != nil {
		fmt.Printf("Failed to delete directory. (path: %s)", exportDir)
	}
}

func GetExportFilePath(dir, name, ext string) string {
	return filepath.Join(
		dir,
		fmt.Sprintf("%s.%s", name, ext),
	)
}

func WriteFile(dir, fileName, ext string, metrics exporter.CsvRecords, withFriendly bool) error {
	file := GetExportFilePath(dir, fileName, ext)
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// [TODO] The standard Go encoding/csv library may not be able to handle CSV escape, therefore, it may be a good idea to introduce a library for this purpose.
	cw := csv.NewWriter(f)
	defer cw.Flush()

	for _, metric := range metrics {
		if err := cw.Write(metric.ToStringArray(withFriendly)); err != nil {
			fmt.Printf("write csv file failed (reason: %s)", err)
		}
	}
	return nil
}
