package rpcserver

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	. "rpcdef"
)

type Server struct {
	srv      *rpc.Server
	listener net.Listener //needed for shutdown
	users    *Users
}

func (s *Server) Launch(addr string, data string) error { //launch server and start listening
	log.Println("Initializing server")
	s.users = new(Users)
	err := s.users.Init(data) //load users
	if err != nil {
		log.Println("Error initializing users: " + err.Error())
		return err
	}

	s.srv = new(rpc.Server)
	s.srv.Register(s.users)

	s.srv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Println("Listening error: " + err.Error())
		return err
	}
	go s.Serve() //main server loop
	return nil
}

func (s *Server) Serve() { //main server loop
	for {
		conn, err := s.listener.Accept()
		if err != nil { //specifically when listener is closed
			log.Println(err.Error())
			return
		} else {
			log.Println("Accepted connection")
			go func() {
				s.srv.ServeCodec(jsonrpc.NewServerCodec(conn))
			}()
		}
	}
}

func (s *Server) Shutdown(data string) { //shut down server by closing listener
	s.listener.Close()
	s.users.Finalize(data) //save users
}
