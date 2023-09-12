package utils

import "github.com/hitoshi44/go-uid64"

var generator *uid64.Generator

func InitGenerator() {
	var err error
	generator, err = uid64.NewGenerator(0)
	if err != nil {
		panic(err)
	}
}

func GenerateId() int64 {
	id, err := generator.Gen()
	if err != nil {
		panic(err)
	}
	return id.ToInt()
}
