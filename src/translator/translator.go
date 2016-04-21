package translator

import (
	"github.com/go-playground/lars"
	"log"
	"os"
)

type Translator struct {
	Logger	*log.Logger
}

func NewTranslator() *Translator {
	return &Translator{
		Logger: log.New(os.Stdout, "[translator] ", log.Ldate|log.Lmicroseconds),
	}
}

func Routes(routes lars.IRouteGroup) {
	routes.Any("",func(context lars.Context) {
		value, _ := context.Get("translator")
		translator := value.(*Translator)
		translator.Logger.Printf("new connection %v", context.Request().RemoteAddr)

		client := NewClient(translator)
		client.Listen()

		context.Done()
	})
}


