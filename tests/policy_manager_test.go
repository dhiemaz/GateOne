package gate_test

import (
	"context"
	"fmt"
	"github.com/ndv6/gate"
	"github.com/ory/ladon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

func seedPolicies(n int) []*gate.DefaultPolicy {
	var policies = make([]*gate.DefaultPolicy, 0)
	for i := 0; i < n; i++ {
		policies = append(policies, &gate.DefaultPolicy{
			Description: fmt.Sprintf("description #%d", i),
			Subjects:    []string{"groups:administrators"},
			Effect:      ladon.AllowAccess,
			Resources:   []string{fmt.Sprintf("room:%d", i)},
			Actions:     []string{"create", "update", "delete"},
			Conditions: gate.Conditions{
				"va": &gate.StringPrefixCondition{
					Prefix:        fmt.Sprintf("PRE-%d", i),
					CaseSensitive: true,
				},
			},
		})
	}
	return policies
}

func TestMongoPolicyManagerSuite(t *testing.T) {
	db, cb := initTest()
	defer cb()

	mp := gate.NewMongoPolicyManager("eliving", db)
	policyCollection := db.Collection("eliving_policies")
	var cleanup = func() {
		policyCollection.DeleteMany(context.TODO(), bson.D{})
	}

	t.Run("MongoPolicyManager_Create", func(t *testing.T) {
		defer cleanup()
		p := seedPolicies(1)[0]
		if err := mp.Create(p); err != nil {
			t.Error(err)
		}
	})

	t.Run("MongoPolicyManager_Get", func(t *testing.T) {
		defer cleanup()

		id := primitive.NewObjectID().Hex()
		mp := gate.NewMongoPolicyManager("eliving", db)
		p := seedPolicies(1)[0]
		p.ID = id
		if err := mp.Create(p); err != nil {
			t.Error(err)
		}
		ip, err := mp.Get(id)
		if err != nil {
			t.Errorf("%+v", err)
		}
		if ip.GetID() != id {
			t.Errorf("expected id %s got %s", id, ip.GetID())
		}
	})

	t.Run("MongoPolicyManager_Update", func(t *testing.T) {
		defer cleanup()

		id := primitive.NewObjectID().Hex()
		p := seedPolicies(1)[0]
		p.ID = id
		if err := mp.Create(p); err != nil {
			t.Error(err)
		}

		var (
			updatedDescription = "Updated description"
			updatedDocument = *p
		)
		updatedDocument.Description = updatedDescription
		if err := mp.Update(&updatedDocument); err != nil {
			t.Errorf("%+v", err)
		}

		ip, err := mp.Get(id)
		if err != nil {
			t.Errorf("%+v", err)
		}
		if ip.GetID() != id {
			t.Errorf("expected id %s got %s", id, ip.GetID())
		}
		if ip.GetDescription() != updatedDescription {
			t.Errorf("expected description: %s got %s", updatedDescription, ip.GetDescription())
		}
		if !reflect.DeepEqual(p.GetSubjects(), ip.GetSubjects()) {
			t.Errorf("subject shouldn't be updated; %+v vs. %+v", p.GetSubjects(), ip.GetSubjects())
		}
	})

	t.Run("MongoPolicyManager_Delete", func(t *testing.T) {
		defer cleanup()

		id := primitive.NewObjectID().Hex()
		p := seedPolicies(1)[0]
		p.ID = id
		if err := mp.Create(p); err != nil {
			t.Error(err)
		}

		if err := mp.Delete(id); err != nil {
			t.Error(err)
		}
	})

	t.Run("MongoPolicyManager_GetAll", func(t *testing.T) {
		defer cleanup()

		ps := seedPolicies(10)
		for _,p := range ps {
			p.ID = primitive.NewObjectID().Hex()
			if err := mp.Create(p); err != nil {
				t.Error(err)
			}
		}
		list, err := mp.GetAll(10, 0)
		if err != nil {
			t.Errorf("%+v", err)
		}
		if len(list) != 10 {
			t.Errorf("expected %d got %d", 10, len(list))
		}

		list, err = mp.GetAll(5, 4)
		if err != nil {
			t.Errorf("%+v", err)
		}
		if len(list) != 5 {
			t.Errorf("expected %d got %d", 5, len(list))
		}
	})

	t.Run("MongoPolicyManager_FindPoliciesForResource", func(t *testing.T) {
		defer cleanup()

		ps := seedPolicies(10)
		for _,p := range ps {
			p.ID = primitive.NewObjectID().Hex()
			if err := mp.Create(p); err != nil {
				t.Error(err)
			}
		}
		list, err := mp.FindPoliciesForResource("room:5")
		if err != nil {
			t.Errorf("%+v", err)
		}
		if len(list) != 1 {
			t.Errorf("expected %d got %d", 1, len(list))
		}
	})

	t.Run("MongoPolicyManager_FindRequestCandidates", func(t *testing.T) {
		defer cleanup()

		ps := seedPolicies(10)
		for _,p := range ps {
			p.ID = primitive.NewObjectID().Hex()
			if err := mp.Create(p); err != nil {
				t.Error(err)
			}
		}
		list, err := mp.FindRequestCandidates(&ladon.Request{
			Resource: "room:1",
			Action:   "update",
			Subject:  "groups:administrators",
		})
		if err != nil {
			t.Errorf("%+v", err)
		}
		if len(list) != 1 {
			t.Errorf("expected %d got %d", 1, len(list))
		}
	})

	t.Run("MongoPolicyManager_FindPoliciesForSubject", func(t *testing.T) {
		defer cleanup()

		ps := seedPolicies(10)
		for _,p := range ps {
			p.ID = primitive.NewObjectID().Hex()
			if err := mp.Create(p); err != nil {
				t.Error(err)
			}
		}
		list, err := mp.FindPoliciesForSubject("groups")
		if err != nil {
			t.Errorf("%+v", err)
		}
		if len(list) != 10 {
			t.Errorf("expected %d got %d", 10, len(list))
		}
	})
}
