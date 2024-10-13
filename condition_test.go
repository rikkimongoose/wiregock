package wiregock

import (
	"encoding/json"
	"testing"
)

type TrueRule struct {
}

type FalseRule struct {
}

func (rule *TrueRule) check(str string) (bool, error) {
	return true, nil
}

func (rule *FalseRule) check(str string) (bool, error) {
	return false, nil
}

func TestMarshaling(t *testing.T) {
	body := []byte(`
{
    "request": {
        "urlPath": "/everything",
        "method": "ANY",
        "headers": {
            "Accept": {
                "contains": "xml"
            }
        },
        "queryParameters": {
            "search_term": {
                "equalTo": "WireMock"
            }
        },
        "cookies": {
            "session": {
                "matches": ".*12345.*"
            }
        },
        "bodyPatterns": [
            {
                "equalToXml": "<search-results />"
            },
            {
                "matchesXPath": "//search-results"
            }
        ],
        "multipartPatterns": [
            {
                "matchingType": "ANY",
                "headers": {
                    "Content-Disposition": {
                        "contains": "name=\"info\""
                    },
                    "Content-Type": {
                        "contains": "charset"
                    }
                },
                "bodyPatterns": [
                    {
                        "equalToJson": "{}"
                    }
                ]
            }
        ],
        "basicAuthCredentials": {
            "username": "jeff@example.com",
            "password": "jeffteenjefftyjeff"
        }
    },
    "response": {
        "status": 200
    }
}`)

	var mockData MockData
	err := json.Unmarshal(body, &mockData)
	if err != nil {
    	t.Fatalf(`Error parsing JSON format: %s`, err)
	}
	if *mockData.Request.QueryParameters["search_term"].EqualTo != "WireMock" {
    	t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.QueryParameters[\"search_term\"].EqualTo")
	}
}