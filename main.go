package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	path := os.Args[1]
	suffix := os.Args[2]
	wordAfter := os.Args[3]

	files, err := fetchFiles(path, suffix)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	for _, file := range files {
		if err := os.Chmod(file, 0755); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		data, err := ioutil.ReadFile(file)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		res, err := regexp.Compile(`\b` + wordAfter + `\s+\S+`)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}

		name := strings.TrimSpace(strings.TrimPrefix(string(res.Find(data)), wordAfter))

		if name != "" {
			newPath := filepath.Dir(file)

			if err := os.Rename(file, newPath+"/"+name+"."+suffix); err != nil {
				os.Stderr.WriteString(err.Error())
				os.Exit(1)
			}
		}
	}

}

func fetchFiles(path, suffix string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			files = append(files, file)
		}

		return nil
	})

	return files, err
}
