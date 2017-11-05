package main

import (
	"flag"

	"github.com/fallenhitokiri/hibpnotify"
)

func main() {
	init := flag.Bool("init", false, "initialise a new, empty config file")
	cfg := flag.String("config", "./hibpnotify.json", "path to config file")
	flag.Parse()

	if *init {
		hibpnotify.InitConfig(*cfg)
		return
	}

	hibpn, err := hibpnotify.New(*cfg)

	if err != nil {
		panic(err)
	}

	hibpn.Run()
}
