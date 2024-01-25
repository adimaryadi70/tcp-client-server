package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net"
)

const key = "SDJIWJIADJIAJDIAJDIJIW@I)@)@)DKIWJDIJDIAJDI"

func handleConnection(conn net.Conn, messages chan<- string) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s\n", clientAddr)

	//Send a welcome message to the client
	conn.Write([]byte("Selamat Datang Gateway Komunikasi!\n"))

	// Listen for incoming messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()

		fmt.Printf("Received from %s: %s\n", clientAddr, message)

		// Send the received message to all connected clients
		messages <- message
	}
}

func broadcastMessages(messages <-chan string, clients map[net.Conn]bool) {
	for message := range messages {
		// Broadcast the message to all connected clients
		for client := range clients {
			_, err := client.Write([]byte(message + "\n"))
			if err != nil {
				// Handle error, e.g., client disconnected
				delete(clients, client)
				fmt.Printf("Client %s disconnected\n", client.RemoteAddr())
			}
		}
	}
}

func main() {
	clients := make(map[net.Conn]bool)
	messages := make(chan string, 10)

	listenAddr := ":8081"
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("TCP server is listening on %s\n", listenAddr)

	go broadcastMessages(messages, clients)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		clients[conn] = true
		go handleConnection(conn, messages)
	}
}

func AES256Decrypt(key, crypt string) (string, error) {
	block, err := aes.NewCipher(getPaddedKey(key, 256))
	if err != nil {
		return "", err
	}
	dcrypt, err := hex.DecodeString(crypt)
	if err != nil {
		return "", err
	}
	ecb := cipher.NewCBCDecrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	decrypted := make([]byte, len(dcrypt))
	ecb.CryptBlocks(decrypted, dcrypt)

	return string(PKCS5Trimming(decrypted)), nil
}
func getPaddedKey(key string, bit int) []byte {

	paddedLen := bit/8 - len(key)
	if paddedLen < 0 {
		return []byte(key[0 : bit/8])
	}

	for i := 0; i < paddedLen; i++ {
		key = key + "f"
	}

	return []byte(key)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
