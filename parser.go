package wiregock

import (
    "time"
    "regexp"
    "strings"
    "mime/multipart"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
)

type DataContext struct {
	Body func() string
	Get func(key string) string
	Params func(key string) string
	Cookies func(key string) string
	MultipartForm func() multipart.Form
}

func ParseCondition(request *MockRequest, context *DataContext) (Condition, error) {
	conditions := []Condition{}
	if request.Headers != nil {
		for key, value := range request.Headers {
			newCondition, err := createCondition(&value, func() string { return context.Get(key) })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, *newCondition)
		}
	}

	if request.QueryParameters != nil {
		for key, value := range request.QueryParameters {
			newCondition, err := createCondition(&value, func() string { return context.Params(key) })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, *newCondition)
		}
	}

	if request.Cookies != nil {
		for key, value := range request.Cookies {
			newCondition, err := createCondition(&value, func() string { return context.Cookies(key) })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, *newCondition)
		}
	}

	if len(request.BodyPatterns) > 0 {
		for _, value := range request.BodyPatterns {
			newCondition, err := createCondition(&value, func() string { return context.Body() })
			if err != nil {
				return nil, err
			}
			conditions = append(conditions, *newCondition)
		}
	}
	return AndCondition{conditions}, nil
}

/*func createMultipartCondition(multipartPatternsData []MultipartPatternsData, loaderMethod func() multipart.Form) ([]MultipartDataCondition, error) {
	//TODO - init multipart condition
	//checkAny := multipartPatternsData.MatchingType == nil || strings.Compare(*multipartPatternsData.MatchingType, "ALL") == 0
	
	//return &MultipartDataCondition{checkAny: checkAny, loaderMethod: loaderMethod}, nil
}*/

func createCondition(filter *Filter, loaderMethod func() string) (*DataCondition, error) {
	rules, err := parseRules(filter, true)
	return &DataCondition{loaderMethod, rules}, err
}

func parseRules(filter *Filter, defaultAnd bool) (*BlockRule, error) {
    rules, err := parseRule(filter)
    if err != nil {
    	return nil, err
    }
    rulesAnd := []Rule{}
    rulesOr := []Rule{}
    if defaultAnd {
    	rulesAnd = append(rulesAnd, rules...)
    } else {
		rulesOr = append(rulesOr, rules...)
    }
    if len(filter.And) > 0 {
    	for _, filterAnd := range filter.And {
    		parsedRules, err := parseRule(&filterAnd)
    		if err != nil {
    			return nil, err
    		}
    		rulesAnd = append(rulesAnd, parsedRules...)
    	}
    }
    if len(filter.Or) > 0 {
    	for _, filterOr := range filter.Or {
    		parsedRules, err := parseRule(&filterOr)
    		if err != nil {
    			return nil, err
    		}
    		rulesOr = append(rulesOr, parsedRules...)
    	}
    }
    return &BlockRule{rulesAnd, rulesOr}, nil
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

	if filter.EqualToJson != nil {
		node, err := jsonquery.Parse(strings.NewReader(*filter.EqualToJson))
		if err != nil {
			return nil, err
		}
		ignoreArrayOrder := true
    	ignoreExtraElements := true
    	if filter.IgnoreArrayOrder != nil {
    		ignoreArrayOrder = *filter.IgnoreArrayOrder
    	}
    	if filter.IgnoreExtraElements != nil {
    		ignoreExtraElements = *filter.IgnoreExtraElements
    	}
    	rules = append(rules, EqualToJsonRule{node, ignoreArrayOrder, ignoreExtraElements})
	}

	if filter.EqualToXml != nil {
		node, err := xmlquery.Parse(strings.NewReader(*filter.EqualToXml))
		if err != nil {
			return nil, err
		}
		ignoreArrayOrder := true
    	ignoreExtraElements := true
    	if filter.IgnoreArrayOrder != nil {
    		ignoreArrayOrder = *filter.IgnoreArrayOrder
    	}
    	if filter.IgnoreExtraElements != nil {
    		ignoreExtraElements = *filter.IgnoreExtraElements
    	}
		rules = append(rules, EqualToXmlRule{node, ignoreArrayOrder, ignoreExtraElements})
	}

	if filter.MatchesJsonPath != nil {
		xPath, err := generateXPath(filter.MatchesJsonPath.Expression, filter.MatchesJsonPath.XPathNamespaces)
		if err != nil {
			return nil, err
		}
		rule := MatchesJsonXPathRule{xPath: xPath}
		//TODO - inner rule
		if filter.MatchesJsonPath.Contains != nil {
			rule.contains = filter.MatchesJsonPath.Contains
		}
		if filter.MatchesJsonPath.EqualToJson != nil {
			node, err := jsonquery.Parse(strings.NewReader(*filter.MatchesJsonPath.EqualToXml))
			if err != nil {
				return nil, err
			}
			rule.node = node
		}
		rules = append(rules, rule)
	}

	if filter.MatchesXPath != nil {
		xPath, err := generateXPath(filter.MatchesXPath.Expression, filter.MatchesXPath.XPathNamespaces)
		if err != nil {
			return nil, err
		}
		rule := MatchesXmlXPathRule{xPath: xPath}
		//TODO - inner rule
		if filter.MatchesXPath.Contains != nil {
			rule.contains = filter.MatchesXPath.Contains
		}
		if filter.MatchesXPath.EqualToXml != nil {
			node, err := xmlquery.Parse(strings.NewReader(*filter.MatchesXPath.EqualToXml))
			if err != nil {
				return nil, err
			}
			rule.node = node
		}
		rules = append(rules, rule)
	}

	return rules, nil
}