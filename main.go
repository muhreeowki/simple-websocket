package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("Incomming connection from client: ", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read err: ", err.Error())
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))

		s.broadcast(msg)
	}
}

func (s *Server) broadcast(msg []byte) {
	for conn := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := conn.Write(msg); err != nil {
				fmt.Println("write err: ", err.Error())
			}
		}(conn)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	fmt.Println("Listening at localhost:3000")

	http.ListenAndServe(":3000", nil)
}
