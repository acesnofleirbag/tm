package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"tm/guard"
)

type CopyRunner struct {
	Replace bool
	Name    string
}

func (self *CopyRunner) CpFile(srcPath string, dstPath string) {
	src, err := os.Open(srcPath)
	defer src.Close()
	guard.Err(err)

	dst, err := os.Create(dstPath)
	guard.Err(err)

	if self.Replace {
		data, err := io.ReadAll(src)
		guard.Err(err)

		content := strings.ReplaceAll(string(data), "<:AGGREGATE:>", self.Name)
		reader := strings.NewReader(content)

		_, err = io.Copy(dst, reader)
		guard.Err(err)
	} else {
		_, err = io.Copy(dst, src)
		guard.Err(err)
	}
}

func (self *CopyRunner) CpDir(src string, dst string) {
	entries, err := os.ReadDir(src)
	guard.Err(err)

	err = os.MkdirAll(dst, 0755)
	guard.Err(err)

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			self.CpDir(srcPath, dstPath)
		} else {
			self.CpFile(srcPath, dstPath)
		}
	}
}
