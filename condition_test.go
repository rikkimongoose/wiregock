package wiregock

import (
	"encoding/json"
	"testing"
)

func TestUnmarshaling(t *testing.T) {
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
    if *mockData.Request.UrlPath != "/everything" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.urlPath")
    }
    if *mockData.Request.Method != "ANY" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.method")
    }
    if *mockData.Request.Headers["Accept"].Contains != "xml" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.Headers[\"Accept\"].Contains")
    }
    if *mockData.Request.QueryParameters["search_term"].EqualTo != "WireMock" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.QueryParameters[\"search_term\"].EqualTo")
    }
	if *mockData.Request.Cookies["session"].Matches != ".*12345.*" {
    	t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.Cookies.Session[\"matches\"].EqualTo")
	}
    if *mockData.Request.BodyPatterns[0].EqualToXml != "<search-results />" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BodyPatterns[0].EqualToXml")
    }
    if mockData.Request.BodyPatterns[1].MatchesXPath.Expression != "//search-results" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BodyPatterns[1].MatchesXPath.Expression")
    }
    if *mockData.Request.BasicAuthCredentials.Username != "jeff@example.com" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BasicAuthCredentials.Username")
    }
    if *mockData.Request.BasicAuthCredentials.Password != "jeffteenjefftyjeff" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BasicAuthCredentials.Password")
    }
}

func TestAndCondition(t *testing.T) {
    conditionTrue := AndCondition{[]Condition{DataCondition{rulesAnd:[]Rule{TrueRule{}}}}}
    res, err := conditionTrue.check()
    if err != nil {
        t.Fatalf(`Error conditionTrue: %s`, err)
    }
    if !res {
        t.Fatalf(`Wrong execution conditionTrue`)
    }
    conditionTrueTrue := AndCondition{[]Condition{DataCondition{rulesAnd:[]Rule{TrueRule{}, TrueRule{}}}}}
    res, err = conditionTrueTrue.check()
    if err != nil {
        t.Fatalf(`Error conditionTrueTrue: %s`, err)
    }
    if !res {
        t.Fatalf(`Wrong execution conditionTrueTrue`)
    }
    conditionTrueFalse := AndCondition{[]Condition{DataCondition{rulesAnd:[]Rule{TrueRule{}, FalseRule{}}}}}
    res, err = conditionTrueFalse.check()
    if err != nil {
        t.Fatalf(`Error conditionTrueFalse: %s`, err)
    }
    if res {
        t.Fatalf(`Wrong execution conditionTrueFalse`)
    }
    conditionTrueFalseOr := AndCondition{[]Condition{DataCondition{rulesOr:[]Rule{TrueRule{}, FalseRule{}}}}}
    res, err = conditionTrueFalseOr.check()
    if err != nil {
        t.Fatalf(`Error conditionTrueFalseOr: %s`, err)
    }
    if !res {
        t.Fatalf(`Wrong execution conditionTrueFalseOr`)
    }
    conditionTrueTrueOr := AndCondition{[]Condition{DataCondition{rulesOr:[]Rule{TrueRule{}, TrueRule{}}}}}
    res, err = conditionTrueTrueOr.check()
    if err != nil {
        t.Fatalf(`Error conditionTrueTrueOr: %s`, err)
    }
    if !res {
        t.Fatalf(`Wrong execution conditionTrueTrueOr`)
    }
}


func TestUnmarshalingXPathFilter(t *testing.T) {
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
        "bodyPatterns": [
            {
                "equalToXml": "<search-results />"
            },
            {
                "matchesXPath": {
                    "expression": "//search-results",
                    "contains": "wash",
                    "equalToXml": "<todo-item>Do the washing</todo-item>",
                    "equalToJson": "{}"
                }
                
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
    if mockData.Request.BodyPatterns[1].MatchesXPath.Expression != "//search-results" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BodyPatterns[1].MatchesXPath.Expression")
    }
    if *mockData.Request.BodyPatterns[1].MatchesXPath.Contains != "wash" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BodyPatterns[1].MatchesXPath.Contains")
    }
    if *mockData.Request.BodyPatterns[1].MatchesXPath.EqualToXml != "<todo-item>Do the washing</todo-item>" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BodyPatterns[1].MatchesXPath.EqualToXml")
    }
    if *mockData.Request.BodyPatterns[1].MatchesXPath.EqualToJson != "{}" {
        t.Fatalf(`Unable to load from parsed JSON: %s`, "mockData.Request.BodyPatterns[1].MatchesXPath.EqualToJson")
    }
}