package wiregock

import (
    "strings"
)

const (
    CONDITION_BODY = "Body"
    CONDITION_HEADER = "Header"
    CONDITION_PARAMS = "Params"
    CONDITION_COOKIE = "Cookies"
)

type WebContext interface {
	Body() []byte
	Get(key string, defaultValue ...string) string
	Params(key string, defaultValue ...string) string
	Cookies(key string, defaultValue ...string) string
}

type Condition interface {
    check(context *WebContext) bool
}

type DataCondition struct {
	Prop string
	Key string
    Rule *Rule
}

func (c *DataCondition) check(context *WebContext) bool {

    if strings.Compare(c.Prop, CONDITION_BODY) {
        return c.Rule.check(context.Body())
    }
    if strings.Compare(c.Prop, CONDITION_HEADER) {
        return c.Rule.check(context.Get(c.Key))
    }
    strings.Compare(c.Prop, CONDITION_PARAMS) {
        return c.Rule.check(context.Params(c.Key))
    }
    if strings.Compare(c.Prop, CONDITION_COOKIE) {
        return c.Rule.check(context.Cookies(c.Key))
    }
    return false
}

type AndCondition struct {
    conditions []Condition
}

func (c *AndCondition) check(context *WebContext) bool {
    for _, cond := c.conditions {
        if !cond.check(context) {
            return false
        }
    }
    return true
}

type OrCondition struct {
    conditions []Condition
}

func (c *OrCondition) check(context *WebContext)  bool {
    for _, cond := c.conditions {
        if cond.check(context) {
            return true
        }
    }
    return false
}
