package config

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/vertex-center/vertex/common"
)

var (
	commit = kingpin.Flag("commit", "Print commit hash and quit.").Bool()
	date   = kingpin.Flag("date", "Print build date and quit.").Bool()
	port   = kingpin.Flag("port", "Port to listen on.").Default("8080").String()
)

func RegisterPort(id string, def string) {
	kingpin.Flag(id+"-port", "Port for "+id+".").Default(def).String()
}

func ParseArgs(about common.About) {
	kingpin.Version(about.Version)
	kingpin.Parse()

	if commit != nil && *commit {
		fmt.Println(about.Commit)
		os.Exit(1)
	}

	if date != nil && *date {
		fmt.Println(about.Date)
		os.Exit(1)
	}

	if port != nil {
		Current.Port = *port
	}
}
