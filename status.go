package silverlining

import (
	"github.com/go-www/silverlining/h1"
)

func DefineStatusLine(status int, statusText string) {
	h1.DefineStatusLine(status, statusText)
}

func GetStatusLine(status int) []byte {
	return h1.GetStatusLine(status)
}
