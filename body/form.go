package body

import (
	"io"
	"net/url"
	"strings"

	goquery "github.com/google/go-querystring/query"
)

const formContentType = "application/x-www-form-urlencoded"

func NewForm(body interface{}) Provider {
	return &formProvider{payload: body}
}

func NewURLForm(body url.Values) Provider {
	return &urlFormProvider{payload: body}
}

type formProvider struct {
	payload interface{}
}

func (f *formProvider) Provide() (io.Reader, string, error) {
	values, err := goquery.Values(f.payload)
	if err != nil {
		return nil, "", err
	}

	return strings.NewReader(values.Encode()), formContentType, nil
}

type urlFormProvider struct {
	payload url.Values
}

func (f *urlFormProvider) Provide() (io.Reader, string, error) {
	return strings.NewReader(f.payload.Encode()), formContentType, nil
}
