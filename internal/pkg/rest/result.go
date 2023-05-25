package rest

import (
	"fmt"

	"github.com/changaolee/skeleton/pkg/runtime"
	"github.com/changaolee/skeleton/pkg/third_party/gorequest"
)

// Result 包含调用 Request.Do() 的返回结果.
type Result struct {
	response *gorequest.Response
	err      error
	body     []byte
	decoder  runtime.Decoder
}

// Raw 返回原始结果.
func (r Result) Raw() ([]byte, error) {
	return r.body, r.err
}

// Into 将结果存储到对象 v 中.
func (r Result) Into(v interface{}) error {
	if r.err != nil {
		return r.Error()
	}

	if r.decoder == nil {
		return fmt.Errorf("serializer doesn't exist")
	}

	if err := r.decoder.Decode(r.body, &v); err != nil {
		return err
	}

	return nil
}

// Error 实现了 error 接口.
func (r Result) Error() error {
	return r.err
}
