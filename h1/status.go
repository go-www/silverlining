package h1

import (
	"net/http"
)

type Status uint16

const statusMask = 1<<10 - 1

var statusLine [1024][]byte

var _ = func() int {
	Default := []byte("HTTP/1.1 200 OK\r\n")
	for i := range statusLine {
		statusLine[i] = Default
	}

	statusLine[http.StatusContinue] = []byte("HTTP/1.1 100 Continue\r\n")
	statusLine[http.StatusSwitchingProtocols] = []byte("HTTP/1.1 101 Switching Protocols\r\n")
	statusLine[http.StatusProcessing] = []byte("HTTP/1.1 102 Processing\r\n")
	statusLine[http.StatusEarlyHints] = []byte("HTTP/1.1 103 Early Hints\r\n")

	statusLine[http.StatusOK] = []byte("HTTP/1.1 200 OK\r\n")
	statusLine[http.StatusCreated] = []byte("HTTP/1.1 201 Created\r\n")
	statusLine[http.StatusAccepted] = []byte("HTTP/1.1 202 Accepted\r\n")
	statusLine[http.StatusNonAuthoritativeInfo] = []byte("HTTP/1.1 203 Non-Authoritative Information\r\n")
	statusLine[http.StatusNoContent] = []byte("HTTP/1.1 204 No Content\r\n")
	statusLine[http.StatusResetContent] = []byte("HTTP/1.1 205 Reset Content\r\n")
	statusLine[http.StatusPartialContent] = []byte("HTTP/1.1 206 Partial Content\r\n")
	statusLine[http.StatusMultiStatus] = []byte("HTTP/1.1 207 Multi-Status\r\n")
	statusLine[http.StatusAlreadyReported] = []byte("HTTP/1.1 208 Already Reported\r\n")
	statusLine[http.StatusIMUsed] = []byte("HTTP/1.1 226 IM Used\r\n")

	statusLine[http.StatusMultipleChoices] = []byte("HTTP/1.1 300 Multiple Choices\r\n")
	statusLine[http.StatusMovedPermanently] = []byte("HTTP/1.1 301 Moved Permanently\r\n")
	statusLine[http.StatusFound] = []byte("HTTP/1.1 302 Found\r\n")
	statusLine[http.StatusSeeOther] = []byte("HTTP/1.1 303 See Other\r\n")
	statusLine[http.StatusNotModified] = []byte("HTTP/1.1 304 Not Modified\r\n")
	statusLine[http.StatusUseProxy] = []byte("HTTP/1.1 305 Use Proxy\r\n")
	statusLine[306] = []byte("HTTP/1.1 306 Switch Proxy")
	statusLine[http.StatusTemporaryRedirect] = []byte("HTTP/1.1 307 Temporary Redirect\r\n")
	statusLine[http.StatusPermanentRedirect] = []byte("HTTP/1.1 308 Permanent Redirect\r\n")

	statusLine[http.StatusBadRequest] = []byte("HTTP/1.1 400 Bad Request\r\n")
	statusLine[http.StatusUnauthorized] = []byte("HTTP/1.1 401 Unauthorized\r\n")
	statusLine[http.StatusPaymentRequired] = []byte("HTTP/1.1 402 Payment Required\r\n")
	statusLine[http.StatusForbidden] = []byte("HTTP/1.1 403 Forbidden\r\n")
	statusLine[http.StatusNotFound] = []byte("HTTP/1.1 404 Not Found\r\n")
	statusLine[http.StatusMethodNotAllowed] = []byte("HTTP/1.1 405 Method Not Allowed\r\n")
	statusLine[http.StatusNotAcceptable] = []byte("HTTP/1.1 406 Not Acceptable\r\n")
	statusLine[http.StatusProxyAuthRequired] = []byte("HTTP/1.1 407 Proxy Authentication Required\r\n")
	statusLine[http.StatusRequestTimeout] = []byte("HTTP/1.1 408 Request Timeout\r\n")
	statusLine[http.StatusConflict] = []byte("HTTP/1.1 409 Conflict\r\n")
	statusLine[http.StatusGone] = []byte("HTTP/1.1 410 Gone\r\n")
	statusLine[http.StatusLengthRequired] = []byte("HTTP/1.1 411 Length Required\r\n")
	statusLine[http.StatusPreconditionFailed] = []byte("HTTP/1.1 412 Precondition Failed\r\n")
	statusLine[http.StatusRequestEntityTooLarge] = []byte("HTTP/1.1 413 Request Entity Too Large\r\n")
	statusLine[http.StatusRequestURITooLong] = []byte("HTTP/1.1 414 Request-URI Too Long\r\n")
	statusLine[http.StatusUnsupportedMediaType] = []byte("HTTP/1.1 415 Unsupported Media Type\r\n")
	statusLine[http.StatusRequestedRangeNotSatisfiable] = []byte("HTTP/1.1 416 Range Not Satisfiable\r\n")
	statusLine[http.StatusExpectationFailed] = []byte("HTTP/1.1 417 Expectation Failed\r\n")
	statusLine[http.StatusTeapot] = []byte("HTTP/1.1 418 I'm a teapot\r\n")
	statusLine[http.StatusMisdirectedRequest] = []byte("HTTP/1.1 421 Misdirected Request\r\n")
	statusLine[http.StatusUnprocessableEntity] = []byte("HTTP/1.1 422 Unprocessable Entity\r\n")
	statusLine[http.StatusLocked] = []byte("HTTP/1.1 423 Locked\r\n")
	statusLine[http.StatusFailedDependency] = []byte("HTTP/1.1 424 Failed Dependency\r\n")
	statusLine[http.StatusTooEarly] = []byte("HTTP/1.1 425 Too Early\r\n")
	statusLine[http.StatusUpgradeRequired] = []byte("HTTP/1.1 426 Upgrade Required\r\n")
	statusLine[http.StatusPreconditionRequired] = []byte("HTTP/1.1 428 Precondition Required\r\n")
	statusLine[http.StatusTooManyRequests] = []byte("HTTP/1.1 429 Too Many Requests\r\n")
	statusLine[http.StatusRequestHeaderFieldsTooLarge] = []byte("HTTP/1.1 431 Request Header Fields Too Large\r\n")
	// Unofficial status code 444
	statusLine[444] = []byte("HTTP/1.1 444 Connection Closed Without Response\r\n")
	statusLine[http.StatusUnavailableForLegalReasons] = []byte("HTTP/1.1 451 Unavailable For Legal Reasons\r\n")

	statusLine[http.StatusInternalServerError] = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	statusLine[http.StatusNotImplemented] = []byte("HTTP/1.1 501 Not Implemented\r\n")
	statusLine[http.StatusBadGateway] = []byte("HTTP/1.1 502 Bad Gateway\r\n")
	statusLine[http.StatusServiceUnavailable] = []byte("HTTP/1.1 503 Service Unavailable\r\n")
	statusLine[http.StatusGatewayTimeout] = []byte("HTTP/1.1 504 Gateway Timeout\r\n")
	statusLine[http.StatusHTTPVersionNotSupported] = []byte("HTTP/1.1 505 HTTP Version Not Supported\r\n")
	statusLine[http.StatusVariantAlsoNegotiates] = []byte("HTTP/1.1 506 Variant Also Negotiates\r\n")
	statusLine[http.StatusInsufficientStorage] = []byte("HTTP/1.1 507 Insufficient Storage\r\n")
	statusLine[http.StatusLoopDetected] = []byte("HTTP/1.1 508 Loop Detected\r\n")
	statusLine[http.StatusNotExtended] = []byte("HTTP/1.1 510 Not Extended\r\n")
	statusLine[http.StatusNetworkAuthenticationRequired] = []byte("HTTP/1.1 511 Network Authentication Required\r\n")

	return 0
}()

/*
func DefineStatusLine(status int, statusText string) {
	status = status & statusMask
	statusLine[status] = []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", status, statusText))
}
*/

func GetStatusLine(status int) []byte {
	status = status & statusMask
	return statusLine[status]
}
