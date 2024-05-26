package free

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

func CurrentPath() (string, error) {
	exePath, err := executableDir()
	if err != nil {
		return "", err
	}
	temPath, err := getTempDir()
	if err != nil {
		return "", err
	}
	if strings.Contains(exePath, temPath) {
		return currentDir()
	}
	return exePath, nil
}

func getTempDir() (string, error) {
	temp := os.Getenv("TEMP")
	if temp == "" {
		temp = os.Getenv("TMP")
	}
	return filepath.EvalSymlinks(temp)
}

func executableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(filepath.Dir(exePath))
}

func currentDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return path.Dir(filename), nil
	}
	return "", errors.New("runtime caller path error")
}
