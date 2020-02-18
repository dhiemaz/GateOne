package policies

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ory/ladon"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// PolicyTableSuffix ...
	PolicyTableSuffix = "_policies"
)

var (
	// ErrPolicyInvalidParameter is ...
	ErrPolicyInvalidParameter = errors.New("request contains invalid parameter")

	// ErrPolicyNotFound is ...
	ErrPolicyNotFound = errors.New("requested policy could not be found")

	// ErrNoPolicy is ...
	ErrNoPolicy = errors.New("not policy found matching criteria")
)

// MongoPolicyManager is ...
type MongoPolicyManager struct {
	db *mongo.Collection
}

// NewMongoPolicyManager is ...
func NewMongoPolicyManager(merchant string, db *mongo.Database) *MongoPolicyManager {
	return &MongoPolicyManager{db: db.Collection(fmt.Sprintf("%s%s", merchant, PolicyTableSuffix))}
}

// Create policy
func (pm *MongoPolicyManager) Create(policy ladon.Policy) error {
	pp := policy.(*DefaultPolicy)
	if pp.ID == "" {
		pp.ID = primitive.NewObjectID().Hex()
	}

	if _, err := pm.db.InsertOne(context.TODO(), pp); err != nil {
		return errors.Wrap(err, "failed creating new policy")
	}
	return nil
}

// Update existing policy
func (pm *MongoPolicyManager) Update(policy ladon.Policy) error {
	if policy.GetID() == "" {
		return errors.Wrap(ErrPolicyInvalidParameter, "update request requires id attribute")
	}

	updated := bson.D{{"$set", policy}}
	if _, err := pm.db.UpdateOne(context.TODO(), bson.D{{"_id", policy.GetID()}}, updated); err != nil {
		return errors.Wrapf(err, "failed updating policy #%s", policy.GetID())
	}
	return nil
}

// Get policy by id
func (pm *MongoPolicyManager) Get(id string) (ladon.Policy, error) {
	r := pm.db.FindOne(context.TODO(), bson.D{{"_id", id}})
	if err := r.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed retrieving policy #%s", id)
	}

	var p = new(DefaultPolicy)
	if err := r.Decode(&p); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.Wrapf(ErrPolicyNotFound, "policy #%s does not exists", id)
		}
		return nil, errors.Wrapf(err, "failed decoding policy #%s", id)
	}

	return p, nil
}

// Delete is ...
func (pm *MongoPolicyManager) Delete(id string) error {
	r, err := pm.db.DeleteOne(context.TODO(), bson.D{{"_id", id}})
	if err != nil {
		return errors.Wrap(err, "failed deleting policy")
	}

	if r.DeletedCount == 0 {
		return errors.Wrap(ErrPolicyNotFound, "requested policy not found")
	}

	return nil
}

// GetAll policies stored
func (pm *MongoPolicyManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	c, err := pm.db.Find(context.TODO(), bson.M{}, options.Find().SetLimit(limit).SetSkip(offset))
	if err != nil {
		return nil, errors.Wrap(err, "failed retrieving all policies")
	}

	return pm.policiesListFromCursor(c)
}

// FindRequestCandidates is ...
func (pm *MongoPolicyManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	opt := options.Find().SetLimit(0)
	qp := bson.A{}
	if r.Subject != "" {
		qp = append(qp, bson.D{{
			"subjects", bson.D{
				{"$regex", primitive.Regex{
					Pattern: fmt.Sprintf("^%s", regexp.QuoteMeta(r.Subject)),
					Options: "i",
				}},
			},
		}})
	}

	if r.Resource != "" {
		qp = append(qp, bson.D{{
			"resources", bson.D{
				{"$regex", primitive.Regex{
					Pattern: fmt.Sprintf("^%s", regexp.QuoteMeta(r.Resource)),
					Options: "i",
				}},
			},
		}})
	}

	if r.Action != "" {
		qp = append(qp, bson.D{{"actions", r.Action}})
	}

	query := bson.D{{"$and", qp}}
	c, err := pm.db.Find(context.TODO(), query, opt)
	if err != nil {
		return nil, errors.Wrap(err, "failed retrieving policies by request")
	}

	return pm.policiesListFromCursor(c)
}

// FindPoliciesForSubject is to search policies stored for specified subject
func (pm *MongoPolicyManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	query := bson.M{"subjects": bson.M{"$regex": primitive.Regex{
		Pattern: fmt.Sprintf("^%s", regexp.QuoteMeta(subject)),
		Options: "i",
	}}}

	opt := options.Find().SetLimit(0)
	c, err := pm.db.Find(context.TODO(), query, opt)
	if err != nil {
		return nil, errors.Wrap(err, "failed retrieving policies by subject")
	}

	return pm.policiesListFromCursor(c)
}

// FindPoliciesForResource is ...
func (pm *MongoPolicyManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	opt := options.Find().SetLimit(0)
	query := bson.D{
		{"resources", bson.D{
			{"$regex", primitive.Regex{
				Pattern: fmt.Sprintf("^%s", regexp.QuoteMeta(resource)),
				Options: "i",
			}},
		}},
	}

	c, err := pm.db.Find(context.TODO(), query, opt)
	if err != nil {
		return nil, errors.Wrap(err, "failed retrieving policies by resource")
	}

	return pm.policiesListFromCursor(c)
}

func (pm *MongoPolicyManager) policiesListFromCursor(c *mongo.Cursor) (ladon.Policies, error) {
	var (
		dp []*DefaultPolicy
		p  ladon.Policies
	)

	if err := c.All(context.TODO(), &dp); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoPolicy
		}
		return nil, errors.Wrap(err, "failed decoding all policies")
	}

	for _, v := range dp {
		p = append(p, v)
	}

	return p, nil
}
