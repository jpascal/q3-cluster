package cluster

import (
	"log"
	"os"
	"server"
	"storage"
)

type Cluster struct {
	Servers []*server.Server `json:"servers"`
	Logger  *log.Logger
	Storage *storage.Storage
}

func NewCluster(storage *storage.Storage) *Cluster {
	var cluster Cluster
	cluster.Logger = log.New(os.Stdout, "[cluster] ", log.Ldate|log.Lmicroseconds)
	cluster.Storage = storage
	return &cluster
}

func (self *Cluster) AddServer(server *server.Server) *server.Server {
	self.Servers = append(self.Servers, server)
	return server
}

func (self *Cluster) Startup() {

	result := self.Storage.Redis.Keys("servers:?*")

	for _, key := range result.Val() {
		data := self.Storage.Redis.Get(key)

		if server, err := server.NewServerFromJSON(data.Val()); err != nil {
			self.Logger.Printf("can't load server by key %v because: %s", key, err)
		} else {
			self.AddServer(server)
		}
	}

	self.Logger.Print("startup all servers")
	for _, server := range self.Servers {
		if err := server.Startup(); err != nil {
			server.Shutdown()
		}
	}
}

func (self *Cluster) Shutdown() {
	self.Logger.Print("shutdown all servers")
	for _, server := range self.Servers {
		server.Shutdown()
	}
}
