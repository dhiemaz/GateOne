package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConnect is ...
func MongoConnect(uri string) (c *mongo.Client, err error) {
	var (
		opt                   = make([]*options.ClientOptions, 0)
		defaultConnectTimeout = 10 * time.Second
		defaultPingTimeout    = 2 * time.Second
	)
	opt = append(opt, options.Client().SetAppName("on-gate").ApplyURI(uri))
	ctx, _ := context.WithTimeout(context.Background(), defaultConnectTimeout)
	c, err = mongo.Connect(ctx, opt...)
	if err != nil {
		err = errors.Wrap(err, "failed to create mongodb client")
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), defaultPingTimeout)
	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		err = errors.Wrap(err, "failed to establish connection to mongodb server")
	}

	return
}

// MongoMustConnect is ...
func MongoMustConnect(uri string) *mongo.Client {
	c, err := MongoConnect(uri)
	if err != nil {
		panic(err)
	}
	return c
}
