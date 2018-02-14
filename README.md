# inc [![Build Status](https://travis-ci.org/Urethramancer/inc.svg)](https://travis-ci.org/Urethramancer/inc)

Embed binaries in Go programs the way *I* like it.

## Why

Other solutions weren't to my liking, and I've been using embedded HTML/CSS/JS in a special way. I wanted to include a default option, and export those as templates for customisation by the user, so I found myself writing saving code a lot. This little utility embeds all that data and can optionally include a function to save it to a configurable path.

## How

Run it with any number of files as arguments:

```go
inc one.html two.css three.js
```

Or include save code:

```go
inc -s one.html two.css three.js
```

Or make a list of files to include, perhaps generated from another pre-processor:
```go
inc -l files.txt
```

Sav the resulting file to something other than `embed.go`:
```go
inc -l files.txt -o files.go
```

Get the name and version of the program:

```sh
$ inc -V
inc v0.2.4
```
