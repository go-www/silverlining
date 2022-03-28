package silverlining

import "github.com/goccy/go-json"

//var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (r *Context) ReadJSON(v any) error {
	bodyReader := r.BodyReader()
	return json.NewDecoder(bodyReader).Decode(v)
}

func (r *Context) WriteJSON(status int, v any) error {
	encoded, err := json.Marshal(v)
	if err != nil {
		return err
	}
	r.ResponseHeaders().Set("Content-Type", "application/json")
	return r.WriteFullBody(status, encoded)
}

func (r *Context) WriteJSONIndent(status int, v any, prefix string, indent string) error {
	encoded, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return err
	}
	r.ResponseHeaders().Set("Content-Type", "application/json")
	return r.WriteFullBody(status, encoded)
}

func (r *Context) WriteJSONStream(status int, v any) error {
	r.ResponseHeaders().Set("Content-Type", "application/json")
	chunked := r.ChunkedBodyWriter()
	defer chunked.Close()
	encoder := json.NewEncoder(chunked)
	err := encoder.Encode(v)
	return err
}
