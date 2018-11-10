package main

import (
	"io"
	"net"
)
import "fmt"
import "bufio"
import "os"


var PRINT_DEBUG = true

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

	client := Client{
		conn,
		make(chan string),
		bufio.NewReader(conn),
		make(chan string),
		bufio.NewReader(os.Stdin),
		true,
	}

	go client.ListenForIncomingMessage()
	go client.printIncomingMessage()
	go client.SendOutgoingMessage()

	for {
		//in need to figure out a way to alert user to input message..
		//fmt.Print("Input>")
		text, _ := client.OutgoingMessageReader.ReadString('\n')
		if !client.Connected {
			fmt.Println("Chat server disconnected.. Goodbye!")
			return
		}
		if PRINT_DEBUG {
			fmt.Println("read inputted message")
		}
		client.OutgoingMessageChannel <- text
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


func (c *Client) ListenForIncomingMessage() {
	for {
		message, err := c.IncomingMessageReader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading incoming message, stopping to listen")
			if err == io.EOF {
				fmt.Println("connection closed normally")
				c.Connected = false
				c.Connection.Close()
			} else {
				fmt.Println("got different error: ", err)
			}
			return
		}
		if PRINT_DEBUG {
			fmt.Println("read incoming message")
		}
		c.IncomingMessageChannel <- message
	}
}

func (c *Client) printIncomingMessage() {
	for {
		message := <- c.IncomingMessageChannel
		if !c.Connected {
			return
		}
		fmt.Print("Message received: ", message)


	}
}

