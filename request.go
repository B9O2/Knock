package knock

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

type Request interface {
	Method() Method
	URI() string
	Headers() map[string][]string
	Body() []byte
}

// BaseRequest 基础请求
type BaseRequest struct {
	method  Method
	uri     string
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
