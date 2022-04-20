package install

import (
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

var curlie = &Z.Cmd{
	Name:     `curlie`,
	Summary:  `install github.com/rs/curlie`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		if err := exeCheck("go"); err == nil {
			err = goInstall("go", "github.com/rs/curlie@latest")
			if err != nil {
				return err
			}
		}

		if err := exeCheck("brew"); err == nil {
			err = goInstall("brew", "github.com/rs/curlie@latest")
			if err != nil {
				return err
			}
		}

		if err := exeCheck("scoop"); err == nil {
			err = goInstall("scoop", "github.com/rs/curlie@latest")
			if err != nil {
				return err
			}
		}
		fmt.Println("curlie installed successfully")
		return nil
	},
}
