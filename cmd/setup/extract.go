package main

import (
	"compress/gzip"
	"embed"
	"io"
	"os"
	"path/filepath"

	"examen/pkg/logging"
)

//go:embed opengl32.dll.gz
var embedFS embed.FS

func extractEmbeddedGZ(folder, fileName_gz string) string {
	logging.Infof("Extract embedded %s to %s", fileName_gz, folder)
	file, err := embedFS.Open(fileName_gz)
	if err != nil {
		logging.Errorf("Open(\"%s\"): %v", fileName_gz, err)
		panic(err)
	}
	defer file.Close()
	logging.Debugf("Opened %s", fileName_gz)
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		logging.Errorf("Open(\"%s\"): %v", fileName_gz, err)
		panic(err)
	}
	defer gzipReader.Close()
	logging.Debugf("Crated gzip reader for %s", fileName_gz)
	targetFileName := fileName_gz[:len(fileName_gz)-3]
	logging.Debugf("Target file name %s", targetFileName)
	targetPath := filepath.Join(folder, targetFileName)
	logging.Debugf("Target path %s", targetPath)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		logging.Errorf("Error creating %s: %v", targetPath, err)
		panic(err)
	}
	logging.Debugf("Create %s", targetPath)
	size, err := io.Copy(targetFile, gzipReader)
	if err != nil {
		logging.Errorf("Error creating %s: %v", targetPath, err)
		panic(err)
	}
	logging.Debugf("Extracted %s, %d bytes", targetPath, size)
	return targetPath
}
