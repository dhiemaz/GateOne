package conditions

import "github.com/ory/ladon"

// StringList match conditions where given value match one of predefined options in case sensitive fashion
type StringList struct {
	Options []string `json:"options"`
}

func init() {
	ladon.ConditionFactories[new(StringList).GetName()] = func() ladon.Condition {
		return new(StringList)
	}
}

// Fulfills checking condition rule
func (c *StringList) Fulfills(value interface{}, _ *ladon.Request) bool {
	var rv = make([]string, 0)
	if s, ok := value.(string); ok {
		rv = append(rv, s)
	} else if s, ok := value.([]string); ok {
		rv = s
	}
	if len(rv) == 0 {
		return false
	}
	for _, x := range c.Options {
		var f = false
		for _, y := range rv {
			if y == x {
				f = true
				break
			}
		}
		if !f {
			return false
		}
	}

	return true
}

// GetName condition
func (c *StringList) GetName() string {
	return "StringListCondition"
}
