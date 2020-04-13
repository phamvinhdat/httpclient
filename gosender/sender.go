package gosender

import (
	"context"
	"net/http"
	"time"
)

const defaultTimeOut = 30 * time.Second

type goClient struct {
	c http.Client
}

func New(opts ...Option) *goClient {
	optsArg := getOptionsArg(opts)
	return &goClient{
		http.Client{
			Transport: optsArg.transport,
			Timeout:   optsArg.timeout,
		},
	}
}

func (g *goClient) Send(ctx context.Context,
	rq *http.Request) (*http.Response, error) {
	return g.c.Do(rq)
}

func getOptionsArg(opts []Option) options {
	// Init default options arg
	optsArgs := options{
		timeout:   defaultTimeOut,
		transport: http.DefaultTransport,
	}

	for _, opt := range opts {
		opt.apply(&optsArgs)
	}
	return optsArgs
}
