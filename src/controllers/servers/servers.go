package servers

import (
	"cluster"
	"context"
	"github.com/go-playground/lars"
	"strconv"
	"log"
)

func Routes(routes lars.IRouteGroup) {
	routes.Get("", Index)
	routes.Post("", Create)
	routes.Patch("/:id", Update)
	routes.Delete("/:id", Destroy)
	routes.Get("/:id", Show)
}

func Index(context *context.Context) {
	value, _ := context.Get("cluster")
	c := value.(cluster.Cluster)
	context.JSON(200, c.Servers)
}

func Create(context *context.Context) {

}

func Update(context *context.Context) {

}

func Show(context *context.Context) {
	value, _ := context.Get("cluster")
	c := value.(cluster.Cluster)
	param := context.Param("id")
	var id int64
	var err error
	if id, err = strconv.ParseInt(param,10,32); err != nil {
		log.Print(err)
		context.Response().WriteHeader(500)
		return
	}
	server := c.Servers[id+2]
	if status, err := server.GetStatus(); err != nil {
		log.Print(err)
		context.Response().WriteHeader(500)
	} else {
		context.JSON(200, status)
	}
}

func Destroy(context *context.Context) {

}

func Start(context *context.Context) {

}

func Stop(context *context.Context) {

}

func Rcon(context *context.Context) {

}
