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
	"storage"
	"time"
	"translator"
)

func main() {

	storage := storage.NewStorage()
	cluster := cluster.NewCluster(storage)

	//s1 := server.NewServer("127.0.0.1", 27961)
	//s2 := server.NewServer("127.0.0.1", 27963)
	//
	//cluster.AddServer(s1)
	//cluster.AddServer(s2)

	cluster.Startup()

	//s1.Console("map q3dm6")
	//s2.Console("map q3dm6")
	//
	router := lars.New()

	router.SetRedirectTrailingSlash(false)

	router.RegisterContext(context.NewContext)
	router.RegisterCustomHandler(func(*context.Context) {}, context.CastContext)

	router.Use(func(context lars.Context) {
		context.Set("cluster", cluster)
		context.Set("storage", storage)
		context.Set("translator", translator.NewTranslator())
		context.Next()
	})

	translator.Routes(router.Group("/translator"))

	router.Use(func(context lars.Context) {

		logger := log.New(os.Stdout, "[rest] ", log.Ldate|log.Lmicroseconds)

		t1 := time.Now()
		defer func() {
			if err := recover(); err != nil {
				trace := make([]byte, 1<<16)
				n := runtime.Stack(trace, true)
				logger.Printf(" recovering from panic: %+v\nStack Trace:\n %s", err, trace[:n])
				context.Response().WriteHeader(500)
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

	servers.Routes(router.Group("/servers"))

	http.ListenAndServe(config.Config().General.Listen, router.Serve())

	cluster.Shutdown()
}
