// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package runtime

// Encoder 将对象序列化.
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

// Decoder 尝试从 data 加载对象.
type Decoder interface {
	Decode(data []byte, v interface{}) error
}

type ClientNegotiator interface {
	Encoder() (Encoder, error)
	Decoder() (Decoder, error)
}
