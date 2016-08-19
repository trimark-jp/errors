package errors

import (
	"path/filepath"
	"strings"
)

func trimGOPATHProbably(path string, function string) string {
	maybePackage := ""
	lastSlashIndex := strings.LastIndex(function, "/")
	if 0 <= lastSlashIndex {
		maybePackage = function[:lastSlashIndex]
		return trimGOPATHByPackage(path, maybePackage)
	}
	return parentFolderAndFileName(path)
}

func trimGOPATHByPackage(path string, pkg string) string {
	index := strings.Index(path, pkg)
	if index <= 0 {
		return parentFolderAndFileName(path)
	}
	return path[index:]
}

func parentFolderAndFileName(path string) string {
	return filepath.Join(filepath.Base(filepath.Dir(path)),
		filepath.Base(path))
}
