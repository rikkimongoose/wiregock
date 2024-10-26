package wiregock

import (
	"testing"
	"regexp"
	"time"
)

func TestEqualToRuleCheck(t *testing.T) {
	ruleCaseSensitive := EqualToRule{"test", false}
	res, err := ruleCaseSensitive.check("test")
	if err != nil || !res {
        t.Fatalf(`EqualToRule failed checking: test`)
    }
	ruleCaseInsensitive := EqualToRule{"test", true}
	res, err = ruleCaseInsensitive.check("tEst")
    if err != nil || !res {
        t.Fatalf(`EqualToRule failed checking: tEst`)
    }
}

func TestEqualToBinaryRuleCheck(t *testing.T) {
	rule := EqualToBinaryRule{[]byte("test")}
	res, err := rule.check("test")
	if err != nil || !res {
        t.Fatalf(`EqualToBinaryRule failed checking: test`)
    }
}

func TestNotRuleRuleCheck(t *testing.T) {
	ruleNotTrue := NotRule{TrueRule{}}
	res, err := ruleNotTrue.check("test")
	if err != nil || res {
        t.Fatalf(`NotRule(TrueRule) failed checking`)
    }
	ruleNotFalse := NotRule{FalseRule{}}
	res, err = ruleNotFalse.check("test")
	if err != nil || !res {
        t.Fatalf(`NotRule(FalseRule) failed checking`)
    }
}

func TestContainsRuleCheck(t *testing.T) {
	ruleCaseSensitive := ContainsRule{"test", false}
	res, err := ruleCaseSensitive.check("testing")
	if err != nil || !res {
        t.Fatalf(`ContainsRule failed checking: test`)
    }
	ruleCaseInsensitive := ContainsRule{"test", true}
	res, err = ruleCaseInsensitive.check("tEsting")
    if err != nil || !res {
        t.Fatalf(`ContainsRule failed checking: tEsting`)
    }
}

func TestDateTimeRuleCheck(t *testing.T) {
	sourceData := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

    ruleEqualToDateTime := DateTimeRule{equalToDateTime: &sourceData, timeFormat: time.RFC3339}
	res, err := ruleEqualToDateTime.check("2009-11-10T23:00:00Z")
	if err != nil || !res {
        t.Fatalf(`Rule doesn't check that %s is equal %s Error: %s`, "2009-11-10T23:00:00Z", sourceData, err)
    }
	ruleBefore := DateTimeRule{before: &sourceData, timeFormat: time.RFC3339}
	res, err = ruleBefore.check("2009-11-09T23:00:00Z")
	if err != nil || !res {
        t.Fatalf(`Rule doesn't check that %s is before %s. Error: %s`, "2009-11-09T23:00:00Z", sourceData, err)
    }
    ruleAfter := DateTimeRule{after: &sourceData, timeFormat: time.RFC3339}
	res, err = ruleAfter.check("2009-11-11T23:00:00Z")
	if err != nil || !res {
        t.Fatalf(`Rule doesn't check that %s is after %s Error: %s`, "2009-11-11T23:00:00Z", sourceData, err)
    }
}

func TestAbsentRuleCheck(t *testing.T) {
	absentRule := AbsentRule{}
	res, err := absentRule.check("test")
	if err != nil && res {
        t.Fatalf(`Absent rule catches existing. Error: %s`, err)
    }
	res, err = absentRule.check("")
	if err != nil && !res {
        t.Fatalf(`Absent rule doesn't catch non-existing. Error: %s`, err)
    }
}

func TestWildcardsRuleCheck(t *testing.T) {
	checkWildcard("test", "test", false, t)
	checkWildcard("?a*da*d.?*", "daaadabadmanda", false, t)
	checkWildcard("?a*da*d.?*", "DaaadAbadmanda", true, t)
}

func checkWildcard(wildcard string, value string, caseInsensitive bool, t *testing.T) {
	ruleWildcards := WildcardsRule{wildcard, caseInsensitive}
	res, err := ruleWildcards.check(value)
	if err != nil || !res {
        t.Fatalf(`WildcardsRule %s failed checking: %s`, wildcard, value)
    }
}

func TestRegExRuleCheck(t *testing.T) {
	regEx := `00-[a-f\d]{32}-[a-f\d]{16}-01`
	value := "/00-0af7651916cd43dd8448eb211c80319c-b9c7c989f97918e1-01/"
	ruleRegEx := RegExRule{regexp.MustCompile(regEx)}
	res, err := ruleRegEx.check(value)
	if err != nil || !res {
        t.Fatalf(`RegExRule %s failed checking: %s`, regEx, value)
    }
}

//TODO - UnitTests for MatchesJsonXPathRule, MatchesXmlXPathRule, EqualToJsonRule, EqualToXmlRule