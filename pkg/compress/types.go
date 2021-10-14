package compress

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Compress is a common compress interface
type Compress interface {
	ExtractFiles(sourceFile, targetName string) error
}

func extraFile(name, targetName, tarFile string, header *tar.Header, tarReader *tar.Reader) (err error) {
	if name != targetName && !strings.HasSuffix(name, "/"+targetName) {
		return
	}
	var targetFile *os.File
	if targetFile, err = os.OpenFile(fmt.Sprintf("%s/%s", filepath.Dir(tarFile), targetName),
		os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)); err != nil {
		return
	}
	if _, err = io.Copy(targetFile, tarReader); err != nil {
		return
	}
	_ = targetFile.Close()
	return
}

// GetCompressor gets the compressor base on file extension
func GetCompressor(extension string, additionBinaries []string) Compress {
	// Select the right decompressor based on file type
	switch extension {
	case ".xz":
		return NewXz(additionBinaries)
	case ".zip":
		return NewZip(additionBinaries)
	default:
		return NewGZip(additionBinaries)
	}
}