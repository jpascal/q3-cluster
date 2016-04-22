package servers

import (
	"context"
	"github.com/go-playground/lars"
	"io/ioutil"
)

func Routes(routes lars.IRouteGroup) {
	routes.Get("", Index)
	routes.Get("/:id", Show)
	routes.Get("/:id/startup", Startup)
	routes.Get("/:id/shutdown", Shutdown)
	routes.Get("/:id/status", Status)
	routes.Post("/:id/console", Console)
}

type HolderServer struct {
	Id	string		`json:"id"`
	Address	string	`json:"address"`
	Port	int		`json:"port"`
	Started	bool	`json:"started"`
}

func Index(context *context.Context) {
	cluster := context.Cluster()

	var holder []HolderServer

	for _, server := range cluster.Servers {
		holder = append(holder, HolderServer{
			Id: "asdf",
			Address: server.Address,
			Port: server.Port,
			Started: server.Started,
		})
	}
	context.JSON(200, holder)
}

func Show(context *context.Context) {
	id := context.Param("id")

	if server := context.Cluster().ServerByID(id); server != nil {
		context.JSON(200, HolderServer{
			Id: id,
			Address: server.Address,
			Port: server.Port,
			Started: server.Started,
		})
	} else {
		context.Response().WriteHeader(404)
	}
}

func Status(context *context.Context) {
	if server := context.Cluster().ServerByID(context.Param("id")); server != nil {
		context.JSON(200, server.GetStatus())
	} else {
		context.Response().WriteHeader(404)
	}
}

func Console(context *context.Context) {
	if server := context.Cluster().ServerByID(context.Param("id")); server != nil {
		var buffer []byte
		if context.Request().Body != nil {
			buffer, _ = ioutil.ReadAll(context.Request().Body)
		}
		if err := server.Console(string(buffer)); err != nil {
			panic(err)
		}
	} else {
		context.Response().WriteHeader(404)
	}
}

func Startup(context *context.Context) {
	if server := context.Cluster().ServerByID(context.Param("id")); server != nil {
		if err := server.Startup(); err != nil {
			panic(err)
		}
	} else {
		context.Response().WriteHeader(404)
	}
}

func Shutdown(context *context.Context) {
	if server := context.Cluster().ServerByID(context.Param("id")); server != nil {
		if err := server.Shutdown(); err != nil {
			panic(err)
		}
	} else {
		context.Response().WriteHeader(404)
	}
}
