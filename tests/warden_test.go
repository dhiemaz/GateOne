package gate_test

import (
	"github.com/ndv6/gate"
	"github.com/ory/ladon"
	"log"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	db, cb := initTest()
	defer cb()

	mp := gate.NewMongoPolicyManager("eliving", db)
	warden := &ladon.Ladon{
		Manager:     mp,
		Matcher:     nil,
		AuditLogger: new(ladon.AuditLoggerInfo),
	}

	policies := seedPolicies(10)
	for _, p := range policies {
		if err := mp.Create(p); err != nil {
			t.Error(err)
		}
	}

	type Pair struct {
		result  bool
		request ladon.Request
	}

	var payloads = []Pair{
		{true, ladon.Request{
			Resource: "room:5",
			Action:   "create",
			Subject:  "groups:administrators",
			Context: ladon.Context{
				"va": "PRE-5",
			},
		}},
		{false, ladon.Request{
			Resource: "room:5",
			Action:   "get", // no get request allowed
			Subject:  "groups:administrators",
			Context: ladon.Context{
				"va": "PRE-5",
			},
		}},
		{false, ladon.Request{
			Resource: "room:1", // room 1 vs. VA PRE-5
			Action:   "create",
			Subject:  "groups:administrators",
			Context: ladon.Context{
				"va": "PRE-5",
			},
		}},
	}
	for i, p := range payloads {
		log.Printf("testing %d", i)
		if err := warden.IsAllowed(&p.request); err != nil && p.result {
			t.Error(err)
		}
	}
}
