package get

import "strings"

var templateFuncs = map[string]interface{}{
	"HasPrefix": func(s, prefix string) bool { return strings.HasPrefix(s, prefix) },
	"ToLower":   strings.ToLower,
}

type Tool struct {
	// Name of the tool
	Name string

	// Repo is the GitHub repo
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

	// BinaryTemplate is the naming convention for a binary from GitHub.
	// Using runtime.GOOS and runtime.GOARCH it is possible to determine the
	// binary name, and BinaryTemplate must match it.
	BinaryTemplate string

	// URLTemplate specifies a Go template for the download URL
	// override the OS, architecture and extension
	// All whitespace will be trimmed
	URLTemplate string
}

// IsArchive determines if a binary is in archive format from the download URL.
func (tool Tool) IsArchive(downloadURL string) (bool, error) {
	return strings.HasSuffix(downloadURL, "tar.gz") ||
		strings.HasSuffix(downloadURL, "zip") ||
		strings.HasSuffix(downloadURL, "tgz"), nil
}

type Tools []Tool

func (t Tools) Len() int { return len(t) }

func (t Tools) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

func (t Tools) Less(i, j int) bool {
	var ti = t[i]
	var tj = t[j]
	var tiNameLower = strings.ToLower(ti.Name)
	var tjNameLower = strings.ToLower(tj.Name)
	if tiNameLower == tjNameLower {
		return ti.Name < tj.Name
	}
	return tiNameLower < tjNameLower
}

func MakeTools() Tools {
	var tools []Tool

	// arkade
	// curlie
	// gh cli
	// glab
	// k9s
	// rclone
	// kubectl
	// ds
	// cmd-zet

	tools = append(tools,
		Tool{
			Name:        "ds",
			Repo:        "ds",
			Owner:       "danielmichaels",
			Version:     "v0.1.7",
			Description: "A command box for a danielmichaels things.",
			NonBinary:   false,
			BinaryTemplate: `
		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "Windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "Linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "Darwin"}}
		{{- end -}}

		{{$archStr := .Arch}}
		{{- if eq .Arch "armv7l" -}}
		{{$archStr = "arm"}}
		{{- else if eq .Arch "aarch64" -}}
		{{$archStr = "arm64"}}
		{{- end -}}

		{{.Name}}_{{.VersionNumber}}_{{$osStr}}_{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Name:  "k9s",
			Repo:  "k9s",
			Owner: "derailed",
			//Version:     "v0.26.7",
			Description: "A kubernetes TUI.",
			NonBinary:   false,
			BinaryTemplate: `
		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "Windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "Linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "Darwin"}}
		{{- end -}}

		{{$archStr := .Arch}}
		{{- if eq .Arch "armv7l" -}}
		{{$archStr = "arm"}}
		{{- else if eq .Arch "aarch64" -}}
		{{$archStr = "arm64"}}
		{{- end -}}

		{{.Name}}_{{$osStr}}_{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Name:        "arkade",
			Repo:        "arkade",
			Owner:       "alexellis",
			Description: "Portable marketplace for downloading your favourite devops CLIs and installing helm charts, with a single command.",
			BinaryTemplate: `
			{{ if HasPrefix .OS "ming" -}}
			{{.Name}}.exe
			{{- else if eq .OS "darwin" -}}
				{{ if eq .Arch "arm64" -}}
				{{.Name}}-darwin-arm64
				{{- else -}}
				{{.Name}}-darwin
				{{- end -}}
			{{- else if eq .Arch "armv6l" -}}
			{{.Name}}-armhf
			{{- else if eq .Arch "armv7l" -}}
			{{.Name}}-armhf
			{{- else if eq .Arch "aarch64" -}}
			{{.Name}}-arm64
			{{- else -}}
			{{.Name}}
			{{- end -}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "cli",
			Repo:        "cli",
			Name:        "gh",
			Description: "GitHub’s official command line tool.",
			BinaryTemplate: `

	{{$extStr := "tar.gz"}}
	{{ if HasPrefix .OS "ming" -}}
	{{$extStr = "zip"}}
	{{- end -}}

	{{$osStr := ""}}
	{{ if HasPrefix .OS "ming" -}}
	{{$osStr = "windows"}}
	{{- else if eq .OS "linux" -}}
	{{$osStr = "linux"}}
	{{- else if eq .OS "darwin" -}}
	{{$osStr = "macOS"}}
	{{- end -}}

	{{$archStr := .Arch}}
	{{- if eq .Arch "aarch64" -}}
	{{$archStr = "arm64"}}
	{{- else if eq .Arch "x86_64" -}}
	{{$archStr = "amd64"}}
	{{- end -}}

	gh_{{.VersionNumber}}_{{ ToLower $osStr}}_{{ ToLower $archStr}}.{{$extStr}}`,
		})

	return tools
}
