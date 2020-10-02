package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Server struct {
	listener    net.Listener
	exit        chan bool
	connections map[int]net.Conn
}

func InitServer() *Server {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("ERROR : failed to create listener", err.Error())
	}
	srv := &Server{
		listener:    l,
		exit:        make(chan bool),
		connections: map[int]net.Conn{},
	}
	go srv.serve()
	return srv
}

func (srv *Server) serve() {
	var id int
	log.Println("INFO : Listening for clients")
	for {
		select {
		case <-srv.exit:
			log.Println("INFO : Shutting down the server")
			err := srv.listener.Close()
			if err != nil {
				log.Println("ERROR : Could not close listener", err.Error())
			}
			if len(srv.connections) > 0 {
				srv.closeConnections()
			}
			return

		default:
			tcpListener := srv.listener
			conn, err := tcpListener.Accept()
			if err != nil {
				log.Println("ERROR : Failed to accept connection:", err.Error())
			}
			srv.connections[id] = conn
			write(conn, "Welcome to TCP server")
			go func(connID int) {
				log.Println("INFO : Client with id:", connID, "joined")
				srv.handleConnection(conn)
				delete(srv.connections, connID)
				log.Println("INFO : Client with id:", connID, "left")
			}(id)
			id++
		}
	}

}

func (srv *Server) handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		value := scanner.Text()
		switch value {
		case "input":
			write(conn, "output: "+value)
		case "exit":
			err := conn.Close()
			if err != nil {
				log.Println("ERROR : Could not close connection", err.Error())
			}
		}
	}
}

func write(c net.Conn, s string) {
	_, err := fmt.Fprintf(c, "%s\n-> ", s)
	if err != nil {
		log.Fatal(err)
	}
}

func (srv *Server) closeConnections() {
	log.Println("INFO : Closing all connections")
	for id, conn := range srv.connections {
		err := conn.Close()
		if err != nil {
			log.Println("ERROR : Could not close connection with id:", id, err.Error())
		}

	}
}

func (srv *Server) Stop() {
	close(srv.exit)
}
