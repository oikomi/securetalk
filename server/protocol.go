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
	//"log"
	"fmt"
	"errors"
)

type CmdFuc func(server *SecureServer, client *SecureClient, args []string)

var (
	commands map[string]CmdFuc
)

func init() {
	commands = map[string]CmdFuc{
		"name": setName,
	}
}

type Cmd struct {	
	cmd string
	args []string
}

func CreateCmd() *Cmd {
	return &Cmd {
		
	}
}

//func (self *Cmd)addCmd(cmd string) {
//	commands[cmd]
//}

func (self *Cmd)parseCmd(msglist []string) {
	self.cmd = msglist[1]
	self.args = msglist[2:]
}

func (self *Cmd) executeCommand(server *SecureServer, client *SecureClient) (err error) {
	fmt.Println(self.cmd)
	if f, ok := commands[self.cmd]; ok {
		fmt.Println("-1--------")
		f(server, client, self.args)
		return
	}

	err = errors.New("Unsupported command: " + self.cmd)
	return
}

func setName(server *SecureServer, client *SecureClient, args []string) {
	fmt.Println("setName")
	oldname := client.GetName()
	client.SetName(args[0])
	server.messageProcess(fmt.Sprintf("Notification: %s changed its name to %s", oldname, args[0]))
}

