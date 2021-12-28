package main

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/lestrrat-go/jwx/jwk"
)

func main() {
	rsamp, err := rsa.GenerateMultiPrimeKey(rand.Reader, 3, 2048)
	if err != nil {
		panic(err)
	}
	jwk.New()
}
