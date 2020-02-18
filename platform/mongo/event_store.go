package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EventStore list functions
type EventStore interface {
	Emit(ctx context.Context, e *Event, multiple bool) error
	FindUserMerchants(ctx context.Context, userId string) ([]string, error)
	Retrieve(ctx context.Context, userId, merchantId string, e []string, meta map[string]interface{}, limit, skip int64) (events []Event, err error) /**/
}

// EventMongoStore is ...
type EventMongoStore struct {
	db *mongo.Collection
}

// NewEventMongoStore is ...
func NewEventMongoStore(db *mongo.Collection) *EventMongoStore {
	return &EventMongoStore{db: db}
}

// Emit is ...
func (store *EventMongoStore) Emit(ctx context.Context, e *Event, multiple bool) error {
	t := time.Now().Unix()
	if time.Unix(e.EventTime, 0).IsZero() {
		e.EventTime = t
	}

	// store all events - allow similar events co-exists
	if multiple {
		e.ID = primitive.NewObjectID().Hex()
		e.CreatedAt = t
		if _, err := store.db.InsertOne(ctx, e); err != nil {
			return errors.Wrap(err, "failed creating new persistent event")
		}
		return nil
	}

	// allow only 1 similar event exists by user/merchant/action
	filter := bson.M{"$and": bson.A{
		bson.D{{"user_id", e.UserId}},
		bson.D{{"merchant_id", e.MerchantId}},
		bson.D{{"action", e.Action}},
	}}

	values := bson.M{"user_id": e.UserId, "merchant_id": e.MerchantId, "action": e.Action}
	values["notes"] = e.Notes
	values["meta"] = e.Meta
	values["created_at"] = t

	update := bson.D{{"$set", bson.D{
		{"event_time", e.EventTime},
		{"updated_at", t},
	}}, {"$setOnInsert", values}}

	opt := options.Update().SetUpsert(true)

	if _, err := store.db.UpdateOne(ctx, filter, update, opt); err != nil {
		return errors.Wrap(err, "failed updating event")
	}

	return nil
}

// FindUserMerchants is ...
func (store *EventMongoStore) FindUserMerchants(ctx context.Context, userId string) (mid []string, err error) {
	var t []interface{}
	t, e := store.db.Distinct(ctx, "merchant_id", bson.D{{"user_id", userId}})
	if e != nil {
		err = errors.Wrap(e, "failed retrieving unique merchants")
		return
	}
	for _, v := range t {
		mid = append(mid, v.(string))
	}
	return
}

// Retrieve is ...
func (store *EventMongoStore) Retrieve(ctx context.Context, userId, merchantId string, e []string, meta map[string]interface{}, limit, skip int64) (events []Event, err error) {
	filters := bson.A{
		bson.M{"user_id": userId},
		bson.M{"merchant_id": merchantId},
	}
	if e != nil && len(e) > 0 {
		filters = append(filters, bson.M{"action": bson.M{"$in": e}})
	}
	if meta != nil && len(meta) > 0 {
		for k, v := range meta {
			filters = append(filters, bson.M{fmt.Sprintf("meta.%s", k): v})
		}
	}
	filter := bson.M{
		"$and": filters,
	}
	opt := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.M{"_id": -1})

	var c *mongo.Cursor
	c, err = store.db.Find(ctx, filter, opt)
	if err != nil {
		err = errors.Wrap(err, "failed retrieving all events")
		return
	}

	if err = c.All(ctx, &events); err != nil {
		err = errors.Wrap(err, "failed formatting events result")
	}
	return
}
