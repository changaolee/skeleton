// Copyright 2022 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package main

import (
	"os"

	//_ "go.uber.org/automaxprocs".

	"github.com/changaolee/skeleton/internal/skeleton"
)

func main() {
	command := skeleton.NewSkeletonCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
