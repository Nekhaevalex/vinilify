package utils

import (
	"io"
	"log"
	"net/http"
	"os"
)

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
