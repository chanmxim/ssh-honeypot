package ssh

import (
	"fmt"
	"log"
	"net"
)

func StartServer(port string) error{
	address := fmt.Sprintf("0.0.0.0:%s", port)

	// Open listener and resource cleanup
	listener, err := net.Listen("tcp", address)
	if err != nil{
		return fmt.Errorf("[-] Failed to bind to port: %w", err)
	}
	defer listener.Close()

	log.Printf("Honeypot listening for prey on %s...\n", address)

	// Accept loop
	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Fatal(err)
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn){
	defer conn.Close()

	remoteIP := conn.RemoteAddr().String()
	log.Printf("[+] New connection from: %s\n", remoteIP)

	// Send fake ssh banner
	conn.Write([]byte("Fake ssh banner\n"))
}