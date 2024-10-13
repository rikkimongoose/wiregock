package wiregock

import (
	"testing"
	"regexp"
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