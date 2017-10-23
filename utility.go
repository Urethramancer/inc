package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

func mangleName(name string) string {
	s := strings.Replace(name, ".", "", -1)
	s = strings.Replace(s, "_", "", -1)
	s = strings.Replace(s, "/", "", -1)
	return s
}
