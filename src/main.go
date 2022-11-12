package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"server/src/proto"
)

func process(conns map[int]net.Conn, conn net.Conn) {
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
		for _, v := range conns {
			if v != conn && conn != nil {
				fmt.Println(1)
				_, err := v.Write(data)

				if err != nil {
					fmt.Println("Error:", err.Error())
				}
			}
		}
	}
}

func main() {

	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	conns := make(map[int]net.Conn, 1024)
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
		conns[i] = conn
		go process(conns, conn)
		i++
	}
}
