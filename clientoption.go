package httpclient

import (
	"net/http"

	"github.com/phamvinhdat/httpclient/gosender"
)

type clientOption struct {
	sender  Sender
	header  http.Header
	hookFns []HookFn
}

type ClientOption interface {
	apply(*clientOption)
}

type clientOptFunc func(*clientOption)

func (f clientOptFunc) apply(c *clientOption) {
	f(c)
}

func getClientOption(opts ...ClientOption) clientOption {
	opt := clientOption{
		sender: gosender.New(),
		header: http.Header{},
	}

	for _, f := range opts {
		f.apply(&opt)
	}

	return opt
}

func WithSender(sender Sender) ClientOption {
	return clientOptFunc(func(c *clientOption) {
		c.sender = sender
	})
}

func WithClientHookFn(hookFn HookFn) ClientOption {
	return clientOptFunc(func(c *clientOption) {
		c.hookFns = append(c.hookFns, hookFn)
	})
}

// WithClientHeader sets the header entries associated with key to the single
// element value. It replaces any existing values associated with key. If
// isAdding[0] == true (default is false) then It appends to any existing values
// associated with key
func WithClientHeader(key, value string, isAdding ...bool) ClientOption {
	return clientOptFunc(func(c *clientOption) {
		fn := c.header.Set
		if isAdding != nil && isAdding[0] == true {
			fn = c.header.Add
		}
		fn(key, value)
	})
}
