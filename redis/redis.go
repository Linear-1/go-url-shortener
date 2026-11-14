package redis

import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context,) (*redis.Client,error){
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		MaxRetries: -1,
		DialTimeout: 500* time.Millisecond,
		DialerRetries: 1,
	})

	err := rdb.Ping(ctx).Err()

	if err != nil {
		rdb.Close()
		return nil , err
	}

	return rdb , nil
}

