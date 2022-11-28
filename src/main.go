package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"server/src/controller"
	"server/src/proto"
)

func process(clients map[int]client, conn net.Conn, id int) {

	welcomeData, err := proto.Encode(pWelcome, clients[id].uuid.String())
	newPlayer(clients, id)

	fmt.Println(clients[id].uuid.String())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, e := conn.Write(welcomeData)
	if e != nil {
		fmt.Println("Error:", err)
		return
	}

	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		msg, pID, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("Полученные данные от клиента:", msg, pID)
		data, err := proto.Encode(pMessage, msg)
		m := controller.Message{msg, len(msg)}
		go controller.Save(&m)
		for _, v := range clients {
			if v.conn != conn && v.conn != nil {
				fmt.Println(1)
				_, err := v.conn.Write(data)

				if err != nil {
					fmt.Println("Error:", err.Error())
				}
			}
		}
	}
}

func newPlayer(clients map[int]client, id int) {
	for _, v := range clients {
		if v != clients[id] {
			data, err := proto.Encode(pNewPlayer, clients[id].uuid.String())
			if err != nil {
				fmt.Println("Error:", err.Error())
			}

			//Говорим всем, кто на сервере, что мы зашли
			_, err = v.conn.Write(data)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}

			//Подгружаем всех игроков, которые уже на сервере
			aboutPlayer, err := proto.Encode(pNewPlayer, v.uuid.String())
			_, err = clients[id].conn.Write(aboutPlayer)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
		}
	}
}

type client struct {
	uuid uuid.UUID
	conn net.Conn
}

const (
	pWelcome   int32 = 101
	pNewPlayer       = 102
	pMessage         = 103
	pMove            = 104
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	clients := make(map[int]client, 1024)
	i := 0

	fmt.Println("Server started...")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		clients[i] = client{uuid.New(), conn}
		fmt.Println("New Connection: " + conn.LocalAddr().String())
		go process(clients, conn, i)
		i++
	}
}
