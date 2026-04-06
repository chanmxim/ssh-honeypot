package main

import (
	"log"
	"github.com/chanmxim/ssh-honeypot/internal/ssh"
)

func main(){
	log.Println("Initiating a honeypot...")

	err := ssh.StartServer("2222")
	if err != nil{
		log.Fatalf("[-] Fatal server error: %v\n", err)
	}
}
