package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gookit/properties"
)

var (
	_, b, _, _ = runtime.Caller(0)
	RootPath   = filepath.Dir(b)
)

var PMan *properties.Parser

func NewPman(fileName string) *properties.Parser {
	pman := properties.NewParser(properties.ParseEnv, properties.ParseInlineSlice)
	fileData, err := os.ReadFile(fmt.Sprintf("%s/../../resources/properties/%s", RootPath, fileName))
	if err != nil {
		panic("Properties file does not found")
	}
	err = pman.Parse(string(fileData))
	if err != nil {
		panic(err)
	}
	return pman
}
