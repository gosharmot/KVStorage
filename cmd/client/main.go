package main

import (
	kvs "KVStorage/pkg/KVStorage"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8085")
	if err != nil {
		log.Fatal(err)
	}
	for {
		var key, value, method string
		fmt.Print("Введите метод: ")
		fmt.Scanln(&method)
		switch method {
		case "get":
			method = kvs.METHODGET
			fmt.Print("Введите ключ: ")
			fmt.Scanln(&key)
		case "put":
			method = kvs.METHODPUT
			fmt.Print("Введите ключ: ")
			fmt.Fscanln(os.Stdin, &key)
			fmt.Print("Введите значение: ")
			fmt.Fscan(os.Stdin, &value)
		case "delete":
			method = kvs.METHODDELETE
			fmt.Print("Введите ключ: ")
			fmt.Scanln(&key)
		default:
			continue
		}
		msg := &kvs.Message{
			Method: method,
			Key:    key,
			Value:  value,
		}
		err := json.NewEncoder(conn).Encode(msg)
		if err != nil {
			fmt.Println(err)
		}
		if method == kvs.METHODGET {
			err = json.NewDecoder(conn).Decode(msg)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(msg)
		}
	}
}
