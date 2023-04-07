// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package app

import "github.com/spf13/pflag"

// CliOptions 定义了从命令行读取参数的配置选项接口.
type CliOptions interface {
	Flags() (fss NamedFlagSets) // 通过命令行参数解析出的 FlagSets
	Validate() []error          // 用于校验参数是否合法
}

// SubCliOptions 定义了每个子命令行选项需要实现的接口，为 CliOptions 提供支持.
type SubCliOptions interface {
	AddFlags(fs *pflag.FlagSet) // 向指定 FlagSet 中添加标志
	Validate() []error          // 用于校验参数是否合法
}

// CompletableOptions 定义了完整选项需要实现的接口.
type CompletableOptions interface {
	Complete() error
}

// PrintableOptions 定义了可被打印的选项需要实现的接口.
type PrintableOptions interface {
	String() string
}
