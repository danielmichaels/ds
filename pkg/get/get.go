package get

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"text/template"
	"time"
)

var (
	httpTimeout = 10 * time.Second
)

func httpClient(timeout *time.Duration) http.Client {
	cl := http.Client{}
	cl.Timeout = *timeout
	return cl
}

var Cmd = &Z.Cmd{
	Name:    `get`,
	Summary: `install executables and applications on the host system [requires internet]`,
	Description: `
		The *get* command downloads a tools or applications from that providers releases or
		downloads page. Typically, tools are downloaded as a binary for fast and efficient access
		on the host platform.
		`,
	Other: []Z.Section{
		{
			Title: "Examples",
			Body: `
			ds get - list all available tools

			ds get arkade - download the Arkade binary`,
		},
	},
	Commands: []*Z.Cmd{
		// imported commands
		help.Cmd,
		// local
	},
	Call: func(_ *Z.Cmd, args ...string) error {
		tools := MakeTools()
		arch, opSystem := GetClientArch()
		sort.Sort(tools)
		if len(args) == 0 {
			ListToolsTable(tools)
			return nil
		}
		tool := args[0]
		log.Printf("Looking up version for %q\n", tool)
		t, err := getTool(tool, tools)
		if err != nil {
			return err
		}

		version := t.Version
		if version == "" {
			version = "latest"
		}
		_, err = Download(&t, arch, opSystem, version)
		if err != nil {
			return err
		}

		err = PrintPostInstallMessage(t)
		if err != nil {
			return err
		}

		return nil
	},
}

// FindGithubRelease retrieves a response from GitHub's API for any valid repository
// in JSON format.
func FindGithubRelease(owner, repo string) ([]*GithubAPIReleasesResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	cl := httpClient(&httpTimeout)

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := cl.Do(r)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	var release []*GithubAPIReleasesResponse
	//var r2 any
	err = json.NewDecoder(res.Body).Decode(&release)
	if err != nil {
		return nil, fmt.Errorf("failed to decode release with err: %s", err)
	}
	return release, nil
}

func PrintPostInstallMessage(t Tool) error {
	lt := ToolLocal{
		Name:    t.Name,
		Path:    fmt.Sprintf("%s/%s/%s", os.Getenv("HOME"), toolFilePath, t.Name),
		BinPath: fmt.Sprintf("%s", toolFilePath),
	}
	msg, err := PostInstallationMessage(lt)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", msg)
	return nil
}

// ListToolsTable returns a list of all supported tools in tabular format.
func ListToolsTable(tools Tools) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetColWidth(60)
	table.SetHeader([]string{"Tool", "Description"})
	count := 0
	for _, tool := range tools {
		table.Append([]string{tool.Name, tool.Description})
		count++
	}
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.Normal},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Normal, tablewriter.Normal},
	)
	table.SetRowLine(true)
	table.SetCaption(true, fmt.Sprintf("%d tools are currently supported.\n", count))
	table.Render()
}

// getTool retrieves tool information from the all the available Tool structs
// and if a valid entry is found returns it.
func getTool(tool string, tools Tools) (Tool, error) {
	for _, t := range tools {
		if tool == t.Name {
			return Tool{
				Name:           t.Name,
				Repo:           t.Repo,
				Owner:          t.Owner,
				Version:        t.Version,
				Description:    t.Description,
				NonBinary:      t.NonBinary,
				BinaryTemplate: t.BinaryTemplate,
			}, nil
		}
	}
	return Tool{}, fmt.Errorf("error: %q not found", tool)
}

// GetBinaryName returns the name of a binary for the given tool or an
// error if the tool's template cannot be parsed or executed.
func GetBinaryName(tool *Tool, os, arch, version string) (string, error) {
	if len(tool.BinaryTemplate) > 0 {
		var err error
		t := template.New(tool.Name + "_binaryname")
		t = t.Funcs(templateFuncs)

		t, err = t.Parse(tool.BinaryTemplate)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		ver := toolVersion(tool, version)
		if err := t.Execute(&buf, map[string]string{
			"OS":            os,
			"Arch":          arch,
			"Name":          tool.Name,
			"Version":       ver,
			"VersionNumber": strings.TrimPrefix(ver, "v"),
		}); err != nil {
			return "", err
		}

		res := strings.TrimSpace(buf.String())
		fmt.Printf("[DEBUG] binaryName %q\n", res)
		return res, nil
	}

	return "", errors.New("BinaryTemplate is not set")
}

func toolVersion(tool *Tool, version string) string {
	ver := tool.Version
	if len(version) > 0 {
		ver = version
	}
	return ver
}

// GetClientArch retrieves the host systems architecture and operating system.
// It will change the naming to a consistent format for linux/amd64. Further
// transformations of the architecture and operating system are then done elsewhere
// from a common base.
func GetClientArch() (arch, os string) {
	os = runtime.GOOS
	if runtime.GOOS == "linux" {
		os = "linux"
	}
	arch = runtime.GOARCH
	if runtime.GOARCH == "amd64" {
		arch = "x86_64"
	}
	return arch, os
}

// GetDownloadURL returns the downloadable assets from GitHub for use in other
// functions.
func GetDownloadURL(tool Tool, arch, opSystem, version string) (string, error) {
	releases, err := FindGithubRelease(tool.Owner, tool.Repo)
	if err != nil {
		return "", err
	}
	if version == "latest" {
		// get then latest tag which is always index 0
		version = releases[0].TagName
	}

	binaryName, err := GetBinaryName(&tool, opSystem, arch, version)
	if err != nil {
		return "", err
	}

	defer func() {
		log.Printf("Found version %q\n", version)
	}()

	for _, release := range releases {
		if release.Name == version || release.TagName == version {
			for _, asset := range release.Assets {
				if asset.Name == binaryName {
					return asset.BrowserDownloadUrl, nil
				}
			}
		}
	}
	return "", fmt.Errorf("no download URL found for %s", tool.Name)
}

// GithubAPIReleasesResponse is taken from the GitHub Releases API
// ref: https://docs.github.com/en/rest/releases/releases#list-releases
type GithubAPIReleasesResponse struct {
	Url       string `json:"url"`
	AssetsUrl string `json:"assets_url"`
	UploadUrl string `json:"upload_url"`
	HtmlUrl   string `json:"html_url"`
	Id        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeId          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []struct {
		Url      string `json:"url"`
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadUrl string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballUrl    string `json:"tarball_url"`
	ZipballUrl    string `json:"zipball_url"`
	Body          string `json:"body"`
	MentionsCount int    `json:"mentions_count,omitempty"`
}
