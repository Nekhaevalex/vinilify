package utils

import (
	"path/filepath"
	"testing"
)

func TestDirExists(t *testing.T) {
	usersPath := filepath.Join(GetRoot(), "utils")
	exists, _ := DirExists(usersPath)
	t.Log(exists)

	userPath1 := filepath.Join(GetRoot(), "utils", "123")
	exists, _ = DirExists(userPath1)
	t.Log(exists)
}

func TestUserFilepath(t *testing.T) {
	t.Log(GetUserPath(123))
}
