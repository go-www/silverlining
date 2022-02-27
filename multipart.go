package silverlining

import (
	"errors"
	"mime"
	"mime/multipart"
	"strings"
)

var ErrContentTypeInvalid = errors.New("content-type invalid")

func (r *Context) MultipartReader() (*multipart.Reader, error) {
	contentType, ok := r.RequestHeaders().Get("Content-Type")
	if !ok {
		return nil, ErrContentTypeInvalid
	}

	mediatype, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(mediatype, "multipart/") {
		return nil, ErrContentTypeInvalid
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, ErrContentTypeInvalid
	}

	return multipart.NewReader(r.BodyReader(), boundary), nil
}
