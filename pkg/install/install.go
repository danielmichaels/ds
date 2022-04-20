package install

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"os/exec"
)

var Cmd = &Z.Cmd{
	Name:        `install`,
	Summary:     `install executables and applications onto the host system`,
	Description: `Commands under this branch are used to install common executables and applications on the host system.`,
	Other:       []Z.Section{{`Application Options`, `curlie, gh`}},
	Commands: []*Z.Cmd{
		// imported commands
		help.Cmd,
		// local
		curlie, gh,
	},
}

func exeCheck(exe string) error {
	_, err := exec.LookPath(exe)
	if err != nil {
		return err
	}
	return nil
}
func goInstall(exe, pkg string) error {
	err := Z.Exec(exe, "install", pkg)
	if err != nil {
		return err
	}
	return nil
}
