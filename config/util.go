package config

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log/slog"
	"os"
	"path"
	"path/filepath"
)

func init() {
	//fmt.Printf("=== main init\n")
	//InitLog()
	InitDotEnv()
}

func InitDotEnv() {
	file := path.Join(os.Getenv("HOME"), ".env")
	if godotenv.Load(file) != nil {
		slog.Info("File not found", slog.String("path", file))
	}
}

func MustGetFromEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		slog.Error("Config error", slog.String(name, "Env.var. not found"))
		panic("Config error")
	}
	return v
}

func GetFromEnv(name string) string {
	return os.Getenv(name)
}
func FindConfigFile() (string, error) {
	return FindFileUpToHome(FileConfig)
}

func CurrentDir() string {
	dir, _ := os.Getwd()
	return dir
}

func HomeDir() string {
	return os.Getenv("HOME")
}

func FindFileUpToHome(fileName string) (string, error) {
	// Get the current working directory
	currentDir := CurrentDir()
	// Start searching from the current directory and go up recursively
	for {
		filePath := filepath.Join(currentDir, fileName)
		_, err := os.Stat(filePath)
		// If the file exists, return its path
		if err == nil {
			return filePath, nil
		}
		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		// Check if we have reached the home directory
		if parentDir == HomeDir() {
			break
		}
		currentDir = parentDir
	}
	// If the file is not found, return an error
	return "", fmt.Errorf("file '%s' not found", fileName)
}

func LogError(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.Any("error", err))
	}
}

//func LogErr(err error) {
//	slog.Error("ERROR", slog.String("error", err.Error()))
//}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
