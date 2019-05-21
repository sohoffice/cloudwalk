package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const lineFeed = byte('\n')

func main() {
	cmdPtr := flag.String("cmd", "heroku", "Heroku command to execute, default to heroku.")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		fmt.Println("heroku-gcp-migrate-tool <heroku app name>")
		os.Exit(1)
	}
	env := flag.Arg(0)

	buf := readConfig(*cmdPtr, env)
	data := parseConfig(&buf)

	output(env, "Secret declaration, save the below and run with `kubectl apply -f ./secret.yaml`",
		fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: secret-%s
type: Opaque
stringData:`, env),
		data, genSecret)

	output(env,
		"Add the below to your k8s environment yml", "",
		data, genEnv)
}

func readConfig(executable string, env string) bytes.Buffer {
	log.Println(fmt.Sprintf("Running `%s config -a %s`", executable, env))
	cmd := exec.Command(executable, "config", "-a", env)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Heroku configs loaded.")
	return out
}

func parseConfig(buffer *bytes.Buffer) []KeyValue {
	var data []KeyValue
	for str, err := buffer.ReadString(lineFeed); err == nil; str, err = buffer.ReadString(lineFeed) {
		pos := strings.Index(str, ":")
		if pos != -1 {
			kv := KeyValue{
				key: str[:pos],
				val: strings.TrimSpace(str[pos+1:]),
			}
			data = append(data, kv)
		}
	}
	return data
}

// Print pretty output of GCP yml
func output(env string, title string, headerText string, data []KeyValue, processor KeyValueConvert) {
	printSquare(title, "BEGIN")
	fmt.Println(headerText)
	for _, s := range traverse(env, data, processor) {
		fmt.Println(s)
	}
	printSquare("END", "")
	fmt.Print("\n\n")
}

// Print the msg and surround it with a pretty square, optionally add a decor on the top.
func printSquare(msg string, decor string) {
	if len(decor) > 0 {
		decor = fmt.Sprintf(" %s ", decor)
	}
	var length int
	if len(msg) > len(decor) {
		length = len(msg) + 2
	} else {
		length = len(decor)
	}

	fmt.Printf("┌%s%s┐\n", decor, strings.Repeat("─", length-len(decor)))
	fmt.Printf("│ %s %s│\n", msg, strings.Repeat(" ", length-len(msg)-2))
	fmt.Printf("└%s┘\n", strings.Repeat("─", length))
}

// Traverse through data and process each of them with the processor.
func traverse(env string, data []KeyValue, processor KeyValueConvert) []string {
	var out []string
	for _, ent := range data {
		s := processor(env, ent)
		out = append(out, s)
	}
	return out
}

func genEnv(env string, data KeyValue) string {
	return fmt.Sprintf(`      - name: %s
        valueFrom:
          secretKeyRef:
            name: secret-%s
            key: %s`, data.key, env, data.key)
}

func genSecret(_ string, data KeyValue) string {
	return fmt.Sprintf(`  %s: "%s"`, data.key, data.val)
}

type KeyValue struct {
	key string
	val string
}

type KeyValueConvert func(string, KeyValue) string
