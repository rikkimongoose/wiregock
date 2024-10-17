package wiregock

import (
    "time"
    "regexp"
    "strings"
)

type DataContext struct {
	Body func() string
	Get func(key string) string
	Params func(key string) string
	Cookies func(key string) string
}

func parseCondition(request *MockRequest, context *DataContext) (Condition, error) {
	conditions := []Condition{}
	if request.Headers != nil {
		for key, value := range request.Headers {
			newCondition, err := createCondition(value, func() string { return context.Get(key) })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, newCondition)
		}
	}

	if request.QueryParameters != nil {
		for key, value := range request.QueryParameters {
			newCondition, err := createCondition(value, func() string { return context.Params(key) })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, newCondition)
		}
	}

	if request.Cookies != nil {
		for key, value := range request.Cookies {
			newCondition, err := createCondition(value, func() string { return context.Cookies(key) })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, newCondition)
		}
	}

	if len(request.BodyPatterns) > 0 {
		for _, value := range request.BodyPatterns {
			newCondition, err := createCondition(value, func() string { return context.Body() })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, newCondition)
		}
	}
	if strings.Compare(*request.Method, "ANY") != 0 {
		return AndCondition{conditions}, nil
	}
	return OrCondition{conditions}, nil
}

func createCondition(filter Filter, loaderMethod func() string) (DataCondition, error) {
	rulesAnd, rulesOr, err := parseRules(&filter)
	return DataCondition{loaderMethod, rulesAnd, rulesOr}, err
}

func parseRules(filter *Filter) ([]Rule, []Rule, error) {
    rulesAnd, err := parseRule(filter)
    if err != nil {
    	return nil, nil, err
    }
    rulesOr := []Rule{}
    if len(filter.And) > 0 {
    	for _, filterAnd := range filter.And {
    		parsedRules, err := parseRule(&filterAnd)
    		if err != nil {
    			return nil, nil, err
    		}
    		rulesAnd = append(rulesAnd, parsedRules...)
    	}
    }
    if len(filter.Or) > 0 {
    	for _, filterOr := range filter.Or {
    		parsedRules, err := parseRule(&filterOr)
    		if err != nil {
    			return nil, nil, err
    		}
    		rulesOr = append(rulesOr, parsedRules...)
    	}
    }
    return rulesAnd, rulesOr, nil
}

func parseRule(filter *Filter) ([]Rule, error) {
	rules := []Rule{}
	caseInsensitive := false
	if filter.CaseInsensitive != nil {
		caseInsensitive = *filter.CaseInsensitive
	}

	if filter.Contains != nil {
		val := *filter.Contains
		rules = append(rules, ContainsRule{val, caseInsensitive})
	}

	if filter.EqualTo != nil {
		val := *filter.EqualTo
		rules = append(rules, EqualToRule{val, caseInsensitive})
	}

	if filter.BinaryEqualTo != nil {
		val := *filter.BinaryEqualTo
		rules = append(rules, EqualToBinaryRule{[]byte(val)})
	}

	if filter.DoesNotContain != nil {
		val := *filter.DoesNotContain
		rules = append(rules, NotRule{ContainsRule{val, caseInsensitive}})
	}

	if filter.Matches != nil {
		regexStr := *filter.Matches
		regex, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, err
		}
		rules = append(rules, RegExRule{regex})
	}

	if filter.DoesNotMatch != nil {
		regexStr := *filter.DoesNotMatch
		regex, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, err
		}
		rules = append(rules, NotRule{RegExRule{regex}})
	}

	if filter.Absent != nil {
		absent := *filter.Absent
		if absent {
			rules = append(rules, AbsentRule{})
		}
	}
	actualFormat := time.RFC3339
	if filter.ActualFormat != nil {
		actualFormat = *filter.ActualFormat
	}
	if filter.Before != nil || filter.After != nil || filter.EqualToDateTime != nil {
		rules = append(rules, DateTimeRule{filter.Before, filter.After, filter.EqualToDateTime, actualFormat})
	}
	return rules, nil
}