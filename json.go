package silverlining

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (r *RequestContext) ReadJSON(v any) error {
	bodyReader := r.BodyReader()
	return json.NewDecoder(bodyReader).Decode(v)
}

func (r *RequestContext) WriteJSON(status int, v any) error {
	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.ResponseHeaders().Set("Content-Type", "application/json")
	return r.WriteFullBody(status, encoded)
}

func (r *RequestContext) WriteJSONIndent(status int, v any, prefix string, indent string) error {
	encoded, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return err
	}
	r.ResponseHeaders().Set("Content-Type", "application/json")
	return r.WriteFullBody(status, encoded)
}
