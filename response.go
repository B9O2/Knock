package knock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// Response 响应
type Response struct {
	*http.Response
}

func (r *Response) ReadBody() ([]byte, error) {
	if body, err := io.ReadAll(r.Body); err != nil {
		return nil, err
	} else {
		r.Body = io.NopCloser(bytes.NewReader(body))
		return body, nil
	}
}

func (r *Response) Raw() []byte {
	lines := [][]byte{
		[]byte(fmt.Sprintf("HTTP/%d.%d %s", r.ProtoMajor, r.ProtoMinor, r.Status)),
	}
	for k, values := range r.Header {
		for _, v := range values {
			lines = append(lines, []byte(fmt.Sprintf("%s: %s", k, v)))
		}
	}
	lines = append(lines, []byte{'\r', '\n'})
	if body, err := r.ReadBody(); err != nil {
		lines = append(lines, []byte("<Knock::ReadBody>"+err.Error()))
	} else {
		lines = append(lines, body)
	}

	return bytes.Join(lines, []byte{'\r', '\n'})
}

func (r *Response) String() string {
	return string(r.Raw())
}
