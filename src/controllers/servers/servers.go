package servers

import (
	"context"
	"cluster"
	"github.com/go-playground/lars"
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
}

func Create(context *context.Context) {

}

func Update(context *context.Context) {

}

func Show(context *context.Context) {
	value, d := context.Get("cluster")
	c := value.(cluster.Cluster)
	log.Printf("Get %v -> %v", d, len(c.Servers))
}

func Destroy(context *context.Context) {

}

func Start(context *context.Context) {

}

func Stop(context *context.Context) {

}

func Rcon(context *context.Context) {

}
