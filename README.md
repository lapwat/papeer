# Papeer

Papeer is a powerful **ereader internet vacuum**. It can scrape any website, removing ads and keeping only the relevant content (formatted text and images). You can export the content to Markdown, EPUB or MOBI files.

# Table of contents

- [Usage](#usage)
  * [Scrape a web page](#scrape-a-web-page)
  * [Scrape a whole website](#scrape-a-whole-website)
    + [`depth` option](#-depth--option)
    + [`selector` option](#-selector--option)
    + [Display the table of contents](#display-the-table-of-contents)
    + [Scrape time](#scrape-time)
- [Installation](#installation)
  * [From source](#from-source)
  * [From binary](#from-binary)
    + [Linux / MacOS](#linux---macos)
    + [Windows](#windows)
  * [MOBI support](#mobi-support)
- [Autocompletion](#autocompletion)
- [Dependencies](#dependencies)

# Usage

## Scrape a web page

The `get` command lets you retrieve the content of any web page.

```
Scrape URL content

Usage:
  papeer get URL [flags]

Examples:
papeer get https://www.eff.org/cyberspace-independence

Flags:
  -a, --author string      book author
      --delay int          time in milliseconds to wait before downloading next chapter, use with depth/selector (default -1)
  -d, --depth int          scraping depth
  -f, --format string      file format [stdout, md, epub, mobi] (default "md")
  -h, --help               help for get
      --images             retrieve images only
  -i, --include            include URL as first chapter, use with depth/selector
  -l, --limit int          limit number of chapters, use with depth/selector (default -1)
  -n, --name string        book name (default: page title)
  -o, --offset int         skip first chapters, use with depth/selector
      --output string      file name (default: book name)
  -q, --quiet              hide progress bar
  -s, --selector strings   table of contents CSS selector
  -t, --threads int        download concurrency, use with depth/selector (default -1)
      --use-link-name      use link name for chapter title
```

## Scrape a whole website

If a navigation menu is present on a website, you can scrape the content of each page.

You can activate this mode by using the `depth` or `selector` options.

### `depth` option

This option defaults to 0, `papeer` will grab only the main page.

If you specify a value greater than 0, `papeer` will grab pages as deep as the value you specify.

> Using `include` option will include all intermediary levels into the book.

### `selector` option

If this option is not specified, `papeer` will grab only the one page.

If this option is specified, `papeer` will select the links (a HTML tag) present on the main page, then grab each one of them.

You can chain this option to grab several level of pages with diferent selectors for each level.

### Display the table of contents

Before actually scraping a whole website, it is a good idea to use the `list` command. This command is like a **dry run**, which lets you vizualize the content before actually retrieving it. You can use several options to customize the table of contents extraction, such as `selector`, `limit`, `offset` and `include`. Type `papeer list --help` for more information about those options.

```sh
papeer list https://12factor.net/ -s 'section.concrete>article>h2>a'
```
```
 #  NAME                    URL                                    
 1  I. Codebase             https://12factor.net/codebase          
 2  II. Dependencies        https://12factor.net/dependencies      
 3  III. Config             https://12factor.net/config            
 4  IV. Backing services    https://12factor.net/backing-services  
 5  V. Build, release, run  https://12factor.net/build-release-run 
 6  VI. Processes           https://12factor.net/processes         
 7  VII. Port binding       https://12factor.net/port-binding      
 8  VIII. Concurrency       https://12factor.net/concurrency       
 9  IX. Disposability       https://12factor.net/disposability     
10  X. Dev/prod parity      https://12factor.net/dev-prod-parity   
11  XI. Logs                https://12factor.net/logs              
12  XII. Admin processes    https://12factor.net/admin-processes
```

### Scrape time

Once you are satisfied with the table of contents listed by the `ls` command, you can actually scrape the content of those pages. You can use the same options that you specified for the `ls` command. You can specify `delay` and `threads` options when using `selector` or `depth` options.

```sh
papeer get https://12factor.net/ --selector='section.concrete>article>h2>a'
```
```
[======================================>-----------------------------] Chapters 7 / 12
[====================================================================] 1. I. Codebase
[====================================================================] 2. II. Dependencies
[====================================================================] 3. III. Config
[====================================================================] 4. IV. Backing services
[====================================================================] 5. V. Build, release, run
[====================================================================] 6. VI. Processes
[====================================================================] 7. VII. Port binding
[--------------------------------------------------------------------] 8. VIII. Concurrency
[--------------------------------------------------------------------] 9. IX. Disposability
[--------------------------------------------------------------------] 10. X. Dev/prod parity
[--------------------------------------------------------------------] 11. XI. Logs
[--------------------------------------------------------------------] 12. XII. Admin processes
Markdown saved to "The_Twelve-Factor_App.md"
```

# Installation

## From source

```sh
go get -u github.com/lapwat/papeer
```

## From binary

### Linux / MacOS

```sh
# use platform=darwin for MacOS
platform=linux
release=0.4.0

# download and extract
curl -L https://github.com/lapwat/papeer/releases/download/v$release/papeer-v$release-$platform-amd64.tar.gz > papeer.tar.gz
tar xzvf papeer.tar.gz
rm papeer.tar.gz

# move to user binaries
sudo mv papeer /usr/local/bin
```

### Windows

Download [latest release](https://github.com/lapwat/papeer/releases/download/v0.4.0/papeer-v0.4.0-windows-amd64.exe.zip).

## MOBI support

Install kindlegen to convert websites, Linux only

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

You can replace `bash` by your own shell (zsh, fish or powershell).

Type `papeer completion bash -h` for more information.

# Dependencies

- `cobra` command line interface
- `go-readability` extract content from HTML
- `html-to-markdown` convert HTML to Markdown
- `go-epub` convert HTML to EPUB
- `colly` query HTML trees
- `uiprogress` display progress bars
