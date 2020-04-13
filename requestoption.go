package httpclient

import (
	"net/http"

	"github.com/phamvinhdat/httpclient/body"
)

type requestOption struct {
	bodyProvider body.Provider
	header       http.Header
	query        interface{}

	hookFns []HookFn
}

type RequestOption interface {
	apply(*requestOption)
}

func getRequestOption(opts ...RequestOption) requestOption {
	opt := requestOption{
		header: http.Header{},
	}

	for _, f := range opts {
		f.apply(&opt)
	}

	return opt
}

type reqOptFunc func(*requestOption)

func (f reqOptFunc) apply(r *requestOption) {
	f(r)
}

func WithHeader(key, value string) RequestOption {
	return reqOptFunc(func(r *requestOption) {
		r.header.Set(key, value)
	})
}

func WithBodyProvider(bProvider body.Provider) RequestOption {
	return reqOptFunc(func(r *requestOption) {
		r.bodyProvider = bProvider
	})
}

func WithQuery(query interface{}) RequestOption {
	return reqOptFunc(func(r *requestOption) {
		r.query = query
	})
}

func WithHookFn(hookFn HookFn) RequestOption {
	return reqOptFunc(func(r *requestOption) {
		r.hookFns = append(r.hookFns, hookFn)
	})
}
