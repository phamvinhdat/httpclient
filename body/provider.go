package body

import "io"

type Provider interface {
	Provide() (body io.Reader, contentType string, err error)
}
