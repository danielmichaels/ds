package get

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// decompressArchive will decompress any valid archive file.
func decompressArchive(tool *Tool, dlURL, outFilePath, opSystem, arch, version string) (string, error) {
	file, err := os.Open(outFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	target := filepath.Dir(outFilePath)
	outFilePath = path.Join(target, tool.Name)

	switch {
	case strings.HasSuffix(dlURL, "tar.gz"):
		err := Untar(file, target, true)
		if err != nil {
			return "", err
		}
	case strings.HasSuffix(dlURL, "tgz"):
		err := Untar(file, target, true)
		if err != nil {
			return "", err
		}
	case strings.HasSuffix(dlURL, "zip"):
		info, err := file.Stat()
		if err != nil {
			return "", err
		}

		err = Unzip(file, info.Size(), target)
		if err != nil {
			return "", err
		}
	default:
		log.Printf("tool %s appears th have an unsupported binary format\n", tool.Name)
	}
	return outFilePath, nil
}

// Untar untar a file to a target directory on the host it is running on.
func Untar(r io.Reader, target string, gzip bool) error {
	return untar(r, target, gzip)
}

// untar is the private interface for Untar. untar accepts any valid archive
// and will open it for reading, extract the files and place them into the hosts
// /tmp directory.
func untar(r io.Reader, target string, gz bool) error {
	tarReader := tar.NewReader(r)
	if gz {
		gzr, err := gzip.NewReader(r)
		if err != nil {
			return fmt.Errorf("failed to decompress gzip: %v", err)
		}
		tarReader = tar.NewReader(gzr)
	}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error reading tarfile: %w", err)
			return fmt.Errorf("tar error: %w", err)
		}

		fpath := filepath.Base(header.Name)
		fpath = path.Join(target, fpath)
		info := header.FileInfo()
		mode := info.Mode()
		switch {
		case mode.IsDir():
			break
		case mode.IsRegular():
			f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode.Perm())
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(f, tarReader)
			if err != nil {
				return err
			}
		default:
		}
	}
	return nil
}

// Unzip unzips files and puts them on the host system.
func Unzip(r io.ReaderAt, size int64, target string) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return fmt.Errorf("error creating zip reader: %w", err)
	}
	return unzip(zr, target)
}

// unzip calls extractAndWrite to unzip and write files to the hosts target directory.
func unzip(r *zip.Reader, target string) error {
	for _, f := range r.File {
		err := extractAndWrite(f, target)
		if err != nil {
			return err
		}
	}
	return nil
}

// extractAndWrite extracts the zip file and writes its contents to the target
// directory. Nested directories are flattened into the target.
func extractAndWrite(zf *zip.File, target string) error {
	r, err := zf.Open()
	if err != nil {
		return err
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	fpath := filepath.Base(zf.Name)
	fpath = path.Join(target, fpath)
	info := zf.FileInfo()
	mode := info.Mode()

	switch {
	case mode.IsDir():
		break
	case mode.IsRegular():
		f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zf.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(f, r)
		if err != nil {
			return err
		}
	default:
	}

	return nil
}
