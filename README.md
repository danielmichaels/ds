# 🌳 *Do Something*

[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

## Install

`ds` is meant to be portable with a simple installation process.

### Recommended method (easy mode)

*Requires Go 1.18+*
```
go install github.com/danielmichaels/cmd/ds@latest
```

### Go Install doesn't work

If the device does not have Go installed, or only has access to an older version, then
you will need to use one of the following installation methods. It's just as easy but has
some extra steps and requires you to manually select your operating system.

**Curl the Binary**

*Requires [jq](https://stedolan.github.io/jq/)*

Running the following shell command will download the latest release to your current working directory. Once downloaded, you can run the
binary with `./ds` or install it into your path as described in the next section.

`OS=linux_x86 URL=$(curl api.github.com/repos/danielmichaels/ds/releases/latest -Ls | jq -r '.assets[].browser_download_url' | grep -i $OS) &&  curl -L -s $URL | tar xvz -C . ds`

*Notes*: setting the `OS=` is important and must match your current operating system. The most common options are:

- `linux_x86`
- `linux_arm64`
- `windows_x86`
- `darwin_x86`
- `darwin_arm64`

**Browser Based Download**

Go to the [releases](https://github.com/danielmichaels/ds/releases) page and 
download the appropriate binary for your operating system.

For linux you can then drop the binary into your `PATH` or call the binary directly.

To drop it into your path simply copy the `ds` binary to `$HOME/.local/bin`. If that directory is 
in your `$PATH` it will now be accessible by calling `ds` from your terminal.

## Usage

To use this binary simply execute `ds` in your terminal. All the commands available are listed in 
the `COMMANDS` section of the output.

As a general rule any command that has an alias will contain several subcommands. For instance,
`ds scripts` contains various subcommands which must be called using either the alias `s` or the
full command of `scripts`.

Alias' can be identified in the `COMMANDS` list by the pipe operator such as `s|scripts`. 

## Example Usage

`ds` comes with many commands, some are local but others are imported from 
other branches, or packages such as `y2j` and `yq`.

**A few examples:**

To get your current external IP address use `ds scripts ipify`, or its shorthand alias of `ds s ipify`.

From your external IP address, print its IP information from [ipinfo.io] 
using two `ds` commands together; 

```shell
ds scripts ipinfo $(ds scripts ip)
```

Retrieving the current weather for a given location:

```shell
ds scripts weather london
```

Parsing JSON or YAML by using `ds yq`.

```shell
echo '{"arr":{"index1":"one","index2":"two"},"nested":{"arr":{"indexArr1": {"idx":["one","two"]}}},"test":"string"}' | ds yq .arr                     
# outputs: {"index1": "one", "index2": "two"}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C ds ds
```

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.

Example: `ds y2j help`

## Other Examples

* <https://github.com/rwxrob/z> - *the one that started it all* by [rwxrob]
* <https://github.com/rwxrob/bonzai-example> - a template to use when creating your own

[rwxrob]: https://github.com/rwxrob
[ipinfo.io]: https://ipinfo.io
