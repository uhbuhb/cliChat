package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type ChatServer struct {
	IncomingConnectionChannel chan net.Conn
	IncomingMessageChannel chan string
	Clients []ChatClient

}

type ChatClient struct {
	Connection net.Conn
	Reader *bufio.Reader
	Connected bool
}


func main() {
	//init chatObject
	//forever: listen for new connections
		//on new connection add it to chatObject

	//chatObject
		//holds an array of clients
		//listens on incoming connection channel
		//listens on incoming message channel
		//broadcasts incoming messages

	incomingConnectionChannel := make(chan net.Conn)

	server := ChatServer {IncomingConnectionChannel: incomingConnectionChannel}
	server.launchServer()

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("error starting to listen")
	}
	defer listener.Close()

	for {
		fmt.Println("waiting for client..")
		//waits here until client connects
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection")
		}
		incomingConnectionChannel <- conn
		fmt.Println("client received")
	}
	
}

func (s *ChatServer) launchServer() {
	s.Clients = make([]ChatClient, 0)

	go s.ListenForIncomingConnections()
	//go s.BroadcastIncomingMessages()

}


func (s *ChatServer) ListenForIncomingConnections() {
	for {
		incomingConnection := <- s.IncomingConnectionChannel
		fmt.Println("connection came")
		reader := bufio.NewReader(incomingConnection)
		client := ChatClient{incomingConnection, reader, true}
		go s.ListenForIncomingMessage(client)
		s.Clients = append(s.Clients, client)
		s.Broadcast("Welcome new user, type a message to send to group\n")

	}
}

func (s *ChatServer) ListenForIncomingMessage(client ChatClient) {
	for {
		msg, err := client.Reader.ReadString('\n')
		if err!= nil {
			if err == io.EOF {
				fmt.Println("user disconnected")
			} else {
				fmt.Println("error reading from reader: ", err)
			}
			client.Connected = false
			return
		}
		fmt.Println("message came: ", msg)
		if err != nil {
			fmt.Println("error reading from reader", err)
		}
		//s.IncomingMessageChannel <- msg //this doesnt work for some reason..
		go s.Broadcast(msg)
	}
}


func (s *ChatServer) BroadcastIncomingMessages(){
	for {
		fmt.Println("waiting for incoming message on channel")
		message := <- s.IncomingMessageChannel
		fmt.Println("broadcasting message: ", message)
		for _, client := range s.Clients {
			fmt.Fprintf(client.Connection, message)
		}
	}
}


func (s *ChatServer) Broadcast(message string) {
	fmt.Println("broadcasting message: ", message)
	for i, client := range s.Clients {
		if !client.Connected {
			s.RemoveClient(i, client)
			continue
		}
		fmt.Fprintf(client.Connection, message)
	}
}


func (s *ChatServer) RemoveClient(i int, client ChatClient) {
	client.Connection.Close()
	s.Clients = append(s.Clients[:i], s.Clients[i+1:]...)
}






