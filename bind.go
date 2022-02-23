package silverlining

import (
	"errors"
	"reflect"
	"strconv"
)

var ErrBindPtrError = errors.New("bind function's parameter must be a pointer")
var ErrBindType = errors.New("bind type error")

func bindStruct(v any, stag string, src func(key string) (value string, found bool)) error {
	s := reflect.ValueOf(v)
	if s.Type().Kind() != reflect.Ptr {
		return ErrBindPtrError
	}
	se := s.Elem()
	if se.Type().Kind() != reflect.Struct {
		return ErrBindType
	}

	for i := 0; i < se.NumField(); i++ {
		f := se.Field(i)
		if f.CanSet() {
			tag := se.Type().Field(i).Tag.Get(stag)
			if tag == "" {
				continue
			}

			v, ok := src(tag)
			if !ok {
				continue
			}

			switch f.Kind() {
			case reflect.String:
				f.SetString(string(v))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				IntValue, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				f.SetInt(IntValue)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				UintValue, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return err
				}
				f.SetUint(UintValue)
			case reflect.Float32:
				FloatValue, err := strconv.ParseFloat(v, 32)
				if err != nil {
					return err
				}
				f.SetFloat(FloatValue)
			case reflect.Float64:
				FloatValue, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				f.SetFloat(FloatValue)
			case reflect.Bool:
				BoolValue, err := strconv.ParseBool(v)
				if err != nil {
					return err
				}
				f.SetBool(BoolValue)
			}
		}
	}

	return nil
}

func (rctx *RequestContext) BindQuery(v any) error {
	return bindStruct(v, "query", func(key string) (value string, found bool) {
		value, err := rctx.GetParam([]byte(key))
		if err != nil {
			found = false
			return
		}
		found = true
		return
	})
}

func (rctx *RequestContext) BindJSON(v any) error {
	return rctx.ReadJSON(v)
}
