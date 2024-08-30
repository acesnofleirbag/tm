package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"tm/guard"
)

type Config struct {
	Template  string `json:"template"`
	Aggregate string `json:"aggregate"`
}

var config Config

func startConfig(home string) {
	configFile, err := os.Open(fmt.Sprintf("%v/.config/tmrc", home))
	guard.Err(err)
	defer configFile.Close()

	configData, err := io.ReadAll(configFile)
	guard.Err(err)

	err = json.Unmarshal(configData, &config)
	guard.Err(err)
}

func checkArgsSize(args []string, size int) {
	if len(args) < size {
		fmt.Println("ERR: missing args")
		os.Exit(1)
	}
}

func main() {
	args := os.Args
    checkArgsSize(args, 2)
	cmd := args[1]

	home := os.Getenv("HOME")

	startConfig(home)

	switch cmd {
	case "init":
		checkArgsSize(args, 3)
		path := args[2]

		templatePath := fmt.Sprintf("%v/%v", home, config.Template)

		runner := CopyRunner{Replace: false, Name: ""}
		runner.CpDir(templatePath, path)

		break
	case "aggregate":
		checkArgsSize(args, 3)
		subcmd := args[2]

		switch subcmd {
		case "new":
			checkArgsSize(args, 5)
			aggregate, dst := args[3], args[4]

			aggregatePath := fmt.Sprintf("%v/%v", home, config.Aggregate)

			runner := CopyRunner{Replace: true, Name: aggregate}
			runner.CpDir(aggregatePath, dst)

			break
		}

		break
	default:
		fmt.Println("Usage:")
		fmt.Println()
		fmt.Println("tm init <path>")
		fmt.Println("tm aggregate new <name> <path>")
	}
}
