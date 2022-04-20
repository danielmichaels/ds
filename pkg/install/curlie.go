package install

import (
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"runtime"
)

var curlie = &Z.Cmd{
	Name:    `curlie`,
	Summary: `install github.com/rs/curlie`,
	Description: `Installs *curlie* a drop-in replacement for cURL. It offers all the features
			of curl but with *jq* output by default. 

			The installer uses *brew* and Mac and *scoop* on Windows to install the application. Linux requires
			Go to be installed. In future binary releases will be available making the need for Go to be
			installed redundant.`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		var success = func() { fmt.Println("curlie successfully installed") }
		switch runtime.GOOS {
		case "darwin":
			if err := exeCheck("brew"); err == nil {
				err = goInstall("brew", "rs/tap/curlie")
				if err != nil {
					return err
				}
				success()
				return nil
			}
		case "windows":
			if err := exeCheck("scoop"); err == nil {
				err = goInstall("scoop", "curlie")
				if err != nil {
					return err
				}
				success()
				return nil
			}
		case "linux":
			if err := exeCheck("go"); err == nil {
				err = goInstall("go", "github.com/rs/curlie@latest")
				if err != nil {
					return err
				}
				success()
				return nil
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
