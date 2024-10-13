package wiregock

import (
    "strings"
)

type WebContext interface {
	Body() []byte
	Get(key string, defaultValue ...string) string
	Params(key string, defaultValue ...string) string
	Cookies(key string, defaultValue ...string) string
}

func loaderGet(context WebContext, key string) func() (string, bool) {
	return func() (string, bool) {
		data := context.Get(key, "")
		return data, strings.EqualFold(data, "")
	}
}

func loaderParams(context WebContext, key string) func() (string, bool) {
	return func() (string, bool) {
		data := context.Params(key, "")
		return data, strings.EqualFold(data, "")
	}
}

func loaderCookies(context WebContext, key string) func() (string, bool) {
	return func() (string, bool) {
		data := context.Cookies(key, "")
		return data, strings.EqualFold(data, "")
	}
}

func loaderBody(context WebContext) func() (string, bool) {
	return func() (string, bool) {
		body := string(context.Body()[:])
		return body, len(body) > 0
	}
}

func parseCondition(request *MockRequest, context WebContext) Condition {
	conditions := []Condition{}
	if request.Headers != nil {
		for key, value := range request.Headers {
			conditions = append(conditions, createConditions(value, loaderGet(context, key)))
		}
	}

	if request.QueryParameters != nil {
		for key, value := range request.QueryParameters {
			conditions = append(conditions, createConditions(value, loaderParams(context, key)))
		}
	}

	if request.Cookies != nil {
		for key, value := range request.Cookies {
			conditions = append(conditions, createConditions(value, loaderCookies(context, key)))
		}
	}

	if len(request.BodyPatterns) > 0 {
		for _, value := range request.BodyPatterns {
			conditions = append(conditions, createConditions(value, loaderBody(context)))
		}
	}
	return &AndCondition{conditions}
}

type ConditionRules struct {
    rulesAnd []Rule
    rulesOr []Rule
}

func createConditions(filter Filter, loaderMethod func() (string, bool)) DataCondition {
	rulesAnd, rulesOr := parseRules(&filter)
	return DataCondition{loaderMethod, rulesAnd, rulesOr}
}

func parseRules(filter *Filter) ([]Rule, []Rule) {
    rulesAnd := []Rule{parseRule(filter)}
    rulesOr := []Rule{}
    if len(filter.And) > 0 {
    	for _, filterAnd := range filter.And {
    		rulesAnd = append(rulesAnd, parseRule(&filterAnd))
    	}
    }
    if len(filter.Or) > 0 {
    	for _, filterOr := range filter.Or {
    		rulesOr = append(rulesOr, parseRule(&filterOr))
    	}
    }
    return rulesAnd, rulesOr
}

func parseRule(filter *Filter) Rule {
	return TrueRule{}
}