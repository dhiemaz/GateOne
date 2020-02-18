package policies

import (
	"encoding/json"

	"github.com/ory/ladon"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

// Conditions ladon
type Conditions ladon.Conditions

// MarshalBSON change bson object into byte object
func (cs Conditions) MarshalBSON() ([]byte, error) {
	out := make(map[string]*jsonCondition, len(cs))
	for k, c := range cs {
		raw, err := bson.Marshal(c)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}

		out[k] = &jsonCondition{
			Type:    c.GetName(),
			Options: bson.Raw(raw),
		}
	}

	return bson.Marshal(out)
}

// UnmarshalBSON change byte object into BSON object
func (cs Conditions) UnmarshalBSON(data []byte) error {
	if cs == nil {
		return errors.New("Can not be nil")
	}

	var jcs map[string]jsonCondition
	var dc ladon.Condition

	if err := bson.Unmarshal(data, &jcs); err != nil {
		return errors.WithStack(err)
	}

	for k, jc := range jcs {
		var found bool
		for name, c := range ladon.ConditionFactories {
			if name == jc.Type {
				found = true
				dc = c()

				if len(jc.Options) == 0 {
					cs[k] = dc
					break
				}

				if err := bson.Unmarshal(jc.Options, dc); err != nil {
					return errors.WithStack(err)
				}

				cs[k] = dc
				break
			}
		}

		if !found {
			return errors.Errorf("Could not find condition type %s", jc.Type)
		}
	}

	return nil
}

type jsonCondition struct {
	Type    string   `json:"type" bson:"type"`
	Options bson.Raw `json:"options" bson:"options"`
}

// DefaultPolicy is the default implementation of the policy interface.
type DefaultPolicy struct {
	ID          string     `json:"id" bson:"_id"`
	Description string     `json:"description" bson:"description"`
	Subjects    []string   `json:"subjects" bson:"subjects"`
	Effect      string     `json:"effect" bson:"effect"`
	Resources   []string   `json:"resources" bson:"resources"`
	Actions     []string   `json:"actions" bson:"actions"`
	Conditions  Conditions `json:"conditions" bson:"conditions"`
	Meta        []byte     `json:"meta" bson:"meta"`
}

// UnmarshalBSON overwrite own policy with values of the given in policy in JSON format
func (p *DefaultPolicy) UnmarshalBSON(data []byte) error {
	var pol = struct {
		ID          string     `json:"id" bson:"_id"`
		Description string     `json:"description" bson:"description"`
		Subjects    []string   `json:"subjects" bson:"subjects"`
		Effect      string     `json:"effect" bson:"effect"`
		Resources   []string   `json:"resources" bson:"resources"`
		Actions     []string   `json:"actions" bson:"actions"`
		Conditions  Conditions `json:"conditions" bson:"conditions"`
		Meta        []byte     `json:"meta" bson:"meta"`
	}{
		Conditions: Conditions{},
	}

	if err := bson.Unmarshal(data, &pol); err != nil {
		return errors.WithStack(err)
	}

	*p = *&DefaultPolicy{
		ID:          pol.ID,
		Description: pol.Description,
		Subjects:    pol.Subjects,
		Effect:      pol.Effect,
		Resources:   pol.Resources,
		Actions:     pol.Actions,
		Conditions:  pol.Conditions,
		Meta:        pol.Meta,
	}
	return nil
}

// UnmarshalMeta parses the policies []byte encoded metadata and stores the result in the value pointed to by v.
func (p *DefaultPolicy) UnmarshalMeta(v interface{}) error {
	if err := json.Unmarshal(p.Meta, &v); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// GetID returns the policies id.
func (p *DefaultPolicy) GetID() string {
	return p.ID
}

// GetDescription returns the policies description.
func (p *DefaultPolicy) GetDescription() string {
	return p.Description
}

// GetSubjects returns the policies subjects.
func (p *DefaultPolicy) GetSubjects() []string {
	return p.Subjects
}

// AllowAccess returns true if the policy effect is allow, otherwise false.
func (p *DefaultPolicy) AllowAccess() bool {
	return p.Effect == ladon.AllowAccess
}

// GetEffect returns the policies effect which might be 'allow' or 'deny'.
func (p *DefaultPolicy) GetEffect() string {
	return p.Effect
}

// GetResources returns the policies resources.
func (p *DefaultPolicy) GetResources() []string {
	return p.Resources
}

// GetActions returns the policies actions.
func (p *DefaultPolicy) GetActions() []string {
	return p.Actions
}

// GetConditions returns the policies conditions.
func (p *DefaultPolicy) GetConditions() ladon.Conditions {
	return ladon.Conditions(p.Conditions)
}

// GetMeta returns the policies arbitrary metadata set by the user.
func (p *DefaultPolicy) GetMeta() []byte {
	return p.Meta
}

// GetEndDelimiter returns the delimiter which identifies the end of a regular expression.
func (p *DefaultPolicy) GetEndDelimiter() byte {
	return '>'
}

// GetStartDelimiter returns the delimiter which identifies the beginning of a regular expression.
func (p *DefaultPolicy) GetStartDelimiter() byte {
	return '<'
}
