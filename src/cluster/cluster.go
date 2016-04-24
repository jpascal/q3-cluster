package cluster

import (
	"log"
	"os"
	"server"
	"storage"
)

type Cluster struct {
	Logger  *log.Logger
	Storage *storage.Storage
	Servers map[string]*server.Server `json:"servers"`
}

func NewCluster(storage *storage.Storage) *Cluster {
	var cluster Cluster
	cluster.Logger = log.New(os.Stdout, "[cluster] ", log.Ldate|log.Lmicroseconds)
	cluster.Storage = storage
	cluster.Servers = make(map[string]*server.Server)
	return &cluster
}

func (self *Cluster) AddServer(server *server.Server) string {
	id := self.Storage.Servers().Create(storage.ServerRecord{Address: server.Address, Port: server.Port})
	self.Servers[id] = server
	return id
}

func (self *Cluster) DelServer(id string) {
	server := self.Servers[id]
	server.Shutdown()
	self.Storage.Servers().Delete(id)
	delete(self.Servers, id)
}

func (self *Cluster) ServerByID(id string) *server.Server {
	return self.Servers[id]
}

func (self *Cluster) Startup() {

	for _, record := range self.Storage.Servers().All() {
		self.Logger.Printf("loading server from %v", record.Id)
		self.Servers[record.Id] = server.NewServer(record.Address, record.Port)
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
