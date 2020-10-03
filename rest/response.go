package rest

import (
	"github.com/tidwall/gjson"
)

type Response struct {
	StatusCode int
	Data       []byte
}

func (r *Response) Map(path string) map[string]interface{} {
	m, ok := r.JSON(path).Value().(map[string]interface{})
	if !ok {
		return nil
	}
	return m
}

func (r *Response) JSON(path string) gjson.Result {
	if len(path) == 0 {
		return gjson.ParseBytes(r.Data)
	}
	return gjson.ParseBytes(r.Data).Get(path)
}
