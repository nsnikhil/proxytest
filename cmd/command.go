package main

import (
	"log"
	"proxytest/pkg/app"
)

const (
	httpServeCommand = "http-serve"
)

func commands() map[string]func(configFile string) {
	return map[string]func(configFile string){
		httpServeCommand: app.StartHTTPServer,
	}
}

func execute(cmd string, configFile string) {
	run, ok := commands()[cmd]
	if !ok {
		log.Fatal("invalid command")
	}

	run(configFile)
}
