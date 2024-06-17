package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func GetRoot() string {
	abs, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("failed to retrieve absolute path")
	}
	return strings.SplitAfter(abs, "vinilify")[0]
}

// Returns absolute path to Assets directory
func GetAssets() string {
	root := GetRoot()
	return filepath.Join(root, "assets")
}

// Checks if directory exists
func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}

// Downloads and saves attachement from url to file with filepath
func DownloadAttachment(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		log.Panic(err)
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panic(err)
	}

	return err
}
