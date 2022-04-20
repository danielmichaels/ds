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

var gh = &Z.Cmd{
	Name:    `gh`,
	Summary: `install github.com/cli/cli`,
	Description: `Installs github's cli *gh* and supports mac, windows and linux. Only Centos and 
		Debian-based linux derivatives are supported. Mac and windows rely upon *brew* and *scoop*
		respectively. No binary installations are supported as yet. Linux requires sudo privileges.`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		var success = func() { fmt.Println("gh successfully installed") }
		switch runtime.GOOS {
		case "darwin":
			if err := exeCheck("brew"); err == nil {
				// todo add upgrade option
				err = goInstall("brew", "gh")
				if err != nil {
					return err
				}
				success()
				return nil
			}
		case "windows":
			if err := exeCheck("scoop"); err == nil {
				err = goInstall("scoop", "gh")
				if err != nil {
					return err
				}
				success()
				return nil
			}
		case "linux":
			var s string
			fmt.Printf("options: 1. ubuntu/debian/rpi 2. centos/rhel/fedora\nEnter a number: ")
			_, err := fmt.Scanln(&s)
			if err != nil {
				return err
			}

			s = strings.TrimSpace(s)
			s = strings.ToLower(s)

			switch s {
			case "1":
				script, err := scripts.Retriever("files/gh-ubuntu")
				if err != nil {
					return err
				}
				defer func() { _ = os.Remove(script) }()
				return Z.Exec("bash", script)
			case "2":
				script, err := scripts.Retriever("files/gh-centos")
				if err != nil {
					return err
				}
				defer func() { _ = os.Remove(script) }()
				return Z.Exec("bash", script)
			default:
				fmt.Println("did not supply a valid option - if your os is not available please raise an issue")
				return err
			}
		default:
			if err := exeCheck("go"); err == nil {
				err = goInstall("go", "github.com/rs/curlie@latest")
				if err != nil {
					return err
				}
				success()
				return nil
			}
		}
		return nil
	},
}
