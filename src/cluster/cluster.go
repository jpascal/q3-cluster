package cluster

import (
	"server"
	"log"
	"os"
)

type Cluster struct {
	Servers []*server.Server `json:"servers"`
	Logger	*log.Logger
}


func NewCluster() *Cluster {
	var cluster Cluster
	cluster.Logger = log.New(os.Stdout, "[cluster] ", log.Ldate|log.Lmicroseconds)
	return &cluster
}

func (self *Cluster) AddServer(server *server.Server) {
	self.Logger.Printf("add server to cluster %v:%v", server.Address, server.Port)
	self.Servers = append(self.Servers, server)
}

func (self *Cluster) Startup() {
	self.Logger.Print("startup")
	for _, server := range self.Servers {
		if err := server.Startup(); err != nil {
			server.Shutdown()
		}
	}
}

func (self *Cluster) Shutdown() {
	self.Logger.Print("shutdown")
	for _, server := range self.Servers {
		server.Shutdown()
	}
}
