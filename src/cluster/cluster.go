package cluster

import (
	"server"
)

type Cluster struct {
	Servers []*server.Server
}

func (self *Cluster) AddServer(server *server.Server) {
	self.Servers = append(self.Servers, server)
}

func (self *Cluster) Startup() {
	for _, server := range self.Servers {
		server.Startup()
	}
}

func (self *Cluster) Shutdown() {
	for _, server := range self.Servers {
		server.Shutdown()
	}
}
