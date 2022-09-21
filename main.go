package main

import (
	"fmt"
	"os"

	"github.com/jpetrucciani/poglets/cmd"
	"github.com/sirupsen/logrus"
)

var (
	Version   string
	GitCommit string
)

func main() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	if err := cmd.Execute(Version, GitCommit); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
