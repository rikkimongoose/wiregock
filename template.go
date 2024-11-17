package wiregock

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type RequestData map[string]interface{}

func ParseQuery(values url.Values) map[string]map[string]string {
	response := map[string]map[string]string{}
	for key, value := range values {
		dataMap := map[string]string{}
		for index, valueItem := range value {
			dataMap[fmt.Sprintf("[%d]", index)] = valueItem
		}
		response[key] = dataMap
	}
	return response
}

func ToSingleValueMap(values map[string][]string) map[string]string {
	resp := map[string]string{}
	for key, value := range values {
		resp[key] = value[0]
	}
	return resp
}

func CookiesToMap(cookies []*http.Cookie) map[string]string {
	resp := map[string]string{}
	for _, cookie := range cookies {
		resp[cookie.Name] = cookie.Value
	}
	return resp
}

func LoadRequestData(req *http.Request) (*RequestData, error) {
	body := ""
	bodyBase64 := ""
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		body = string(b[:])
		bodyBase64 = b64.URLEncoding.EncodeToString(b)
	}
	return &RequestData{
		"request": RequestData{
			"id":           uuid.New().String(),
			"url":          req.URL.RequestURI(),
			"queryFull":    req.URL.Query(),
			"query":        ToSingleValueMap(req.URL.Query()),
			"method":       req.Method,
			"host":         req.Host,
			"port":         req.URL.Port(),
			"scheme":       req.URL.Scheme,
			"baseUrl":      req.URL.Host,
			"headersFull":  req.Header,
			"headers":      ToSingleValueMap(req.Header),
			"cookies":      CookiesToMap(req.Cookies()),
			"body":         body,
			"bodyAsBase64": bodyBase64,
		},
	}, nil
}
