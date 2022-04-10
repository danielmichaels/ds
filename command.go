// Copyright 2022 ds Authors
// SPDX-License-Identifier: Apache-2.0

// Package ds provides the Bonzai command branch of the same name.
package ds

import (
	"github.com/rwxrob/bonzai/comp"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/config"
	"github.com/rwxrob/help"
	"log"
)

var Cmd = &Z.Cmd{

	Name:      `ds`,
	Summary:   `*Do Something* is a single binary to rule them all`,
	Version:   `v0.0.1`,
	Copyright: `Copyright 2022 Daniel Michaels`,
	License:   `Apache-2.0`,
	Commands:  []*Z.Cmd{help.Cmd, config.Cmd, file, s},
	Description: `
		This is my multi-call library which is set to replace the need for many
		small bash scripts. Instead everything here will be portable and requires
		a single curl command to pull it down on to any box. This functionality
		contained inside is custom to my needs and may not work for anyone else.
		`,
}

var file = &Z.Cmd{
	Name:      `file`,
	Commands:  []*Z.Cmd{help.Cmd},
	Completer: comp.File,
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			return x.UsageError()
		}
		log.Printf("would show file information about %v", args[0])
		return nil
	},
}
