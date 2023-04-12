// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package main

import (
	"fmt"

	"github.com/changaolee/skeleton/internal/pkg/code"
	"github.com/changaolee/skeleton/pkg/errors"
)

func main() {
	if err := bindUser(); err != nil {
		// %s: 打印对外部用户的错误字符串.
		fmt.Printf("====================> %%s <====================\n")
		fmt.Printf("%s\n\n", err)

		// %v: 等同于 %s.
		fmt.Printf("====================> %%v <====================\n")
		fmt.Printf("%v\n\n", err)

		// %-v: 打印调用详情，用于错误定位.
		fmt.Printf("====================> %%-v <====================\n")
		fmt.Printf("%-v\n\n", err)

		// %+v: 打印完整的错误堆栈信息，用于调试.
		fmt.Printf("====================> %%+v <====================\n")
		fmt.Printf("%+v\n\n", err)

		// %#-v: 以 JSON 格式打印调用详情.
		fmt.Printf("====================> %%#-v <====================\n")
		fmt.Printf("%#-v\n\n", err)

		// %#+v: 以 JSON 格式打印完整的错误堆栈信息.
		fmt.Printf("====================> %%#+v <====================\n")
		fmt.Printf("%#+v\n\n", err)

		// 业务方进行错误码判定.
		if errors.IsCode(err, code.ErrEncodingFailed) {
			fmt.Println("this is a ErrEncodingFailed error")
		}
		if errors.IsCode(err, code.ErrDatabase) {
			fmt.Println("this is a ErrDatabase error")
		}

		// 打印错误根因.
		fmt.Println(errors.Cause(err))
	}
}

func bindUser() error {
	if err := getUser(); err != nil {
		// 第三步：用指定 message 和错误码对 err 进行包装.
		return errors.WrapC(err, code.ErrEncodingFailed, "encoding user 'xxx' failed.")
	}

	return nil
}

func getUser() error {
	if err := queryDatabase(); err != nil {
		// 第二步：用指定 message 对 err 进行包装.
		return errors.Wrap(err, "get user failed.")
	}

	return nil
}

func queryDatabase() error {
	// 第一步：创建一个指定错误码的 Error.
	return errors.WithCode(code.ErrDatabase, "user 'xxx' not found.")
}
