# inc [![Build Status](https://travis-ci.org/Urethramancer/inc.svg)](https://travis-ci.org/Urethramancer/inc)

Embed binaries in Go programs the way *I* like it.

## Why

Other solutions weren't to my liking, and I've been using embedded HTML/CSS/JS in a special way. I wanted to include a default option, and export those as templates for customisation by the user, so I found myself writing saving code a lot. This little utility embeds all that data and can optionally include a function to save it to a configurable path.

## How to embed

Run it with any number of files and/or directories as arguments:

```sh
inc one.html two.css three.js tpl/
```

Or include save code:

```sh
inc -s one.html two.css three.js
```

Or make a list of files to include, perhaps generated from another pre-processor:
```sh
inc -l files.txt
```

Save the resulting file to something other than `embed.go`:
```sh
inc -l files.txt -o files.go
```

## Other options
Get the name and version of the program:

```sh
$ inc -V
inc v0.4.2
```

## Code
Compile the generated `embed.go` into your program and set the base path if you want to load physical files from a particular location:

```go
SetBasePath("/var/www/html")
```

Decompress embedded data like this:
```go
data, err := GetData("one.html")
```

If a file exists with the same subpath in the configured base path, it will be loaded instead of the embedded version.

Individual files from the embedded data can be saved to disk:
```go
err := SaveData("one.html")
```

To save everything into the configured base path, simply do this:

```go
err := SaveAllData()
```
