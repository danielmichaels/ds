package get

import (
	"bytes"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

var (
	//toolFilePath = ".local/bin"
	toolFilePath = ".ds/bin"
)

func mkdirp(path string) error {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return err
	}
	return err
}

// InitUserDir will establish a location for the binaries to be stored.
func InitUserDir() (string, error) {
	home := os.Getenv("HOME")
	binPath := fmt.Sprintf("%s/%s", home, toolFilePath)
	err := mkdirp(binPath)
	if err != nil {
		return "", err
	}
	if len(home) == 0 {
		return home, fmt.Errorf("$HOME, not set")
	}
	return binPath, nil
}

// LocalBinary returns the filepath for the binary in the users home directory.
func LocalBinary(name, subdir string) (string, error) {
	home := os.Getenv("HOME")

	val := path.Join(home, toolFilePath)
	if len(subdir) > 0 {
		val = path.Join(val, subdir)
	}

	return path.Join(val, name), nil
}

// Download is a public interface for downloading a file from a provided URL.
func Download(tool *Tool, arch, opSystem, version string) (string, error) {
	dlURL, err := GetDownloadURL(*tool, arch, opSystem, version)
	if err != nil {
		return "", err
	}
	log.Printf("Downloading %q", dlURL)

	outputPath, err := downloadFile(dlURL)
	if err != nil {
		return "", err
	}

	if isArchive, err := tool.IsArchive(dlURL); isArchive {
		if err != nil {
			return "", err
		}

		out, err := decompressArchive(tool, dlURL, outputPath, opSystem, arch, version)
		if err != nil {
			return "", err
		}
		outputPath = out
		log.Printf("Extracted %q\n", outputPath)
	}

	_, err = InitUserDir()
	if err != nil {
		return "", err
	}

	localPath, err := LocalBinary(tool.Name, "")
	if err != nil {
		return "", err
	}

	_, err = CopyFile(outputPath, localPath, 0700)
	log.Printf("Copied %q to %q\n", outputPath, localPath)
	if err != nil {
		return "", err
	}
	return outputPath, nil
}

// downloadFile retrieves a file from a given URL and downloads it to the local
// machine returning the path of that file.
// A file length is required to render the download progress bar.
func downloadFile(url string) (string, error) {
	cl := httpClient(&httpTimeout)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", nil
	}

	res, err := cl.Do(r)
	if err != nil {
		return "", nil
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("unexpected status code during download: %d", res.StatusCode)
	}

	_, file := path.Split(url)
	tmp := os.TempDir()
	outFilePath := path.Join(tmp, file)
	out, err := os.Create(outFilePath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	progBar := progressbar.DefaultBytes(
		res.ContentLength,
		"downloading",
	)
	_, err = io.Copy(io.MultiWriter(out, progBar), res.Body)
	if err != nil {
		return "", err
	}
	return outFilePath, nil
}

// CopyFile copies a source to a destination and applies permissions to that file.
func CopyFile(src, dst string, permissions int) (int64, error) {
	_, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	dest, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(permissions))
	if err != nil {
		return 0, err
	}
	defer dest.Close()
	bites, err := io.Copy(dest, source)
	if err != nil {
		return 0, err
	}
	return bites, err
}

type ToolLocal struct {
	Name    string
	Path    string
	BinPath string
}

// PostInstallationMessage generates installation message after tool has been downloaded
func PostInstallationMessage(localToolsStore ToolLocal) ([]byte, error) {

	t := template.New("Installation Instructions")

	t.Parse(`
# Add arkade binary directory to your PATH variable
export PATH=$PATH:$HOME/{{.BinPath}}

# Test the binary:
{{.Path}}

# Or install with:
sudo mv {{.Path}} /usr/local/bin/
`)

	var tpl bytes.Buffer

	err := t.Execute(&tpl, localToolsStore)
	if err != nil {
		return nil, err
	}

	return tpl.Bytes(), err
}
