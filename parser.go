package wiregock

import (
    "time"
    "regexp"
    "strings"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
)

type MultipartFormData struct {
	Headers map[string]string
	Data []byte
}

type DataContext struct {
	Body func() string
	Get func(key string) string
	GetMulti func(key string) []string
	Params func(key string) string
	ParamsMulti func(key string) []string
	Cookies func(key string) string
	FormValue func(key string) string
	MultipartForm func() func (yield func(MultipartFormData) bool)
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

	if request.FormParameters != nil {
		for key, value := range request.FormParameters {
			newCondition, err := createCondition(&value, func() string { return context.FormValue(key) })
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

func createCondition(filter *Filter, loaderMethod func() string) (*DataCondition, error) {
	rules, err := parseRules(filter, true)
	return &DataCondition{loaderMethod, rules}, err
}

type XPathFilterProps struct {
	caseInsensitive bool
	ignoreArrayOrder bool
	ignoreExtraElements bool
}

func loadFilterProps(filter *Filter) XPathFilterProps {
	caseInsensitive := false
	if filter.CaseInsensitive != nil {
		caseInsensitive = *filter.CaseInsensitive
	}
	ignoreArrayOrder := true
    if filter.IgnoreArrayOrder != nil {
    	ignoreArrayOrder = *filter.IgnoreArrayOrder
    }
    ignoreExtraElements := true
    if filter.IgnoreExtraElements != nil {
    	ignoreExtraElements = *filter.IgnoreExtraElements
    }
    return XPathFilterProps{
    	caseInsensitive,
		ignoreArrayOrder,
		ignoreExtraElements,
    }
}

func loadXPathFilterProps(filter *XPathFilter, xPathFilterPropsDefault *XPathFilterProps) XPathFilterProps {
	caseInsensitive := false
	if xPathFilterPropsDefault != nil {
		caseInsensitive = xPathFilterPropsDefault.caseInsensitive
	}
	if filter.CaseInsensitive != nil {
		caseInsensitive = *filter.CaseInsensitive
	}
	ignoreArrayOrder := true
	if xPathFilterPropsDefault != nil {
		ignoreArrayOrder = xPathFilterPropsDefault.ignoreArrayOrder
	}
    if filter.IgnoreArrayOrder != nil {
    	ignoreArrayOrder = *filter.IgnoreArrayOrder
    }
    ignoreExtraElements := true
	if xPathFilterPropsDefault != nil {
		ignoreExtraElements = xPathFilterPropsDefault.ignoreExtraElements
	}
    if filter.IgnoreExtraElements != nil {
    	ignoreExtraElements = *filter.IgnoreExtraElements
    }
    return XPathFilterProps{
    	caseInsensitive,
		ignoreArrayOrder,
		ignoreExtraElements,
    }
}

func loadXPathFilterRules(filterPath *XPathFilter, caseInsensitive bool) []Rule {
	xPathRules := []Rule{}
	if filterPath.EqualTo != nil {
		xPathRules = append(xPathRules, EqualToRule{*filterPath.EqualTo, caseInsensitive})
	}
	if filterPath.Contains != nil {
		xPathRules = append(xPathRules, ContainsRule{*filterPath.Contains, caseInsensitive})
	}
	actualFormat := time.RFC3339
	if filterPath.Before != nil || filterPath.After != nil || filterPath.EqualToDateTime != nil {
		if filterPath.ActualFormat != nil {
			actualFormat = *filterPath.ActualFormat
		}
		xPathRules = append(xPathRules, DateTimeRule{filterPath.Before, filterPath.After, filterPath.EqualToDateTime, actualFormat})
	}
	return xPathRules
}

type XPathFactory interface {
    generateXPathRule(query string, xPathFilterProps *XPathFilterProps) (Rule, error)
    generateMatchesXPathRule(filterPath *XPathFilter, xPathFilterPropsDefault *XPathFilterProps) (Rule, *XPathFilterProps, error)
}

type XPathJsonFactory struct {}

func (xPathFactory XPathJsonFactory) generateXPathRule(query string, xPathFilterProps *XPathFilterProps) (Rule, error) {
	node, err := jsonquery.Parse(strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	rule := EqualToJsonRule{node: node}
	rule.IgnoreArrayOrder = xPathFilterProps.ignoreArrayOrder
	rule.IgnoreExtraElements = xPathFilterProps.ignoreExtraElements
	return rule, err
}

func (xPathFactory XPathJsonFactory) generateMatchesXPathRule(filterPath *XPathFilter, xPathFilterPropsDefault *XPathFilterProps) (Rule, error) {
	xPath, err := generateXPath(filterPath.Expression, filterPath.XPathNamespaces)
	if err != nil {
		return nil, err
	}
	xPathFilterPropsLocal := loadXPathFilterProps(filterPath, xPathFilterPropsDefault)
	xPathRules := loadXPathFilterRules(filterPath, xPathFilterPropsLocal.caseInsensitive)
	if filterPath.EqualToJson != nil {
		ruleJson, err := xPathFactory.generateXPathRule(*filterPath.EqualToJson, &xPathFilterPropsLocal)
		if err != nil {
			return nil, err
		}
    	xPathRules = append(xPathRules, ruleJson)
	}
	if filterPath.EqualToXml != nil {
		ruleXml, err := xPathFactory.generateXPathRule(*filterPath.EqualToXml, &xPathFilterPropsLocal)
		if err != nil {
			return nil, err
		}
		xPathRules = append(xPathRules, ruleXml)
	}
	for _, filterPathSub := range filterPath.And {
		ruleSub, err := xPathFactory.generateMatchesXPathRule(&filterPathSub, &xPathFilterPropsLocal)
		if err != nil {
			return nil, err
		}
		xPathRules = append(xPathRules, ruleSub)
	}
	rule := MatchesJsonXPathRule{}
	rule.xPath = xPath
	rule.innerRule = BlockRule{rulesOr: xPathRules}
	return rule, nil
}

type XPathXmlFactory struct {}

func (xPathFactory XPathXmlFactory) generateXPathRule(query string, xPathFilterProps *XPathFilterProps) (Rule, error) {
	node, err := xmlquery.Parse(strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	rule := EqualToXmlRule{node: node}
	rule.IgnoreArrayOrder = xPathFilterProps.ignoreArrayOrder
	rule.IgnoreExtraElements = xPathFilterProps.ignoreExtraElements
	return rule, err
}

func (xPathFactory XPathXmlFactory) generateMatchesXPathRule(filterPath *XPathFilter, xPathFilterPropsDefault *XPathFilterProps) (Rule, error) {
	xPath, err := generateXPath(filterPath.Expression, filterPath.XPathNamespaces)
	if err != nil {
		return nil, err
	}
	xPathFilterPropsLocal := loadXPathFilterProps(filterPath, xPathFilterPropsDefault)
	xPathRules := loadXPathFilterRules(filterPath, xPathFilterPropsLocal.caseInsensitive)
	if filterPath.EqualToJson != nil {
		ruleJson, err := xPathFactory.generateXPathRule(*filterPath.EqualToJson, &xPathFilterPropsLocal)
		if err != nil {
			return nil, err
		}
    	xPathRules = append(xPathRules, ruleJson)
	}
	if filterPath.EqualToXml != nil {
		ruleXml, err := xPathFactory.generateXPathRule(*filterPath.EqualToXml, &xPathFilterPropsLocal)
		if err != nil {
			return nil, err
		}
		xPathRules = append(xPathRules, ruleXml)
	}
	for _, filterPathSub := range filterPath.And {
		ruleSub, err := xPathFactory.generateMatchesXPathRule(&filterPathSub, &xPathFilterPropsLocal)
		if err != nil {
			return nil, err
		}
		xPathRules = append(xPathRules, ruleSub)
	}
	rule := MatchesXmlXPathRule{}
	rule.xPath = xPath
	rule.innerRule = BlockRule{rulesOr: xPathRules}
	return rule, nil
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
	xPathJsonFactory := XPathJsonFactory{}
	xPathXmlFactory := XPathXmlFactory{}

	rules := []Rule{}

	xPathFilterProps := loadFilterProps(filter)
	caseInsensitive := xPathFilterProps.caseInsensitive

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
		rule, err := xPathJsonFactory.generateXPathRule(*filter.EqualToJson, &xPathFilterProps)
		if err != nil {
			return nil, err
		}
    	rules = append(rules, rule)
	}
	if filter.EqualToXml != nil {
		rule, err := xPathXmlFactory.generateXPathRule(*filter.EqualToXml, &xPathFilterProps)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if filter.MatchesJsonPath != nil {
		rule, err := xPathJsonFactory.generateMatchesXPathRule(filter.MatchesJsonPath, &xPathFilterProps)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	if filter.MatchesXPath != nil {
		rule, err := xPathXmlFactory.generateMatchesXPathRule(filter.MatchesXPath, &xPathFilterProps)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}