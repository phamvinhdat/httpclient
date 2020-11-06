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

// UnmarshalResponse will be unmarshal response body to target
func UnmarshalResponse(target interface{},
	unmarshaller ...Unmarshaller) httpclient.HookFn {
	u := getUnmarshaller(unmarshaller...)

	return func(ctx context.Context,
		reqChain httpclient.Chain) (*http.Response, error) {
		res, err := reqChain.Proceed(ctx, reqChain.GetRequest(ctx))
		if err != nil {
			return nil, err
		}

		if target == nil {
			return res, nil
		}
		err = unmarshalRes(target, res, u)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func getUnmarshaller(unmarshaller ...Unmarshaller) Unmarshaller {
	if unmarshaller == nil || unmarshaller[0] == nil {
		return &jsonUnmarshaller{}
	}

	return unmarshaller[0]
}

func unmarshalRes(target interface{}, res *http.Response,
	unmarshaller Unmarshaller) error {
	var bodyBuffer bytes.Buffer
	_, err := bodyBuffer.ReadFrom(res.Body)
	if err != nil {
		return err
	}

	res.Body = ioutil.NopCloser(&bodyBuffer)

	// unmarshall response
	err = unmarshaller.Unmarshal(bodyBuffer.Bytes(), target)
	if err != nil {
		return err
	}

	return nil
}