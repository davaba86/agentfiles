package files

import (
	"errors"
	"os"
	"path/filepath"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func WriteNew(path string, data []byte, perm os.FileMode) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

func Write(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func Copy(src, dst string, perm os.FileMode) error {
	data, err := Read(src)
	if err != nil {
		return err
	}
	return WriteNew(dst, data, perm)
}

func Join(dir, name string) string {
	return filepath.Join(dir, name)
}
