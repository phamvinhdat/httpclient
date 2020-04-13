package hook

import "encoding/json"

type Unmarshaller interface {
	Unmarshal(data []byte, target interface{}) error
}

type jsonUnmarshaller struct{}

func (*jsonUnmarshaller) Unmarshal(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
