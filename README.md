# 🌳 *Do Something*

[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

## Install

`ds` is meant to be portable with a simple install process.

### Recommended method

Go to the [releases](https://github.com/danielmichaels/ds/releases) page and 
download the appropriate binary for your operating system.

For linux you can then drop the binary into your `PATH` or call the binary directly.

To drop it into your path simply copy the `ds` binary to `$HOME/.local/bin`. If that directory is 
in your `$PATH` it will now be accessible by calling `ds` from your terminal.

### Standalone

```
go install github.com/danielmichaels/ds/ds@latest
```

## Usage

To use this binary simply execute `ds` in your terminal. All the commands available are listed in 
the `COMMANDS` section of the output.

As a general rule any command that has an alias will contain several subcommands. For instance,
`ds scripts` contains various subcommands which must be called using either the alias `s` or the
full command of `scripts`.

Alias' can be identified in the `COMMANDS` list by the pipe operator such as `s|scripts`. 

To get your current external IP address use `ds scripts ipify`, or its shorthand alias of `ds s ipify`.


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
