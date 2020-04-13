package gosender

import (
	"net/http"
	"net/url"
	"time"
)

type options struct {
	timeout   time.Duration
	transport http.RoundTripper
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(args *options) {
	f(args)
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(args *options) {
		args.timeout = timeout
	})
}

func WithProxyUrl(proxyUrl string) Option {
	return optionFunc(func(args *options) {
		args.transport = getTransport(proxyUrl)
	})
}

func getTransport(proxyUrl string) http.RoundTripper {
	if len(proxyUrl) == 0 {
		return http.DefaultTransport
	}
	proxyURL, _ := url.Parse(proxyUrl)
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return transport
}
