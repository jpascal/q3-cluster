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
			Addr:     config.GetConfig().Storage.Address,
			Password: config.GetConfig().Storage.Password,
			DB:       config.GetConfig().Storage.Database,
		}),
	}
}

func (self *Storage) Servers() *ServersStore {
	return GetServersStore(self.Redis)
}
