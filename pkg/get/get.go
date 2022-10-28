package get

import (
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"net/http"
	"strings"
	"time"
)

var (
	httpTimeout = 10 * time.Second
)

type Tool struct {
	// Name of the tool
	Name string

	// Repo is the github repo
	Repo string

	// Owner is the tool Repo owner, such as
	// derailed/k9s
	Owner string

	// Version to pull. An empty string means "latest"
	Version string

	// Description of what this tool does/is.
	Description string

	// NonBinary is used to determine if the tool is not a binary such as
	// kubetail which is a bash script
	NonBinary bool
}

func (t Tool) GetURL(ver string, quiet bool) (string, error) {
	if len(ver) == 0 {
		v, err := FindGithubRelease(t.Owner, t.Repo)
		if err != nil {
			return "", err
		}
		ver = v
	}
	return binaryURL(t.Owner, t.Repo, ver), nil
}

func binaryURL(owner string, repo string, ver string) string {
	return fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/%s", owner, repo, ver,
	)
	// determine if the url is:
	// github.com/OWNER/REPO/releases/download/VERSION/NAME
	// or
	// github.com/OWNER/REPO/releases/download/NAME (no version tag)
}

func FindGithubRelease(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://github.com/%s/%s/releases/latest", owner, repo)

	cl := httpClient(&httpTimeout)
	cl.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	r, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}

	res, err := cl.Do(r)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 302 {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	location := res.Header.Get("Location")
	if len(location) == 0 {
		return "", fmt.Errorf("unable to find location of release")
	}

	version := location[strings.LastIndex(location, "/")+1:]
	return version, nil
}

func GetDownloadURL(tool *Tool, version string, quiet bool) (string, error) {
	return "", nil
}

// httpClient w/ timeouts
func httpClient(timeout *time.Duration) http.Client {
	cl := http.Client{}
	cl.Timeout = *timeout
	return cl
}

// get tool url (github)
func DownloadURL(t *Tool, version string, quiet bool) (string, error) {
	v := getToolVersion(t, version)

	dl, err := t.GetURL(v, quiet)
	if err != nil {
		return "", err
	}
	return dl, nil
}

func getToolVersion(t *Tool, version string) string {
	v := t.Version
	if len(v) > 0 {
		v = version
	}
	return v
}

// download binary
// check tool exists in list of tools
// post install message
// display table of supported binaries

var Cmd = &Z.Cmd{
	Name:    `get`,
	Summary: `get executables and applications onto the host system [requires internet]`,
	Description: `
		The *get* command downloads a tools or applications from that providers releases or
		downloads page. Typically, tools are downloaded as a binary for fast and efficient access
		on the host platform.
		`,
	Commands: []*Z.Cmd{
		// imported commands
		help.Cmd,
		// local
	},
	Call: func(_ *Z.Cmd, args ...string) error {
		t := Tool{}
		dl, err := GetDownloadURL(&t, "v1.0", false)
		if err != nil {
			return err
		}
		fmt.Printf("get: '%s'\n", dl)
		return nil
	},
}
