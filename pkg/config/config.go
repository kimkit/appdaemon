package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/go-yaml/yaml"
)

func Load(dict interface{}) error {
	defaultConfigFile := "config.yaml"
	file := flag.String("c", defaultConfigFile, "Config file")
	help := flag.Bool("h", false, "Print this message")

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*file); err != nil {
		if os.IsNotExist(err) && *file == defaultConfigFile {
			return nil
		} else {
			return err
		}
	}

	conf, err := ioutil.ReadFile(*file)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(conf, dict); err != nil {
		return err
	}

	dir := path.Dir(*file)
	fp, err := os.Open(dir)
	if err != nil {
		return err
	}

	fs, err := fp.Readdir(0)
	if err != nil {
		return err
	}

	var files []string
	for _, fi := range fs {
		if strings.HasPrefix(fi.Name(), path.Base(*file)+".inc.") {
			files = append(files, path.Join(dir, fi.Name()))
		}
	}

	sort.Strings(files)
	for _, file := range files {
		conf, err := ioutil.ReadFile(file)
		if err != nil {
			// ignore err
			continue
		}
		if err := yaml.Unmarshal(conf, dict); err != nil {
			// ignore err
			continue
		}
	}

	return nil
}
