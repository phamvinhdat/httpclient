package httpclient

import (
	"context"
	"io"
	"net/http"

	goquery "github.com/google/go-querystring/query"
)

const contentTypeKey = "Content-Type"

func buildHttpRequest(ctx context.Context, url, method string,
	opt requestOption) (*http.Request, error) {
	var (
		body        io.Reader
		contentType string
		err         error
	)
	if opt.bodyProvider != nil {
		body, contentType, err = opt.bodyProvider.Provide()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// set header
	req.Header = opt.header
	if len(contentType) > 0 {
		req.Header.Set(contentTypeKey, contentType)
	}

	// set query
	err = setQuery(req, opt.query)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func setQuery(req *http.Request, query interface{}) error {
	v, err := goquery.Values(query)
	if err != nil {
		return err
	}
	req.URL.RawQuery = v.Encode()
	return nil
}
