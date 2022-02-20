package silverlining

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (r *RequestContext) ReadAsJSON(v any) error {
	bodyReader := r.BodyReader()
	return json.NewDecoder(bodyReader).Decode(v)
}

func (r *RequestContext) WriteAsJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.SetHeader("Content-Type", "application/json")
	r.SetContentLength(len(data))

	return nil
}
