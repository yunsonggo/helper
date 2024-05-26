package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/yunsonggo/helper/types"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

type GoRedis struct {
	Client *redis.Client
	Mutex  *redsync.Redsync
	Keys   map[string]string
}

func NewRedisClient(conf *types.Redis) (*GoRedis, error) {
	addr := fmt.Sprintf("%s:%d", conf.RedisAddr, conf.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     conf.RedisUser,
		Password:     conf.RedisPass,
		DB:           conf.RedisDB,
		MaxRetries:   conf.WithTriesTimes,
		DialTimeout:  time.Duration(conf.WithExpirySecond) * time.Second,
		ReadTimeout:  time.Duration(conf.ReadTimeoutSecond) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeoutSecond) * time.Second,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConn,
		MaxConnAge:   time.Duration(conf.MaxConnAgeMinute) * time.Minute,
		PoolTimeout:  time.Duration(conf.PoolTimeoutMinute) * time.Minute,
		IdleTimeout:  time.Duration(conf.IdleTimeoutMinute) * time.Minute,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	// 连接池
	pool := goredis.NewPool(rdb)
	rs := redsync.New(pool)
	return &GoRedis{
		Client: rdb,
		Mutex:  rs,
		Keys:   make(map[string]string),
	}, nil
}

func (r *GoRedis) AddKey(name, key string) {
	r.Keys[name] = key
}

func (r *GoRedis) GetKey(name string) (string, bool) {
	if key, ok := r.Keys[name]; ok {
		return key, true
	}
	return "", false
}

func (r *GoRedis) RemoveKey(name string) {
	if _, ok := r.GetKey(name); ok {
		delete(r.Keys, name)
	}
}

func (r *GoRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *GoRedis) Get(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key:%s does not exist or has expired", key)
	} else if err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func (r *GoRedis) Delete(ctx context.Context, key string) error {
	if _, err := r.Client.Get(ctx, key).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		} else {
			return err
		}
	}
	return r.Client.Del(ctx, key).Err()
}

func (r *GoRedis) ExpireTime(key string) (time.Duration, error) {
	return r.Client.TTL(context.Background(), key).Result()
}

func (r *GoRedis) HashSet(ctx context.Context, key, field, value string) error {
	return r.Client.HSet(ctx, key, field, value).Err()
}

func (r *GoRedis) HashGet(ctx context.Context, key, field string) (string, error) {
	value, err := r.Client.HGet(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("key:%s,field:%s does not exist or has expired", key, field)
	} else if err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func (r *GoRedis) HashList(ctx context.Context, key string) (map[string]string, error) {
	value, err := r.Client.HGetAll(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("key:%s does not exist or has expired", key)
		} else {
			return nil, err
		}
	}
	return value, nil
}

func (r *GoRedis) HashDelete(ctx context.Context, key, field string) error {
	if _, err := r.Client.HGet(ctx, key, field).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		} else {
			return err
		}
	}
	return r.Client.HDel(ctx, key, field).Err()
}
