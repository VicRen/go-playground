package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	src  = flag.String("src", "", "Source root")
	repo = flag.String("repo", "", "Repo for protobuf generation, such as galaxy")
	tool = flag.String("tool", "protoc", "protoc tool for gen protobuf, default is protoc")
)

func main() {
	flag.Parse()

	if *src == "" {
		panic("please specify src")
	}

	if *repo == "" {
		panic("please specify repo")
	}

	protofiles := make(map[string][]string)
	reporoot := filepath.Join(*src, *repo)
	protoc := *tool

	if err := filepath.Walk(reporoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		filename := filepath.Base(path)
		if strings.HasSuffix(filename, ".proto") {
			protofiles[dir] = append(protofiles[dir], path)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	fmt.Println("proto files:", protofiles)

	for _, files := range protofiles {
		args := []string{"--proto_path", *src, "--go_out=plugins=grpc:" + *src}
		args = append(args, files...)
		fmt.Println("args", args)
		cmd := exec.Command(protoc, args...)
		cmd.Env = append(cmd.Env, os.Environ()...)
		output, err := cmd.CombinedOutput()
		if len(output) > 0 {
			fmt.Println(string(output))
		}
		if err != nil {
			panic(err)
		}
	}

	if err := filepath.Walk(reporoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".pb.go") {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content = bytes.Replace(content, []byte("\"go-playground"), []byte("\"github.com/VicRen/go-playground"), 1)

		//pos := bytes.Index(content, []byte("\npackage"))
		//if pos > 0 {
		//	content = content[pos+1:]
		//}

		return ioutil.WriteFile(path, content, info.Mode())
	}); err != nil {
		panic(err)
	}
}
