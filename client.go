package main

import (
	"io"
	"net"
)
import "fmt"
import "bufio"
import "os"

var PRINT_DEBUG = false

type Client struct {
	Connection net.Conn
	IncomingMessageChannel chan string
	IncomingMessageReader *bufio.Reader
	OutgoingMessageChannel chan string
	OutgoingMessageReader *bufio.Reader
	Connected bool

}


func main() {
	conn, _ := net.Dial("tcp", "localhost:8081")
	closeChannel := make(chan bool)

	client := Client{
		conn,
		make(chan string, 1),
		bufio.NewReader(conn),
		make(chan string, 1),
		bufio.NewReader(os.Stdin),
		true,
	}

	go client.ListenForIncomingMessage(closeChannel)
	go client.PrintIncomingMessage()
	go client.SendOutgoingMessage()
	go client.WaitforInput()

	for {
		//client stays open until listenForIncomingMessages gets an error
		closeClient := <- closeChannel
		if closeClient {
			fmt.Println("Chat server disconnected.. Goodbye!")
			return
		}
	}

}


func (c *Client) WaitforInput() {
	for {
		text, _ := c.OutgoingMessageReader.ReadString('\n')
		if !c.Connected {
			return
		}
		if PRINT_DEBUG {
			fmt.Println("read inputted message")
		}
		c.OutgoingMessageChannel <- text
	}
}


func (c *Client) SendOutgoingMessage() {
	for {
		outgoingMessageText := <- c.OutgoingMessageChannel
		if !c.Connected {
			return
		}
		if PRINT_DEBUG {
			fmt.Println("sending outgoing string")
		}
		fmt.Fprintf(c.Connection, outgoingMessageText)
	}
}


func (c *Client) ListenForIncomingMessage(close chan bool ) {
	for {
		message, err := c.IncomingMessageReader.ReadString('\n')
		if err != nil {
			fmt.Println("connection error on listenForIncomingMessages, stopping to listen")
			if err == io.EOF {
				fmt.Println("connection closed normally")
				c.Connected = false
				c.Connection.Close()
				close <- true
			} else {
				fmt.Println("connection closed strangely.. error: ", err)
			}
			return
		}
		if PRINT_DEBUG {
			fmt.Println("read incoming message")
		}
		c.IncomingMessageChannel <- message
	}
}


func (c *Client) PrintIncomingMessage() {
	for {
		message := <- c.IncomingMessageChannel
		if !c.Connected {
			return
		}
		fmt.Print("> ", message)
	}
}



