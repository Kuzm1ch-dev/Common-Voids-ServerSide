package main

import (
	"bufio"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"server/src/common"
	"server/src/proto"
)

const (
	UUIDPackage int32 = 101
	pNewPlayer        = 102
	pMessage          = 103
	pBroadcast        = 104
)

func process(clients map[int]common.Client, conn net.Conn, id int) {

	defer closeConnection(clients, id)

	//Отправляем UUID игроку
	welcomeData, err := proto.Encode(common.Package{UUIDPackage, "", clients[id].Uuid.String()})
	fmt.Println("welcome: ", string(welcomeData))
	fmt.Println("Игроку присвоен uuid: " + clients[id].Uuid.String())

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, e := conn.Write(welcomeData)
	if e != nil {
		fmt.Println("Error:", err)
		return
	}

	newPlayer(clients, id)

	for {
		reader := bufio.NewReader(conn)
		pack, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		//fmt.Println("Полученные данные от клиента:", string(pack.Marshal()))
		data, err := proto.Encode(pack)
		//m := controller.Message{pack, len(msg)}
		//go controller.Save(&m)

		if pack.Code == pBroadcast {
			fmt.Println(string(data))
			for _, v := range clients {
				if v.Conn != conn && v.Conn != nil {
					_, err := v.Conn.Write(data)
					if err != nil {
						fmt.Println("Error:", err.Error())
					}
				}
			}
		}
	}
}

func closeConnection(clients map[int]common.Client, id int) {
	fmt.Println(clients[id].Uuid.String() + " Отключился.")
	clients[id].Conn.Close()
	delete(clients, id)
}

func newPlayer(clients map[int]common.Client, id int) {
	for _, v := range clients {
		if v != clients[id] {
			data, err := proto.Encode(common.Package{pNewPlayer, "", clients[id].Uuid.String()})
			fmt.Println("Отсылаем ", v.Uuid.String(), " игроку пакет с новым игроком: ", clients[id].Uuid.String())
			if err != nil {
				fmt.Println("Error:", err.Error())
			}

			//Говорим всем, кто на сервере, что мы зашли
			_, err = v.Conn.Write(data)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}

			//Подгружаем всех игроков, которые уже на сервере
			aboutPlayer, err := proto.Encode(common.Package{pNewPlayer, "", v.Uuid.String()})
			fmt.Println("Отсылаем ", clients[id].Uuid.String(), " игроку пакет с другим игроком: ", v.Uuid.String())
			_, err = clients[id].Conn.Write(aboutPlayer)
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
		}
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		log.Println("No .env file found")
	}

	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	clients := make(map[int]common.Client, 1024)
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
		clients[i] = common.Client{uuid.New(), conn}
		fmt.Println("New Connection: " + conn.LocalAddr().String())
		go process(clients, conn, i)
		i++
	}
}
