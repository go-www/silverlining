package h1

type Method uint8

const (
	MethodInvalid Method = iota
	MethodGET
	MethodHEAD
	MethodPOST
	MethodPUT
	MethodDELETE
	MethodCONNECT
	MethodOPTIONS
	MethodTRACE
	MethodPATCH

	MethodBREW // HTCPCP/1.0 (https://datatracker.ietf.org/doc/html/rfc2324)
)

/*
	According to RFC7231 Section 4.1,
	All general purpose HTTP/1.1 servers MUST support the GET, HEAD.
*/

func (m Method) String() string {
	switch m {
	case MethodGET:
		return "GET"
	case MethodHEAD:
		return "HEAD"
	case MethodPOST:
		return "POST"
	case MethodPUT:
		return "PUT"
	case MethodDELETE:
		return "DELETE"
	case MethodCONNECT:
		return "CONNECT"
	case MethodOPTIONS:
		return "OPTIONS"
	case MethodTRACE:
		return "TRACE"
	case MethodPATCH:
		return "PATCH"
	case MethodBREW:
		return "BREW"
	default:
		return "UNKNOWN"
	}
}
