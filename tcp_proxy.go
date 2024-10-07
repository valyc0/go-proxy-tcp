package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func handleConnection(src net.Conn, remoteAddress string) {
	defer src.Close()

	// Connessione al server remoto
	dst, err := net.Dial("tcp", remoteAddress)
	if err != nil {
		log.Printf("Impossibile connettersi a %s: %v", remoteAddress, err)
		return
	}
	defer dst.Close()

	// Copia i dati dal client al server remoto e viceversa
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Printf("Errore durante il forwarding dal client: %v", err)
		}
	}()
	if _, err := io.Copy(src, dst); err != nil {
		log.Printf("Errore durante il forwarding dal server: %v", err)
	}
}

func startProxy(localPort string, remoteAddress string) {
	listener, err := net.Listen("tcp", ":"+localPort)
	if err != nil {
		log.Fatalf("Impossibile avviare il listener sulla porta %s: %v", localPort, err)
	}
	defer listener.Close()

	log.Printf("Proxy avviato sulla porta %s, inoltrando verso %s", localPort, remoteAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Errore durante l'accettazione della connessione: %v", err)
			continue
		}
		go handleConnection(conn, remoteAddress)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./tcp_proxy <local_port> <remote_address:remote_port>")
		return
	}

	localPort := os.Args[1]
	remoteAddress := os.Args[2]

	startProxy(localPort, remoteAddress)
}

