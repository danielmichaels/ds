// Copyright 2022 ds Daniel Michaels
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/danielmichaels/ds/pkg/install"
	"github.com/danielmichaels/ds/pkg/scripts"
	"github.com/danielmichaels/zet-cmd"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/uniq"
	"github.com/rwxrob/vars"
	"github.com/rwxrob/y2j"
	"github.com/rwxrob/yq"
)

func init() {
	err := Z.Conf.SoftInit()
	if err != nil {
		panic("system error initializing conf")
	}
	err = Z.Vars.SoftInit()
	if err != nil {
		panic("system error initializing vars")
	}
}

func main() {
	Z.AllowPanic = true
	Cmd.Run()
}

var Cmd = &Z.Cmd{
	Name:      `ds`,
	Summary:   `*Do Something* is a single binary to rule them all`,
	Version:   `v0.1.7`,
	Copyright: `Copyright 2022 Daniel Michaels`,
	License:   `Apache-2.0`,
	Shortcuts: Z.ArgMap{
		"uuid":   {"uniq", "uuid"},
		"isosec": {"uniq", "isosec"},
		"env":    {"scripts", "env-check"},
	},
	Commands: []*Z.Cmd{
		// imported
		help.Cmd, conf.Cmd, yq.Cmd, vars.Cmd, y2j.Cmd, vars.Cmd, uniq.Cmd, zet.Cmd,
		// internal
		scripts.Cmd, install.Cmd,
	},
	Issues: `github.com/danielmichaels/ds/issues`,
	Site:   `danielms.site`,
	Source: `git@github.com:danielmichaels/ds.git`,
	Description: `
		This is my multi-call library which is set to replace the need for many
		small bash scripts. Instead everything here will be portable and requires
		a single curl command to pull it down on to any box. 

		This functionality contained inside is custom to my needs and may not work for anyone else.

		Intentions for **{{ .Name }}** 

		* All-in-one binary
		* Install often used applications/binaries
		* Replace shell scripts on host 

		`,
}
