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
	"os/signal"
	"syscall"
)

var signals = make(chan os.Signal, 1)

func main() {

	storage := storage.NewStorage()
	cluster := cluster.NewCluster(storage)

	signal.Notify(signals, syscall.SIGQUIT)
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for {
			switch <-signals {
			case syscall.SIGQUIT, os.Interrupt:
				cluster.Shutdown()
				os.Exit(1)
			case syscall.SIGTERM:
				os.Exit(1)
			}
		}
	}()

	cluster.Startup()

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

		logger := log.New(os.Stdout, "[http] ", log.Ldate|log.Lmicroseconds)

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
