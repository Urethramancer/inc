package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

type incFile struct {
	Name string
	Path string
	Data []byte
}

// Convert binary data to an embeddable variable.
func Convert(content []byte, name string) ([]byte, error) {
	var b bytes.Buffer

	c := fmt.Sprintf("var %s = []byte(\"", name)
	b.WriteString(c)
	for _, in := range content {
		c = fmt.Sprintf("\\x%02x", in)
		b.WriteString(c)
	}
	b.WriteString("\")\n\n")
	return b.Bytes(), nil
}

// ConvertFiles takes a list of files and runs Convert() on each file.
func ConvertFiles(files []string) (string, error) {
	var list []incFile
	for _, path := range files {
		name := mangleName(path)
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return "", err
		}

		gz, err := Compress(in)
		if err != nil {
			return "", err
		}

		out, err := Convert(gz, name)
		if err != nil {
			return "", err
		}

		file := incFile{path, name, out}
		list = append(list, file)
	}

	var b bytes.Buffer
	header := `package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var basePath = ""

// EmbeddedFileList gets the byte slices from original paths.
type EmbeddedFileList map[string]*[]byte

var embeddedFiles EmbeddedFileList

// SetBasePath sets the path prepended to file paths when checking for actual files.
func SetBasePath(path string) {
	basePath = path
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// GetData decompresses an embedded file.
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

	save := `// SaveData saves the specified embedded file relative to the specified path.
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
	if err != nil {
		return err
	}

	return nil
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

	init := `func init() {
	embeddedFiles = make(EmbeddedFileList)
`

	b.WriteString(header)
	if opts.Save {
		b.WriteString(save)
	}
	b.WriteString(init)
	for _, v := range list {
		s := fmt.Sprintf("\tembeddedFiles[\"%s\"] = &%s\n", v.Name, v.Path)
		b.WriteString(s)
	}
	b.WriteString("}\n\n")

	for _, f := range list {
		_, err := b.Write(f.Data)
		if err != nil {
			return "", err
		}
	}

	return b.String(), nil
}
