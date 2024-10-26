package wiregock

import (
    "time"
    "encoding/json"
    "mime/multipart"
    "go.mongodb.org/mongo-driver/bson"
)

type Filter struct {
    Contains            *string       `json:"contains,omitempty" bson:"contains,omitempty"`
    EqualTo             *string       `json:"equalTo,omitempty" bson:"equalTo,omitempty"`
    CaseInsensitive     *bool         `json:"caseInsensitive,omitempty" bson:"caseInsensitive,omitempty"`
    BinaryEqualTo       *string       `json:"binaryEqualTo,omitempty" bson:"binaryEqualTo,omitempty"`
    DoesNotContain      *string       `json:"doesNotContain,omitempty" bson:"doesNotContain,omitempty"`
    Matches             *string       `json:"matches,omitempty" bson:"matches,omitempty"`
    DoesNotMatch        *string       `json:"doesNotMatch,omitempty" bson:"doesNotMatch,omitempty"`
    Absent              *bool         `json:"absent,omitempty" bson:"absent,omitempty"`
    And                 []Filter      `json:"and,omitempty" bson:"and,omitempty"`
    Or                  []Filter      `json:"or,omitempty" bson:"or,omitempty"`
    Before              *time.Time    `json:"before,omitempty" bson:"before,omitempty"` // "2021-05-01T00:00:00Z"
    After               *time.Time    `json:"after,omitempty" bson:"after,omitempty"` // "2021-05-01T00:00:00Z"
    EqualToDateTime     *time.Time    `json:"equalToDateTime,omitempty" bson:"equalToDateTime,omitempty"`
    ActualFormat        *string       `json:"actualFormat,omitempty" bson:"actualFormat,omitempty"`
    EqualToJson         *string       `json:"equalToJson,omitempty" bson:"equalToJson,omitempty"`
    IgnoreArrayOrder    *bool         `json:"ignoreArrayOrder,omitempty" bson:"ignoreArrayOrder,omitempty"`
    IgnoreExtraElements *bool         `json:"ignoreExtraElements,omitempty" bson:"ignoreExtraElements,omitempty"`
    MatchesJsonPath     *XPathFilter  `json:"matchesJsonPath,omitempty" bson:"matchesJsonPath,omitempty"`
    EqualToXml          *string       `json:"equalToXml,omitempty" bson:"equalToXml,omitempty"`
    MatchesXPath        *XPathFilter  `json:"matchesXPath,omitempty" bson:"matchesXPath,omitempty"`
    Includes            []MultiFilter `json:"includes,omitempty" bson:"includes,omitempty"`
    HasExactly          []MultiFilter `json:"hasExactly,omitempty" bson:"hasExactly,omitempty"`
}

type XPathFilter struct {
    Expression          string            `json:"-" bson:"-"`
    EqualTo             *string           `json:"equalTo,omitempty" bson:"equalTo,omitempty"`
    EqualToJson         *string           `json:"equalToJson,omitempty" bson:"equalToJson,omitempty"`
    EqualToXml          *string           `json:"equalToXml,omitempty" bson:"equalToXml,omitempty"`
    Contains            *string           `json:"contains,omitempty" bson:"contains,omitempty"`
    CaseInsensitive     *bool             `json:"caseInsensitive,omitempty" bson:"caseInsensitive,omitempty"`
    IgnoreArrayOrder    *bool             `json:"ignoreArrayOrder,omitempty" bson:"ignoreArrayOrder,omitempty"`
    IgnoreExtraElements *bool             `json:"ignoreExtraElements,omitempty" bson:"ignoreExtraElements,omitempty"`
    And                 []XPathFilter     `json:"and,omitempty" bson:"and,omitempty"`
    Before              *time.Time        `json:"before,omitempty" bson:"before,omitempty"` // "2021-05-01T00:00:00Z"
    After               *time.Time        `json:"after,omitempty" bson:"after,omitempty"` // "2021-05-01T00:00:00Z"
    EqualToDateTime     *time.Time        `json:"equalToDateTime,omitempty" bson:"equalToDateTime,omitempty"`
    ActualFormat        *string           `json:"actualFormat,omitempty" bson:"actualFormat,omitempty"`
    XPathNamespaces     map[string]string `json:"xPathNamespaces,omitempty" bson:"xPathNamespaces,omitempty"`
}

type MultiFilter struct {
    EqualTo             *string           `json:"equalTo,omitempty" bson:"equalTo,omitempty"`
    Contains            *string           `json:"contains,omitempty" bson:"contains,omitempty"`
    DoesNotContain      *string           `json:"doesNotContain,omitempty" bson:"doesNotContain,omitempty"`    
}

type MockRequest struct {
    UrlPath              *string                 `json:"urlPath,omitempty" bson:"urlPath,omitempty"`
    UrlPattern           *string                 `json:"urlPattern,omitempty" bson:"urlPattern,omitempty"`
    Method               *string                 `json:"method,omitempty" bson:"method,omitempty"`
    Headers              map[string]Filter       `json:"headers,omitempty" bson:"headers,omitempty"`
    QueryParameters      map[string]Filter       `json:"queryParameters,omitempty" bson:"queryParameters,omitempty"`
    FormParameters       map[string]Filter       `json:"formParameters,omitempty" bson:"formParameters,omitempty"`    
    Cookies              map[string]Filter       `json:"cookies,omitempty" bson:"cookies,omitempty"`
    BodyPatterns         []Filter                `json:"bodyPatterns,omitempty" bson:"bodyPatterns,omitempty"`
    MultipartPatterns    []MultipartPatternsData `json:"multipartPatterns,omitempty" bson:"multipartPatterns,omitempty"`
    BasicAuthCredentials *struct {
        Username *string `json:"username,omitempty" bson:"username,omitempty"`
        Password *string `json:"password,omitempty" bson:"password,omitempty"`
    } `json:"basicAuthCredentials,omitempty" bson:"basicAuthCredentials,omitempty"`
}

