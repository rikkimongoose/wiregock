package wiregock

import (
    "encoding/json"
    "encoding/xml"
    "regexp"
    "strings"
    "bytes"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	"github.com/IGLOU-EU/go-wildcard/v2"
)

type Rule interface {
    check(str string) bool, error
}

type EqualToRule struct {
    val string
    caseInsensitive false
}

type EqualToBinaryRule struct {
    val []byte
}

type ContainsRule struct {
    val string
}

type WildcardsRule struct {
    match string
}
type RegExRule struct {
    match *regexp.Regexp
}

type MatchesXPathRule struct {
    xPath *xpath.Expr
}

type EqualToXmlRule struct {
    node *xpath.Node
}

type EqualToJsonRule struct {
    node *xpath.Node
}

func (rule *EqualToRule) check(str string) (bool, error) {
	if rule.caseInsensitive {
		return strings.EqualFold(rule.val, str), nil
	}
	return strings.Compare(rule.val, str), nil
}

func (rule *EqualToBinaryRule) check(str string) (bool, error) {
	return bytes.Compare(rule.val, []byte(str)), nil
}

func (rule *ContainsRule) check(str string) (bool, error) {
	return strings.Contains(str, rule.val), nil
}

func (rule *WildcardsRule) check(str string) (bool, error) {
	return wildcard.Match(str, rule.val), nil
}

func (rule *EqualToXmlRule) check(str string) (bool, error) {
	node, err := xmlquery.Parse(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(&rule.node, data)
}

func (rule *EqualToJsonRule) check(str string) (bool, error) {
	node, err := jsonquery.Parse(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(&rule.node, node)
}

func (rule *MatchesXPathRule) check(str string) (bool, error) {
	node, err := jsonquery.Parse(strings.NewReader(str))
	if err != nil {
		node, err := xmlquery.Parse(strings.NewReader(str))
		if err != nil {
			return false, err
		}
		return xmlquery.QuerySelector(node, rule.xPath) != nil
	}
	return jsonquery.QuerySelector(node, rule.xPath) != nil
}

func (rule *RegExRule) check(str string) (bool, error) {
	return rule.MatchString(str)
}