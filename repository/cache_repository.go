package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)


type CacheRepository struct{
	client *redis.Client
}


func NewCacheRepository (client *redis.Client) *CacheRepository{
	return &CacheRepository{
		client: client , 
	}
}

func (r *CacheRepository) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := r.client.Set(ctx , key , value , expiration).Err()

	if err != nil {
		return err
	}

	return nil

}


func (r *CacheRepository) Get(ctx context.Context , key string) (string, error) {
	originalURL , err := r.client.Get(ctx,key).Result() 

	if err != nil {
		return "" , err
	}

	return originalURL , nil
}


func (r *CacheRepository) Delete(ctx context.Context , key string) error {
	return r.client.Del(ctx, key).Err()

}