package storage

import (
	"fmt"
	"gopkg.in/redis.v4"
	"log"
	"os"
	"strings"
	"uuid"
)

type ServersStore struct {
	Redis  *redis.Client
	Logger *log.Logger
}

type ServerRecord struct {
	Id      string `json:"-"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

const SERVERS_STORE_PATH = "servers::"

func (self *ServerRecord) Key() string {
	return SERVERS_STORE_PATH + self.Id
}

func GetServersStore(redis *redis.Client) *ServersStore {
	return &ServersStore{Redis: redis, Logger: log.New(os.Stdout, "[storage.servers] ", log.Ldate|log.Lmicroseconds)}
}

func (self *ServersStore) Exist(key string) bool {
	if result := self.Redis.Exists(key); result.Err() != nil {
		panic(result.Err())
	} else {
		return result.Val()
	}
}

func (self *ServersStore) Find(id string) *ServerRecord {
	self.Logger.Printf("find %v", id)
	var record ServerRecord

	record.Id = id

	if !self.Exist(record.Key()) {
		return nil
	}

	if result := self.Redis.Get(record.Key()); result.Err() != nil {
		panic(result.Err())
	} else {
		Deserialize(result.Val(), &record)
	}

	return &record
}

func (self *ServersStore) Create(record ServerRecord) string {
	record.Id = uuid.NewUUID()
	if self.Exist(record.Key()) {
		panic(fmt.Sprintf("record already exists (id=%v)", record.Id))
	}

	if result := self.Redis.Set(record.Key(), Serialize(record), 0); result.Err() != nil {
		panic(result.Err())
	}

	self.Logger.Printf("created %v", record.Id)

	return record.Id
}

func (self *ServersStore) Delete(id string) int64 {
	self.Logger.Printf("delete %v", id)
	var record ServerRecord

	record.Id = id

	if !self.Exist(record.Key()) {
		return 0
	}

	if result := self.Redis.Del(record.Key()); result.Err() != nil {
		panic(result.Err())
	} else {
		return result.Val()
	}
}

func (self *ServersStore) Save(record ServerRecord) {
	if record.Id == "" {
		panic("can't save without id")
	}

	self.Logger.Printf("save %v", record.Id)

	if !self.Exist(record.Key()) {
		panic(fmt.Sprintf("record not exists (id=%v)", record.Id))
	}
	if result := self.Redis.Set(record.Key(), Serialize(record), 0); result.Err() != nil {
		panic(result.Err())
	}
}

func (self *ServersStore) All() (records []*ServerRecord) {

	self.Logger.Printf("get all")

	result := self.Redis.Keys(SERVERS_STORE_PATH + "*")

	if result.Err() != nil {
		panic(result.Err())
	}

	for _, key := range result.Val() {
		result := self.Redis.Get(key)
		if result.Err() != nil {
			panic(result.Err())
		}
		record := ServerRecord{}
		Deserialize(result.Val(), &record)
		record.Id = strings.Replace(key, SERVERS_STORE_PATH, "", 1)
		records = append(records, &record)
	}
	return
}
