package wiregock

import (
    "fmt"
    "strings"
)

const anyMathods = [...]string{ "GET", "HEAD", "OPTIONS", "TRACE", "PUT", "DELETE", "POST", "PATCH", "CONNECT" }

func loadUrl(wiremockRequest bson.M) string {
	url, ok = wiremockRequest["url"]
	if ok {
		return url
	}
	urlPath, ok = wiremockRequest["urlPath"]
	if ok {
		return urlPath
	}
	urlPattern, ok = wiremockRequest["urlPattern"]
	if ok {
		return fmt.Sprintf("regex(%s)", urlPattern) 
	}
	return nil
}

func loadMethods(methodNames string) string[] {
    if strings.Compare(methodNames, "ANY") {
        return anyMathods
    }
    return strings.Split(methodNames, ",")
}
