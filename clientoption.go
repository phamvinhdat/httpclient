package httpclient

import "github.com/phamvinhdat/httpclient/gosender"

type clientOption struct {
	sender  Sender
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
