package wiregock

import (
    "regexp"
    "strings"
    "bytes"
    "reflect"
    "time"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/IGLOU-EU/go-wildcard/v2"
	"github.com/antchfx/xpath"
)

type Rule interface {
    check(str string) (bool, error)
}

type NotRule struct {
    base *Rule
}

type EqualToRule struct {
    val string
    caseInsensitive bool
}

type EqualToBinaryRule struct {
    val []byte
}

type DateTimeRule struct {
    before          *time.Time
    after           *time.Time
    equalToDateTime *time.Time
    timeFormat string //default: time.RFC3339
}

type ContainsRule struct {
    val string
    caseInsensitive bool
}

type WildcardsRule struct {
    val string
    caseInsensitive bool
}
type RegExRule struct {
    regex *regexp.Regexp
}

type MatchesXPathRule struct {
    xPath *xpath.Expr
}

type EqualToXmlRule struct {
    node *xmlquery.Node
}

type EqualToJsonRule struct {
    node *jsonquery.Node
}

type TrueRule struct {
}

type FalseRule struct {
}

func (rule NotRule) check(str string) (bool, error) {
	res, err := (*rule.base).check(str)
	return !res, err
}

func (rule EqualToRule) check(str string) (bool, error) {
	if rule.caseInsensitive {
		return strings.EqualFold(rule.val, str), nil
	}
	return strings.Compare(rule.val, str) == 0, nil
}

func (rule EqualToBinaryRule) check(str string) (bool, error) {
	return bytes.Compare(rule.val, []byte(str)) == 0, nil
}

func (rule DateTimeRule) check(str string) (bool, error) {
	sourceTime, error := time.Parse(rule.timeFormat, str)
	if error != nil {
		return false, error
	}
	if rule.equalToDateTime != nil && !sourceTime.Equal(*rule.equalToDateTime) {
		return false, nil
	}
	if rule.before != nil && !sourceTime.Before(*rule.before) {
		return false, nil
	}
	if rule.after != nil && !sourceTime.After(*rule.after) {
		return false, nil
	}
	return true, nil
}

func (rule ContainsRule) check(str string) (bool, error) {
	if rule.caseInsensitive {
		return strings.Contains(strings.ToLower(str), strings.ToLower(rule.val)), nil
	}
	return strings.Contains(str, rule.val), nil
}

func (rule WildcardsRule) check(str string) (bool, error) {
	if rule.caseInsensitive {
		return wildcard.Match(strings.ToLower(rule.val), strings.ToLower(str)), nil
	}
	return wildcard.Match(rule.val, str), nil
}

func (rule RegExRule) check(str string) (bool, error) {
	return rule.regex.MatchString(str), nil
}

func (rule EqualToXmlRule) check(str string) (bool, error) {
	node, err := xmlquery.Parse(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(&rule.node, node), nil
}

func (rule EqualToJsonRule) check(str string) (bool, error) {
	node, err := jsonquery.Parse(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(&rule.node, node), nil
}

func (rule MatchesXPathRule) check(str string) (bool, error) {
	node, err := jsonquery.Parse(strings.NewReader(str))
	if err != nil {
		node, err := xmlquery.Parse(strings.NewReader(str))
		if err != nil {
			return false, err
		}
		return (xmlquery.QuerySelector(node, rule.xPath) != nil), nil
	}
	return (jsonquery.QuerySelector(node, rule.xPath) != nil), nil
}

func (rule TrueRule) check(str string) (bool, error) {
	return true, nil
}

func (rule FalseRule) check(str string) (bool, error) {
	return false, nil
}
