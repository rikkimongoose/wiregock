package wiregock

import (
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
    "go.mongodb.org/mongo-driver/bson"
    "strings"
    "regexp"
)

const (
    KEY_HEADERS = "headers"
    KEY_QUERY = "queryParameters"
    KEY_COOKIES = "cookies"
    KEY_BODY = "bodyPatterns"
    KEY_MATCHING_TYPE = "matchingType"
    KEY_MULTIPART_HEADERS = "multipartPatterns" //TODO - will be supported later

    CMD_EQUAL = "equalTo"
    CMD_EQUAL_BINARY = "binaryEqualTo" // TODO - реализовать
    CMD_CONTAINS = "contains"
    CMD_MATCHES = "matches"
    CMD_MATCHES_WILDCARDS = "wildcards"
    CMD_MATCHES_JSON = "equalToJson"
    CMD_MATCHES_XML = "equalToXml"
    CMD_XPATH = "matchesXPath"

    MATCHING_ANY = "ANY"
    MATCHING_ALL = "ALL"

    MATCHING_DEFAULT = MATCHING_ANY
)

func parseCondition(m bson.M) *Condition {
	conditions := []Condition{}
	headers, ok = m[KEY_HEADERS]
	if ok {
		for key, val := range headers {
			rule = parseRule(key, val)
			if rule == nil {
				logger.Warn("Wrong header")
				continue
			}
			append(conditions, DataCondition {
				Prop: CONDITION_HEADER
				Key: key
			    Rule: rule
			})
		}
	}

	queries, ok = m[KEY_QUERY]
	if ok {
		for key, val := range queries {
			rule = parseRule(key, val)
			if rule == nil {
				logger.Warn("Wrong query")
				continue
			}
			append(conditions, DataCondition {
				Prop: CONDITION_PARAMS
				Key: key
			    Rule: rule
			})
		}
	}

	cookies, ok = m[KEY_COOKIES]
	if ok {
		for key, val := range cookies {
			rule = parseRule(key, val)
			if rule == nil {
				logger.Warn("Wrong cookie")
				continue
			}
			append(conditions, DataCondition {
				Prop: CONDITION_COOKIE
				Key: key
			    Rule: rule
			})
		}
	}

	body, ok = m[KEY_BODY]
	if ok {
		for _, val := range body {
			conditionsBody := []Condition{}
			for keyBody, valueBody := range val {
				rule = parseRule(key, val)
				if rule == nil {
					logger.Warn("Wrong cookie")
					continue
				}
				append(conditionsBody, DataCondition {
					Prop: CONDITION_BODY
					Key: nil
				    Rule: rule
				})
			}
			append(conditions, AndCondition {conditionsBody})
		}
	}

	matchingType, ok = m[KEY_MATCHING_TYPE]
	if !ok {
	   matchingType = MATCHING_DEFAULT
	}
	if strings.Compare(matchingType, MATCHING_ALL) {
		return &AndCondition{conditions}
	}
	return &OrCondition{conditions}
}

func parseRule(cmd string, value string) *Rule {
	if strings.Compare(cmd, CMD_EQUAL) {
		return &EqualToRule{value}
	}
	if strings.Compare(cmd, CMD_EQUAL_BINARY) {
		return &EqualToBinaryRule{byte[](value)}
	}
	if strings.Compare(cmd, CMD_CONTAINS) {
		regex, err := regexp.Compile(value)
		if err != nil {
			logger.Error(err)
			return nil
		}
		return &ContainsRule{regex}
	}
	if strings.Compare(cmd, CMD_MATCHES) {
		return &RegExRule{value}
	}
	if strings.Compare(cmd, CMD_MATCHES_WILDCARDS) {
		return &WildcardsRule{value}
	}
	if strings.Compare(cmd, CMD_XPATH) {
		expr, err := xpath.Compile(value)
		if err != nil {
			logger.Error(err)
			return nil
		}
		return &MatchesXPathRule{expr}
	}
	if strings.Compare(cmd, CMD_MATCHES_XML) {
		node, err := xmlquery.Parse(strings.NewReader(value))
		if err != nil {
			logger.Error(err)
			return nil
		}
		return &EqualToXmlRule{value}
	}
	if strings.Compare(cmd, CMD_MATCHES_JSON) {
		node, err := jsonquery.Parse(strings.NewReader(value))
		if err != nil {
			logger.Error(err)
			return nil
		}
		return &EqualToJsonRule{value}
	}
	return nil
}