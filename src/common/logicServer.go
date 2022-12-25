package common

import (
	"bufio"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"server/src/game"
)

const (
	UUIDPackage int32 = 101
	pNewPlayer        = 102
	pBroadcast        = 104
	pDisconnect       = 105
)

type LogicServer struct {
	Addr      string
	Clients   map[int]Client
	CurrentID int
}

func (s *LogicServer) Init(addr string, maxPlayers int) {
	s.Addr = addr
	s.Clients = make(map[int]Client, maxPlayers)
}

func (s *LogicServer) ListenAndServe() error {
	Listener, err := net.Listen("tcp", s.Addr+":30000")
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	log.Println("Logic server started ", s.Addr)

	for {
		conn, err := Listener.Accept()
		if err != nil {
			log.Fatal("accept failed, err:", err)
			continue
		}
		s.Clients[s.CurrentID] = Client{uuid.New(), conn, game.Player{}}
		log.Println("New Connection: " + conn.LocalAddr().String())
		go s.ConnectionHandler(s.CurrentID)
		s.CurrentID++
	}
}

func (s LogicServer) ConnectionHandler(id int) error {

	defer s.closeConnection(id)

	//Отправляем UUID игроку
	welcomeData, err := Encode(Package{UUIDPackage, "", s.Clients[id].Uuid.String()})
	log.Println("welcome: ", string(welcomeData))
	log.Println("Игроку присвоен uuid: " + s.Clients[id].Uuid.String())

	if err != nil {
		log.Fatal("Error:", err)
		return err
	}

	_, err = s.Clients[id].Conn.Write(welcomeData)
	if err != nil {
		log.Fatal("Error:", err)
		return err
	}

	err = s.NewPlayer(id)
	if err != nil {
		log.Fatal("Error:", err)
		return err
	}

	for {
		reader := bufio.NewReader(s.Clients[id].Conn)
		pack, err := Decode(reader)
		if err == io.EOF {
			log.Fatal("Error:", err)
			return err
		}
		if err != nil {
			log.Fatal("decode msg failed, err:", err)
			return err
		}

		err = s.handleReceivedPacket(pack, id)
		if err != nil {
			log.Fatal("Error:", err)
			return err
		}
	}
}

func (s LogicServer) handleReceivedPacket(pack Package, id int) error {
	if pack.Code == pBroadcast {
		s.broadCastWithout(pack, id)
	}
	return nil
}

func (s LogicServer) broadCast(pack Package) error {
	data, err := Encode(pack)
	if err != nil {
		log.Fatal("Error:", err)
		return err
	}
	for _, client := range s.Clients {
		if client.Conn != nil {
			_, err := client.Conn.Write(data)
			if err != nil {
				log.Fatal("Error:", err.Error())
				return err
			}
		}
	}
	return nil
}

func (s LogicServer) broadCastWithout(pack Package, id int) error {
	data, err := Encode(pack)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	for _, client := range s.Clients {
		if client.Conn != s.Clients[id].Conn && client.Conn != nil {
			_, err := client.Conn.Write(data)
			if err != nil {
				log.Println("Error:", err.Error())
				return err
			}
		}
	}
	return nil
}

func (s LogicServer) closeConnection(id int) {
	err := s.Clients[id].Conn.Close()
	if err != nil {
		log.Fatal("Error:", err)
	}
	s.broadCastWithout(Package{pDisconnect, "", s.Clients[id].Uuid.String()}, id)
	log.Println(s.Clients[id].Uuid.String() + " Отключился.")
	delete(s.Clients, id)
}

func (s LogicServer) NewPlayer(id int) error {
	for _, v := range s.Clients {
		if v != s.Clients[id] {
			data, err := Encode(Package{pNewPlayer, "", s.Clients[id].Uuid.String()})
			log.Println("Отсылаем ", v.Uuid.String(), " игроку пакет с новым игроком: ", s.Clients[id].Uuid.String())
			if err != nil {
				log.Println("Error:", err.Error())
			}

			//Говорим всем, кто на сервере, что мы зашли
			_, err = v.Conn.Write(data)
			if err != nil {
				log.Println("Error:", err.Error())
			}

			//Подгружаем всех игроков, которые уже на сервере
			aboutPlayer, err := Encode(Package{pNewPlayer, "", v.Uuid.String()})
			log.Println("Отсылаем ", s.Clients[id].Uuid.String(), " игроку пакет с другим игроком: ", v.Uuid.String())
			_, err = s.Clients[id].Conn.Write(aboutPlayer)
			if err != nil {
				log.Println("Error:", err.Error())
			}
		}
	}
	return nil
}
