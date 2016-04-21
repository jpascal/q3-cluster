package server

import (
	"bufio"
	"config"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	Address  string         `json:"address"`
	Port     int            `json:"port"`
	Instance *exec.Cmd      `json:"-"`
	Stdin    io.WriteCloser `json:"-"`
	Stdout   io.ReadCloser  `json:"-"`
	Password string         `json:"-"`
	Logger   *log.Logger    `json:"-"`
}

func NewServer(address string, port int) *Server {
	var server Server
	server.Address = address
	server.Port = port

	server.Logger = log.New(os.Stdout, fmt.Sprintf("[%v:%v] ", server.Address, server.Port), log.Ldate|log.Lmicroseconds)

	server.Password = randStringRunes(16)

	arguments := config.Config().Cluster.Arguments

	arguments = strings.Replace(arguments, "$address", server.Address, 1)
	arguments = strings.Replace(arguments, "$port", fmt.Sprint(server.Port), 1)

	server.Instance = exec.Command(config.Config().Cluster.Server, arguments)

	log := config.Config().Cluster.Log
	log = strings.Replace(log, "$address", server.Address, 1)
	log = strings.Replace(log, "$port", fmt.Sprint(server.Port), 1)

	server.Stdin, _ = server.Instance.StdinPipe()
	server.Stdout, _ = server.Instance.StdoutPipe()

	return &server
}

func (self *Server) send(data []byte) ([]byte, error) {
	var connection net.Conn
	var err error
	if connection, err = net.Dial("udp", fmt.Sprintf("%v:%v", self.Address, self.Port)); err != nil {
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
	buffer, err := self.send([]byte("\xff\xff\xff\xffgetstatus\n"))
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
	self.Logger.Printf("remote command: '%s'", command)
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
		return err
	}

	self.Logger.Printf("started with pid %v", self.Instance.Process.Pid)

	scanner := bufio.NewScanner(self.Stdout)
	go func() {
		for scanner.Scan() {
			self.Logger.Printf("console: %v", scanner.Text())
		}
	}()

	return nil
}

func (self *Server) HasInstance() bool {
	return self.Instance.Process != nil
}

func (self *Server) Shutdown() {
	if self.HasInstance() {
		self.Logger.Printf("stopping process %v", self.Instance.Process.Pid)
		self.Instance.Process.Kill()
		self.Instance.Process = nil
	}
}

func (self *Server) Console(command string) error {
	if n, err := self.Stdin.Write([]byte(fmt.Sprintf("%s\n", command))); err != nil {
		return err
	} else {
		self.Logger.Printf("send to server (%d) command: %v", n, command)
	}
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
