package httpclient

import (
	"context"
	"net/http"
)

type Chain interface {
	Proceed(ctx context.Context, req *http.Request) (*http.Response, error)
	GetRequest(ctx context.Context) *http.Request
}

type chain struct {
	index   int
	req     *http.Request
	hookFns []HookFn
}

func (c *chain) GetRequest(ctx context.Context) *http.Request {
	return c.req
}

func (c *chain) Proceed(ctx context.Context, req *http.Request) (*http.Response, error) {
	ch := &chain{
		index:   c.index + 1,
		req:     req,
		hookFns: c.hookFns,
	}
	hookFn := c.hookFns[c.index]
	return hookFn(ctx, ch)
}
