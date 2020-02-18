package conditions

import (
	"strings"

	"github.com/ory/ladon"
)

// StringPrefix match given value prefixed with pre-defined prefix
// CaseSensitive an option whether comparison done in case sensitive or not
type StringPrefix struct {
	Prefix        string `json:"prefix" bson:"prefix"`
	CaseSensitive bool   `json:"case_sensitive" bson:"case_sensitive"`
}

func init() {
	ladon.ConditionFactories[new(StringPrefix).GetName()] = func() ladon.Condition {
		return new(StringPrefix)
	}
}

// Fulfills checking condition rule
func (c *StringPrefix) Fulfills(value interface{}, _ *ladon.Request) bool {
	s, ok := value.(string)
	if !ok {
		return false
	}
	if c.CaseSensitive {
		return strings.HasPrefix(s, c.Prefix)
	}
	return strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(c.Prefix))
}

// GetName condition
func (c *StringPrefix) GetName() string {
	return "StringPrefixCondition"
}