type MultipartPatternsData struct {
    MatchingType     *string           `json:"matchingType,omitempty" bson:"matchingType,omitempty"`
    FileName         []Filter          `json:"fileName,omitempty" bson:"fileName,omitempty"`
    Headers          map[string]Filter `json:"headers,omitempty" bson:"headers,omitempty"`
    BodyPatterns     []Filter          `json:"bodyPatterns,omitempty" bson:"bodyPatterns,omitempty"`
}

type MockData struct {
    Request *MockRequest  `json:"request" bson:"request"`
    Response *struct {
        Status          *int               `json:"status,omitempty" bson:"status,omitempty"`
        Body            *string            `json:"body,omitempty" bson:"body,omitempty"`
        BodyFileName    *string            `json:"bodyFileName,omitempty" bson:"bodyFileName,omitempty"`
        JsonBody        *interface{}       `json:"jsonBody,omitempty" bson:"jsonBody,omitempty"`
        Headers         map[string]string  `json:"headers,omitempty" bson:"headers,omitempty"`
        Cookies         map[string]string  `json:"cookies,omitempty" bson:"cookies,omitempty"`
    } `json:"response" bson:"response"`
}

type Condition interface {
    Check() (bool, error)
}

type DataCondition struct {
    loaderMethod func() string
    blockRule Rule
}

type DataConditionMulti struct {
    loaderMethod func() []string
    rulesAnd []Rule
    rulesOr []Rule
}

type MultipartDataCondition struct {
    checkAny bool
    loaderMethod func() multipart.Form
    rulesFileName *Rule
    rulesHeader map[string]Rule
    rulesBody *Rule
}

func (c DataCondition) Check() (bool, error) {
    data := ""
    if c.loaderMethod != nil {
        data = c.loaderMethod()
    }
    return c.blockRule.check(data)
}

type AndCondition struct {
    conditions []Condition
}

func (c AndCondition) Check() (bool, error) {
    for _, cond := range c.conditions {
        res, err := cond.Check() 
        if err != nil {
            return false, err
        }
        if !res {
            return false, nil
        }
    }
    return true, nil
}

type OrCondition struct {
    conditions []Condition
}

func (c OrCondition) Check() (bool, error) {
    for _, cond := range c.conditions {
        res, err := cond.Check()
        if err != nil {
            return false, err
        }
        if res {
            return true, nil
        }
    }
    return false, nil
}

func (xPathFilter *XPathFilter) UnmarshalJSON(data []byte) error {
    switch data[0] {
    case '"':
        var expression string
        if err := json.Unmarshal(data, &expression); err != nil {
            return err
        }
        xPathFilter.Expression = expression
    case '{':
        var fieldsData interface{}
        if err := json.Unmarshal(data, &fieldsData); err != nil {
            return err
        }

        var objmap map[string]json.RawMessage
        err := json.Unmarshal(data, &objmap)
        if err != nil {
            return err
        }

        fields := fieldsData.(map[string]interface{})
        expression, ok := fields["expression"].(string)
        if ok {
            xPathFilter.Expression = expression
        }
        equalToJson, ok := fields["equalToJson"].(string)
        if ok {
            xPathFilter.EqualToJson = &equalToJson
        }
        equalToXml, ok := fields["equalToXml"].(string)
        if ok {
            xPathFilter.EqualToXml = &equalToXml
        }
        contains, ok := fields["contains"].(string)
        if ok {
            xPathFilter.Contains = &contains
        }
        xPathNamespacesBytes, ok := objmap["xPathNamespaces"]
        if ok {
            xPathNamespaces := make(map[string]string)
            if err := json.Unmarshal(xPathNamespacesBytes, &xPathNamespaces); err != nil {
                return err
            }
            xPathFilter.XPathNamespaces = xPathNamespaces
        }
    }
    return nil
}

func (xPathFilter *XPathFilter) UnmarshalBSON(data []byte) error {
    switch data[0] {
    case '"':
        var expression string
        if err := bson.Unmarshal(data, &expression); err != nil {
            return err
        }
        xPathFilter.Expression = expression
    case '{':
        var fieldsData interface{}
        if err := bson.Unmarshal(data, &fieldsData); err != nil {
            return err
        }

        var objmap map[string]bson.Raw
        err := json.Unmarshal(data, &objmap)
        if err != nil {
            return err
        }

        fields := fieldsData.(map[string]interface{})
        expression, ok := fields["expression"].(string)
        if ok {
            xPathFilter.Expression = expression
        }
        equalToJson, ok := fields["equalToJson"].(string)
        if ok {
            xPathFilter.EqualToJson = &equalToJson
        }
        equalToXml, ok := fields["equalToXml"].(string)
        if ok {
            xPathFilter.EqualToXml = &equalToXml
        }
        contains, ok := fields["contains"].(string)
        if ok {
            xPathFilter.Contains = &contains
        }
        xPathNamespacesBytes, ok := objmap["xPathNamespaces"]
        if ok {
            xPathNamespaces := make(map[string]string)
            if err := bson.Unmarshal(xPathNamespacesBytes, &xPathNamespaces); err != nil {
                return err
            }
            xPathFilter.XPathNamespaces = xPathNamespaces
        }
    }
    return nil
}