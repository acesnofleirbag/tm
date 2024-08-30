package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"tm/guard"
)

func capitalize(str string) string {
	if len(str) == 0 {
		return str
	}

	first := strings.ToUpper(string(str[0]))
	rest := strings.ToLower(str[1:])

	return first + rest
}

type CopyRunner struct {
	Replace bool
	Name    string
}

func (self *CopyRunner) CpFile(srcPath string, dstPath string) {
	src, err := os.Open(srcPath)
	defer src.Close()
	guard.Err(err)

	dstPath = strings.ReplaceAll(string(dstPath), "<:A:>", self.Name)
	dstPath = strings.ReplaceAll(string(dstPath), "<:AFU:>", capitalize(self.Name))
	dst, err := os.Create(dstPath)
	guard.Err(err)

	if self.Replace {
		data, err := io.ReadAll(src)
		guard.Err(err)

		content := strings.ReplaceAll(string(data), "<:A:>", self.Name)
		content = strings.ReplaceAll(string(content), "<:AFU:>", capitalize(self.Name))
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
