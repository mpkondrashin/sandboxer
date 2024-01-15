package extract

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"examen/pkg/logging"
)

func FileGZ(fs fs.FS, folder string, fileName_gz string) (string, error) {
	logging.Debugf("Extract embedded %s to %s", fileName_gz, folder)
	file, err := fs.Open(fileName_gz)
	if err != nil {
		return "", fmt.Errorf("Open(\"%s\"): %w", fileName_gz, err)

	}
	defer func() {
		logging.Debugf("Close embed.FS file")
		file.Close()
	}()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("Open(\"%s\"): %w", fileName_gz, err)
	}
	defer func() {
		logging.Debugf("Close GZip reader")
		gzipReader.Close()
	}()
	targetFileName := filepath.Base(fileName_gz)[:len(fileName_gz)-3]
	targetPath := filepath.Join(folder, targetFileName)
	logging.Debugf("Target path %s", targetPath)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("error creating %s: %w", targetPath, err)
	}
	defer func() {
		logging.Debugf("Close %s", targetPath)
		targetFile.Close()
	}()
	logging.Debugf("Created %s", targetPath)
	size, err := io.Copy(targetFile, gzipReader)
	if err != nil {
		return "", fmt.Errorf("error creating %s: %w", targetPath, err)
	}
	logging.Debugf("Extracted %s, %d bytes", targetPath, size)
	return targetPath, nil
}
