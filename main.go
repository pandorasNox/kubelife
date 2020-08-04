/*
Copyright 2020 The Kubelife Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"log"

	"github.com/pandorasnox/kubelife/pkg/ssh"
)

func main() {
	user := flag.String("user", "", "remote server login user")
	addr := flag.String("addr", "", "remote server address (ip/dns)")
	flag.Parse()

	ssh, err := ssh.New(*user, *addr, ssh.AgentAuth())
	if err != nil {
		log.Fatal(err)
	}
	_ = ssh
}
