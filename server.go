package main

import (
	"bufio"
	"fmt"
	"net"
)

type ChatServer struct {
	IncomingConnectionChannel chan net.Conn
	IncomingMessageChannel chan string
	Clients []ChatClient

}

type ChatClient struct {
	Connection net.Conn
}


func main() {
	//init chatObject
	//forever: listen for new connections
		//on new connection add it to chatObject

	//chatObject
		//listens on incoming connection channel
		//listens on incoming message channel
		//broadcast message feature
		//holds an array of clients

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
	go s.BroadcastIncomingMessages()



}


func (s *ChatServer) ListenForIncomingConnections() {
	for {
		incomingConnection := <- s.IncomingConnectionChannel
		fmt.Println("connection came")
		reader := bufio.NewReader(incomingConnection)
		go s.ListenForIncomingMessage(reader)
		s.Clients = append(s.Clients, ChatClient{incomingConnection})


	}
}

func (s *ChatServer) ListenForIncomingMessage(reader *bufio.Reader) {
	for {
		msg, err := reader.ReadString('\n')
		fmt.Println("message came: ", msg)
		if err != nil {
			fmt.Println("error reading from reader", err)
		}
		//s.IncomingMessageChannel <- msg //this doesnt work for some reason..
		go s.BroadcastIncoming(msg)
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
func (s *ChatServer) BroadcastIncoming(message string) {
	fmt.Println("broadcasting message: ", message)
	for _, client := range s.Clients {
		fmt.Fprintf(client.Connection, message)
	}
}








//func (s ChatServer) handle(conn net.Conn) {
//	//this starts right after connection was opened
//	fmt.Println("handling connection..")
//	s.Connections = append(s.Connections, conn)
//
//	for {
//		message, err := bufio.NewReader(conn).ReadString('\n')
//		if err != nil {
//			s.RemoveConnection(conn)
//		}
//		fmt.Print("Message Received: ", string(message))
//
//		s.broadcast(message)
//		fmt.Println("broadcasted")
//	}
//}
//
//func (s ChatServer) broadcast(message string){
//	fmt.Print("Broadcasting message: ", message)
//	for _, value := range s.Connections {
//		value.Write([]byte(message))
//
//	}
//}
//func (s ChatServer) RemoveConnection(conn net.Conn) {
//	for i, v := range s.Connections {
//		if v == conn {
//			s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
//		}
//	}
//}






