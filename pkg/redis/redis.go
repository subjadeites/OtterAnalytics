package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func ConnectRedisClient(addr string, port int, db int, pwd ...string) *RedisClient {
	password := ""
	if len(pwd) > 0 {
		password = pwd[0]
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr + ":" + strconv.Itoa(port),
		Password: password,
		DB:       db,
	})
	return &RedisClient{Client: client, Ctx: context.Background()}
}

func (r *RedisClient) Ping() error {
	return r.Client.Ping(r.Ctx).Err()
}

func (r *RedisClient) Set(key string, value interface{}, expiration int) error {
	return r.Client.Set(r.Ctx, key, value, time.Duration(expiration)*time.Second).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *RedisClient) Del(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

func (r *RedisClient) HSet(key, field string, value interface{}) error {
	return r.Client.HSet(r.Ctx, key, field, value).Err()
}

func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.Client.HGet(r.Ctx, key, field).Result()
}

func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.Client.HGetAll(r.Ctx, key).Result()
}

func (r *RedisClient) HDel(key string, fields ...string) error {
	return r.Client.HDel(r.Ctx, key, fields...).Err()
}
