package repository

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"konami_backend/internal/pkg/csrf"
)

type TokenManager struct {
	redisPool *redis.Pool
}

func NewRedisTokenManager(conn *redis.Pool) csrf.Repository {
	return &TokenManager{redisPool: conn}
}

func (tm *TokenManager) Add(token string, expire int64) error {
	conn := tm.redisPool.Get()
	result, err := redis.String(conn.Do("SET", token, 1, "EX", expire))
	if err != nil {
		return err
	}
	if result != "OK" {
		return errors.New("unable to add token")
	}
	return nil
}

func (tm *TokenManager) Validate(token string) error {
	conn := tm.redisPool.Get()
	_, err := redis.String(conn.Do("GET", token))

	if err != nil {
		if err == redis.ErrNil {
			return nil
		}
		return err
	}
	return errors.New("invalid token")
}
