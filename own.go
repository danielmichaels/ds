// Copyright 2022 foo Authors
// SPDX-License-Identifier: Apache-2.0

package ds

// Go treats all files as if they are, more or less, in the same large
// file. Create separate files to help you and others find the code you
// need quickly.

import (
	"errors"
	"fmt"
	"github.com/danielmichaels/ds/scripts"
	"github.com/rwxrob/help"
	json2 "github.com/rwxrob/json"
	"net/http"
	"os"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
)

var s = &Z.Cmd{
	Name:     `scripts`,
	Summary:  `call custom scripts`,
	Commands: []*Z.Cmd{help.Cmd, weather, ipify, til, hugo, envCheck, ipinfo, pfsenseManager},
}

var weather = &Z.Cmd{
	Name:     `weather`,
	Summary:  `the current weather for Canberra`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, none ...string) error {
		return Z.Exec("curl", "v2.wttr.in/Canberra")
	},
}

var ipify = &Z.Cmd{
	Name:     `ipify`,
	Summary:  `print out the current external IP address`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, none ...string) error {
		script, err := scripts.Retriever("files/ipify")

		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(script) }()

		return Z.Exec("sh", script)
	},
}

var hugo = &Z.Cmd{
	Name:     `hugo`,
	Summary:  `run the hugo docker image`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		cmdlineArgs := strings.Join(args, " ")

		script, err := scripts.Retriever("files/hugo")
		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(script) }()

		return Z.Exec("sh", script, cmdlineArgs)
	},
}

var envCheck = &Z.Cmd{
	Name:     `env-check`,
	Summary:  `prompt user before showing environment variables`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(caller *Z.Cmd, args ...string) error {
		var s string

		fmt.Printf("Show 'env'? y/N")
		_, err := fmt.Scanln(&s)
		if err != nil {
			return err
		}
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)

		if s == "y" || s == "yes" {
			return Z.Exec("env")
		}
		return nil
	},
}

var ipinfo = &Z.Cmd{
	Name:     `ipinfo`,
	Summary:  `return information about an IP address`,
	Commands: []*Z.Cmd{help.Cmd},
	Usage:    `enter a valid IP address e.g. 1.1.1.1`,
	Call: func(caller *Z.Cmd, args ...string) error {
		if len(args) < 1 {
			fmt.Println("must provide an IP address")
			return caller.UsageError()
		}
		token := Z.Conf.Query(".ipinfo")
		bearer := fmt.Sprintf("Bearer %s", token)

		var result map[string]interface{}

		cl := json2.Client
		cl.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			for k, v := range via[0].Header {
				r.Header[k] = v
			}
			return nil
		}
		headers := map[string]string{}
		headers["Authorization"] = bearer
		req := json2.Request{
			Method: "GET",
			URL:    fmt.Sprintf("https://ipinfo.io/%s", args[0]),
			Query:  nil,
			Header: headers,
			Body:   nil,
			Into:   &result,
		}
		err := json2.Fetch(&req)
		if err != nil {
			return err
		}
		marshal, err := json2.MarshalIndent(&result, " ", " ")
		if err != nil {
			return err
		}
		fmt.Println(string(marshal))
		return nil
	},
}

var til = &Z.Cmd{
	Name:    `til`,
	Usage:   `til <filename>`,
	Params:  []string{"filename"},
	Summary: `create a new blog post about something I learned`,
	Call: func(caller *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			fmt.Println(caller.UsageError())
			return errors.New("no args supplied. must supply a name for the file")
		}
		cmdlineArgs := strings.Join(args, " ")

		script, err := scripts.Retriever("files/til")
		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(script) }()

		return Z.Exec("bash", script, cmdlineArgs)
	},
}

var pfsenseManager = &Z.Cmd{
	Name:     `pfsense-vm-manager`,
	MinArgs:  1,
	Params:   []string{"start", "stop"},
	Summary:  `pfsense-vm-manager starts or stops multiple pfsense virtual machines for local testing`,
	Commands: []*Z.Cmd{help.Cmd},
	Description: `
	**pfsense-vm-manager** is a shortcut to stop or start multiple pfsense
	virtual machines for testing locally.`,
	Call: func(caller *Z.Cmd, args ...string) error {
		cmdlineArgs := strings.Join(args, " ")

		script, err := scripts.Retriever("files/pfsense-vm-manager")

		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(script) }()

		return Z.Exec("bash", script, cmdlineArgs)
	},
}
