package main

import (
	"cluster"
	"context"
	"controllers/servers"
	"github.com/go-playground/lars"
	"log"
	"net/http"
	"server"
)

func main() {

	var cluster cluster.Cluster

	s1 := server.NewServer("localhost", 1)
	s2 := server.NewServer("localhost", 2)

	cluster.AddServer(s1)
	cluster.AddServer(s2)

	cluster.Startup()

	if err := s1.Console("test"); err != nil {
		log.Print(err)
	}

	router := lars.New()

	router.SetRedirectTrailingSlash(false)

	router.RegisterContext(context.NewContext)
	router.RegisterCustomHandler(func(*context.Context) {}, context.CastContext)

	router.Use(func(context lars.Context) {
		context.Set("cluster", cluster)
		context.Next()
	})

	servers.Routes(router.Group("/servers"))

	http.ListenAndServe(":3007", router.Serve())

	cluster.Shutdown()
}
