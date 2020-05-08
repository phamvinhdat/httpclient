package hook

import (
	"net/http"
)

type UnmarshalConditionFn func(response *http.Response) bool
type unmarshalResponse struct {
	target interface{}
	conFn  UnmarshalConditionFn
}

type UnmarshalOption interface {
	apply(opt *unmarshalOption)
}

type unmarshalOption struct {
	unmarshaller       Unmarshaller
	unmarshalResponses []unmarshalResponse
}

type unmarshalOptionFn func(opt *unmarshalOption)

func (fn unmarshalOptionFn) apply(opt *unmarshalOption) {
	fn(opt)
}

func WithUnmashaller(unmarshaller Unmarshaller) UnmarshalOption {
	return unmarshalOptionFn(func(opt *unmarshalOption) {
		opt.unmarshaller = unmarshaller
	})
}

// WithUnmarshalResponse add unmarshalResponse to unmarshalResponses
// if more than UnmarshalConditionFn is true, only first case in unmarshalResponses
// will be execute. Default target use when all UnmarshalConditionFns are false
func WithUnmarshalResponse(uRes unmarshalResponse) UnmarshalOption {
	return unmarshalOptionFn(func(opt *unmarshalOption) {
		opt.unmarshalResponses = append(opt.unmarshalResponses, uRes)
	})
}

func getUnmarshalOption(opts ...UnmarshalOption) unmarshalOption {
	opt := unmarshalOption{
		unmarshaller: &jsonUnmarshaller{},
	}
	for _, fn := range opts {
		fn.apply(&opt)
	}

	return opt
}
