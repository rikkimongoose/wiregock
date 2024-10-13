package wiregock

import (
    "testing"
    "regexp"
)

// TestGenerateTraceparent calls generateTraceparent, checking
// for a valid return value.
func TestGenerateTraceparent(t *testing.T) {
	traceRegex := regexp.MustCompile(`^00-[a-f\d]{32}-[a-f\d]{16}-01$`)
    msg := generateTraceparent()
    if !traceRegex.MatchString(msg) {
        t.Fatalf(`Wrong traceparent id: %s`, msg)
    }
}