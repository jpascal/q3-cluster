package servers

import (
	"context"
	"encoding/json"
	"github.com/go-playground/lars"
	"io/ioutil"
	"server"
)

func Routes(routes lars.IRouteGroup) {
	routes.Get("", Index)
	routes.Get("/:id", Show)
	routes.Get("/:id/startup", Startup)
	routes.Get("/:id/shutdown", Shutdown)
	routes.Get("/:id/status", Status)
	routes.Post("/:id/console", Console)
	routes.Post("", Create)
	routes.Delete("/:id", Delete)
}

type ServerBaseFields struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type ResponseServer struct {
	ServerBaseFields
	Id      string `json:"id"`
	Started bool   `json:"started"`
}

type RequestServer struct {
	ServerBaseFields
}

func Delete(context *context.Context) {
	if server := context.Cluster().ServerByID(context.Param("id")); server != nil {
		context.Cluster().DelServer(context.Param("id"))
	} else {
		context.Response().WriteHeader(404)
	}
}

func Create(context *context.Context) {
	cluster := context.Cluster()

	request := RequestServer{}

	body, err := ioutil.ReadAll(context.Request().Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		panic(err)
	}

	var id string

	if request.Address == "" || request.Port == 0 {
		context.Response().WriteHeader(406)
	} else {
		id = cluster.AddServer(server.NewServer(request.Address, request.Port))
	}

	context.JSON(200, ResponseServer{ Id: id, ServerBaseFields:ServerBaseFields{ Address: request.Address, Port: request.Port }})
}

func Index(context *context.Context) {
	cluster := context.Cluster()

	var holder []ResponseServer

	for id, server := range cluster.Servers {
		holder = append(holder, ResponseServer{
			ServerBaseFields:ServerBaseFields{ Address: server.Address, Port: server.Port },
			Id:      id,
			Started: server.Started,
		})
	}
	context.JSON(200, holder)
}

func Show(context *context.Context) {
	id := context.Param("id")

	if server := context.Cluster().ServerByID(id); server != nil {
		context.JSON(200, ResponseServer{
			ServerBaseFields:ServerBaseFields{ Address: server.Address, Port: server.Port },
			Id:      id,
			Started: server.Started,
		})
	} else {
		context.Response().WriteHeader(404)
	}
}

func Status(context *context.Context) {
	if server := context.Cluster().ServerByID(context.Param("id")); server != nil {
		if status, err := server.GetStatus(); err != nil {
			panic(err)
		} else {
			context.JSON(200, status)
		}

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
