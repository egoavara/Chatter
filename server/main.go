package main

import (
	"chatter-server/chq"
	"log"
	"os"
)

func main() {
	cfg, err := chq.MapReadCloser(os.Open("./privates/config.json"))
	if err != nil {
		log.Fatal(err)
	}
	chq, err := cfg.ChatterQ()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = chq.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	if err := chq.Run(); err != nil {
		log.Fatal(err)
	}
}
