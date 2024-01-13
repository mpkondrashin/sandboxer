package main

import (
	"log"
	"time"

	ipc "github.com/james-barrow/golang-ipc"
)

func main() {
	go server()
	c, err := ipc.StartClient("example1", nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		log.Println("client status: ", c.Status())
		if c.Status() == "Connected" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	for {
		log.Println("client Sending Hello")
		err := c.Write(1, []byte("Hello"))
		if err != nil {
			log.Println("client Write error: ", err)
		}
		time.Sleep(1 * time.Second)
		log.Println("client attempt to read")
		message, err := c.Read()
		if err != nil {
			log.Println("client read error: ", err)
			break
		}
		if message.MsgType == -1 {
			log.Println("client status", c.Status())
			if message.Status == "Reconnecting" {
				c.Close()
				return
			}
		} else {
			log.Println("Client received: "+string(message.Data)+" - Message type: ", message.MsgType)
			c.Write(5, []byte("Message from client - PONG"))
		}

	}
}

func server() {
	s, err := ipc.StartServer("example1", nil)
	if err != nil {
		log.Println("server error", err)
		return
	}
	log.Println("server status", s.Status())
	for {
		message, err := s.Read()
		if err != nil {
			log.Print("server read error: ", err)
			break
		}
		if message.MsgType == -1 {
			if message.Status == "Connected" {
				log.Println("server status", s.Status())
				s.Write(1, []byte("server - PING"))
			}
		} else {
			log.Println("Server received: "+string(message.Data)+" - Message type: ", message.MsgType)
			s.Close()
			return
		}
	}
}
