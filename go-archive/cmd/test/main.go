package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	goyaml "gopkg.in/yaml.v2"
)

type Config struct {
	Point *Point `yaml:"point"`
}

type Point struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

func main() {
	sampleYAML := `
---
point:
  x: 5
  y: 6
`

	sampleYAMLBytes := []byte(sampleYAML)
	dec := goyaml.NewDecoder(bytes.NewReader(sampleYAMLBytes))

	var config Config
	err := dec.Decode(&config)
	if err == io.EOF {
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config.Point)

	// configBytes, err := goyaml.Marshal(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(configBytes))
}
