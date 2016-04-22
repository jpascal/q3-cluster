package context

import (
	"cluster"
	"github.com/go-playground/lars"
	"storage"
	"translator"
)

type Context struct {
	*lars.Ctx
}

func (self *Context) Cluster() *cluster.Cluster {
	value, _ := self.Get("cluster")
	return value.(*cluster.Cluster)
}

func (self *Context) Translator() *translator.Translator {
	value, _ := self.Get("translator")
	return value.(*translator.Translator)
}

func (self *Context) Storage() *storage.Storage {
	value, _ := self.Get("storage")
	return value.(*storage.Storage)
}

func NewContext(l *lars.LARS) lars.Context {
	return &Context{
		Ctx: lars.NewContext(l),
	}
}

// casts custom context and calls you custom handler so you don;t have to type cast lars.Context everywhere
func CastContext(c lars.Context, handler lars.Handler) {
	// could do it in all one statement, but in long form for readability
	h := handler.(func(*Context))
	ctx := c.(*Context)
	h(ctx)
}
