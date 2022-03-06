package main

import (
	kvs "KVStorage/pkg/KVStorage"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	message = make(chan *RequestMessage)
)

type RequestMessage struct {
	msg        *kvs.Message
	clientChan chan *kvs.Message
}

func main() {
	storage := kvs.NewStorage()
	mutex := sync.Mutex{}

	listener, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster(storage, &mutex)
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster(store *kvs.Storage, m *sync.Mutex) {
	for {
		select {
		case req := <-message:
			switch req.msg.Method {
			case kvs.METHODPUT:
				m.Lock()
				store.Put(req.msg.Ctx, req.msg.Key, req.msg.Value)
				m.Unlock()
			case kvs.METHODDELETE:
				m.Lock()
				store.Delete(req.msg.Ctx, req.msg.Key)
				m.Unlock()
			case kvs.METHODGET:
				m.Lock()
				res, err := store.Get(req.msg.Ctx, req.msg.Key)
				m.Unlock()
				if err != nil {
					log.Fatal(err)
				}
				req.clientChan <- &kvs.Message{
					Ctx:    req.msg.Ctx,
					Method: req.msg.Method,
					Key:    req.msg.Key,
					Value:  res,
				}
			default:
				log.Println("неверный метод")
				continue
			}
		}
	}
}

func handleConn(conn net.Conn) {
	fmt.Printf("%s connect\n", conn.LocalAddr().String())
	defer conn.Close()
	for {
		ch := make(chan *kvs.Message)
		go func(conn net.Conn, ch <-chan *kvs.Message) {
			for msg := range ch {
				err := json.NewEncoder(conn).Encode(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
		}(conn, ch)
		msg := &kvs.Message{}
		err := json.NewDecoder(conn).Decode(msg)
		if err != nil {
			log.Println(err)
			break
		}
		message <- &RequestMessage{
			msg:        msg,
			clientChan: ch,
		}
	}
}
