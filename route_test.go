package wiregock

import (
	"testing"
	"slices"
)

func TestLoadMethodsCheck(t *testing.T) {
	methods := loadMethods("ANY")
	for _, method := range anyMethods {
		if !slices.Contains(methods, method) {
        	t.Fatalf(`%s method isn't loaded from ANY`, method)
    	}
	}
	methodsGetPost := loadMethods("GET, POST")
	if !slices.Contains(methodsGetPost, "GET") {
    	t.Fatalf(`%s method isn't loaded from "GET, POST"`, "GET")
	}
	if !slices.Contains(methodsGetPost, "POST") {
    	t.Fatalf(`%s method isn't loaded by ANY`, "POST")
	}
}