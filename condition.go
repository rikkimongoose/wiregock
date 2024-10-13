package wiregock

import (
    "time"
)

type Filter struct {
    Contains            *string    `json:"contains,omitempty" bson:"contains,omitempty"`
    EqualTo             *string    `json:"equalTo,omitempty" bson:"equalTo,omitempty"`
    CaseInsensitive     *bool      `json:"caseInsensitive,omitempty" bson:"caseInsensitive,omitempty"`
    BinaryEqualTo       *string    `json:"binaryEqualTo,omitempty" bson:"binaryEqualTo,omitempty"`
    DoesNotContain      *string    `json:"doesNotContain,omitempty" bson:"doesNotContain,omitempty"`
    Matches             *string    `json:"matches,omitempty" bson:"matches,omitempty"`
    DoesNotMatch        *string    `json:"doesNotMatch,omitempty" bson:"doesNotMatch,omitempty"`
    Absent              *bool      `json:"absent,omitempty" bson:"absent,omitempty"`
    And                 []Filter  `json:"and,omitempty" bson:"and,omitempty"`
    Or                  []Filter  `json:"or,omitempty" bson:"or,omitempty"`
    Before              *time.Time `json:"before,omitempty" bson:"before,omitempty"` // "2021-05-01T00:00:00Z"
    After               *time.Time `json:"after,omitempty" bson:"after,omitempty"` // "2021-05-01T00:00:00Z"
    EqualToDateTime     *time.Time `json:"equalToDateTime,omitempty" bson:"equalToDateTime,omitempty"`
    ActualFormat        *string    `json:"equalToDateTime,omitempty" bson:"equalToDateTime,omitempty"`
    EqualToJson         *string   `json:"equalToJson,omitempty" bson:"equalToJson,omitempty"`
    IgnoreArrayOrder    *bool `json:"ignoreArrayOrder,omitempty" bson:"ignoreArrayOrder,omitempty"`
    IgnoreExtraElements *bool `json:"ignoreExtraElements,omitempty" bson:"ignoreExtraElements,omitempty"`
    MatchesJsonPath     *string   `json:"matchesJsonPath,omitempty" bson:"matchesJsonPath,omitempty"`
    EqualToXml          *string   `json:"equalToXml,omitempty" bson:"equalToXml,omitempty"`
    MatchesXPath        *string   `json:"matchesXPath,omitempty" bson:"matchesXPath,omitempty"`
}

type MockData struct {
    Request struct {
        UrlPath              *string            `json:"urlPath,omitempty" bson:"urlPath,omitempty"`
        UrlPattern           *string            `json:"urlPattern,omitempty" bson:"urlPattern,omitempty"`
        Method               *string            `json:"method,omitempty" bson:"method,omitempty"`
        Headers              map[string]Filter `json:"headers,omitempty" bson:"headers,omitempty"`
        QueryParameters      map[string]Filter `json:"queryParameters,omitempty" bson:"queryParameters,omitempty"`
        Cookies              map[string]Filter `json:"cookies,omitempty" bson:"cookies,omitempty"`
        BodyPatterns         []Filter          `json:"bodyPatterns,omitempty" bson:"bodyPatterns,omitempty"`
        BasicAuthCredentials *struct {
            Username string `json:"username,omitempty" bson:"username,omitempty"`
            Password string `json:"password,omitempty" bson:"password,omitempty"`
        } `json:"basicAuthCredentials,omitempty" bson:"basicAuthCredentials,omitempty"`
    } `json:"request" bson:"request"`
    Response struct {
        Status  *int               `json:"status,omitempty" bson:"status,omitempty"`
        Body    *string            `json:"body,omitempty" bson:"body,omitempty"`
        Headers map[string]string `json:"headers,omitempty" bson:"headers,omitempty"`
    } `json:"response" bson:"response"`
}

type Condition interface {
    check() (bool, error)
}

type DataCondition struct {
    loaderData func() (string, bool)
    rule *Rule
}

func (c *DataCondition) check() (bool, error) {
    rule := *c.rule
    data, ok := c.loaderData()
    if !ok {
        return false, nil
    }
    return rule.check(data)
}

type ExistingCondition struct {
    loaderData func() (string, bool)
    rule *Rule
}

func (c *ExistingCondition) check() (bool, error) {
    _, ok := c.loaderData()
    return ok, nil
}

type AndCondition struct {
    conditions []Condition
}

func (c *AndCondition) check() bool {
    for _, cond := range c.conditions {
        res, err := cond.check() 
        if err != nil || !res {
            return false
        }
    }
    return true
}

type OrCondition struct {
    conditions []Condition
}

func (c *OrCondition) check()  bool {
    for _, cond := range c.conditions {
        res, err := cond.check() 
        if err == nil && res {
            return true
        }
    }
    return false
}