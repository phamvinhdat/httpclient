package body

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

func NewMultiPartFile(data, files map[string]string) Provider {
	return &multipartFileProvider{
		files: files,
		data:  data,
	}
}

type multipartFileProvider struct {
	files map[string]string
	data  map[string]string
}

func (m *multipartFileProvider) Provide() (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fileName, filePath := range m.files {
		err := appendFileToForm(writer, fileName, filePath)
		if err != nil {
			return nil, "", err
		}
	}

	for key, val := range m.data {
		_ = writer.WriteField(key, val)
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

func appendFileToForm(writer *multipart.Writer, name, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	var fw io.Writer
	if fw, err = writer.CreateFormFile(name, path); err != nil {
		return err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return err
	}
	return nil
}
