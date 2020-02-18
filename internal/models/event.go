package model

import "time"

// Event data model
type Event struct {
	ID         string `json:"id" bson:"_id" db:"id"`
	UserId     string `json:"user_id" bson:"user_id" db:"user_id"`
	MerchantId string `json:"merchant_id" bson:"merchant_id" db:"merchant_id"`
	Action     string `json:"action" bson:"action" db:"action"`             // action i.e. register; login etc.
	EventTime  int64  `json:"event_time" bson:"event_time" db:"event_time"` // unix timestamp

	Notes string                 `json:"notes" bson:"notes"  db:"notes"` // verbatim notes
	Meta  map[string]interface{} `json:"meta" bson:"meta"  db:"meta"`    // search-able context

	CreatedAt int64 `json:"created_at" bson:"created_at" db:"created_at"` // created at - unix timestamp
	UpdatedAt int64 `json:"updated_at" bson:"updated_at" db:"updated_at"` // updated at - unix timestamp
}

// NewEvent created new event
func NewEvent(userId, merchantId, action, notes string, meta map[string]interface{}) *Event {
	t := time.Now().Unix()
	return &Event{
		UserId:     userId,
		MerchantId: merchantId,
		Action:     action,
		EventTime:  t,
		Notes:      notes,
		Meta:       meta,
		CreatedAt:  t,
	}
}
