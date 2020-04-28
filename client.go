package httpclient

import (
	"context"
	"net/http"
)

type Client interface {
	Get(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Head(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Post(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Put(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Patch(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Delete(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Connect(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Options(ctx context.Context, url string, opts ...RequestOption) (int, error)
	Trace(ctx context.Context, url string, opts ...RequestOption) (int, error)
}

type RequestHookFn func(context.Context, *http.Request) error
type ResponseHookFn func(context.Context, *http.Response) error

type HookFn func(ctx context.Context, reqChain Chain) (*http.Response, error)

type client struct {
	sender  Sender
	header  http.Header
	hookFns []HookFn
}

func NewClient(opts ...ClientOption) Client {
	opt := getClientOption(opts...)

	return &client{
		sender:  opt.sender,
		header:  opt.header,
		hookFns: opt.hookFns,
	}
}

func (c *client) Get(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodGet, opts...)
}

func (c *client) Head(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodHead, opts...)
}

func (c *client) Post(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodPost, opts...)
}

func (c *client) Put(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodPut, opts...)
}

func (c *client) Patch(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodPatch, opts...)
}

func (c *client) Delete(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodDelete, opts...)
}

func (c *client) Connect(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodConnect, opts...)
}

func (c *client) Options(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodOptions, opts...)
}

func (c *client) Trace(ctx context.Context, url string, opts ...RequestOption) (int, error) {
	return c.execute(ctx, url, http.MethodTrace, opts...)
}

func (c *client) execute(ctx context.Context, url, method string,
	opts ...RequestOption) (int, error) {
	// get requestOption
	opt := getRequestOption(opts...)

	// convert to http.Request
	req, err := c.buildHttpRequest(ctx, url, method, opt)
	if err != nil {
		return 0, err
	}

	// send request
	res, err := c.proceed(ctx, req, opt)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	return res.StatusCode, nil
}

func (c *client) proceed(ctx context.Context, req *http.Request,
	opt requestOption) (*http.Response, error) {
	hookFns := append(opt.hookFns, sendingRequestHook(c.sender))
	reqChain := &chain{
		index:   0,
		req:     req,
		hookFns: append(c.hookFns, hookFns...),
	}

	return reqChain.Proceed(ctx, req)
}

func sendingRequestHook(sender Sender) HookFn {
	return func(ctx context.Context, reqChain Chain) (*http.Response, error) {
		req := reqChain.GetRequest(ctx)
		return sender.Send(ctx, req)
	}
}
