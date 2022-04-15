// Copyright 2022 ds Daniel Michaels
// SPDX-License-Identifier: Apache-2.0

// Package ds provides the Bonzai command branch of the same name.
package ds

import (
	"github.com/danielmichaels/ds/scripts"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/config"
	"github.com/rwxrob/help"
	"github.com/rwxrob/uniq"
	"github.com/rwxrob/vars"
	"github.com/rwxrob/y2j"
	"github.com/rwxrob/yq"
)

var Cmd = &Z.Cmd{
	Name:      `ds`,
	Summary:   `*Do Something* is a single binary to rule them all`,
	Version:   `v0.0.2`,
	Copyright: `Copyright 2022 Daniel Michaels`,
	License:   `Apache-2.0`,
	Commands:  []*Z.Cmd{help.Cmd, config.Cmd, scripts.Cmd, yq.Cmd, y2j.Cmd, vars.Cmd, uniq.Cmd},
	Issues:    `github.com/danielmichaels/ds/issues`,
	Site:      `danielms.site`,
	Source:    `github.com/danielmichaels/ds`,
	Description: `
		This is my multi-call library which is set to replace the need for many
		small bash scripts. Instead everything here will be portable and requires
		a single curl command to pull it down on to any box. This functionality
		contained inside is custom to my needs and may not work for anyone else.
		`,
}
