package main

import (
	"io/ioutil"
	"os"

	"github.com/Joker/jade"
)

func main() {
	str, err := jade.ParseFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(os.Args[2], []byte(str), 0644)
}
