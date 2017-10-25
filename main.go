package main

import (
	"sort"

	flags "github.com/jessevdk/go-flags"
)

// Version is filled in from Git tags by the build script.
var Version = "0.0.0"

var opts struct {
	Info    bool   `short:"V" description:"Show program version and exit."`
	Verbose bool   `short:"v" description:"Show progress."`
	List    string `short:"l" long:"list" description:"Name of a text file listing files to embed, one per line."`
	Output  string `short:"o" long:"output" description:"Output file name." default:"embed.go"`
	Save    bool   `short:"s" description:"Include save code in the output."`
	Args    struct {
		Files []string `positional-arg-name:"FILE" description:"File to embed."`
	} `positional-args:"true"`
}

func main() {
	var err error
	_, err = flags.Parse(&opts)
	if err != nil {
		return
	}

	if opts.Info {
		pr("inc %s\n", Version)
		return
	}

	var list []string

	list = append(list, opts.Args.Files...)

	if opts.List != "" {
		var l []string
		l, err = loadList(opts.List)
		if err != nil {
			pr("Error loading list: %s\n", err.Error())
			return
		}
		list = append(list, l...)
	}

	if len(list) == 0 {
		return
	}
	sort.Strings(list)
	var s string
	s, err = ConvertFiles(list)
	if err != nil {
		pr("Error converting list: %s\n", err.Error())
		return
	}

	err = saveString(s, opts.Output)
	if err != nil {
		pr("Error saving output: %s\n", err.Error())
		return
	}
}
