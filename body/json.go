package body

import (
	"encoding/json"
	"io"
	"strings"
)

const jsonContentType = "application/json"

func NewJson(body interface{}) Provider {
	return &jsonProvider{payload: body}
}

type jsonProvider struct {
	payload interface{}
}

func (j *jsonProvider) Provide() (io.Reader, string, error) {
	dataBytes, err := json.Marshal(j.payload)
	if err != nil {
		return nil, "", err
	}

	return strings.NewReader(string(dataBytes)), jsonContentType, nil
}
