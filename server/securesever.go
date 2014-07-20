//
// Copyright 2014 Hong Miao. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Message chan string
type ClientTable map[net.Conn]*SecureClient

type SecureServer struct {
	listener net.Listener
	clients  ClientTable
	clientcoming chan net.Conn
	incoming Message
	
}

func CreateServer() *SecureServer {	
	server := &SecureServer {
		clients : make(ClientTable),
		clientcoming : make(chan net.Conn),
		incoming: make(Message),
	}
	
	return server
}

func (self *SecureServer)serverEvent() {
	for {
		select {
		case message := <-self.incoming:
			self.messageProcess(message)
		case conn := <- self.clientcoming:
			self.processClient(conn)
		}
	}
}

func (self *SecureServer)Listen(port string) {
	log.Printf("connport %s \n", port)
	self.listener, _ = net.Listen("tcp", port)	
	go self.serverEvent()
	for {
		conn, err := self.listener.Accept()
		log.Printf("Accept")
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		self.clientcoming <- conn
	}
}

func (self *SecureServer)sendToEveryClient(message string) {
	for _, client := range self.clients {
		client.outgoing <- message
	}
}

func (self *SecureServer)messageProcess(message string) {
	
	self.sendToEveryClient(message)
}

func (self *SecureServer)processClient(conn net.Conn) {
	log.Printf("processClient")
	client := CreateClient(conn)
	
	self.clients[conn] = client
	
	client.ClientEvent()
	go func() {
		for {
			msg := <-client.incoming
			
			if strings.HasPrefix(msg, "::") {
				cmd := CreateCmd()
				fmt.Println(cmd)
				msglist := strings.Split(msg, " ")
				fmt.Println(msglist)
				cmd.parseCmd(msglist)
				cmd.executeCommand(self, client)
				
			} else {
				log.Printf("Got message: %s\n", msg)
				self.incoming <- msg
			}
		}
	}()

		
}