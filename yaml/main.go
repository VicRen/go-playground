package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type TestYAML struct {
	MySql MySql
	Cache Cache
}

type MySql struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
}

type Cache struct {
	Enable bool     `yaml:"enable"`
	List   []string `yaml:"list,flow"`
}

func main() {
	var conf *TestYAML
	f, err := ioutil.ReadFile("./yaml/test.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("--- t file:\n%s\n\n", string(f))

	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("conf: %v\n\n", conf)

	d, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
	err = ioutil.WriteFile("./yaml/test_dump.yaml", d, 0644)
	if err != nil {
		fmt.Println(err)
	}

	s, err := readFile()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}

// readFile returns the file name, or error when error occur.
func readFile() (s string, err error) {
	file, cErr := os.Open("./yaml/test.yaml")
	if cErr != nil {
		return "", cErr
	}
	defer func() {
		if cErr := file.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()
	return file.Name(), nil
}
