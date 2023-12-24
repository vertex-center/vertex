package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/alecthomas/kingpin/v2"
	"github.com/vertex-center/vertex/common"
)

var (
	commit = kingpin.Flag("commit", "Print commit hash and quit.").Bool()
	date   = kingpin.Flag("date", "Print build date and quit.").Bool()
	host   = kingpin.Flag("host", "Host to listen on.").Envar("VERTEX_HOST").Default("127.0.0.1").String()
	port   = kingpin.Flag("port", "Port to listen on.").Envar("VERTEX_PORT").Default("8080").String()

	mu     sync.RWMutex
	hosts  = map[string]*string{}
	fields = map[string]*string{}
)

// RegisterHost registers a host flag with the given id and default port.
func RegisterHost(id, defaultPort string) {
	idUpper := strings.ToUpper(id)
	idUpper = strings.ReplaceAll(idUpper, "-", "_")

	mu.Lock()
	defer mu.Unlock()
	hosts[id] = kingpin.
		Flag(id+"-addr", "Address for "+id+".").
		Envar("VERTEX_" + idUpper + "_ADDR").
		Default(Current.DefaultApiAddr(defaultPort)).
		String()
	Current.RegisterAPIAddr(id, defaultPort)
}

func (c *Config) RegisterDBArgs() {
	mu.Lock()
	defer mu.Unlock()
	fields["VERTEX_DB_DRIVER"] = kingpin.Flag("db-driver", "Database driver.").Envar("VERTEX_DB_DRIVER").Default("postgres").Enum("postgres", "sqlite")
	fields["VERTEX_DB_HOST"] = kingpin.Flag("db-host", "Database host.").Envar("VERTEX_DB_HOST").Default("localhost").String()
	fields["VERTEX_DB_PORT"] = kingpin.Flag("db-port", "Database port.").Envar("VERTEX_DB_PORT").Default("5432").String()
	fields["VERTEX_DB_USER"] = kingpin.Flag("db-user", "Database user.").Envar("VERTEX_DB_USER").Default("postgres").String()
	fields["VERTEX_DB_PASS"] = kingpin.Flag("db-pass", "Database pass.").Envar("VERTEX_DB_PASS").Default("postgres").String()
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

	if host != nil {
		Current.host = *host
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

	for id, val := range fields {
		if val == nil {
			continue
		}
		Current.fields[id] = *val
	}
}
