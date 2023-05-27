package v1

import (
	"github.com/changaolee/skeleton/internal/pkg/rest"
	"github.com/changaolee/skeleton/pkg/runtime"
)

type APIV1Interface interface {
	RESTClient() rest.Interface
	UsersGetter
}

type APIV1Client struct {
	restClient rest.Interface
}

func (c *APIV1Client) Users() UserInterface {
	return newUsers(c)
}

func NewForConfig(c *rest.Config) (*APIV1Client, error) {
	config := *c
	setConfigDefaults(&config)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &APIV1Client{client}, nil
}

func NewForConfigOrDie(c *rest.Config) *APIV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}

	return client
}

func New(c rest.Interface) *APIV1Client {
	return &APIV1Client{c}
}

func setConfigDefaults(config *rest.Config) {
	gv := SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = ""
	config.Negotiator = runtime.NewSimpleClientNegotiator()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultUserAgent()
	}
}

func (c *APIV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}

	return c.restClient
}
