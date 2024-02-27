/*
Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
Software is distributed under MIT license as stated in LICENSE file

extract.go

Extract files and copy them to destination folder
*/
package extract

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"sandboxer/pkg/logging"
)

func ExtractFile(fs fs.FS, folder string, filePath string) (string, error) {
	if filepath.Ext(filePath) == ".gz" {
		return FileGZ(fs, folder, filePath)
	}
	return CopyFile(fs, folder, filePath)
}

func FileGZ(fs fs.FS, folder string, fileName_gz string) (string, error) {
	//logging.Debugf("Extract embedded %s to %s", fileName_gz, folder)
	file, err := fs.Open(fileName_gz)
	if err != nil {
		return "", fmt.Errorf("Open(\"%s\"): %w", fileName_gz, err)

	}
	defer func() {
		//logging.Debugf("Close embed.FS file")
		file.Close()
	}()
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("Open(\"%s\"): %w", fileName_gz, err)
	}
	defer func() {
		//logging.Debugf("Close GZip reader")
		gzipReader.Close()
	}()
	fileName := filepath.Base(fileName_gz)
	targetFileName := fileName[:len(fileName)-3]
	targetPath := filepath.Join(folder, targetFileName)
	//logging.Debugf("Target path %s", targetPath)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("create File: %w", err)
	}
	defer func() {
		//logging.Debugf("Close %s", targetPath)
		targetFile.Close()
	}()
	//logging.Debugf("Created %s", targetPath)
	if _, err := io.Copy(targetFile, gzipReader); err != nil {
		return "", fmt.Errorf("error creating %s: %w", targetPath, err)
	}
	//logging.Debugf("Extracted %s, %d bytes", targetPath, size)
	return targetPath, nil
}

func CopyFile(fs fs.FS, folder string, filePath string) (string, error) {
	//logging.Debugf("Extract embedded %s to %s", fileName_gz, folder)
	sourceFile, err := fs.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Open(\"%s\"): %w", filePath, err)
	}
	defer func() {
		//logging.Debugf("Close embed.FS file")
		sourceFile.Close()
	}()
	fileName := filepath.Base(filePath)
	targetPath := filepath.Join(folder, fileName)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer func() {
		targetFile.Close()
	}()
	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		return "", fmt.Errorf("copy error %s: %w", targetPath, err)
	}
	//logging.Debugf("Extracted %s, %d bytes", targetPath, size)
	return targetPath, nil
}

func CopyFolder(f fs.FS, root, sourceFolder, targetFolder string) error {
	return fs.WalkDir(f, filepath.Join(root, sourceFolder),
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			targetPath := filepath.Join(targetFolder, path[len(root)+1:]) // chop "embed/" prefix
			if d.IsDir() {
				logging.Debugf("MkdirAll(%s)", targetPath)
				if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
					return err
				}
				return nil
			}
			logging.Debugf("MkdirAll(%s)", filepath.Dir(targetPath))
			if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
				return err
			}

			logging.Debugf("Open(%s)", path)
			sourceFile, err := f.Open(path)
			if err != nil {
				return err
			}
			defer sourceFile.Close()

			logging.Debugf("Open(%s)", targetPath)
			targetFile, err := os.Open(targetPath)
			if err != nil {
				return err
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, sourceFile); err != nil {
				return err
			}
			return nil
		})
}

func ExtractFileTGZ(targetFolder string, fileSystem fs.FS, filePath string) error {
	logging.Debugf("ExtractFileTGZ(%s,%v,%s)", targetFolder, fileSystem, filePath)
	f, err := fileSystem.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return UntarReader(targetFolder, f)
}

/*
func ExtractGz(fileSystem fs.FS, filePath string, callback func(reader io.Reader) error) error {
	return ExtractFileTGZ(fileSystem, filePath, func(reader io.Reader) error {
		gzipFile, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		return callback(gzipFile)
	})
}

func ExtractTar(folder string, reader io.Reader) error {
	return nil
}
*/

func Untar(fs fs.FS, folder string, filePath string) error {
	logging.Debugf("Untar %s to %s", folder, filePath)
	sourceFile, err := fs.Open(filePath)
	if err != nil {
		return fmt.Errorf("Open(\"%s\"): %w", filePath, err)
	}
	defer func() {
		//logging.Debugf("Close embed.FS file")
		sourceFile.Close()
	}()
	return UntarReader(folder, sourceFile)
}

// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func UntarReader(targetFolder string, r io.Reader) error {
	logging.Debugf("Untar(%s)", targetFolder)

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		case strings.HasPrefix(header.Name, "._"):
			continue
		}
		target := filepath.Join(targetFolder, header.Name)
		//logging.Debugf("Untar(%s): taget: %s", targetFolder, target)
		switch header.Typeflag {
		case tar.TypeDir:
			//if _, err := os.Stat(target); err != nil {
			logging.Debugf("MkdirAll(%s)", target)
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
			//}
		case tar.TypeReg:
			logging.Debugf("MkdirAll(%s)", filepath.Dir(target))
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			logging.Debugf("OpenFile(%s)", target)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
	}
}

/*
func FileSize(seeker io.Seeker) (int64, error) {
	size, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return size, err
}

// UnzipSource - unpack zip file to folder.
// Taken from https://gosamples.dev/unzip-file/
func UnzipSource(fs fs.FS, source, destination string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	size, err := FileSize(file)
	reader, err := zip.NewReader(file)
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}
	return nil
}

func unzipFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}
*/
