package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func pr(f string, v ...interface{}) {
	fmt.Printf(f, v...)
}

func loadList(path string) ([]string, error) {
	var list []string
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	l := strings.Split(string(in), "\n")
	for _, s := range l {
		if len(s) > 0 {
			list = append(list, s)
		}
	}
	return list, nil
}

func saveString(s, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			pr("Error closing %s: %s", path, err.Error())
		}
	}()

	_, err = f.WriteString(s)
	return err
}

// mangleName turns a file name into alphanumerics only.
func mangleName(name string) string {
	r, _ := regexp.Compile("[^0-9a-zA-Z]+")
	return r.ReplaceAllString(name, "")
}
