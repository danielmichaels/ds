// Copyright 2022 ds Daniel Michaels
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/danielmichaels/ds/pkg/install"
	"github.com/danielmichaels/ds/pkg/scripts"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/compfile"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/uniq"
	"github.com/rwxrob/vars"
	"github.com/rwxrob/y2j"
	"github.com/rwxrob/yq"
)

func main() {
	Z.AllowPanic = true
	// Global shortcuts for frequently used commands
	Z.Shortcuts = map[string][]string{
		"uuid":   {"uniq", "uuid"},
		"isosec": {"uniq", "isosec"},
		"env":    {"scripts", "env-check"},
	}

	Cmd.Run()
}

var Cmd = &Z.Cmd{
	Name:      `ds`,
	Summary:   `*Do Something* is a single binary to rule them all`,
	Version:   `v0.1.0`,
	Copyright: `Copyright 2022 Daniel Michaels`,
	License:   `Apache-2.0`,
	Comp:      compfile.New(),
	Commands: []*Z.Cmd{
		// imported
		help.Cmd, conf.Cmd, yq.Cmd, vars.Cmd, y2j.Cmd, vars.Cmd, uniq.Cmd,
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

		1. All-in-one binary
		2. Install often used applications/binaries
		3. Replace shell scripts on host

		`,
	Other: []Z.Section{
		{"Shortcuts", `
			For brevity the are several helpful shortcuts provided. These enable top level execution of commands that
			may be several subcommands deep in the node tree.

			**Shortcut**: Command Equivalent

			* **env**:  "{{ .Name }} scripts env-check"
			* **uuid**: "uniq uuid"
			* **isosec**:  "uniq isosec"


`},
	},
}
