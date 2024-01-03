package knock

import (
	"bytes"
	"fmt"
	"net/http"
)

// Response 响应
type Response struct {
	*http.Response
	body []byte
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
	lines = append(lines, r.body)

	return bytes.Join(lines, []byte{'\r', '\n'})
}

func (r *Response) ReadBody() []byte {
	return r.body
}

func (r *Response) String() string {
	return string(r.Raw())
}
