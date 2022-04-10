// Copyright 2022 foo Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/danielmichaels/ds"
	Z "github.com/rwxrob/bonzai/z"
)

func main() {
	Z.AllowPanic = true
	ds.Cmd.Run()
}
