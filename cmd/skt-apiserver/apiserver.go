package main

import (
	"math/rand"
	"time"

	"github.com/changaolee/skeleton/internal/apiserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	apiserver.NewApp("skt-apiserver").Run()
}