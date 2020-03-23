package mediary

import (
	"net/http"
)

// Handler is just an alias to http.RoundTriper.RoundTrip function
type Handler func(*http.Request) (*http.Response, error)

// Interceptor is a user defined function that can alter a request before it's sent
// and/or alter a response before it's returned to the caller
type Interceptor func(*http.Request, Handler) (*http.Response, error)

type options struct {
	predefinedClient *http.Client
	interceptors     []Interceptor
}

// Builder is a builder interface to help guide and create a mediary instance
type Builder interface {
	AddInterceptors(...Interceptor) Builder
	WithPreconfiguredClient(*http.Client) Builder
	Build() *http.Client
}

type builderImpl struct {
	apply func(*options)
	next  *builderImpl
}

// Init initializes mediary instance configuration
func Init() Builder {
	return new(builderImpl)
}

func (impl *builderImpl) AddInterceptors(interceptors ...Interceptor) Builder {
	return &builderImpl{
		apply: func(opt *options) {
			opt.interceptors = append(opt.interceptors, interceptors...)
		},
		next: impl,
	}
}

func (impl *builderImpl) WithPreconfiguredClient(client *http.Client) Builder {
	return &builderImpl{
		apply: func(opt *options) {
			opt.predefinedClient = client
		},
		next: impl,
	}
}

func (impl *builderImpl) Build() *http.Client {
	var client = &http.Client{}

	if impl != nil {
		opts := new(options)
		iterator := impl
		for iterator.next != nil {
			iterator.apply(opts)
			iterator = iterator.next
		}
		if opts.predefinedClient != nil {
			client = opts.predefinedClient
		}
		if client.Transport == nil {
			client.Transport = http.DefaultTransport
		}
		client.Transport = prepareCustomRoundTripper(client.Transport, opts.interceptors...)
	}

	return client
}
