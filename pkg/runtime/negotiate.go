package runtime

import (
	"encoding/json"
	"fmt"
)

type NegotiateError struct {
	ContentType string
	Stream      bool
}

func (e NegotiateError) Error() string {
	if e.Stream {
		return fmt.Sprintf("no stream serializers registered for %s", e.ContentType)
	}
	return fmt.Sprintf("no serializers registered for %s", e.ContentType)
}

type apimachineryClientNegotiator struct{}

var _ ClientNegotiator = &apimachineryClientNegotiator{}

func (n *apimachineryClientNegotiator) Encoder() (Encoder, error) {
	return &apimachineryClientNegotiatorSerializer{}, nil
}

func (n *apimachineryClientNegotiator) Decoder() (Decoder, error) {
	return &apimachineryClientNegotiatorSerializer{}, nil
}

type apimachineryClientNegotiatorSerializer struct{}

var _ Decoder = &apimachineryClientNegotiatorSerializer{}

func (s *apimachineryClientNegotiatorSerializer) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (s *apimachineryClientNegotiatorSerializer) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func NewSimpleClientNegotiator() ClientNegotiator {
	return &apimachineryClientNegotiator{}
}
