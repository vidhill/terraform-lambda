package util

import (
	"os"
	"path"
)

func RemoveFiles(paths ...string) error {
	for _, p := range paths {
		if err := os.Remove(p); err != nil {
			return err
		}
	}
	return nil
}

func IsJpegExtension(p string) bool {
	switch path.Ext(p) {
	case ".jpeg", ".jpg":
		return true
	default:
		return false
	}
}
