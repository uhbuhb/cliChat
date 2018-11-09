package main

import "net"
import "fmt"
import "bufio"
import "os"


var PRINT_DEBUG = false

func main() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "localhost:8081")

	incomingMessageChannel := make(chan string)
	incomingMessageReader := bufio.NewReader(conn)

	outgoingMessageChannel := make(chan string)
	outgoingMessageReader := bufio.NewReader(os.Stdin)


	go ListenForIncomingMessage(incomingMessageReader, incomingMessageChannel)
	//go WaitForOutgoingMessage(outgoingMessageReader, outgoingMessageChannel)

	go printIncomingMessage(incomingMessageChannel)
	go sendOutgoingMessage(conn, outgoingMessageChannel)

	for {
		text, _ := outgoingMessageReader.ReadString('\n')
		if PRINT_DEBUG {
			fmt.Println("read inputted message")
		}
		outgoingMessageChannel <- text
	}

}


func WaitForOutgoingMessage(reader *bufio.Reader, ch chan string) {
	for {
		text, _ := reader.ReadString('\n')
		if PRINT_DEBUG {
			fmt.Println("read inputted message")
		}
		ch <- text
	}
}

func sendOutgoingMessage(conn net.Conn, ch chan string) {
	for {
		outgoingMessageText := <- ch
		if PRINT_DEBUG {
			fmt.Println("sending outgoing string")
		}
		fmt.Fprintf(conn, outgoingMessageText)
	}

}


func ListenForIncomingMessage(reader *bufio.Reader, ch chan string) {
	for {
		message, _ := reader.ReadString('\n')
		if PRINT_DEBUG {
			fmt.Println("read incoming message")
		}
		ch <- message
	}
}

func printIncomingMessage(ch chan string) {
	for {
		message := <- ch
		fmt.Println("Message received: ", message)

	}
}

