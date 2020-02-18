package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// RedisConnect it create a redis client from given URL and attempt to establish connection
// url format based on [rfc1738](https://www.ietf.org/rfc/rfc1738.txt)
// redis://<user>:<password>@<host>:<port>/<database>
func RedisConnect(uri string) (client *redis.Client, err error) {
	var opt *redis.Options
	opt, err = redis.ParseURL(uri)
	if err != nil {
		err = errors.Wrap(err, "failed to parse redis uri")
		return
	}

	client = redis.NewClient(opt)
	_, err = client.Ping().Result()
	return
}

// RedisMustConnect is ...
func RedisMustConnect(uri string) *redis.Client {
	c, err := RedisConnect(uri)
	if err != nil {
		panic(err)
	}
	return c
}
