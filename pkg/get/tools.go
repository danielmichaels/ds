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

	tools = append(tools,
		Tool{
			Name:        "ds",
			Repo:        "ds",
			Owner:       "danielmichaels",
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
			Name:        "zet-cmd",
			Repo:        "zet-cmd",
			Owner:       "danielmichaels",
			Description: "A commander for your Zettelkasten notes.",
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
			Name:        "k9s",
			Repo:        "k9s",
			Owner:       "derailed",
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
			Owner:       "derailed",
			Repo:        "popeye",
			Name:        "popeye",
			Description: "Scans live Kubernetes cluster and reports potential issues with deployed resources and configurations.",
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
			Owner:       "alexellis",
			Repo:        "k3sup",
			Name:        "k3sup",
			Description: "Bootstrap Kubernetes with k3s over SSH < 1 min.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
				{{.Name}}.exe
				{{- else if eq .OS "darwin" -}}
				{{.Name}}-darwin
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
			Owner:       "alexellis",
			Repo:        "hey",
			Name:        "hey",
			Description: "Load testing tool",
			BinaryTemplate: `
				{{$osStr := ""}}
				{{- if eq .OS "linux" -}}
					{{- if eq .Arch "x86_64" -}}
						{{$osStr = ""}}
					{{- else if eq .Arch "aarch64" -}}
						{{$osStr = "-linux-arm64"}}
					{{- else if eq .Arch "armv7l" -}}
						{{$osStr = "-linux-armv7"}}
					{{- end -}}
				{{- else if eq .OS "darwin" -}}
					{{- if eq .Arch "x86_64" -}}
						{{$osStr = "-darwin-amd64"}}
					{{- else if eq .Arch "arm64" -}}
						{{$osStr = "-darwin-arm64"}}
					{{- end -}}
				{{ else if HasPrefix .OS "ming" -}}
					{{$osStr =".exe"}}
				{{- end -}}
				
				{{.Name}}{{$osStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "openfaas",
			Repo:        "faas-cli",
			Name:        "faas-cli",
			Description: "Official CLI for OpenFaaS.",
			BinaryTemplate: `{{ if HasPrefix .OS "ming" -}}
				{{.Name}}.exe
				{{- else if eq .OS "darwin" -}}
				{{.Name}}-darwin
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

				{{.Name}}_{{.VersionNumber}}_{{ ToLower $osStr}}_{{ ToLower $archStr}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "rs",
			Repo:        "curlie",
			Name:        "curlie",
			Description: "The power of curl, the ease of use of httpie.",
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

				{{.Name}}_{{.VersionNumber}}_{{ ToLower $osStr}}_{{ ToLower $archStr}}.{{$extStr}}`,
		})

	tools = append(tools,
		Tool{
			Name:        "rclone",
			Repo:        "rclone",
			Owner:       "rclone",
			Description: "\"rsync for cloud storage\" - Google Drive, S3, Dropbox, Backblaze B2, One Drive, Swift, Hubic, Wasabi, Google Cloud Storage, Yandex Files",
			NonBinary:   false,
			BinaryTemplate: `
				{{$osStr := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$osStr = "Windows"}}
				{{- else if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$osStr = "Darwin"}}
				{{- end -}}

				{{$archStr := .Arch}}
				{{- if eq .Arch "armv7l" -}}
				{{$archStr = "arm"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$archStr = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{$archStr = "amd64"}}
				{{- end -}}

				{{.Name}}-v{{.VersionNumber}}-{{$osStr}}-{{$archStr}}.zip`,
		})

	tools = append(tools,
		Tool{
			Name:        "hugo",
			Repo:        "hugo",
			Owner:       "gohugoio",
			Description: "The world’s fastest framework for building websites.",
			NonBinary:   false,
			BinaryTemplate: `
				{{$osStr := ""}}
				{{- if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$osStr = "darwin"}}
				{{- end -}}

				{{$archStr := .Arch}}
				{{- if eq .Arch "armv7l" -}}
				{{$archStr = "arm"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$archStr = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{$archStr = "amd64"}}
				{{- end -}}

				{{.Name}}_{{.VersionNumber}}_{{$osStr}}-{{$archStr}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "goreleaser",
			Repo:        "goreleaser",
			Name:        "goreleaser",
			Description: "Deliver Go binaries as fast and easily as possible",
			BinaryTemplate: `
		{{$osStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$osStr = "Windows"}}
		{{- else if eq .OS "linux" -}}
		{{$osStr = "Linux"}}
		{{- else if eq .OS "darwin" -}}
		{{$osStr = "Darwin"}}
		{{- end -}}

		{{$archStr := ""}}
		{{- if eq .Arch "x86_64" -}}
		{{$archStr = "x86_64"}}
		{{- else if eq .Arch "aarch64" -}}
        {{$archStr = "arm64"}}
		{{- end -}}

		{{$archiveStr := ""}}
		{{ if HasPrefix .OS "ming" -}}
		{{$archiveStr = "zip"}}
		{{- else -}}
		{{$archiveStr = "tar.gz"}}
		{{- end -}}

		{{.Name}}_{{$osStr}}_{{$archStr}}.{{$archiveStr}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "FiloSottile",
			Repo:        "mkcert",
			Name:        "mkcert",
			Description: "A simple zero-config tool to make locally trusted development certificates with any names you'd like.",
			BinaryTemplate: `
				{{ $osStr := "" }}
				{{ $archStr := "" }}
				{{$archiveStr := ""}}
				{{- if HasPrefix .OS "ming" -}}
				{{$archiveStr = ".exe"}}
				{{- else -}}
				{{$archiveStr = ""}}
				{{- end -}}
	
				{{ if HasPrefix .OS "ming" -}}
				{{ $osStr = "windows" }}
				{{- else if eq .OS "linux" -}}
				{{ $osStr = "linux" }}
				{{- else if eq .OS "darwin" -}}
				{{ $osStr = "darwin" }}
				{{- end -}}
	
				{{- if eq .Arch "x86_64" -}}
				{{ $archStr = "amd64" }}
				{{- else if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $archStr = "arm" }}
				{{- end -}}
	
				{{.Name}}-v{{.VersionNumber}}-{{$osStr}}-{{$archStr}}{{$archiveStr}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "junegunn",
			Repo:        "fzf",
			Name:        "fzf",
			Description: "General-purpose command-line fuzzy finder",
			BinaryTemplate: `
				{{ $osStr := "linux" }}
				{{ $ext := ".tar.gz" }}
				{{ if HasPrefix .OS "ming" -}}
				{{ $osStr = "windows" }}
				{{ $ext = ".zip" }}
				{{- else if eq .OS "darwin" -}}
				{{  $osStr = "darwin" }}
				{{ $ext = ".zip" }}
				{{- end -}}
				{{ $archStr := "amd64" }}
				{{- if eq .Arch "armv6l" -}}
				{{ $archStr = "armv6" }}
				{{- else if eq .Arch "armv7l" -}}
				{{ $archStr = "armv7" }}
				{{- else if eq .Arch "arm64" -}}
				{{ $archStr = "arm64" }}
				{{- else if eq .Arch "aarch64" -}}
				{{ $archStr = "arm64" }}
				{{- end -}}
				{{.Name}}-{{.VersionNumber}}-{{$osStr}}_{{$archStr}}{{$ext}}
				`,
		})
	// k3sup

	tools = append(tools,
		Tool{
			Owner:       "stedolan",
			Repo:        "jq",
			Name:        "jq",
			Description: "jq is a lightweight and flexible command-line JSON processor",
			BinaryTemplate: `{{$arch := "arm"}}
				{{- if eq .Arch "x86_64" -}}
				{{$arch = "64"}}
				{{- else if eq .Arch "arm64" -}}
				{{$arch = "64"}}
				{{- else -}}
				{{$arch = "32"}}
				{{- end -}}

				{{$ext := ""}}
				{{$os := .OS}}

				{{ if HasPrefix .OS "ming" -}}
				{{$ext = ".exe"}}
				{{$os = "win"}}
				{{- else if eq .OS "darwin" -}}
				{{$os = "osx-amd"}}
				{{- end -}}

				jq-{{$os}}{{$arch}}{{$ext}}`,
		})

	// kubectx

	// kubetail
	tools = append(tools,
		Tool{
			Owner:       "stern",
			Repo:        "stern",
			Name:        "stern",
			Description: "Multi pod and container log tailing for Kubernetes.",
			BinaryTemplate: `{{$arch := "arm"}}
				{{- if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- else if eq .Arch "x86_64" -}}
				{{$arch = "amd64"}}
				{{- end -}}

				{{$os := .OS}}
				{{$ext := "tar.gz"}}

				{{ if HasPrefix .OS "ming" -}}
				{{$os = "windows"}}
				{{- end -}}

				{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.tar.gz`,
		})

	tools = append(tools,
		Tool{
			Owner:       "jesseduffield",
			Repo:        "lazygit",
			Name:        "lazygit",
			Description: "A simple terminal UI for git commands.",
			BinaryTemplate: `
				{{$os := ""}}
				{{$ext := "tar.gz" }}
				{{ if HasPrefix .OS "ming" -}}
				{{$os = "Windows"}}
				{{$ext = "zip" }}
				{{- else if eq .OS "linux" -}}
				{{$os = "Linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$os = "Darwin"}}
				{{- end -}}

				{{$arch := .Arch}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "x86_64"}}
				{{- else if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
				{{$arch = "arm64"}}
				{{- end -}}
				{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "jesseduffield",
			Repo:        "lazydocker",
			Name:        "lazydocker",
			Description: "The lazier way to manage everything docker.",
			BinaryTemplate: `
				{{$os := ""}}
				{{$ext := "tar.gz" }}
				{{ if HasPrefix .OS "ming" -}}
				{{$os = "Windows"}}
				{{$ext = "zip" }}
				{{- else if eq .OS "linux" -}}
				{{$os = "Linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$os = "Darwin"}}
				{{- end -}}

				{{$arch := .Arch}}
				{{ if (or (eq .Arch "x86_64") (eq .Arch "amd64")) -}}
				{{$arch = "x86_64"}}
				{{- else if (or (eq .Arch "aarch64") (eq .Arch "arm64")) -}}
				{{$arch = "arm64"}}
				{{- end -}}
				{{.Name}}_{{.VersionNumber}}_{{$os}}_{{$arch}}.{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "docker",
			Repo:        "compose",
			Name:        "docker-compose",
			Description: "Define and run multi-container applications with Docker.",
			BinaryTemplate: `
				{{$arch := .Arch}}

				{{$osStr := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$osStr = "windows"}}
				{{- else if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				  {{- if eq .Arch "armv7l" -}}
				  {{ $arch = "armv7"}}
				  {{- end }}
				{{- else if eq .OS "darwin" -}}
				{{$osStr = "darwin"}}
				{{- end -}}
				{{$ext := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$ext = ".exe"}}
				{{- end -}}

				{{.Name}}-{{$osStr}}-{{$arch}}{{$ext}}`,
		})
	tools = append(tools,
		Tool{
			Owner:       "nats-io",
			Repo:        "natscli",
			Name:        "nats",
			Description: "Utility to interact with and manage NATS.",
			BinaryTemplate: `{{$arch := .Arch}}
				{{ if eq .Arch "x86_64" -}}
				{{$arch = "amd64"}}
				{{- else if eq .Arch "armv6l" -}}
				{{$arch = "arm6"}}
				{{- else if eq .Arch "armv7l" -}}
				{{$arch = "arm7"}}
				{{- else if eq .Arch "aarch64" -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{$osStr := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$osStr = "windows"}}
				{{- else if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$osStr = "darwin"}}
				{{- end -}}

				{{.Name}}-{{.VersionNumber}}-{{$osStr}}-{{$arch}}.zip`,
		})

	tools = append(tools,
		Tool{
			Owner:       "argoproj",
			Repo:        "argo-cd",
			Name:        "argocd",
			Description: "Declarative, GitOps continuous delivery tool for Kubernetes.",
			BinaryTemplate: `
				{{$arch := .Arch}}
				{{- if eq .Arch "x86_64" -}}
				{{$arch = "amd64"}}
				{{- else if or (eq .Arch "aarch64") (eq .Arch "arm64") -}}
				{{$arch = "arm64"}}
				{{- end -}}

				{{$osStr := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$osStr = "windows"}}
				{{- else if eq .OS "linux" -}}
				{{$osStr = "linux"}}
				{{- else if eq .OS "darwin" -}}
				{{$osStr = "darwin"}}
				{{- end -}}

				{{$ext := ""}}
				{{ if HasPrefix .OS "ming" -}}
				{{$ext = ".exe"}}
				{{- end -}}

				argocd-{{$osStr}}-{{$arch}}{{$ext}}`,
		})

	tools = append(tools,
		Tool{
			Owner:       "containerd",
			Repo:        "nerdctl",
			Name:        "nerdctl",
			Description: "Docker-compatible CLI for containerd, with support for Compose",
			BinaryTemplate: `
				{{ $file := "" }}
				{{- if eq .OS "linux" -}}
					{{- if eq .Arch "armv6l" -}}
						{{ $file = "arm-v7.tar.gz" }}
					{{- else if eq .Arch "armv7l" -}}
						{{ $file = "arm-v7.tar.gz" }}
					{{- else if eq .Arch "aarch64" -}}
						{{ $file = "arm64.tar.gz" }}
					{{- else -}}
						{{ $file = "amd64.tar.gz" }}
					{{- end -}}
				{{- end -}}

				{{.Name}}-{{.VersionNumber}}-{{.OS}}-{{$file}}`,
		})

	// packer

	// terraform

	// nomad

	// waypoint

	//vault
	return tools
}
