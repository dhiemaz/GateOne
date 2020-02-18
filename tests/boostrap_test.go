package gate_test

import (
	"context"
	"github.com/ndv6/gate"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

// MONGO_URL=""mongodb://localhost:27017""
func initTest() (*mongo.Database, func()) {
	client := gate.MongoMustConnect(os.Getenv("MONGO_URL"))
	db := client.Database("onelabs")
	cb := func() {
		_ = db.Drop(context.TODO())
		_ = client.Disconnect(context.TODO())
	}
	return db, cb
}
