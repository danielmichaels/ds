// Copyright 2022 ds Daniel Michaels
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/danielmichaels/ds"
	Z "github.com/rwxrob/bonzai/z"
)

func main() {
	Z.AllowPanic = true
	// Global alias' shortcuts for frequently used
	Z.Aliases = map[string][]string{
		"uuid": {"uniq", "uuid"},
		"unix": {"scripts", "date"},
	}
	ds.Cmd.Run()
}
