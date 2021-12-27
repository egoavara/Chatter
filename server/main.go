package main

import (
	"chatter-server/chq/jwk"
	"fmt"
)

func main() {
	s, err := jwk.FetchSet("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		panic(err)
	}

	fmt.Println(s.First())
	fmt.Println(s.First().Raw)
}
