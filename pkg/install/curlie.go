package install

import (
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"runtime"
)

var curlie = &Z.Cmd{
	Name:     `curlie`,
	Summary:  `install github.com/rs/curlie`,
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
