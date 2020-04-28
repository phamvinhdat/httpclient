package httpclient

import (
	"context"
	"io"
	"net/http"

	goquery "github.com/google/go-querystring/query"
)

const contentTypeKey = "Content-Type"

func (c *client) buildHttpRequest(ctx context.Context, url, method string,
	reqOpt requestOption) (*http.Request, error) {
	var (
		body        io.Reader
		contentType string
		err         error
	)
	if reqOpt.bodyProvider != nil {
		body, contentType, err = reqOpt.bodyProvider.Provide()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// set header
	req.Header = c.header.Clone()  // client header
	addHeaders(req, reqOpt.header) // request header

	if len(contentType) > 0 {
		req.Header.Set(contentTypeKey, contentType)
	}

	// set query
	err = setQuery(req, reqOpt.query)
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

func addHeaders(req *http.Request, header http.Header) {
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}
