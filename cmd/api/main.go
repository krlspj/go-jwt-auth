package main

import (
	"log"

	"github.com/krlspj/go-jwt-auth/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
