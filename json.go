package silverlining

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (r *RequestContext) ReadJSON(v any) error {
	bodyReader := r.BodyReader()
	return json.NewDecoder(bodyReader).Decode(v)
}

func (r *RequestContext) WriteJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.SetHeader("Content-Type", "application/json")
	r.SetContentLength(len(data))
	r.WriteHeader(200)
	r.Write(data)

	return nil
}