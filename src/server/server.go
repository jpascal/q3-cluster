package server

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
	"math/rand"
	"os"
	"net"
	"strings"
	"strconv"
)

type ServerStatus struct {
	TimeLimit   int    `json:"time_limit"`
	FragLimit   int    `json:"frag_limit"`
	MaxClients  int    `json:"max_clients"`
	Protocol    int    `json:"protocol"`
	GameVersion string `json:"game_version"`
	Version     string `json:"version"`
	MapName     string `json:"map_name"`
	HostName    string `json:"host_name"`
	ProMode     int    `json:"pro_mode"`
}

type Server struct {
	Address  string
	Port     int
	Instance *exec.Cmd
	Stdout   io.WriteCloser
	Stdin    io.ReadCloser
	Password string
	Logger 	 *log.Logger
}

func NewServer(address string, port int) *Server {
	var server Server
	server.Address = address
	server.Port = port

	server.Password = randStringRunes(16)

	server.Instance = exec.Command("./test.sh",
		"+set net_noudp 0",
		"+set sv_strictAuth 0",
		"+set dedicated 1",
		"+set sv_punkbuster 0",
		"+set sv_lanForceRate 0",
		fmt.Sprintf("+set net_ip %v", server.Address),
		fmt.Sprintf("+set net_port %v", server.Port),
	)

	server.Logger = log.New(os.Stdout, fmt.Sprintf("[%v:%v] ",server.Address, server.Port), log.Ldate | log.Lmicroseconds)

	server.Stdout, _ = server.Instance.StdinPipe()
	server.Stdin, _ = server.Instance.StdoutPipe()
	return &server
}

func (self *Server) send(data []byte) ([]byte, error) {
	var connection net.Conn
	var err error
	if connection, err = net.Dial("udp", fmt.Sprintf("%v:%v",self.Address,self.Port)); err != nil {
		return nil, err
	}
	defer connection.Close()

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
	self.Logger.Printf("get status")
	buffer, err := self.send([]byte("\xff\xff\xff\xffgetstatus"))
	if err != nil {
		return nil, err
	}
	var server_status ServerStatus
	response := strings.Split(strings.Trim(string(buffer), "\\n"), "\\")[1:]
	var name, value string
	for len(response) > 0 {
		value, response = response[len(response)-1], response[:len(response)-1]
		name, response = response[len(response)-1], response[:len(response)-1]

		switch name {
		case "fraglimit":
			server_status.FragLimit, _ = strconv.Atoi(value)
		case "timelimit":
			server_status.TimeLimit, _ = strconv.Atoi(value)
		case "sv_maxclients":
			server_status.MaxClients, _ = strconv.Atoi(value)
		case "com_protocol":
			server_status.Protocol, _ = strconv.Atoi(value)
		case "gameversion":
			server_status.GameVersion = value
		case "mapname":
			server_status.MapName = value
		case "version":
			server_status.Version = value
		case "hostname":
			server_status.HostName = value
		case "server_promode":
			server_status.ProMode, _ = strconv.Atoi(value)
		}
	}
	return &server_status, nil
}

func (self *Server) RemoteConsole(command string) (string, error) {
	log.Printf("rcon(%s): %s", self.Address, command)
	buffer, err := self.send([]byte(fmt.Sprintf("\xff\xff\xff\xffrcon %s %s\n", self.Password, command)))
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}


func (self *Server) Startup() error {

	self.Logger.Printf("executing server with password '%s'", self.Password)

	if err := self.Instance.Start(); err != nil {
		self.Logger.Printf("unable to run: %s", err.Error())
	} else {
		self.Logger.Printf("started with pid %v", self.Instance.Process.Pid)
	}

	return nil
}

func (self *Server) HasInstance() bool {
	return self.Instance.Process != nil
}

func (self *Server) Shutdown() {
	if self.HasInstance() {
		self.Logger.Printf("stopping process %v", self.Instance.Process.Pid)
		self.Instance.Wait()
	}
}

func (self *Server) Console(command string) error {
	self.Logger.Printf("command: %v", command)
	self.Stdout.Write([]byte(fmt.Sprintf("%s\n", command)))
	return nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}