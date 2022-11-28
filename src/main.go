package main

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"server/src/controller"
	"server/src/proto"
)

func process(clients map[int]client, conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		fmt.Println("Полученные данные от клиента:", msg)
		data, err := proto.Encode(msg)
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

type client struct {
	conn net.Conn
}

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
		clients[i] = client{conn}
		go process(clients, conn)
		i++
	}
}
