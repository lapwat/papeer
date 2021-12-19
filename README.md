# Papeer

Papeer is a tool that lets you scrape content from the internet. It can scrape any web page, keeping only relevant content (formatted text and images) and removing ads and menus. You can export the content to Markdown, EPUB or MOBI files.

# Usage

```
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
  -a, --author string     book author
  -d, --delay int         time to wait before downloading next chapter, in milliseconds (default -1)
  -f, --format string     file format [stdout, md, epub, mobi] (default "stdout")
  -h, --help              help for papeer
      --images            retrieve images only
  -i, --include           include URL as first chapter, in resursive mode
  -l, --limit int         limit number of chapters, in recursive mode (default -1)
  -n, --name string       book name (default: page title)
  -o, --offset int        skip first chapters, in recursive mode
      --output string     file name (default: book name)
  -r, --recursive         create one chapter per natigation item
  -s, --selector string   table of content CSS selector, in resursive mode
  -t, --threads int       download concurrency, in recursive mode (default -1)

Use "papeer [command] --help" for more information about a command.
```

# Examples

## Grab a single page

The `get` command lets you retrieve the content of a web page.

```sh
papeer get https://www.eff.org/cyberspace-independence
# A Declaration of the Independence of Cyberspace
# ===============================================

# Governments of the Industrial World, you weary giants of flesh and steel, I come from Cyberspace, the new home of Mind. On behalf of the future, I ask you of the past to leave us alone. You are not welcome among us. You have no sovereignty where we gather...
```

## Grab several pages (recursive mode)

The `recursive` option lets you extract the table of content of a website, then scrape the content of each one of its pages.

### Display table of content

Before trying the `recursive` option, it is a good idea to use the `ls` option, which lets you vizualize the content that will be retrieved. You can use several options to customize the table of content extraction, such as `selector`, `limit`, `offset` and `include`. Type `papeer help` for more information about those options.

```sh
papeer ls https://news.ycombinator.com/ --limit=5
#  #  NAME                                                      URL                                                                                             
#  1  Tailwind CSS v3.0                                         https://tailwindcss.com/blog/tailwindcss-v3                                                     
#  2  A molten salt storage solution using sodium hydroxide     https://sifted.eu/articles/salt-energy-storage-seaborg-hyme/                                    
#  3  HashiCorp IPO today                                       https://www.hashicorp.com/blog/a-new-chapter-for-hashicorp                                      
#  4  Stack Graphs                                              https://github.blog/2021-12-09-introducing-stack-graphs/
#  5  ‘Tipping point’ makes partisan polarization irreversible  https://news.cornell.edu/stories/2021/12/tipping-point-makes-partisan-polarization-irreversible 
```

### Scrape time

Once you are satisfied with the table of content listed by the `ls` command, you can actually scrape the content of those pages. You can use the same options that you specified for the `ls` command. In recursive mode, you also have the possibility to use `delay` and `threads` options.

```sh
papeer get https://news.ycombinator.com/ --recursive --delay=500 --limit=5 --format=md
# [========================================>---------------------------] Chapters 3 / 5
# [====================================================================] 1. Tailwind CSS v3.0
# [====================================================================] 2. A molten salt storage solution using sodium hydroxide
# [====================================================================] 3. HashiCorp IPO today
# [--------------------------------------------------------------------] 4. Stack Graphs
# [--------------------------------------------------------------------] 5. ‘Tipping point’ makes partisan polarization irreversible
# Markdown saved to "Hacker News.md"
```

# Installation

## From source

```sh
go get -u github.com/lapwat/papeer
```

## From binary

### On Linux / MacOS

```sh
platform=linux # use platform=darwin for MacOS
curl -L https://github.com/lapwat/papeer/releases/download/v0.3.0/papeer-v0.3.0-$platform-amd64 > papeer
chmod +x papeer
sudo mv papeer /usr/local/bin
```

### On Windows

Download [latest release](https://github.com/lapwat/papeer/releases/download/v0.3.0/papeer-v0.3.0-windows-amd64.exe).

## Install kindlegen to export websites to MOBI (optional)

```sh
TMPDIR=$(mktemp -d -t papeer-XXXXX)
curl -L https://github.com/lapwat/papeer/releases/download/kindlegen/kindlegen_linux_2.6_i386_v2_9.tar.gz > $TMPDIR/kindlegen.tar.gz
tar xzvf $TMPDIR/kindlegen.tar.gz -C $TMPDIR
chmod +x $TMPDIR/kindlegen
sudo mv $TMPDIR/kindlegen /usr/local/bin
rm -rf $TMPDIR
```

# Autocompletion

Execute this command in your current shell, or add it to your `.bashrc`.

```sh
. <(papeer completion bash)
```

Type `papeer completion bash -h` for more information.

You can replace `bash` by your own shell (zsh, fish or powershell).

# Dependencies

- `cobra` command line interface
- `go-readability` extract content from HTML
- `html-to-markdown` convert HTML to Markdown
- `go-epub` convert HTML to EPUB
- `colly` query HTML trees
- `uiprogress` display progress bars
