package cluster

import (
	"log"
	"os"
	"server"
	"storage"
	"uuid"
)

type Cluster struct {
	_logger  *log.Logger
	_storage *storage.Storage
	Servers  map[string]*server.Server `json:"servers"`
}

func NewCluster(storage *storage.Storage) *Cluster {
	var cluster Cluster
	cluster._logger = log.New(os.Stdout, "[cluster] ", log.Ldate|log.Lmicroseconds)
	cluster._storage = storage
	return &cluster
}

func (self *Cluster) AddServer(server *server.Server) string {
	id := uuid.NewUUID()
	self.Servers[id] = server
	return id
}


func (self *Cluster) DelServer(id string) *server.Server {
	server := self.Servers[id]
	delete(self.Servers, id)
	return server
}

func (self *Cluster) ServerByID(id string) *server.Server {
	return self.Servers[id]
}

func (self *Cluster) Startup() {

	result := self._storage.Redis.Keys("servers:?*")

	for _, key := range result.Val() {
		data := self._storage.Redis.Get(key)

		if server, err := server.NewServerFromJSON(data.Val()); err != nil {
			self._logger.Printf("can't load server by key %v because: %s", key, err)
		} else {
			self.AddServer(server)
		}
	}

	self._logger.Print("startup all servers")
	for _, server := range self.Servers {
		if err := server.Startup(); err != nil {
			server.Shutdown()
		}
	}
}

func (self *Cluster) Shutdown() {
	self._logger.Print("shutdown all servers")
	for _, server := range self.Servers {
		server.Shutdown()
	}
}
