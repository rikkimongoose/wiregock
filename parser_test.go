package wiregock

import (
	"testing"
	"time"
	"reflect"
	"regexp"
)

func TestParseRule(t *testing.T) {
	Contains := "Contains"
    EqualTo := "EqualTo"
    CaseInsensitive := false
    BinaryEqualTo := "BinaryEqualTo"
    DoesNotContain := "DoesNotContain"
    Matches := ".*"
    DoesNotMatch := ".*"
    Absent := true
	Before := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
    After := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
    EqualToDateTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
    ActualFormat := "ActualFormat"
    EqualToJson := "{ \"total_results\": 4 }"
    IgnoreArrayOrder := true
    IgnoreExtraElements := true
    MatchesJsonPath := "MatchesJsonPath"
    EqualToXml := "<thing>Hello</thing>"
    MatchesXPath := "MatchesXPath"

	filter := Filter {
	    Contains: &Contains,
	    EqualTo: &EqualTo,
	    CaseInsensitive: &CaseInsensitive,
	    BinaryEqualTo: &BinaryEqualTo,
	    DoesNotContain: &DoesNotContain,
	    Matches: &Matches,
	    DoesNotMatch: &DoesNotMatch,
	    Absent: &Absent,
	    Before: &Before,
	    After: &After,
	    EqualToDateTime: &EqualToDateTime,
	    ActualFormat: &ActualFormat,
	    EqualToJson: &EqualToJson,
	    IgnoreArrayOrder: &IgnoreArrayOrder,
	    IgnoreExtraElements: &IgnoreExtraElements,
	    MatchesJsonPath: &XPathFilter{Expression:MatchesJsonPath},
	    EqualToXml: &EqualToXml,
	    MatchesXPath: &XPathFilter{Expression:MatchesXPath},
	}

	rules, err := parseRule(&filter)
	if err != nil {
        t.Fatalf(`Error parsing rule: %s`, err)
	}
	if !reflect.DeepEqual(rules[0], ContainsRule{Contains, CaseInsensitive}) {
        t.Fatalf(`Error parsing: %s`, "ContainsRule")
	}
	if !reflect.DeepEqual(rules[1], EqualToRule{EqualTo, CaseInsensitive}) {
        t.Fatalf(`Error parsing: %s`, "EqualToRule")
	}
	if !reflect.DeepEqual(rules[2], EqualToBinaryRule{[]byte(BinaryEqualTo)}) {
        t.Fatalf(`Error parsing: %s`, "EqualToBinaryRule")
	}
	if !reflect.DeepEqual(rules[3], NotRule{ContainsRule{DoesNotContain, CaseInsensitive}}) {
        t.Fatalf(`Error parsing: %s`, "NotRule.ContainsRule")
	}
	if !reflect.DeepEqual(rules[4], RegExRule{regexp.MustCompile(Matches)}) {
        t.Fatalf(`Error parsing: %s`, "RegExRule")
	}
	if !reflect.DeepEqual(rules[5], NotRule{RegExRule{regexp.MustCompile(Matches)}}) {
        t.Fatalf(`Error parsing: %s`, "NotRule.RegExRule")
	}
	if !reflect.DeepEqual(rules[6], AbsentRule{}) {
        t.Fatalf(`Error parsing: %s`, "AbsentRule")
	}
	if !reflect.DeepEqual(rules[7], DateTimeRule{&Before, &After, &EqualToDateTime, ActualFormat}) {
        t.Fatalf(`Error parsing: %s`, "DateTimeRule")
	}
}

func TestParseRules(t *testing.T) {
	Contains := "Contains"
	CaseInsensitive := false
	filter := Filter{
	    Contains: &Contains,
	    And: []Filter{Filter { Contains: &Contains }},
	    Or: []Filter{Filter { Contains: &Contains }},
	}

	rules, err := parseRules(&filter, true)
	rulesAnd := rules.rulesAnd
	rulesOr := rules.rulesOr
	rule := ContainsRule{Contains, CaseInsensitive}
	if err != nil {
        t.Fatalf(`Error parsing rules: %s`, err)
	}
	if !reflect.DeepEqual(rulesAnd[0], rule) {
        t.Fatalf(`Error parsing rule And: %s`, rulesAnd[0])
	}
	if !reflect.DeepEqual(rulesAnd[1], rule) {
        t.Fatalf(`Error parsing rule And: %s`, rulesAnd[1])
	}
	if !reflect.DeepEqual(rulesOr[0], rule) {
        t.Fatalf(`Error parsing rule Or: %s`, rulesOr[0])
	}
}