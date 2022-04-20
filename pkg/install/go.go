package install

import (
	"fmt"
	"github.com/danielmichaels/ds/pkg/scripts"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"os"
	"runtime"
	"strings"
)

var goLinux = &Z.Cmd{
	Name:    `go`,
	Summary: `install the latest Go version onto a linux system`,
	Description: `This will install the latest Go version on to the host system.

		**Must be linux/amd64.**`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		switch runtime.GOOS {
		case "linux":
			err := installGo()
			if err != nil {
				return err
			}
		default:
			fmt.Printf("Operating system: %s\n", runtime.GOOS)
			fmt.Println("Did not detect a linux operating system. Exiting")
			Z.Exit()
		}
		return nil
	},
}

func installGo() error {
	var s string
	fmt.Printf("Install Go 1.18 onto the system? y/N ")
	_, err := fmt.Scanln(&s)
	if err != nil {
		return err
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	switch s {
	case "y":
		script, err := scripts.Retriever("files/go-linux")
		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(script) }()
		return Z.Exec("bash", script)
	default:
		fmt.Println("Did not select 'y'. Aborting install.")
		return err
	}
}
