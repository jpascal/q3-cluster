package main

import (
	"cluster"
	"config"
	"context"
	"controllers/servers"
	"github.com/go-playground/lars"
	"log"
	"net/http"
	"os"
	"runtime"
	"server"
	"time"
)

func main() {

	var cluster cluster.Cluster

	s1 := server.NewServer("0.0.0.0", 27961)
	s2 := server.NewServer("0.0.0.0", 27962)

	cluster.AddServer(s1)
	cluster.AddServer(s2)

	cluster.Startup()

	s1.Console("map q3dm6")
	s2.Console("map q3dm6")

	router := lars.New()

	router.SetRedirectTrailingSlash(false)

	router.RegisterContext(context.NewContext)
	router.RegisterCustomHandler(func(*context.Context) {}, context.CastContext)

	router.Use(func(context lars.Context) {

		logger := log.New(os.Stdout, "[web] ", log.Ldate|log.Lmicroseconds)

		t1 := time.Now()
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logger.Printf(" recovering from panic: %+v\nStack Trace:\n %s", err, trace[:n])
				return
			}
		}()
		context.Next()

		res := context.Response()
		req := context.Request()
		code := res.Status()

		t2 := time.Now()

		logger.Printf("%d [%s] %v %q %v %d\n", code, req.Method, req.RemoteAddr, req.URL, t2.Sub(t1), res.Size())

	})

	router.Use(func(context lars.Context) {
		context.Set("cluster", cluster)
		context.Next()
	})

	servers.Routes(router.Group("/servers"))

	http.ListenAndServe(config.Config().General.Listen, router.Serve())

	cluster.Shutdown()
}
