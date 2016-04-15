package context

import (
	"github.com/go-playground/lars"
)

type Context struct {
	*lars.Ctx
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

