# Installation

```sh
go install github.com/lapwat/papeer
```

# Usage

```sh
Browse the web in the eink era

Usage:
  papeer [flags]
  papeer [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  get         Scrape URL content
  help        Help about any command
  ls          Print table of content
  version     Print the version number of papeer

Flags:
  -d, --delay int         wait before downloading next chapter, in milliseconds (default -1)
  -f, --format string     file format [md, epub, mobi] (default "md")
  -h, --help              help for papeer
  -i, --include           include URL as first chapter, in resursive mode
  -o, --output string     output file
  -q, --quiet             do not show progress bars
  -r, --recursive         create one chapter per natigation item
  -s, --selector string   table of content CSS selector
      --stdout            print to standard output

Use "papeer [command] --help" for more information about a command.
```

# Autocompletion

Execute this command in your current shell, or add it to your `.bashrc`.

```sh
. <(papeer completion [bash|fish|powershell|zsh])
```
