package storage

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/KurosawaAngel/gobackup/helper"
	"github.com/KurosawaAngel/gobackup/logger"
)

// Local storage
//
// type: local
// path: /data/backups
type Local struct {
	Base
	path string
}

func (s *Local) open() error {
	s.path = s.viper.GetString("path")
	return helper.MkdirP(s.path)
}

func (s *Local) close() {}

func (s *Local) upload(fileKey string) (err error) {
	logger := logger.Tag("Local")

	// Related path
	if !path.IsAbs(s.path) {
		s.path = path.Join(s.model.WorkDir, s.path)
	}

	targetPath := path.Join(s.path, fileKey)
	targetDir := path.Dir(targetPath)
	if err := helper.MkdirP(targetDir); err != nil {
		logger.Errorf("failed to mkdir %q, %v", targetDir, err)
	}

	_, err = helper.Exec("cp", "-a", s.archivePath, targetPath)
	if err != nil {
		return err
	}
	logger.Info("Store succeeded", targetPath)
	return nil
}

func (s *Local) delete(fileKey string) (err error) {
	targetPath := filepath.Join(s.path, fileKey)
	logger.Info("Deleting", targetPath)

	return os.Remove(targetPath)
}

// List all files
func (s *Local) list(parent string) ([]FileItem, error) {
	remotePath := filepath.Join(s.path, parent)
	items := []FileItem{}

	files, err := os.ReadDir(remotePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				return nil, err
			}
			items = append(items, FileItem{
				Filename:     info.Name(),
				Size:         info.Size(),
				LastModified: info.ModTime(),
			})
		}
	}

	return items, nil
}

func (s *Local) download(fileKey string) (string, error) {
	return "", fmt.Errorf("Local is not support download")
}
