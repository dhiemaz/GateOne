package gate_test

import (
	"context"
	"fmt"
	"github.com/ndv6/gate"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestEventMongoStoreSuite(t *testing.T) {
	db, cb := initTest()
	defer cb()

	eventCollection := db.Collection("events")
	var cleanup = func() {
		eventCollection.DeleteMany(context.TODO(), bson.D{})
	}
	cleanup()

	storage := gate.NewEventMongoStore(eventCollection)

	var seeds = func(n int, duplicate bool) []*gate.Event {
		var events = make([]*gate.Event, 0)
		for i := 0; i < n; i++ {
			seq := i
			if !duplicate {
				seq = 0
			}
			events = append(events, gate.NewEvent("komang","merchant",
				fmt.Sprintf("action-%d", seq),
				fmt.Sprintf("this is a note of event-%d", i),
				map[string]interface{}{
					"deals_id": seq,
					"agenda_id": seq,
				},
			))
		}
		return events
	}

	t.Run("MongoStore_Write_Historical", func(t *testing.T) {
		defer cleanup()

		events := seeds(10, true)
		for _,e := range events {
			if err := storage.Emit(context.TODO(), e, true); err != nil {
				t.Error(err)
			}
		}
		filter := bson.M{"user_id": "komang", "merchant_id": "merchant"}
		i, err := eventCollection.CountDocuments(context.TODO(), filter)
		if err != nil {
			t.Error(err)
		}
		if i != 10 {
			t.Errorf("expected count 10 got %d", i)
		}
	})

	t.Run("MongoStore_Write_Non_Historical", func(t *testing.T) {
		defer cleanup()

		events := seeds(10, false)
		for _,e := range events {
			if err := storage.Emit(context.TODO(), e, false); err != nil {
				t.Error(err)
			}
		}
		filter := bson.M{"user_id": "komang", "merchant_id": "merchant"}
		i, err := eventCollection.CountDocuments(context.TODO(), filter)
		if err != nil {
			t.Error(err)
		}
		if i != 1 {
			t.Errorf("expected count 1 got %d", i)
		}
	})

	t.Run("MongoStore_FindUserMerchants_Historical", func(t *testing.T) {
		defer cleanup()

		events := seeds(10, true)
		for _,e := range events {
			if err := storage.Emit(context.TODO(), e, true); err != nil {
				t.Error(err)
			}
		}
		mid, err := storage.FindUserMerchants(context.TODO(), "komang")
		if err != nil {
			t.Error(err)
		}
		if len(mid) != 1 {
			t.Error(err)
		}
	})

	t.Run("MongoStore_FindUserMerchants_Non_Historical", func(t *testing.T) {
		defer cleanup()

		events := seeds(10, false)
		for _,e := range events {
			if err := storage.Emit(context.TODO(), e, true); err != nil {
				t.Error(err)
			}
		}
		mid, err := storage.FindUserMerchants(context.TODO(), "komang")
		if err != nil {
			t.Error(err)
		}
		if len(mid) != 1 {
			t.Error(err)
		}
	})

	t.Run("MongoStore_Retrieve", func(t *testing.T) {
		defer cleanup()

		events := seeds(10, true)
		for _,e := range events {
			if err := storage.Emit(context.TODO(), e, true); err != nil {
				t.Error(err)
			}
		}
		t.Run("Retrieve Without Event Names/Meta", func(t *testing.T) {
			mid, err := storage.Retrieve(context.TODO(), "komang", "merchant", nil, nil, 5, 0)
			if err != nil {
				t.Errorf("%+v", err)
			}
			if len(mid) != 5 {
				t.Errorf("expected %d got %d", 5, len(mid))
			}
		})

		t.Run("Retrieve Without Meta", func(t *testing.T) {
			mid, err := storage.Retrieve(context.TODO(), "komang", "merchant", []string{"action-1","action-2","action-3"}, nil, 2, 1)
			if err != nil {
				t.Error(err)
			}
			if len(mid) != 2 {
				t.Errorf("expected %d got %d", 2, len(mid))
			}
		})

		t.Run("Retrieve", func(t *testing.T) {
			mid, err := storage.Retrieve(context.TODO(), "komang", "merchant", []string{"action-1","action-2","action-3"}, map[string]interface{}{"deals_id":1}, 5, 0)
			if err != nil {
				t.Error(err)
			}
			if len(mid) != 1 {
				t.Errorf("expected %d got %d", 1, len(mid))
			}
		})

	})
}
