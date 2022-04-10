package scripts

import (
	"embed"
	"io/ioutil"
	"os"
)

//go:embed "files"
var ScriptFiles embed.FS

// tmpFileCreator creates a temp file containing the script data passed in. The
// file is not removed in this function - it must be done in the caller.
func tmpFileCreator(script []byte) (string, error) {
	tmp, err := ioutil.TempFile("/tmp", "ds-file")
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(tmp.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	_, err = f.Write(script)
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

// Retriever reads the file provided and returns its data. It expects a valid
// shell script. A temp file is created and that filename is then returned to be
// called by Z.Exec
func Retriever(filename string) (string, error) {
	data, err := ScriptFiles.ReadFile(filename)
	if err != nil {
		return "", err
	}

	script, err := tmpFileCreator(data)
	if err != nil {
		return "", err
	}
	return script, nil
}
