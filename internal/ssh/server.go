package ssh

import (
	"golang.org/x/crypto/ssh"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
)

func StartServer(port string) error{
	// ==========================
	// SSH SERVER CONFIGURATION
	// ==========================
	config := &ssh.ServerConfig{
		PasswordCallback: func(metaData ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error){
			
			log.Printf("[!] Login attempt - User: [%s] Password: [%s]\n", metaData.User(), string(pass))

			// Returning (nil, nil) means successful authentication
			// Accept any username/password
			return nil, nil
		},
	}

	// Attach fake identity
	config.AddHostKey(generateSigner())
	

	// ==========================
	// TCP SERVER SETUP
	// ==========================
	// Open listener and clean up resources
	address := fmt.Sprintf("0.0.0.0:%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil{
		return fmt.Errorf("[-] Failed to bind to port: %w", err)
	}
	defer listener.Close()

	log.Printf("[*] Honeypot listening for prey on %s...\n", address)

	// Accept loop
	for {
		conn, err := listener.Accept()
		if err != nil{
			log.Fatal(err)
		}

		go handleConnection(conn, config)

	}
}

func handleConnection(conn net.Conn, config *ssh.ServerConfig){
	defer conn.Close()

	remoteIP := conn.RemoteAddr().String()
	log.Printf("[+] New connection from: %s\n", remoteIP)

	// SSH handshake
	_, _, _, err := ssh.NewServerConn(conn, config)
	if err != nil{
		log.Printf("[-] SSH handshake failed for %s: %v\n", remoteIP, err)
		return
	}

	log.Printf("[+] Attacker %s is authenticated\n", remoteIP)
}

func generateSigner() ssh.Signer{
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil{
		log.Fatalf("[-] Failed to generate RSA key: %v", err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil{
		log.Fatalf("[-] Failed to create signer: %v", err)
	}

	return signer
}