package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed web/*
var web embed.FS

type embeddedFileSystem struct {
	fs http.FileSystem
}

func (b *embeddedFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(setDefault(name))
}

func (b *embeddedFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(setDefault(p)); err != nil {
			return false
		}
		return true
	}
	return false
}

func setDefault(path string) string {
	if len(path) == 0 {
		return "/index.html"
	}
	return path
}

func EmbeddedFileSystem(targetPath string) *embeddedFileSystem {
	fsys, err := fs.Sub(web, targetPath)
	if err != nil {
		panic(err)
	}
	return &embeddedFileSystem{
		http.FS(fsys),
	}
}
