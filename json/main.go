package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

var conf = `log:
  time-format: 2006-01-02T15:04:05.999Z07:00
  mode: 3
  path: /tmp/galaxy/log
  file-prefix: client_
  reserve-hours: 1
  context-skip: 5`

func main() {
	c := new(Conf)
	if err := yaml.Unmarshal([]byte(conf), c); err != nil {
		panic(err)
	}

	for k, v := range c.Log {
		fmt.Println(k, ":", v, reflect.TypeOf(v))
	}

	fmt.Println("yml:", c)

	j, err := json.Marshal(c.Log)
	if err != nil {
		panic(err)
	}

	fmt.Println("json:", j)

	js := new(logConf)
	if err := json.Unmarshal(j, js); err != nil {
		panic(err)
	}

	for k, v := range *js {
		fmt.Println(k, ":", v, reflect.TypeOf(v))
	}
}

type Conf struct {
	Log logConf `yml:"log"`
}

type logConf map[string]interface{}
