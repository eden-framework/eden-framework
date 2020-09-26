package main

import (
	"fmt"
	"github.com/eden-framework/plugins"
	"os"
	"path"
	"reflect"
)

func main() {
	//p, err := plugin.Open("/Users/liyiwen/Documents/golang/src/github.com/eden-framework/plugin-redis/generator.so")
	//if err != nil {
	//	panic(err)
	//}
	//symbol, err := p.Lookup("Plugin")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(reflect.TypeOf(symbol).String())

	cwd, _ := os.Getwd()
	path := path.Join(cwd, "plugins")
	os.MkdirAll(path, 0755)
	ldr := plugins.NewLoader(path)
	p, err := ldr.Load("redis", "https://api.github.com/repos/eden-framework/plugin-redis/zipball/v0.0.11")
	if err != nil {
		panic(err)
	}
	symbol, err := p.Lookup("Plugin")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reflect.TypeOf(symbol).String())
}
