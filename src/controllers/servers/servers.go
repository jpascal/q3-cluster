package servers

import (
	"cluster"
	"context"
	"github.com/go-playground/lars"
	"strconv"
	"log"
	"io/ioutil"
)

func Routes(routes lars.IRouteGroup) {
	routes.Get("", Index)
	routes.Get("/:id", Show)
	routes.Get("/:id/status", Status)
	routes.Post("/:id/console", Console)
}

func Index(context *context.Context) {
	value, _ := context.Get("cluster")
	c := value.(cluster.Cluster)
	context.JSON(200, c.Servers)
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
	if (id > 0) && id <= int64(len(c.Servers)) {
		context.JSON(200, c.Servers[id-1])
	} else {
		context.Response().WriteHeader(404)
	}
}

func Status(context *context.Context) {
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
	if (id > 0) && id <= int64(len(c.Servers)) {
		// Read the content
		server := c.Servers[id-1]
		status, d := server.GetStatus()
		if d!= nil { panic(d) }
		context.JSON(200, status)
	} else {
		context.Response().WriteHeader(404)
	}
}


func Console(context *context.Context) {
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
	if (id > 0) && id <= int64(len(c.Servers)) {
		// Read the content
		var buffer []byte
		if context.Request().Body != nil {
			buffer, _ = ioutil.ReadAll(context.Request().Body)
		}
		server := c.Servers[id-1]
		server.Console(string(buffer))
	} else {
		context.Response().WriteHeader(404)
	}
}
