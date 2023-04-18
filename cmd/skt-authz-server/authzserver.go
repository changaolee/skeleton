package main

import (
	"math/rand"
	"time"

	"github.com/changaolee/skeleton/internal/authzserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	authzserver.NewApp("skt-authz-server").Run()
}
