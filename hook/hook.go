package hook

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/phamvinhdat/httpclient"
)

func Log() httpclient.HookFn {
	return func(ctx context.Context,
		reqChain httpclient.Chain) (*http.Response, error) {
		req := reqChain.GetRequest(ctx)
		reqDump, _ := httputil.DumpRequest(req, true)
		reqDumpStr := string(reqDump)
		log.Println(fmt.Sprintf("--> %s", req.Method),
			"url", req.URL.String(),
			"data", reqDumpStr)

		res, err := reqChain.Proceed(ctx, req)
		if err != nil {
			return nil, err
		}

		resDump, _ := httputil.DumpResponse(res, true)
		resDumpStr := string(resDump)
		log.Println(fmt.Sprintf("<-- END %s", res.Request.Method),
			"url", res.Request.URL.String(),
			"data", resDumpStr)
		return res, nil
	}
}

func UnmarshalResponse(defaultTarget interface{},
	opts ...UnmarshalOption) httpclient.HookFn {
	opt := getUnmarshalOption(opts...)

	return func(ctx context.Context,
		reqChain httpclient.Chain) (*http.Response, error) {
		res, err := reqChain.Proceed(ctx, reqChain.GetRequest(ctx))
		if err != nil {
			return nil, err
		}

		for _, uRes := range opt.unmarshalResponses {
			if uRes.conFn(res) {
				err = unmarshalRes(uRes.target, res, opt.unmarshaller)
				if err != nil {
					return nil, err
				}

				return res, nil
			}
		}

		// handle default case
		if defaultTarget == nil {
			return res, nil
		}
		err = unmarshalRes(defaultTarget, res, opt.unmarshaller)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func unmarshalRes(target interface{}, res *http.Response, unmarshaller Unmarshaller) error {
	var bodyBuffer bytes.Buffer
	_, err := bodyBuffer.ReadFrom(res.Body)
	if err != nil {
		return err
	}

	res.Body = ioutil.NopCloser(&bodyBuffer)

	err = unmarshaller.Unmarshal(bodyBuffer.Bytes(), target)
	if err != nil {
		return err
	}

	return nil
}
