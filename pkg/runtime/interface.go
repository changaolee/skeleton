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
