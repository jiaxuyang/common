package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func Post(uri string, data interface{}, params ...interface{}) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	uri = uriAppendQuery(uri, params...)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return do(req)
}

func Put(uri string, data interface{}, params ...interface{}) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	uri = uriAppendQuery(uri, params...)
	req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return do(req)
}

func do(req *http.Request) (*Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body fail")
	}
	if err := checkStatusCode(resp.StatusCode, body); err != nil {
		return nil, err
	}
	ret := &Response{
		StatusCode: resp.StatusCode,
		Data:       body,
	}
	return ret, nil

}

func checkStatusCode(statusCode int, body []byte) error {
	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		err := fmt.Errorf("status code %v, body: %s", statusCode, string(body))
		return NewUserErr("", err)
	}
	if statusCode >= http.StatusInternalServerError {
		err := fmt.Errorf("status code %v, body: %s", statusCode, string(body))
		return NewServiceErr("", err)
	}
	return nil
}

func Get(uri string, params ...interface{}) (*Response, error) {
	uri = uriAppendQuery(uri, params...)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	return do(req)
}

func Delete(uri string, params ...interface{}) (*Response, error) {
	uri = uriAppendQuery(uri, params...)
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return do(req)
}

func uriAppendQuery(uri string, params ...interface{}) string {
	uriQuery := uriQuery(params...)
	if uriQuery != "" {
		uri = fmt.Sprintf("%s?%s", uri, uriQuery)
	}
	return uri
}

func uriQuery(params ...interface{}) string {
	var queryPairs []string
	for i := 0; i < len(params)-1; i += 2 {
		value := fmt.Sprintf("%v", params[i+1])
		value = url.QueryEscape(value)
		queryPairs = append(queryPairs, fmt.Sprintf("%v=%v", params[i], value))
	}
	uriQuery := strings.Join(queryPairs, "&")
	return uriQuery
}
