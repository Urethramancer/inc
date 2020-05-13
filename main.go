package main

import (
	"bytes"
	"fmt"
	"os"
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
	Files   []string `placeholder:"PATH" help:"File or directory to embed."`
}

func main() {
	a := opt.Parse(&opts)
	if opts.Help || opts.Files == nil {
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
	convlist, err := ConvertFiles(list)
	if err != nil {
		pr("Error converting list: %s\n", err.Error())
		return
	}

	var b bytes.Buffer
	b.WriteString(header)
	if opts.Save {
		b.WriteString(save)
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

const (
	header = `package main

import (
	"bytes"
	"compress/gzip"
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

// GetData decompresses an embedded file.
// If a physical file exists at basePath+path, load that instead of the embedded file.
func GetData(path string) ([]byte, error) {
	p := filepath.Join(basePath, path)
	if Exists(p) {
		return ioutil.ReadFile(p)
	}

	in := embeddedFiles[path]
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

	save = `// SaveData saves the specified embedded file relative to the specified path.
func SaveData(path string, data *[]byte) error {
	var err error
	base := filepath.Dir(path)
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

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	var out bytes.Buffer
	_, err = io.Copy(&out, gz)
	gzerr := gz.Close()
	if err != nil {
		return err
	}

	if gzerr != nil {
		return gzerr
	}

	_, err = f.Write(out.Bytes())
	return err
}

// SaveAllData saves all embedded data, relative to the specified path.
func SaveAllData(dest string) error {
	for path, data := range embeddedFiles {
		out := filepath.Join(dest, path)
		err := SaveData(out, data)
		if err != nil {
			return err
		}
	}
	return nil
}

`

	initfunc = `func init() {
	embeddedFiles = make(EmbeddedFileList)
`
)
