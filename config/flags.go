package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/alecthomas/kingpin/v2"
	"github.com/vertex-center/vertex/common"
)

var (
	commit = kingpin.Flag("commit", "Print commit hash and quit.").Bool()
	date   = kingpin.Flag("date", "Print build date and quit.").Bool()
	port   = kingpin.Flag("port", "Port to listen on.").Default("8080").String()

	mu    sync.RWMutex
	hosts = map[string]*string{}
)

// RegisterHost registers a host flag with the given id and default value.
func RegisterHost(id, def string) {
	mu.Lock()
	defer mu.Unlock()
	hosts[id] = kingpin.
		Flag(id+"-addr", "Address for "+id+".").
		Envar("VERTEX_" + id + "_ADDR").
		Default(def).
		String()
	Current.RegisterAPIAddr(id, def)
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

	for id, val := range hosts {
		if val == nil {
			continue
		}
		Current.SetAPIAddr(id, *val)
	}
}
