// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package main

import (
	"os"

	"github.com/changaolee/skeleton/internal/sktctl/cmd"
)

func main() {
	command := cmd.NewDefaultSKTCtlCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
