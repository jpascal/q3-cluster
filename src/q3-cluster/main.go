package main

import (
	////"gopkg.in/redis.v3"
	//"log"
	//"net"
	//"strings"
	//"strconv"
	//"fmt"
	//"github.com/go-playground/lars"
	//"net/http"
	//"context"

)

import (
	"cluster"
	"server"
	"context"
	"net"
	"strings"
	"strconv"
	"log"
	"fmt"
	servers_controller "controllers/servers"
	wlogger "github.com/go-playground/lars/examples/middleware/logging-recovery"
	"github.com/go-playground/lars"
	"net/http"
)

type Player struct {

}

type Team struct {
	Players	[]Player
}

type ServerStatus struct {
	TimeLimit		int		`json:"time_limit"`
	FragLimit		int		`json:"frag_limit"`
	MaxClients		int		`json:"max_clients"`
	Protocol		int		`json:"protocol"`
	GameVersion		string	`json:"game_version"`
	Version			string	`json:"version"`
	MapName			string	`json:"map_name"`
	HostName		string	`json:"host_name"`
	ProMode			int		`json:"pro_mode"`
}

type Server struct {
	Address string
	Password string

}

func NewServer(address, password string) *Server {
	return &Server{Address: address, Password: password}
}

func (self *Server) send(data []byte) ([]byte, error) {
	var connection net.Conn
	var err error
	if connection, err = net.Dial("udp", self.Address); err != nil {
		return nil, err
	}
	defer  connection.Close()

	if _, err = connection.Write(data); err != nil {
		return nil, err
	}

	buffer := make([]byte, 1024)
	if _, err = connection.Read(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

func (self *Server) GetStatus() (*ServerStatus, error) {
	buffer, err := self.send([]byte("\xff\xff\xff\xffgetstatus"))
	if err != nil {
		return nil, err
	}
	var server_status ServerStatus
	response := strings.Split(strings.Trim(string(buffer),"\\n"),"\\")[1:]
	var name, value string
	for len(response) > 0 {
		value, response = response[len(response)-1], response[:len(response)-1]
		name, response = response[len(response)-1], response[:len(response)-1]

		switch(name) {
			case "fraglimit": server_status.FragLimit, _ = strconv.Atoi(value)
			case "timelimit": server_status.TimeLimit, _ = strconv.Atoi(value)
			case "sv_maxclients": server_status.MaxClients, _ = strconv.Atoi(value)
			case "com_protocol": server_status.Protocol, _ = strconv.Atoi(value)
			case "gameversion": server_status.GameVersion = value
			case "mapname": server_status.MapName = value
			case "version": server_status.Version = value
			case "hostname": server_status.HostName = value
			case "server_promode": server_status.ProMode, _ = strconv.Atoi(value)
		}
	}
	return &server_status, nil
}

func (self *Server) Command(command string) (string, error) {
	log.Printf("rcon(%s): %s", self.Address, command)
	buffer, err := self.send([]byte(fmt.Sprintf("\xff\xff\xff\xffrcon %s %s\n",self.Password, command)))
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}


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

	router.Use(func (context lars.Context){
		context.Set("cluster", cluster)
		context.Next()
	})

	router.Use(wlogger.LoggingAndRecovery)

	servers_controller.Routes(router.Group("/servers"))

	http.ListenAndServe(":3007", router.Serve())

	cluster.Shutdown()

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "", // no password set
	//	DB:       0,  // use default DB
	//})
	//pong, err := client.Ping().Result()
	//log.Println(pong, err)
	//
	//if server_status, err := NewServer("q3.turnir.pro:27961", "1").GetStatus(); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(server_status)
	//}
	//
	//if server_status, err := NewServer("q3.turnir.pro:27962", "1").GetStatus(); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(server_status)
	//}
	//
	//if server_status, err := NewServer("q3.turnir.pro:27963", "1").GetStatus(); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(server_status)
	//}
	//
	//if response, err := NewServer("q3.turnir.pro:27963", "1").Command("ban"); err != nil {
	//	log.Println(err)
	//} else {
	//	log.Printf(response)
	//}
}