// Copyright 2022 ds Daniel Michaels
// SPDX-License-Identifier: Apache-2.0

package scripts

import (
	"embed"
	"errors"
	"fmt"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

//go:embed "files"
var ScriptFiles embed.FS

// tmpFileCreator creates a temp file containing the script data passed in. The
// file is not removed in this function - it must be done in the caller.
func tmpFileCreator(script []byte) (string, error) {
	tmp, err := ioutil.TempFile("/tmp", "ds-file")
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(tmp.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	_, err = f.Write(script)
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	return tmp.Name(), nil
}

// Retriever reads the file provided and returns its data. It expects a valid
// shell script. A temp file is created and that filename is then returned to be
// called by Z.Exec
func Retriever(filename string) (string, error) {
	data, err := ScriptFiles.ReadFile(filename)
	if err != nil {
		return "", err
	}

	script, err := tmpFileCreator(data)
	if err != nil {
		return "", err
	}
	return script, nil
}

var Cmd = &Z.Cmd{
	Name:     `scripts`,
	Summary:  `call custom scripts`,
	Aliases:  []string{"s"},
	Commands: []*Z.Cmd{help.Cmd, weather, ipify, til, hugo, envCheck, ipinfo, pfsenseManager, epochDate},
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
		script, err := Retriever("files/ipify")

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

		script, err := Retriever("files/hugo")
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

		fmt.Printf("Show 'env'? y/N ")
		_, err := fmt.Scanln(&s)
		if err != nil {
			return err
		}

		s = strings.TrimSpace(s)
		s = strings.ToLower(s)

		if s == "y" {
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

		cl := json.Client
		cl.CheckRedirect = func(r *http.Request, via []*http.Request) error {
			for k, v := range via[0].Header {
				r.Header[k] = v
			}
			return nil
		}
		headers := map[string]string{}
		headers["Authorization"] = bearer
		req := json.Request{
			Method: "GET",
			URL:    fmt.Sprintf("https://ipinfo.io/%s", args[0]),
			Query:  nil,
			Header: headers,
			Body:   nil,
			Into:   &result,
		}
		err := json.Fetch(&req)
		if err != nil {
			return err
		}
		marshal, err := json.MarshalIndent(&result, " ", " ")
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

		script, err := Retriever("files/til")
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

		script, err := Retriever("files/pfsense-vm-manager")

		if err != nil {
			return err
		}
		defer func() { _ = os.Remove(script) }()

		return Z.Exec("bash", script, cmdlineArgs)
	},
}

var epochDate = &Z.Cmd{
	Name:     `date`,
	MinArgs:  1,
	Summary:  `convert timestamp to the system local time`,
	Commands: []*Z.Cmd{help.Cmd},
	Usage:    `ds scripts date`,
	Description: `
		Convert a *unix* timestamp to this systems local time.
		
		example: ds scripts date 1647826365`,
	Call: func(caller *Z.Cmd, args ...string) error {
		epoch := args[0]
		return Z.Exec("date", "-d", fmt.Sprintf("@%s", epoch))
	},
}
