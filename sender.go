package httpclient

import (
	"context"
	"net/http"
)

type Sender interface {
	Send(ctx context.Context, req *http.Request) (*http.Response, error)
}
