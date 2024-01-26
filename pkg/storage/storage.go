package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pluveto/flydav/internal/config"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	WriteAll(path string, data []byte) error
	ReadAll(path string) ([]byte, error)
	Write(path string, data []byte, offset int64) error
	Read(path string, offset int64, length int64) ([]byte, error)
	Delete(path string) error
	CreateDirectory(path string) error
	Stat(path string) (Metadata, error)
	List(path string) ([]Metadata, error)
	Size(path string) (int64, error)
	Move(src string, dst string) error
	Copy(src string, dst string) (int64, error)
	Merge(srcs []string, dst string) error
}

func NewStorage(cfg config.StorageBackendConfig) Storage {
	switch cfg.GetEnabledBackend() {
	case "local":
		return NewLocalStorage(cfg.Local)
	default:
		panic(fmt.Errorf("unknown storage backend type"))
	}
}

type Metadata struct {
	Name     string
	FullName string
	IsDir    bool
	Size     int64
}

type LocalStorage struct {
	config config.LocalConfig
}

func NewLocalStorage(cfg config.LocalConfig) *LocalStorage {
	// cfg.BaseDir must exists and be a directory writable.
	if err := os.MkdirAll(cfg.BaseDir, 0755); err != nil {
		panic(err)
	}

	absDir, err := filepath.Abs(cfg.BaseDir)
	if err != nil {
		panic(err)
	}

	cfg.BaseDir = absDir

	return &LocalStorage{
		config: cfg,
	}
}

func (ls *LocalStorage) absPath(path string) string {
	return filepath.Join(ls.config.BaseDir, path)
}

func (ls *LocalStorage) relPath(realPath string) string {
	relPath, err := filepath.Rel(ls.config.BaseDir, realPath)
	if err != nil {
		return ""
	}
	return relPath
}

func (ls *LocalStorage) WriteAll(path string, data []byte) error {
	path = ls.absPath(path)
	return os.WriteFile(path, data, 0644)
}

func (ls *LocalStorage) ReadAll(path string) ([]byte, error) {
	path = ls.absPath(path)
	return os.ReadFile(path)
}

func (ls *LocalStorage) Delete(path string) error {
	path = ls.absPath(path)
	return os.Remove(path)
}

func (ls *LocalStorage) Write(path string, data []byte, offset int64) error {
	path = ls.absPath(path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteAt(data, offset)
	return err
}

func (ls *LocalStorage) Read(path string, offset int64, length int64) ([]byte, error) {
	path = ls.absPath(path)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, length)
	_, err = f.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (ls *LocalStorage) Stat(path string) (Metadata, error) {
	path = ls.absPath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return Metadata{}, err
	}

	return Metadata{
		Name:     fileInfo.Name(),
		FullName: ls.relPath(path),
		IsDir:    fileInfo.IsDir(),
		Size:     fileInfo.Size(),
	}, nil
}

func (ls *LocalStorage) List(path string) ([]Metadata, error) {
	path = ls.absPath(path)
	dirents, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var metadata []Metadata
	for _, dirent := range dirents {
		info, err := dirent.Info()
		if err != nil {
			return nil, err
		}

		size := info.Size()
		metadata = append(metadata, Metadata{
			Name:     dirent.Name(),
			FullName: ls.relPath(filepath.Join(path, dirent.Name())),
			IsDir:    dirent.IsDir(),
			Size:     size,
		})
	}

	return metadata, nil
}

func (ls *LocalStorage) Size(path string) (int64, error) {
	path = ls.absPath(path)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

func (ls *LocalStorage) Move(src string, dst string) error {
	src = ls.absPath(src)
	dst = ls.absPath(dst)
	return os.Rename(src, dst)
}

func (ls *LocalStorage) Merge(srcs []string, dst string) error {
	dst = ls.absPath(dst)
	f, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, src := range srcs {
		src = ls.absPath(src)
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}

		_, err = f.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ls *LocalStorage) CreateDirectory(path string) error {
	path = ls.absPath(path)
	return os.MkdirAll(path, 0755)
}

func (ls *LocalStorage) Copy(src string, dst string) (int64, error) {
	src = ls.absPath(src)
	dst = ls.absPath(dst)

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
