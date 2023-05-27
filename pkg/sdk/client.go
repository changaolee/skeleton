package sdk

import (
	"github.com/changaolee/skeleton/internal/pkg/rest"
	apiv1 "github.com/changaolee/skeleton/pkg/sdk/apiserver/v1"
)

type SKTInterface interface {
	APIV1() apiv1.APIV1Interface
}

type SKTClient struct {
	apiV1 *apiv1.APIV1Client
}

var _ SKTInterface = (*SKTClient)(nil)

func (s *SKTClient) APIV1() apiv1.APIV1Interface {
	return s.apiV1
}

func NewForConfig(c *rest.Config) (*SKTClient, error) {
	configShallowCopy := *c

	var skt SKTClient
	var err error

	skt.apiV1, err = apiv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return &skt, nil
}

func NewForConfigOrDie(c *rest.Config) *SKTClient {
	var skt SKTClient
	skt.apiV1 = apiv1.NewForConfigOrDie(c)

	return &skt
}

func New(c rest.Interface) *SKTClient {
	var skt SKTClient
	skt.apiV1 = apiv1.New(c)

	return &skt
}
