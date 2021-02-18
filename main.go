package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Urethramancer/signor/opt"
)

// version is filled in from Git tags by build scripts.
var version = "undefined"

var opts struct {
	opt.DefaultHelp
	Version bool     `short:"V" help:"Show program version and exit."`
	List    string   `short:"l" long:"list" help:"Name of a text file listing files to embed, one per line."`
	Output  string   `short:"o" long:"output" help:"Output file name." default:"embed.go"`
	Save    bool     `short:"s" long:"save" help:"Include save code in the output."`
	Brotli  bool     `short:"b" long:"brotli" help:"Use Brotli compression instead of gzip."`
	Files   []string `placeholder:"PATH" help:"File or directory to embed."`
}

func main() {
	a := opt.Parse(&opts)
	if opts.Help || (opts.Files == nil && opts.List == "") {
		if opts.Version {
			pr("inc %s\n", version)
			return
		}

		a.Usage()
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
	convlist, err := ConvertFiles(list)
	if err != nil {
		pr("Error converting list: %s\n", err.Error())
		return
	}

	var b bytes.Buffer
	if opts.Brotli {
		header = strings.Replace(header, "compress/gzip", "github.com/andybalholm/brotli", 1)
	}

	if opts.Save {
		b.WriteString(header)
		b.WriteString(save)
	} else {
		header = strings.Replace(header, fmtheader, "", 1)
		b.WriteString(header)
	}

	if opts.Brotli {
		b.WriteString(brotlidec)
	} else {
		b.WriteString(gzipdec)
	}

	b.WriteString(initfunc)
	for _, v := range convlist {
		s := fmt.Sprintf("\tembeddedFiles[\"%s\"] = &%s\n", v.Name, v.Path)
		b.WriteString(s)
	}
	b.WriteString("}\n\n")

	for _, f := range convlist {
		_, err := b.Write(f.Data)
		if err != nil {
			pr("Couldn't write to buffer: %s", err.Error())
			os.Exit(2)
		}
	}

	err = saveString(b.String(), opts.Output)
	if err != nil {
		pr("Error saving output: %s\n", err.Error())
		return
	}
}

var (
	fmtheader = `	"fmt"
`
	header = `package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var basePath = ""

// EmbeddedFileList holds the byte slices from original paths.
type EmbeddedFileList map[string]*[]byte

var embeddedFiles EmbeddedFileList

// SetBasePath sets the path prepended to file paths when checking for non-embedded files.
func SetBasePath(path string) {
	basePath = path
}

// Exists checks if a file or directory exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
`
	gzipdec = `// GetData decompresses an embedded gzip file.
// If a physical file exists at basePath+path, load that instead of the embedded file.
func GetData(fn string) ([]byte, error) {
	p := filepath.Join(basePath, fn)
	if Exists(p) {
		return ioutil.ReadFile(p)
	}

	in, ok := embeddedFiles[fn]
	if !ok {
		return nil, os.ErrNotExist
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(*in))
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	_, err = io.Copy(&out, gz)
	err2 := gz.Close()
	if err != nil {
		return nil, err
	}

	if err2 != nil {
		return nil, err2
	}

	return out.Bytes(), nil
}

	`

	brotlidec = `// GetData decompresses an embedded gzip file.
// If a physical file exists at basePath+path, load that instead of the embedded file.
func GetData(fn string) ([]byte, error) {
	p := filepath.Join(basePath, fn)
	if Exists(p) {
		return ioutil.ReadFile(p)
	}
	
	in, ok := embeddedFiles[fn]
	if !ok {
		return nil, os.ErrNotExist
	}
	
	br := brotli.NewReader(bytes.NewBuffer(*in))
	if br == nil {
		return nil, io.ErrUnexpectedEOF
	}
	
	var out bytes.Buffer
	_, err := io.Copy(&out, br)
	if err != nil {
		return nil, err
	}
	
	return out.Bytes(), nil
}
	
	`

	save = `// SaveData saves the specified embedded file relative to the specified path.
func SaveData(fn string) error {
	data, ok := embeddedFiles[fn]
	if !ok {
		return fmt.Errorf("unknown embedded file %s", fn)
	}

	var err error
	base := filepath.Dir(fn)
	if !Exists(base) {
		err = os.MkdirAll(base, 0755)
		if err != nil {
			return err
		}
	}

	gz, err := gzip.NewReader(bytes.NewBuffer(*data))
	if err != nil {
		return err
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = io.Copy(f, gz)
	gzerr := gz.Close()
	if err != nil {
		return err
	}

	return gzerr
}

// SaveAllData saves all embedded data, relative to the configured base path.
func SaveAllData() error {
	for fn := range embeddedFiles {
		out := filepath.Join(basePath, fn)
		err := SaveData(out)
		if err != nil {
			return err
		}
	}
	return nil
}

`

	initfunc = `// init builds the list of embedded files.
func init() {
	embeddedFiles = make(EmbeddedFileList)
`
)
