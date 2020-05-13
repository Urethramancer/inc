package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
func ConvertFiles(filelist []string) ([]*incFile, error) {
	var list []*incFile
	for _, path := range filelist {
		fi, err := os.Stat(path)
		if err != nil {
			pr("Error getting info about '%s': %s", path, err.Error())
			continue
		}

		if fi.IsDir() {
			filist, err := ioutil.ReadDir(path)
			if err == nil {
				l := []string{}
				for _, x := range filist {
					fn := filepath.Join(path, x.Name())
					l = append(l, fn)
				}
				convlist, err := ConvertFiles(l)
				if err != nil {
					pr("Couldn't convert files in directory '%s': %s", path, err.Error())
					continue
				}

				list = append(list, convlist...)
			}
			continue
		}

		name := mangleName(path)
		in, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		gz, err := Compress(in)
		if err != nil {
			return nil, err
		}

		out, err := Convert(gz, name)
		if err != nil {
			return nil, err
		}

		file := incFile{path, name, out}
		list = append(list, &file)
	}

	return list, nil
}
