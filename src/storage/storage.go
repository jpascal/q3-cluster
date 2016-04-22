package storage

import (
	"config"
	"gopkg.in/redis.v4"
)

type Storage struct {
	Redis *redis.Client
}

func NewStorage() *Storage {
	return &Storage{
		Redis: redis.NewClient(&redis.Options{
			Addr:     config.Config().Storage.Address,
			Password: config.Config().Storage.Password,
			DB:       config.Config().Storage.Database,
		}),
	}
}
