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
	"log"
	"net"
	"bufio"
	"fmt"
)

type SecureClient struct {
	conn net.Conn
	
	reader   *bufio.Reader
	writer   *bufio.Writer
	incoming chan string
	outgoing chan string
	id string
	name string
}


func CreateClient(conn net.Conn) *SecureClient {
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	client := &SecureClient {
		reader: r,
		writer: w,
		incoming : make(chan string),
		outgoing : make(chan string),
	}
	
	return client
}


func (self *SecureClient) ClientEvent() {
	go self.Read()
	go self.Write()
}

func (self *SecureClient) Read() {
	fmt.Println("Read")
	for {
		buf := make([]byte, 1000)
		n , err := self.reader.Read(buf)
		if err != nil {
			fmt.Println("Read error")
			//log.Printf("Read error: %s\n", err)
			return			
		}
		log.Printf("Read %d bytes\n", n)
		self.incoming <- string(buf)
		/*
		if line, _, err := self.reader.ReadLine(); err == nil {
			fmt.Println("###########")
			self.incoming <- string(line)
		} else {
			fmt.Println("Read error")
			//log.Printf("Read error: %s\n", err)
			return
		}
		*/
	}

}

func (self *SecureClient) Write() {
	
	for data := range self.outgoing {
		fmt.Println("Write")
		//_, err := self.writer.WriteString(data + "\n")
		_, err := self.writer.Write([]byte(data))
		
		if err != nil {
			log.Printf("Write error: %s\n", err)
			return
		}
		
		if err := self.writer.Flush(); err != nil {
			log.Printf("Write error: %s\n", err)
			return
		}
	}

}

func (self *SecureClient) GetName() string {
	return self.name
}

func (self *SecureClient) SetName(name string) {
	self.name = name
}
