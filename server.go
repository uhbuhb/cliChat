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
	Clients map[int]ChatClient
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
	s.Clients = make(map[int]ChatClient)
	go s.ListenForIncomingConnections()
}


func (s *ChatServer) ListenForIncomingConnections() {
	clientId := 1
	for {
		incomingConnection := <- s.IncomingConnectionChannel
		fmt.Println("connection came")
		reader := bufio.NewReader(incomingConnection)
		client := ChatClient{incomingConnection, reader, true}
		s.Clients[clientId] = client
		go s.ListenForIncomingMessage(clientId)
		s.Broadcast(-1, fmt.Sprintf("Welcome user%d, there are now %d users\n", clientId, len(s.Clients)))
		clientId++
	}
}

func (s *ChatServer) ListenForIncomingMessage(clientId int) {
	for {
		client := s.Clients[clientId]
		msg, err := client.Reader.ReadString('\n')
		if err!= nil {
			if err == io.EOF {
				fmt.Println("user disconnected")
				s.RemoveClient(clientId)
			} else {
				fmt.Println("error reading from reader: ", err)
			}
			//fmt.Println("setting field to false")
			//client.Connected = false
			//s.Clients[clientId] = client
			return
		}
		msg = fmt.Sprintf("user%d: %s", clientId, msg)
		fmt.Println("message came: ", msg)
		if err != nil {
			fmt.Println("error reading from reader", err)
		}
		go s.Broadcast(clientId, msg)
	}
}

//
//func (s *ChatServer) BroadcastIncomingMessages(){
//	for {
//		fmt.Println("waiting for incoming message on channel")
//		message := <- s.IncomingMessageChannel
//		fmt.Println("broadcasting message: ", message)
//		for _, client := range s.Clients {
//			fmt.Fprintf(client.Connection, message)
//		}
//	}
//}


func (s *ChatServer) Broadcast(senderId int, message string) {
	fmt.Println("broadcasting message: ", message)
	for clientId, client := range s.Clients {
		if !client.Connected {
			fmt.Println("found disconnected client.. removing")
			s.RemoveClient(clientId)
			continue
		}
		if clientId == senderId {
			continue
		}
		fmt.Fprintf(client.Connection, message)
	}
}


func (s *ChatServer) RemoveClient(clientId int) {
	fmt.Println("removing client")
	s.Clients[clientId].Connection.Close()
	delete(s.Clients, clientId)
	s.Broadcast(-1, fmt.Sprintf("User%d left, there are now %d users in chat\n", clientId, len(s.Clients)))
}


