package server

import (
	"log"
	"os/exec"
	"io"
	"fmt"
)

type Server struct {
	Address		string
	Port		int
	Instance	*exec.Cmd
	Stdout		io.WriteCloser
}

func NewServer(address string, port int) *Server {
	var server Server
	server.Address = address
	server.Port = port

	server.Instance = exec.Command("/Users/e.shurmin/Projects/stats_s/test.sh",
		"+set net_noudp 0",
		"+set sv_strictAuth 0",
		"+set dedicated 1",
		"+set sv_punkbuster 0",
		"+set sv_lanForceRate 0",
		fmt.Sprintf("+set net_ip %v",server.Address),
		fmt.Sprintf("+set net_port %v",server.Port),
	)

	server.Stdout, _ = server.Instance.StdinPipe()
	return &server
}

func (self *Server) Startup() (error) {

	log.Printf("[%v:%v] executing %s ", self.Address, self.Port, self.Instance.Path)

	if err := self.Instance.Start(); err != nil {
		log.Printf("[%v:%v] unable to run: %s", self.Address, self.Port, err.Error())
	} else {
		log.Printf("[%v:%v] started with pid %v", self.Address, self.Port, self.Instance.Process.Pid)
	}

	return nil
}

func (self *Server) HasInstance() bool {
	return self.Instance.Process != nil
}

func (self *Server) Shutdown() {
	if self.HasInstance() {
		log.Printf("[%v:%v] stopping process %v", self.Address, self.Port, self.Instance.Process.Pid)
		self.Instance.Wait()
	}
}

func (self *Server) Console(command string) error {
	log.Printf("[%v:%v] command: %v", self.Address, self.Port, command)
	self.Stdout.Write([]byte(fmt.Sprintf("%s\n",command)))
	return nil
}

func (self *Server) RemoteConsole(command string) error {
	return nil
}
