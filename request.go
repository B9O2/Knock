package knock

import (
	"bytes"
	"errors"
	"github.com/B9O2/rawhttp/client"
	"strconv"
)

type Method string

func (m Method) String() string {
	suffix := "(Method)"
	methods := []string{"GET", "POST", "HEAD", "PUT"}
	for _, v := range methods {
		if string(m) == v {
			return v + suffix
		}
	}
	return "Unknown(Method)"
}

var (
	GET  = Method("GET")
	POST = Method("POST")
)

type HTTPVersion client.Version

var (
	HTTP_1_1 = HTTPVersion(client.HTTP_1_1)
	HTTP_1_0 = HTTPVersion(client.HTTP_1_0)
	HTTP_2_0 = HTTPVersion{
		Major: 2,
		Minor: 0,
	}
)

type Request interface {
	Method() Method
	URI() string
	Version() HTTPVersion
	Headers() map[string][]string
	Body() []byte
	Patch(Request)
}

// BaseRequest 基础请求
type BaseRequest struct {
	method  Method
	uri     string
	version HTTPVersion
	headers map[string][]string
	body    []byte
}

func (br *BaseRequest) Headers() map[string][]string {
	return br.headers
}

func (br *BaseRequest) Body() []byte {
	return br.body
}

func (br *BaseRequest) Method() Method {
	return br.method
}

func (br *BaseRequest) URI() string {
	return br.uri
}

func (br *BaseRequest) Version() HTTPVersion {
	return br.version
}

func (br *BaseRequest) Patch(r Request) {
	if len(r.URI()) > 0 {
		br.uri = r.URI()
	}

	if len(r.Method()) > 0 {
		br.method = r.Method()
	}

	if len(r.Headers()) > 0 {
		for k, v := range r.Headers() {
			br.headers[k] = v
		}
	}

	if len(r.Body()) > 0 {
		br.body = r.Body()
	}
}

func NewBaseRequest(method Method, uri string, version HTTPVersion, headers map[string][]string, body []byte) *BaseRequest {
	if !bytes.HasSuffix(body, []byte{'\r', '\n'}) {
		body = append(body, []byte{'\r', '\n'}...)
	}
	return &BaseRequest{
		method:  method,
		uri:     uri,
		version: version,
		headers: headers,
		body:    body,
	}
}

func NewBaseRequestFromRaw(raw []byte) (*BaseRequest, error) {
	lines := bytes.Split(raw, []byte{'\n'})
	if len(lines) <= 0 {
		return nil, errors.New("malformed request. cause: no content")
	}

	//Request Line
	parts := bytes.SplitN(bytes.TrimRight(lines[0], "\r"), []byte{' '}, 3)
	if len(parts) != 3 {
		return nil, errors.New("malformed request. cause: '" + string(lines[0]) + "'")
	}
	lines = lines[1:]

	//Method URI ProtocolVersion
	method := parts[0]
	uri := parts[1]
	_, versionRaw, _ := bytes.Cut(parts[2], []byte{'/'})
	majorRaw, minorRaw, _ := bytes.Cut(versionRaw, []byte{'.'})
	major, err := strconv.Atoi(string(majorRaw))
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(string(minorRaw))
	if err != nil {
		return nil, err
	}
	version := HTTPVersion{
		Major: major,
		Minor: minor,
	}

	//RawHeaders 未来可能有用
	rawHeaders := map[string][][]byte{}
	for i, line := range lines {
		if len(line) <= 0 {
			//body
			lines = lines[i+1:]
			break
		}
		k, v, _ := bytes.Cut(bytes.TrimRight(line, "\r"), []byte{':', ' '})
		key := string(k)
		if _, ok := rawHeaders[key]; ok {
			rawHeaders[key] = append(rawHeaders[key], v)
		} else {
			rawHeaders[key] = [][]byte{v}
		}
	}
	//Headers
	headers := map[string][]string{}
	for k, vs := range rawHeaders {
		if _, ok := headers[k]; !ok {
			headers[k] = []string{}
		}
		for _, v := range vs {
			headers[k] = append(headers[k], string(v))
		}
	}

	//Body

	body := bytes.Join(lines, []byte{'\n'})
	if !bytes.HasSuffix(body, []byte{'\r', '\n'}) {
		body = append(body, []byte{'\r', '\n'}...)
	}
	return &BaseRequest{
		method:  Method(method),
		uri:     string(uri),
		version: version,
		headers: headers,
		body:    body,
	}, nil
}
