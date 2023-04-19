// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

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
