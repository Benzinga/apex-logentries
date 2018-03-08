// +build ignore

package main

import (
	"flag"
	"time"

	logentries "github.com/Benzinga/apex-logentries"
	"github.com/apex/log"
)

var (
	token string
)

func init() {
	// Parse command line.
	flag.StringVar(&token, "token", "", "LogEntries token")
	flag.Parse()
}

func main() {
	// Ensure a token is specified.
	if token == "" {
		panic("You must specify a token with -token.")
	}

	// Create apex-logentries handler.
	le := logentries.New(logentries.Config{
		UseTLS: true,
		Token:  token,
	})

	// Set Apex handler to apex-logentries.
	log.SetHandler(le)

	// Send some test logs.
	log.Debug("Debug Message")
	log.Info("Info Message")
	log.Warn("Warn Message")
	log.Error("Error Message")

	// Wait a few seconds.
	time.Sleep(3 * time.Second)
}
