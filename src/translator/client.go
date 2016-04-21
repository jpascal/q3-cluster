package translator

type Client struct {
	translator *Translator
}

func NewClient(translator *Translator) *Client {
	return &Client{translator: translator}
}


func (self *Client) Listen() {
	go func(){
		self.translator.Logger.Print("...")
	}()
}