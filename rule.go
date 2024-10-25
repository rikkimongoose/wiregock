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
    base Rule
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

type MatchesJsonXPathRule struct {
	xPath *xpath.Expr
    innerRule Rule
}

type MatchesXmlXPathRule struct {
	xPath *xpath.Expr
    innerRule Rule
}

type EqualToXmlRule struct {
    node *xmlquery.Node
    ignoreArrayOrder bool
    ignoreExtraElements bool
}

type EqualToJsonRule struct {
    node *jsonquery.Node
    ignoreArrayOrder bool
    ignoreExtraElements bool
}

type AbsentRule struct {
}

type TrueRule struct {
}

type FalseRule struct {
}

type BlockRule struct {
    rulesAnd []Rule
    rulesOr []Rule
}

func (rule NotRule) check(str string) (bool, error) {
	res, err := rule.base.check(str)
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
	//TODO - comparasion with ignoreArrayOrder, ignoreExtraElements
	return reflect.DeepEqual(&rule.node, node), nil
}

func (rule EqualToJsonRule) check(str string) (bool, error) {
	node, err := jsonquery.Parse(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	return reflect.DeepEqual(&rule.node, node), nil
}

func (rule MatchesXmlXPathRule) check(str string) (bool, error) {
	node, err := xmlquery.Parse(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	if rule.innerRule != nil {
		nodesByXPath := xmlquery.QuerySelectorAll(node, rule.xPath)
		for _, node := range nodesByXPath {
			ok, err := rule.innerRule.check(node.Data)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
	return (xmlquery.QuerySelector(node, rule.xPath) != nil), nil
}

func (rule MatchesJsonXPathRule) check(str string) (bool, error) {
	node, err := jsonquery.Parse(strings.NewReader(str))
	if err != nil {
		return false, err
	}
	if rule.innerRule != nil {
		nodesByXPath := jsonquery.QuerySelectorAll(node, rule.xPath)
		for _, node := range nodesByXPath {
			ok, err := rule.innerRule.check(node.Data)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
	return (jsonquery.QuerySelector(node, rule.xPath) != nil), nil
}

func (rule AbsentRule) check(str string) (bool, error) {
	return len(str) == 0, nil
}

func (rule TrueRule) check(str string) (bool, error) {
	return true, nil
}

func (rule FalseRule) check(str string) (bool, error) {
	return false, nil
}

func (rule BlockRule) check(str string) (bool, error) {
	if rule.rulesAnd != nil {
		for _, ruleAnd := range rule.rulesAnd {
	        res, err := ruleAnd.check(str)
	        if err != nil {
	            return false, err
	        }
	        if !res {
	            return false, nil
	        }
	    }
	}
	if rule.rulesOr != nil {
	    for _, ruleOr := range rule.rulesOr {
	        res, err := ruleOr.check(str)
	        if err != nil {
	            return false, err
	        }
	        if res {
	            return true, nil
	        }
	    }
	}
	return rule.rulesAnd == nil || rule.rulesOr == nil || len(rule.rulesOr) == 0, nil
}


func generateXPath(str string, namespaces map[string]string) (*xpath.Expr, error) {
	if namespaces != nil {
		return xpath.CompileWithNS(str, namespaces)
	}
	return xpath.Compile(str)
}