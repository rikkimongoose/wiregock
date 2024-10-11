package wiregock

import (
    "github.com/gofiber/fiber/v3"
)

type WebContextMock struct {
	body string
	headers map[string]string
	params map[string]string
	cookies map[string]string
}

func (mock *WebContextMock) Body() []byte {
	return []byte{mock.body}
}

func (mock *WebContextMock) Get(key string, defaultValue ...string) string {
	if val, ok := mock.headers[key]; ok {
	    return val
	}
	return defaultValue
}

func (mock *WebContextMock) Params(key string, defaultValue ...string) string {
	if val, ok := mock.params[key]; ok {
	    return val
	}
	return defaultValue
}
func (mock *WebContextMock) Cookies(key string, defaultValue ...string) string {
	if val, ok := mock.cookies[key]; ok {
	    return val
	}
	return defaultValue
}