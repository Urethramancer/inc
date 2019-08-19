package main

import (
	"sort"

	"github.com/Urethramancer/signor/opt"
)

// version is filled in from Git tags by build scripts.
var version = "undefined"

var opts struct {
	opt.DefaultHelp
	Version bool     `short:"V" help:"Show program version and exit."`
	List    string   `short:"l" long:"list" help:"Name of a text file listing files to embed, one per line."`
	Output  string   `short:"o" long:"output" help:"Output file name." default:"embed.go"`
	Save    bool     `short:"s" help:"Include save code in the output."`
	Files   []string `placeholder:"FILE" help:"File to embed."`
}

func main() {
	a := opt.Parse(&opts)
	if opts.Help {
		a.Usage()
		return
	}

	if opts.Version {
		pr("inc %s\n", version)
		return
	}

	var list []string
	var err error
	list = append(list, opts.Files...)

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
