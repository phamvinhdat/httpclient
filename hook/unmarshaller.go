package hook

import (
	"encoding/json"
	"encoding/xml"
)

type Unmarshaller interface {
	Unmarshal(data []byte, target interface{}) error
}

type jsonUnmarshaller struct{}

func (*jsonUnmarshaller) Unmarshal(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

type xmlUnmarshaller struct{}

func (*xmlUnmarshaller) Unmarshal(data []byte, target interface{}) error {
	return xml.Unmarshal(data, target)
}

func NewJsonUnmarshaller() Unmarshaller {
	return &jsonUnmarshaller{}
}

func NewXmlUnmarshaller() Unmarshaller {
	return &xmlUnmarshaller{}
}
