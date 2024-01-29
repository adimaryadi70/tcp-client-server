package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const key = "SDJIWJIADJIAJDIAJDIJIW@I)@)@)DKIWJDIJDIAJDI"

type SendMessage struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Services string `json:"services"`
	Data     string `json:"data"`
	Key      string `json:"key"`
}

func main() {
	host := "localhost"
	port := "8081"
	address := host + ":" + port
	timeSecond := 2 * time.Second
	portInt, _ := strconv.Atoi(port)
	if isPortOpen(host, portInt, timeSecond) {
		fmt.Printf("Port %d on %s is open\n", host, portInt)
	} else {
		fmt.Printf("Ip %d Port %s is Not Open\n", host, portInt)
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go ReadFromServer(conn)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		message := scanner.Text()
		dataSend := SendMessage{
			ID:       generateID(),
			Type:     "Request",
			Services: "authentication.login",
			Data:     message,
			Key:      key,
		}

		sender, err := json.Marshal(dataSend)
		textEncrypt, err := AES256Encrypt(key, string(sender))
		if err != nil {
			log.Println("Error Enryption:", err)
		}
		_, _ = fmt.Fprintf(conn, "%s\n", textEncrypt)
	}
}
func generateID() string {
	// Generate a new UUID
	id := uuid.New()

	// Convert the UUID to a string
	return id.String()
}
func isPortOpen(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func ReadFromServer(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Received From Server => ", message)
	}
}

func AES256Encrypt(key, src string) (string, error) {
	block, err := aes.NewCipher(getPaddedKey(key, 256))
	if err != nil {
		return "", err
	}
	ecb := cipher.NewCBCEncrypter(block, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	return hex.EncodeToString(crypted), nil
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
