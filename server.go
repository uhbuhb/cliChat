package main

import ( 
	"net"
	"fmt"
	"bufio"
)

type ChatServer struct {
	Connections []net.Conn

}


func main() {
	server := ChatServer {make([]net.Conn, 0)}
	server.launchServer()
	
}

func (s ChatServer) launchServer() {
	fmt.Println("Launching server...")
	
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
		go s.handle(conn)
		fmt.Println("client received, finished loop")
  	}
}

func (s ChatServer) handle(conn net.Conn) {
	//this starts right after connection was opened
	fmt.Println("handling connection..")
	s.Connections = append(s.Connections, conn)

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message Received: ", string(message))

	s.broadcast(message)
	fmt.Print("broadcasted")
}

func (s ChatServer) broadcast(message string){
	fmt.Print("Broadcasting message: ", message)
	for _, value := range s.Connections {
		value.Write([]byte(message + "\n"))
		
	}
}






