package cloudwalk

import (
	"bufio"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const lineFeed = byte('\n')

var regexpJavaProperties = regexp.MustCompile("^([^=]+)=(.*)$")

type KeyValue struct {
	Key string
	Val string
}

// Read a java properties file as KeyValue array.
// Be reminded the method can not handle properties that span over multiple line.
func ReadJavaPropertiesLike(filename string) []KeyValue {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	var data []KeyValue
	for s, err := reader.ReadString(lineFeed); err == nil; s, err = reader.ReadString(lineFeed) {
		s = strings.TrimRight(s, "\n")
		if regexpJavaProperties.MatchString(s) {
			found := regexpJavaProperties.FindStringSubmatch(s)
			data = append(data, KeyValue{
				Key: found[1], Val: found[2],
			})
		}
	}
	return data
}

// Parse query string into KeyValue array.
func ReadQueryString(s string) []KeyValue {
	data, err := url.ParseQuery(s)
	if err != nil {
		panic(err)
	}
	var all []KeyValue
	for k, v := range data {
		all = append(all, KeyValue{
			Key: k, Val: v[0],
		})
	}
	return all
}
