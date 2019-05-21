package main

import (
	"bytes"
	"flag"
	"github.com/sohoffice/cloudwalk"
	"io/ioutil"
	"os"
)

func main() {
	configFilePtr := flag.String("config-file", "", "Specify subsitution in a config file. Key is cur value, Value is the new value.")
	configPtr := flag.String("c", "", "Specify substitution in the format of URL query string. Key is cur value, Value is the new value.")
	modePtr := flag.String("mode", "direct", "The mode to match the key. Valid values are: direct (default), db (double braces).")
	filenamePtr := flag.String("f", "", "The generated API filename.")
	outputPtr := flag.String("o", "", "The file should be generated in the output file.")
	flag.Parse()

	if *configPtr == "" && *configFilePtr == "" {
		println("Please supply substitution with: -c 'sub string' or -config-file 'config file' (or both).")
		os.Exit(1)
	}
	if *filenamePtr == "" {
		println("Please supply input filename: -f filename")
		os.Exit(1)
	}
	if *outputPtr == "" {
		println("Please supply output filename: -o filename")
		os.Exit(1)
	}

	mode := parseMode(*modePtr)
	if mode == -1 {
		println("Please specify a valid mode: direct, db.")
		os.Exit(1)
	}

	var substitutes []cloudwalk.KeyValue
	if *configFilePtr != "" {
		substitutes = append(substitutes, cloudwalk.ReadJavaPropertiesLike(*configFilePtr)...)
	}
	if *configPtr != "" {
		substitutes = append(substitutes, cloudwalk.ReadQueryString(*configPtr)...)
	}

	if mode == DoubleBraces {
		for i := range substitutes {
			v := substitutes[i]
			substitutes[i] = cloudwalk.KeyValue{
				Key: "{{" + v.Key + "}}", Val: v.Val,
			}
		}
	}
	var processors []BytesProcessor
	for _, subst := range substitutes {
		processors = append(processors, newStringReplaceProcessor(subst.Key, subst.Val))
	}

	dataPtr := readFile(*filenamePtr)
	data := *dataPtr
	for _, proc := range processors {
		data = proc(data)
	}

	err := ioutil.WriteFile(*outputPtr, data, 0644)
	if err != nil {
		panic(err)
	}
}

func readFile(filename string) *[]byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return &data
}

type BytesProcessor func([]byte) []byte

func newStringReplaceProcessor(old string, new string) BytesProcessor {
	oldb := []byte(old)
	newb := []byte(new)

	f := func(s []byte) []byte {
		return bytes.Replace(s, oldb, newb, -1)
	}
	return f
}

type Mode int

const (
	Direct Mode = iota
	DoubleBraces
)

type ModeInfo struct {
	id   Mode
	name string
	abbr string
}

var allModes = []ModeInfo{
	{
		id:   Direct,
		name: "Direct",
		abbr: "direct",
	},
	{
		id:   DoubleBraces,
		name: "DoubleBraces",
		abbr: "db",
	},
}

func (m Mode) toString() string {
	return allModes[m].name
}

// Parse string into mode.
// Return the mode (int), or -1 if not found.
func parseMode(s string) Mode {
	for i, v := range allModes {
		if v.name == s || v.abbr == s {
			return Mode(i)
		}
	}
	return Mode(-1)
}
