package knock

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/B9O2/knock/rawhttp/client"
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

func (h HTTPVersion) String() string {
	return fmt.Sprintf("HTTP/%d.%d", h.Major, h.Minor)
}

var (
	HTTP_1_1 = HTTPVersion(client.HTTP_1_1)
	HTTP_1_0 = HTTPVersion(client.HTTP_1_0)
	HTTP_2_0 = HTTPVersion{
		Major: 2,
		Minor: 0,
	}
)

type Request interface {
	Raw() []byte
}

// HTTPRequest 基础请求
type HTTPRequest struct {
	method  Method
	uri     string
	version HTTPVersion
	headers map[string][]string
	body    []byte
}

func (br HTTPRequest) Headers() map[string][]string {
	return br.headers
}

func (br HTTPRequest) Body() []byte {
	return br.body
}

func (br HTTPRequest) Method() Method {
	return br.method
}

func (br HTTPRequest) URI() string {
	return br.uri
}

func (br HTTPRequest) Version() HTTPVersion {
	return br.version
}

func (br HTTPRequest) Patch(patchr HTTPRequest) HTTPRequest {
	r := HTTPRequest{
		headers: make(map[string][]string),
	}

	if len(patchr.URI()) > 0 {
		r.uri = patchr.URI()
	}

	if len(patchr.Method()) > 0 {
		r.method = patchr.Method()
	}

	if len(patchr.Headers()) > 0 {
		for k, v := range patchr.Headers() {
			r.headers[k] = v
		}
	}

	if len(patchr.Body()) > 0 {
		r.body = patchr.Body()
	}

	r.version = patchr.version

	return r
}

func (br HTTPRequest) Raw() []byte {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("%s %s %s"+client.NewLine, br.method, br.uri, br.version))

	for key, values := range br.headers {
		if len(values) != 0 {
			for _, v := range values {
				b.WriteString(fmt.Sprintf("%s: %s"+client.NewLine, key, v))
			}
		} else {
			b.WriteString(fmt.Sprintf("%s"+client.NewLine, key))
		}
	}

	b.WriteString(client.NewLine)
	b.Write(br.body)
	return []byte(strings.ReplaceAll(b.String(), "\n", client.NewLine))
}

func NewBaseRequest(method Method, uri string, version HTTPVersion, headers map[string][]string, body []byte) HTTPRequest {
	if !bytes.HasSuffix(body, []byte{'\r', '\n'}) {
		body = append(body, []byte{'\r', '\n'}...)
	}
	return HTTPRequest{
		method:  method,
		uri:     uri,
		version: version,
		headers: headers,
		body:    body,
	}
}

func NewBaseRequestFromRaw(raw []byte) (HTTPRequest, error) {
	lines := bytes.Split(raw, []byte{'\n'})
	if len(lines) <= 0 {
		return HTTPRequest{}, errors.New("malformed request. cause: no content")
	}

	//Request Line
	parts := bytes.SplitN(bytes.TrimRight(lines[0], "\r"), []byte{' '}, 3)
	if len(parts) != 3 {
		return HTTPRequest{}, errors.New("malformed request. cause: '" + string(lines[0]) + "'")
	}
	lines = lines[1:]

	//Method URI ProtocolVersion
	method := parts[0]
	uri := parts[1]
	_, versionRaw, _ := bytes.Cut(parts[2], []byte{'/'})
	majorRaw, minorRaw, _ := bytes.Cut(versionRaw, []byte{'.'})
	major, err := strconv.Atoi(string(majorRaw))
	if err != nil {
		return HTTPRequest{}, err
	}
	minor, err := strconv.Atoi(string(minorRaw))
	if err != nil {
		return HTTPRequest{}, err
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
		} //else {
		//rawHeaders[key] = [][]byte{v}
		//}
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
	return HTTPRequest{
		method:  Method(method),
		uri:     string(uri),
		version: version,
		headers: headers,
		body:    body,
	}, nil
}
